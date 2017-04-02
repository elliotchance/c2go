package ast

import (
	"regexp"
)

type ArraySubscriptExpr struct {
	Address  string
	Position string
	Type     string
	Tags     string
}

func ParseArraySubscriptExpr(line string) ArraySubscriptExpr {
	re := regexp.MustCompile("(?P<address>[0-9a-fx]+) <(?P<position>.*)> '(?P<type>.*?)' (?P<tags>.*)")
	match := re.FindStringSubmatch(line)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return ArraySubscriptExpr{
		Address: result["address"],
		Position: result["position"],
		Type: result["type"],
		Tags: result["tags"],
	}
}
