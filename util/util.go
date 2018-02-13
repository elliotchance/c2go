package util

import (
	"strconv"
	"strings"
)

// InStrings returns true if item exists in items. It must be an exact string
// match.
func InStrings(item string, items []string) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}

	return false
}

// Ucfirst returns the word with the first letter uppercased; none of the other
// letters in the word are modified. For example "fooBar" would return "FooBar".
func Ucfirst(word string) string {
	if word == "" {
		return ""
	}

	if len(word) == 1 {
		return strings.ToUpper(word)
	}

	return strings.ToUpper(string(word[0])) + word[1:]
}

// Atoi converts a string to an integer in cases where we are sure that s will
// be a valid integer, otherwise it will panic.
func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	PanicOnError(err, "bad integer")

	return i
}

// GetExportedName returns a deterministic and Go safe name for a C type. For
// example, "*__foo[]" will return "FooSlice".
func GetExportedName(field string) string {
	if strings.Contains(field, "interface{}") ||
		strings.Contains(field, "Interface{}") {
		return "Interface"
	}

	// Convert "[]byte" into "byteSlice". This also works with multiple slices,
	// like "[][]byte" to "byteSliceSlice".
	for len(field) > 2 && field[:2] == "[]" {
		field = field[2:] + "Slice"
	}

	// NotFunc(int)()
	field = strings.Replace(field, "(", "_", -1)
	field = strings.Replace(field, ")", "_", -1)

	return Ucfirst(strings.TrimLeft(field, "*_"))
}
