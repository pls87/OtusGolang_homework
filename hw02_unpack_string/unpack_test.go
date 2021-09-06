package hw02unpackstring

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func seedRandom() {
	rand.Seed(42)
}

func TestIsRuneInteger(t *testing.T) {
	seedRandom()

	cases := []struct {
		rune           int32
		expectedResult bool
	}{
		{rune: 48, expectedResult: true}, // 0
		{rune: 57, expectedResult: true}, // 9
		{rune: 55, expectedResult: true}, // 7
		{rune: 47, expectedResult: false},
		{rune: rand.Int31n(48), expectedResult: false},
		{rune: rand.Int31n(math.MaxInt32-58) + 58, expectedResult: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(strconv.QuoteRune(tc.rune), func(t *testing.T) {
			require.Equal(t, tc.expectedResult, isRuneInteger(tc.rune))
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
