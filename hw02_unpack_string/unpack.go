package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	var prevRune rune
	escaped := false

	for _, symbol := range input {
		switch {
		case unicode.IsDigit(symbol):
			if prevRune == 0 || escaped {
				return "", errors.New("invalid string")
			}
			count, _ := strconv.Atoi(string(symbol))
			result.WriteString(strings.Repeat(string(prevRune), count))
			prevRune = 0
			escaped = false

		case symbol == '\\' && !escaped:
			escaped = true

		default:
			if prevRune != 0 {
				result.WriteRune(prevRune)
			}
			prevRune = symbol
			escaped = false
		}
	}

	if prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}
