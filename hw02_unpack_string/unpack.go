package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	var prevRune rune
	escaped := false

	for _, symbol := range input {
		switch {
		case escaped:
			if symbol != '\\' && (symbol < '0' || symbol > '9') {
				return "", ErrInvalidString
			}
			prevRune = symbol
			escaped = false

		case symbol == '\\':
			escaped = true
			if prevRune != 0 {
				result.WriteRune(prevRune)
			}
			prevRune = 0

		default:
			count, err := strconv.Atoi(string(symbol))
			if err == nil {
				if prevRune == 0 {
					return "", ErrInvalidString
				}
				result.WriteString(strings.Repeat(string(prevRune), count))
				prevRune = 0
			} else {
				if prevRune != 0 {
					result.WriteRune(prevRune)
				}
				prevRune = symbol
			}
		}
	}
	if prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}
