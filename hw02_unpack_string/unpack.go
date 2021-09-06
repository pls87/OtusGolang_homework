package hw02unpackstring

import (
	"errors"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
	NotADigit        = errors.New("rune is not a digit")
)

const (
	zeroCode   int32 = 48
	escapeCode int32 = 92
)

func multiplyRune(b *strings.Builder, code int32, count int32) {
	var i int32
	for i = 0; i < count; i++ {
		b.WriteRune(code)
	}
}

func isRuneInteger(code int32) bool {
	return code >= zeroCode && code <= zeroCode+9
}

func runeDigitToInt32(code int32) (int32, error) {
	if !isRuneInteger(code) {
		return -1, NotADigit
	}
	return code - zeroCode, nil
}

func isRuneEscape(code int32) bool {
	return code == escapeCode
}

func Unpack(str string) (string, error) {
	builder := strings.Builder{}

	escaped := false
	var current int32 = -1

	for _, code := range str {
		switch {
		case !escaped && isRuneEscape(code):
			escaped = true
		case !escaped && isRuneInteger(code):
			if current == -1 {
				return "", ErrInvalidString
			}
			digit, _ := runeDigitToInt32(code)
			multiplyRune(&builder, current, digit)
			current = -1
		default:
			if current != -1 {
				builder.WriteRune(current)
			}
			current = code
			escaped = false
		}
	}
	// the last character if string isn't ended by digit or backslash
	if current != -1 && !escaped {
		builder.WriteRune(current)
	}
	// the last character is backslash
	if current == -1 && escaped {
		return "", ErrInvalidString
	}

	return builder.String(), nil
}
