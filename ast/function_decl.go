package ast

import "strings"

type FunctionDecl struct {
	Address    string
	Position   string
	Prev       string
	Position2  string
	Name       string
	Type       string
	IsExtern   bool
	IsImplicit bool
	IsUsed     bool
	Children []interface{}
}

func parseFunctionDecl(line string) *FunctionDecl {
	groups := groupsFromRegex(
		"(?P<prev>prev [0-9a-fx]+ )?<(?P<position1>.*)>(?P<position2> [^ ]+)?(?P<implicit> implicit)?(?P<used> used)? (?P<name>\\w+) '(?P<type>.*)'(?P<extern> extern)?",
		line,
	)

	prev := groups["prev"]
	if prev != "" {
		prev = prev[5:len(prev) - 1]
	}

	return &FunctionDecl{
		Address: groups["address"],
		Position: groups["position1"],
		Prev: prev,
		Position2: strings.TrimSpace(groups["position2"]),
		Name: groups["name"],
		Type: groups["type"],
		IsExtern: len(groups["extern"]) > 0,
		IsImplicit: len(groups["implicit"]) > 0,
		IsUsed: len(groups["used"]) > 0,
		Children: []interface{}{},
	}
}
