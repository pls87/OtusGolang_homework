package hw02unpackstring

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func seedRandom() {
	rand.Seed(42)
}

func TestMultiplyRune(t *testing.T) {
	cases := []struct {
		rune           int32
		count          int32
		expectedResult string
	}{
		{rune: 65, count: 5, expectedResult: "AAAAA"},
		{rune: 73, count: 9, expectedResult: "IIIIIIIII"},
		{rune: 70, count: 0, expectedResult: ""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(strconv.QuoteRune(tc.rune), func(t *testing.T) {
			builder := strings.Builder{}
			multiplyRune(&builder, tc.rune, tc.count)
			require.Equal(t, tc.expectedResult, builder.String())
		})
	}
}

func TestIsRuneEscape(t *testing.T) {
	seedRandom()

	cases := []struct {
		rune           int32
		expectedResult bool
	}{
		{rune: 92, expectedResult: true},
		{rune: rand.Int31n(92), expectedResult: false},
		{rune: rand.Int31n(math.MaxInt32-93) + 93, expectedResult: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(strconv.QuoteRune(tc.rune), func(t *testing.T) {
			require.Equal(t, tc.expectedResult, isRuneEscape(tc.rune))
		})
	}
}

func TestRuneDigitToInt32(t *testing.T) {
	seedRandom()
	cases := []struct {
		rune           int32
		expectedResult int32
		success        bool
	}{
		{rune: 48, expectedResult: 0, success: true},
		{rune: 57, expectedResult: 9, success: true},
		{rune: 55, expectedResult: 7, success: true},
		{rune: 47, expectedResult: -1, success: false},
		{rune: 58, expectedResult: -1, success: false},
		{rune: rand.Int31n(48), expectedResult: -1, success: false},
		{rune: rand.Int31n(math.MaxInt32-58) + 58, expectedResult: -1, success: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(strconv.QuoteRune(tc.rune), func(t *testing.T) {
			result, err := runeDigitToInt32(tc.rune)
			require.Equal(t, tc.success, err == nil)
			if tc.success {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			} else {
				require.Truef(t, errors.Is(err, ErrRuneIsNotADigit), "actual error %q", err)
			}
		})
	}
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		// corner cases
		{input: ``, expected: ``},
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
	invalidStrings := []string{`3abc`, `45`, `aaa10b`, `7`, `\`, `ab5\`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
