package ast

import (
	"github.com/elliotchance/c2go/program"
)

type MallocAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseMallocAttr(line string) *MallocAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &MallocAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *MallocAttr) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *MallocAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
