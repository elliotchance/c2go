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
	// for.

	condition, _ := n.Children[2].render(program)

	n.Children[3].(*CompoundStmt).belongsToSwitch = true
	body, _ := n.Children[3].render(program)

	out := fmt.Sprintf("switch %s {\n%s\n}\n", condition, body)

	return out, ""
}

func (n *SwitchStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
