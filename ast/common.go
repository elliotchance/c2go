package ast

import (
	"regexp"
	"strconv"
	"strings"
)

func groupsFromRegex(rx, line string) map[string]string {
	re := regexp.MustCompile("(?P<address>[0-9a-fx]+) " + rx)
	match := re.FindStringSubmatch(line)
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
