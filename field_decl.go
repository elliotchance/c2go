package main

import (
	"strings"
	"fmt"
)

type FieldDecl struct {
	Address    string
	Position   string
	Position2  string
	Name       string
	Type       string
	Referenced bool
	Children   []interface{}
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
		Address:    groups["address"],
		Position:   groups["position"],
		Position2:  strings.TrimSpace(groups["position2"]),
		Name:       groups["name"],
		Type:       groups["type"],
		Referenced: len(groups["referenced"]) > 0,
		Children:   []interface{}{},
	}
}

func (n *FieldDecl) Render() []string {
	fieldType := resolveType(n.Type)
	name := strings.Replace(n.Name, "used", "", -1)

	// Go does not allow the name of a variable to be called "type". For the
	// moment I will rename this to avoid the error.
	if name == "type" {
		name = "type_"
	}

	suffix := ""
	if len(n.Children) > 0 {
		suffix = fmt.Sprintf(" = %s", renderExpression(n.Children[0])[0])
	}

	if suffix == " = (0)" {
		suffix = " = nil"
	}

	return []string{fmt.Sprintf("%s %s%s", name, fieldType, suffix), "unknown3"}
}
