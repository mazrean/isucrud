package analyze

import (
	"strings"
	"unicode"
)

const (
	initializeKeyword = "initialize"
)

func IsInitializeFuncName(name string) bool {
	words := camelCaseSplit(name)
	for _, word := range words {
		if strings.ToLower(word) == initializeKeyword {
			return true
		}
	}

	return false
}

func camelCaseSplit(s string) []string {
	var result []string
	start := 0
	for i, r := range s {
		if unicode.IsUpper(r) {
			result = append(result, s[start:i])
			start = i
		}
	}
	result = append(result, s[start:])

	return result
}
