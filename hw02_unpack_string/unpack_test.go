package hw02unpackstring_test

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
	"unicode"

	hw02unpackstring "github.com/pls87/OtusGolang_homework/hw02_unpack_string"
	"github.com/stretchr/testify/require"
)

const (
	zeroCode   = 48
	escapeCode = 92
)

func seedRandom() int64 {
	seed := time.Now().Unix()
	rand.Seed(seed)

	return seed
}

func randBool() bool {
	return rand.Int31n(2) == 1
}

func randDigit() int32 {
	return rand.Int31n(10)
}

func generateRandomCase(negative bool) (input string, expected string, seedUsed int64) {
	seedUsed = seedRandom()

	var i int32
	inputBuilder, expectedBuilder := strings.Builder{}, strings.Builder{}
	length := rand.Int31n(16) + 1
	alreadyNegative := false
	for i = 0; i < length; i++ {
		code := rand.Int31n(255)
		count := randDigit()
		negativeStep := false
		// if it's asked to generate negative case then extra digit will be written 50% of cases
		if negative && ((!alreadyNegative && i == length-1) || randBool()) {
			inputBuilder.WriteRune(randDigit() + zeroCode)
			negativeStep = true
			alreadyNegative = true
		}
		// digit and backslash should be escaped
		if code == escapeCode || unicode.IsDigit(code) {
			inputBuilder.WriteRune(escapeCode)
		}
		inputBuilder.WriteRune(code)

		// digit is written if count != 1 or if wrong sequence is generated on this step. Otherwise, digit is optional
		if count != 1 || negativeStep || randBool() {
			inputBuilder.WriteRune(count + zeroCode)
		}
		// expected value is generated just for positive case
		if !negative {
			var i int32
			for i = 0; i < count; i++ {
				expectedBuilder.WriteRune(code)
			}
		}
	}
	return inputBuilder.String(), expectedBuilder.String(), seedUsed
}

func TestUnpack(t *testing.T) {
	randomInput, randomExpected, seed := generateRandomCase(false)
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		// positive cases
		{input: "Hh0el2oo0, goo0langg0!3", expected: "Hello, golang!!!"},
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "a0b0c0d0", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: " 5 6", expected: "           "},
		{
			input:    randomInput,
			expected: randomExpected,
			name:     fmt.Sprintf("Random positive case with input:%s and seed:%d", randomInput, seed),
		},
		// cases with escape
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\1\33\52\84\70`, expected: `1333558888`},
		// corner cases
		{input: "", expected: ""},
		{input: `\6`, expected: `6`},
		{input: `\\`, expected: `\`},
		{input: `\a`, expected: `a`},
	}

	for _, tc := range tests {
		tc := tc
		if tc.name == "" {
			tc.name = fmt.Sprintf("Positive case with input:%s", tc.input)
		}
		t.Run(tc.name, func(t *testing.T) {
			result, err := hw02unpackstring.Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	randomInput, _, seed := generateRandomCase(true)
	invalidCases := []struct {
		input string
		name  string
	}{
		{input: `3abc`},
		{input: `45`},
		{input: `aaa10b`},
		{input: `7`},
		{input: `\`},
		{input: `ab5\`},
		{input: `ab\\45`},
		{
			input: `ab\\45`,
			name:  fmt.Sprintf("Random negative case with input:%s and seed:%d", randomInput, seed),
		},
	}

	for _, tc := range invalidCases {
		tc := tc
		if tc.name == "" {
			tc.name = fmt.Sprintf("Negative case with input:%s", tc.input)
		}
		t.Run(tc.name, func(t *testing.T) {
			_, err := hw02unpackstring.Unpack(tc.input)
			require.Truef(t, errors.Is(err, hw02unpackstring.ErrInvalidString), "actual error %q", err)
		})
	}
}
