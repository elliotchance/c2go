package ast

import (
	"strconv"
	"strings"
)

func ucfirst(word string) string {
	if len(word) == 0 {
		return ""
	}

	if len(word) == 1 {
		strings.ToUpper(word)
	}

	return strings.ToUpper(string(word[0])) + word[1:]
}

func removeQuotes(s string) string {
	s = strings.TrimSpace(s)

	if s == `""` {
		return ""
	}
	if s == `''` {
		return ""
	}

	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-2]
	}
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}

	return s
}

func atof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return f
}
