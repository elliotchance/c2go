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
		(?P<position2> col:\d+| line:\d+:\d+)?
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

	// FIXME: There are some cases where the name is empty. We need to
	// investigate this further. For now I will just exclude them.
	if name == "" {
		return []string{"", "unknown72"}
	}

	// Go does not allow the name of a variable to be called "type". For the
	// moment I will rename this to avoid the error.
	if name == "type" {
		name = "type_"
	}

	return []string{fmt.Sprintf("%s %s", name, fieldType), "unknown3"}
}
