package util

import (
	"strconv"
	"strings"
)

func InStrings(item string, items []string) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}

	return false
}

func Ucfirst(word string) string {
	if word == "" {
		return ""
	}

	return strings.ToUpper(string(word[0])) + word[1:]
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

func GetExportedName(field string) string {
	return Ucfirst(strings.TrimLeft(field, "*_"))
}
