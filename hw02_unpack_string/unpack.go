package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString - Error for invalid input string.
var ErrInvalidString = errors.New("invalid string")

type CharInfo struct {
	str          string
	isDigit      bool
	originalRune rune
}

func (c *CharInfo) getInt() (int, error) {
	digit, err := strconv.Atoi(c.str)

	if err != nil {
		return -1, err
	}

	return digit, nil
}

// Unpack - Function for string unpacking
// Use it to convert string like "ab4cd2" to "abbbbcdd".
func Unpack(str string) (string, error) {
	// If input string is empty, immediately return it
	if str == "" {
		return "", nil
	}
	var b strings.Builder
	// Buffer for storing non-digit chunk of string - it will be used to properly populate the result string
	charInfoBuffer := make([]*CharInfo, 0, len(str))

	var prevCharInfo *CharInfo
	// Iterating over the input string
	for _, char := range str {
		if charInfoBuffer == nil {
			charInfoBuffer = make([]*CharInfo, 0, len(str))
		}

		charInfo := CharInfo{
			originalRune: char,
			str:          string(char),
			isDigit:      unicode.IsDigit(char),
		}

		// If first char is digit, then it is invalid string
		if prevCharInfo == nil && charInfo.isDigit {
			return "", ErrInvalidString
		}

		// If there are two digits in a row, then it is invalid string
		if prevCharInfo != nil && charInfo.isDigit && prevCharInfo.isDigit {
			return "", ErrInvalidString
		}

		// If current char is a letter, then append it to buffer
		if !charInfo.isDigit {
			charInfoBuffer = append(charInfoBuffer, &charInfo)
		} else {
			// If current char is digit and previous one is not, then write char buffer to result string
			// and, if digit > 0, multiply the last char of the buffer
			charInfoBufferExcludingLastChar := charInfoBuffer[:len(charInfoBuffer)-1]
			lastBufferChar := charInfoBuffer[len(charInfoBuffer)-1]
			for _, ch := range charInfoBufferExcludingLastChar {
				b.WriteString(ch.str)
			}
			digit, _ := charInfo.getInt()
			if digit > 0 {
				b.WriteString(strings.Repeat(lastBufferChar.str, digit))
			}
			// Clean buffer
			charInfoBuffer = nil
		}

		prevCharInfo = &charInfo
	}

	// If buffer is not empty, the write it all into the result string
	if len(charInfoBuffer) > 0 {
		for _, ch := range charInfoBuffer {
			b.WriteString(ch.str)
		}
	}
	return b.String(), nil
}
