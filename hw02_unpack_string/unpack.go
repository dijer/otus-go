package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString  = errors.New("invalid string")
	ErrInvalidConvert = errors.New("invalid convert string to digit")
)

func Unpack(input string) (string, error) {
	var output strings.Builder
	inputRunes := []rune(input)

	for i, currentChar := range inputRunes {
		isCurrentCharDigit := unicode.IsDigit(currentChar)
		nextChar := ""
		isNextCharDigit := false

		if len(input) > i+1 {
			nextChar = string(input[i+1])
			isNextCharDigit = unicode.IsDigit(inputRunes[i+1])
		}

		if isCurrentCharDigit && (i == 0 || isNextCharDigit) {
			return "", ErrInvalidString
		}

		if isCurrentCharDigit {
			continue
		}

		if isNextCharDigit {
			if nextCharNumber, err := strconv.Atoi(nextChar); err == nil {
				output.WriteString(strings.Repeat(string(currentChar), nextCharNumber))
				continue
			} else if err != nil {
				return "", ErrInvalidConvert
			}
		}

		output.WriteString(string(currentChar))
	}

	return output.String(), nil
}
