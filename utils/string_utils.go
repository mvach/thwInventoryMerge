package utils

import (
	"strconv"
	"unicode"
)

func StartsWithNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsDigit(rune(s[0]))
}

func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}