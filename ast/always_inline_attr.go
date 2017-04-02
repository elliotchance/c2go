package ast

import (
	"regexp"
)

type AlwaysInlineAttr struct {
	Address  string
	Position string
}

func ParseAlwaysInlineAttr(line string) AlwaysInlineAttr {
	re := regexp.MustCompile("(?P<address>[0-9a-fx]+) <(?P<position>.*)> always_inline")
	match := re.FindStringSubmatch(line)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return AlwaysInlineAttr{
		Address: result["address"],
		Position: result["position"],
	}
}
