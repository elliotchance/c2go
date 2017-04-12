package ast

import "strings"

type FieldDecl struct {
	Address    string
	Position   string
	Position2  string
	Name       string
	Type       string
	Referenced bool
	Children []interface{}
}

func parseFieldDecl(line string) *FieldDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<position2> [^ ]+)?
		(?P<referenced> referenced)?
		 (?P<name>\w+?)
		 '(?P<type>.+?)'`,
		line,
	)

	return &FieldDecl{
		Address: groups["address"],
		Position: groups["position"],
		Position2: strings.TrimSpace(groups["position2"]),
		Name: groups["name"],
		Type: groups["type"],
		Referenced: len(groups["referenced"]) > 0,
		Children: []interface{}{},
	}
}
