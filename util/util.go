package util

import (
	"bufio"
	"fmt"
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
	// Convert "[]byte" into "byteSlice". This also works with multiple slices,
	// like "[][]byte" to "byteSliceSlice".
	for field[:2] == "[]" {
		field = field[2:] + "Slice"
	}

	return Ucfirst(strings.TrimLeft(field, "*_"))
}

// StringToLines - convert string to string lines
func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return lines, fmt.Errorf("reading standard input: %v", err)
	}

	return lines, nil
}
