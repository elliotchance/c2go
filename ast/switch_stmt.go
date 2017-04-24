package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
)

type SwitchStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseSwitchStmt(line string) *SwitchStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &SwitchStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *SwitchStmt) render(program *program.Program) (string, string) {
	// The first two children are nil. I don't know what they are supposed to be
	// for. Raise an error so we can work it out.
	if n.Children[0] != nil {
		panic(n.Children[0])
	}
	if n.Children[1] != nil {
		panic(n.Children[1])
	}

	condition, _ := n.Children[2].render(program)

	n.Children[3].(*CompoundStmt).belongsToSwitch = true
	body, _ := n.Children[3].render(program)

	out := fmt.Sprintf("switch %s {\n%s\n}\n", condition, body)

	return out, ""
}

func (n *SwitchStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
