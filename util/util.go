package util

import "strings"

func InStrings(item string, items []string) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}

	return false
}

func Ucfirst(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}
