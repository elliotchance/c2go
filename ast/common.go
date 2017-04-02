package ast

import "regexp"

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
