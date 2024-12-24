package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func runeIsDigit(r rune) bool {
	if r >= 48 && r <= 57 {
		return true
	}

	return false
}

func Unpack(s string) (string, error) {
	srune := []rune(s) // srune LOL))

	lsrune := len(srune)

	if lsrune == 0 { // if string is empty, return string empty
		return "", nil
	}

	if lsrune == 1 {
		if runeIsDigit(srune[0]) { // if string is digit, return error
			return "", ErrInvalidString
		}
		return string(srune), nil
	}

	// if length of string is odd and last symbol is digit, return error, for example a23
	if lsrune%2 == 1 && runeIsDigit(srune[lsrune-1]) {
		return "", ErrInvalidString
	}

	var builder strings.Builder // for result
	builder.WriteString("")

	// check pairs of characters
	l := 0 // left index
	r := 1 // right index
	for r < len(srune) {
		if runeIsDigit(srune[l]) { // if find number > 9, for example ...45...
			return "", ErrInvalidString
		}

		if runeIsDigit(srune[r]) { // if find standard sequence, for example ...4b...
			builder.WriteString(strings.Repeat(string(srune[l]), int(srune[r])-48))
			l += 2 // move to the next pair
			r += 2
		} else { // if find a sequence of non-digites, for example ...aa...
			builder.WriteString(string(srune[l]))
			l++ // move to the next pair
			r++
		}
	}

	if l == len(srune)-1 { // if there is only one character left at the end
		builder.WriteString(string(srune[l]))
	}

	return builder.String(), nil
}
