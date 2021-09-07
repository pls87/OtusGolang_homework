package hw02unpackstring

import (
	"errors"
	"strings"
)

var (
	ErrInvalidString   = errors.New("invalid string")
	ErrRuneIsNotADigit = errors.New("rune is not a digit")
)

const (
	zeroCode   rune = 48
	escapeCode rune = 92
)

func writeRunes2Builder(b *strings.Builder, code rune, count int32) {
	var i int32
	for i = 0; i < count; i++ {
		b.WriteRune(code)
	}
}

func isRuneInteger(code rune) bool {
	return code >= zeroCode && code <= zeroCode+9
}

func runeDigit2Int32(code rune) (int32, error) {
	if !isRuneInteger(code) {
		return -1, ErrRuneIsNotADigit
	}
	return code - zeroCode, nil
}

func isRuneEscape(code rune) bool {
	return code == escapeCode
}

func Unpack(str string) (string, error) {
	builder := strings.Builder{}

	escaped := false
	var prev rune = -1

	for _, code := range str {
		switch {
		case !escaped && isRuneEscape(code):
			escaped = true
		case !escaped && isRuneInteger(code):
			if prev == -1 {
				return "", ErrInvalidString
			}
			digit, _ := runeDigit2Int32(code)
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
