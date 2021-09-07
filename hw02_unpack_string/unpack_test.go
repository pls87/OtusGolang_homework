package hw02unpackstring

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func seedRandom() {
	rand.Seed(time.Now().Unix())
}

func generateRandomPositiveCase() (input string, expected string) {
	seedRandom()

	var i int32
	inputBuilder, expectedBuilder := strings.Builder{}, strings.Builder{}
	for i = 0; i < rand.Int31n(64); i++ {
		code := rand.Int31n(math.MaxInt8-58) + 58
		count := rand.Int31n(10)
		codeCount := count + 48
		inputBuilder.WriteRune(code)
		inputBuilder.WriteRune(codeCount)
		var j int32
		for j = 0; j < count; j++ {
			expectedBuilder.WriteRune(code)
		}
	}
	return inputBuilder.String(), expectedBuilder.String()
}

func TestWriteRunes2Builder(t *testing.T) {
	cases := []struct {
		code     rune
		count    int32
		expected string
	}{
		{code: 65, count: 5, expected: "AAAAA"},
		{code: 73, count: 9, expected: "IIIIIIIII"},
		{code: 70, count: 0, expected: ""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(strconv.QuoteRune(tc.code), func(t *testing.T) {
			builder := strings.Builder{}
			writeRunes2Builder(&builder, tc.code, tc.count)
			require.Equal(t, tc.expected, builder.String())
		})
	}
}

func TestIsRuneEscape(t *testing.T) {
	seedRandom()

	cases := []struct {
		code     rune
		expected bool
	}{
		{code: 92, expected: true},
		{code: rand.Int31n(92), expected: false},
		{code: rand.Int31n(math.MaxInt8-93) + 93, expected: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(strconv.QuoteRune(tc.code), func(t *testing.T) {
			require.Equal(t, tc.expected, isRuneEscape(tc.code))
		})
	}
}

func TestRuneDigit2Int32(t *testing.T) {
	seedRandom()

	cases := []struct {
		code     rune
		expected int32
		success  bool
	}{
		{code: 48, expected: 0, success: true},
		{code: 57, expected: 9, success: true},
		{code: 55, expected: 7, success: true},
		{code: 47, expected: -1, success: false},
		{code: 58, expected: -1, success: false},
		{code: rand.Int31n(48), expected: -1, success: false},
		{code: rand.Int31n(math.MaxInt8-58) + 58, expected: -1, success: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(strconv.QuoteRune(tc.code), func(t *testing.T) {
			result, err := runeDigit2Int32(tc.code)
			require.Equal(t, tc.success, err == nil)
			if tc.success {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			} else {
				require.Truef(t, errors.Is(err, ErrRuneIsNotADigit), "actual error %q", err)
			}
		})
	}
}

func TestUnpack(t *testing.T) {
	randomInput, randomExpected := generateRandomPositiveCase()
	tests := []struct {
		input    string
		expected string
	}{
		// positive cases
		{input: "Hh0el2oo0, goo0langg0!3", expected: "Hello, golang!!!"},
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "a0b0c0d0", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: " 5 6", expected: "           "},
		{input: randomInput, expected: randomExpected},
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
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{`3abc`, `45`, `aaa10b`, `7`, `\`, `ab5\`, `ab\\45`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
