package main

import "strings"

func ucfirst(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}

func getExportedName(field string) string {
	return ucfirst(strings.TrimLeft(field, "_"))
}

func inStrings(item string, items []string) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}

	return false
}
