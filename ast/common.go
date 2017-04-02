package ast

import (
	"regexp"
	"strconv"
	"strings"
)

func Parse(line string) interface{} {
	nodeName := strings.SplitN(line, " ", 2)[0]
	var node interface{}

	switch nodeName {
	case "AlwaysInlineAttr":
		node = ParseAlwaysInlineAttr(line)
	case "ArraySubscriptExpr":
		node = ParseArraySubscriptExpr(line)
	case "AsmLabelAttr":
		node = ParseAsmLabelAttr(line)
	case "AvailabilityAttr":
		node = ParseAvailabilityAttr(line)
	case "BinaryOperator":
		node = ParseBinaryOperator(line)
	case "BreakStmt":
		node = ParseBreakStmt(line)
	default:
		panic(nodeName)
	}

	return node
}

func groupsFromRegex(rx, line string) map[string]string {
	fullRegexp := "(?P<address>[0-9a-fx]+) " + rx
	re := regexp.MustCompile(fullRegexp)

	match := re.FindStringSubmatch(line)
	if len(match) == 0 {
		panic("could not match regexp '" + fullRegexp +
			"' with string '" + line + "'")
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return result
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

	if len(s) >= 2 && s[0] == '"' && s[len(s) - 1] == '"' {
		return s[1:len(s) - 2]
	}

	return s
}
