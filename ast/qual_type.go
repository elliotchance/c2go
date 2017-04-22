package ast

import (
	"github.com/elliotchance/c2go/program"
)

type QualType struct {
	Address  string
	Type     string
	Kind     string
	Children []Node
}

func parseQualType(line string) *QualType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &QualType{
		Address:  groups["address"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

func (n *QualType) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *QualType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
