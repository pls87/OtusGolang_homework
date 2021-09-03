package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

const zeroCode int32 = 48
const escapeCode = 92

func multiplyRune(b *strings.Builder, code int32, count int32) {
	var i int32
	for i = 0; i < count; i++ {
		b.WriteRune(code)
	}
}

func isRuneInteger(code int32) bool {
	return code >= zeroCode && code < zeroCode+10
}

func runeToInt32(code int32) int32 {
	return code - zeroCode
}

func isRuneEscape(code int32) bool {
	return code == escapeCode
}

func Unpack(str string) (string, error) {
	builder := strings.Builder{}
	var current int32 = -1
	var escaped bool = false

	for _, code := range str {
		if isRuneInteger(code) && !escaped {
			// if string starts by digit or two digits-neighbors were found then return error
			if current == -1 {
				return "", ErrInvalidString
			}
			multiplyRune(&builder, current, runeToInt32(code))
			current = -1
		} else if isRuneEscape(code) && !escaped {
			escaped = true
		} else {
			if current != -1 {
				builder.WriteRune(current)
			}
			current = code
			escaped = false
		}
	}
	// the last character if string isn't ended by digit
	if current != -1 {
		builder.WriteRune(current)
	}

	return builder.String(), nil
}
