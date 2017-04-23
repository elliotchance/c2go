package ast

import (
	"github.com/elliotchance/c2go/program"
)

type TransparentUnionAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseTransparentUnionAttr(line string) *TransparentUnionAttr {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &TransparentUnionAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *TransparentUnionAttr) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *TransparentUnionAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
