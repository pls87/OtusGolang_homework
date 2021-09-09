package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func writeRunes2Builder(b *strings.Builder, code rune, count int32) {
	var i int32
	for i = 0; i < count; i++ {
		b.WriteRune(code)
	}
}

func Unpack(str string) (string, error) {
	builder := strings.Builder{}

	escaped := false
	var prev rune = -1

	for _, code := range str {
		switch {
		case !escaped && code == '\\':
			escaped = true
		case !escaped && unicode.IsDigit(code):
			if prev == -1 {
				return "", ErrInvalidString
			}
			digit := code - '0'
			writeRunes2Builder(&builder, prev, digit)
			prev = -1
		default:
			if prev != -1 {
				builder.WriteRune(prev)
			}
			prev = code
			escaped = false
		}
	}
	// the last character is backslash
	if prev == -1 && escaped {
		return "", ErrInvalidString
	}

	// the last character if string isn't ended by digit or backslash
	if prev != -1 && !escaped {
		builder.WriteRune(prev)
	}

	return builder.String(), nil
}
