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
		case escaped:
			if symbol != '\\' && !unicode.IsDigit(symbol) {
				return "", ErrInvalidString
			}
			prevRune = symbol
			escaped = false

		case symbol == '\\':
			escaped = true
			result.WriteRune(prevRune)

		case unicode.IsDigit(symbol):
			if prevRune == 0 {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(symbol))
			result.WriteString(strings.Repeat(string(prevRune), count))
			prevRune = 0

		default:
			if prevRune != 0 {
				result.WriteRune(prevRune)
			}
			prevRune = symbol
		}
	}
	if prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}
