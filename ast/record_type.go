package ast

import (
	"github.com/elliotchance/c2go/program"
)

type RecordType struct {
	Address  string
	Type     string
	Children []Node
}

func parseRecordType(line string) *RecordType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &RecordType{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *RecordType) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *RecordType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
