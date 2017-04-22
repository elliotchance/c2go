package ast

import (
	"strconv"
	"strings"

	"github.com/elliotchance/c2go/util"
)

func getExportedName(field string) string {
	return util.Ucfirst(strings.TrimLeft(field, "_"))
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
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

func unescapeString(s string) string {
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	s = strings.Replace(s, "\\t", "\t", -1)

	return s
}
