package utils

import "strings"

func LowerFirstLetter(text string) string {
	firstLetter := string(text[0])
	return strings.ToLower(firstLetter) + text[1:]
}

func UpperFirstLetter(text string) string {
	firstLetter := string(text[0])
	return strings.ToUpper(firstLetter) + text[1:]
}
