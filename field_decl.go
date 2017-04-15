package main

import (
	"fmt"
	"strings"
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
		(?P<name> \w+?)?
		 '(?P<type>.+?)'`,
		line,
	)

	return &FieldDecl{
		Address:    groups["address"],
		Position:   groups["position"],
		Position2:  strings.TrimSpace(groups["position2"]),
		Name:       strings.TrimSpace(groups["name"]),
		Type:       groups["type"],
		Referenced: len(groups["referenced"]) > 0,
		Children:   []interface{}{},
	}
}

func (n *FieldDecl) Render() []string {
	fieldType := resolveType(n.Type)
	name := n.Name

	if name == "" {
		return []string{"", "unknown71"}
	}

	// Go does not allow the name of a variable to be called "type". For the
	// moment I will rename this to avoid the error.
	if name == "type" {
		name = "type_"
	}

	// It may have a default value.
	suffix := ""
	if len(n.Children) > 0 {
		suffix = fmt.Sprintf(" = %s", renderExpression(n.Children[0])[0])
	}

	// NULL is a macro that one rendered looks like "(0)" we have to be
	// sensitive to catch this as Go would complain that 0 (int) is not
	// compatible with the type we are setting it to.
	if suffix == " = (0)" {
		suffix = " = nil"
	}

	return []string{fmt.Sprintf("%s %s%s", name, fieldType, suffix), "unknown3"}
}
