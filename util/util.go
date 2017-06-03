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

	if len(word) == 1 {
		strings.ToUpper(word)
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
	// Convert "[]byte" into "byteSlice". This also works with multiple slices,
	// like "[][]byte" to "byteSliceSlice".
	for field[:2] == "[]" {
		field = field[2:] + "Slice"
	}

	return Ucfirst(strings.TrimLeft(field, "*_"))
}
