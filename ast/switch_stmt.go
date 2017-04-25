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
	// for. It looks like the number of children is also not reliable, but we
	// know that we need the last two.

	if len(n.Children) < 2 {
		// I don't know what causes this condition. Need to investigate.
		return "", ""
	}

	condition, _ := n.Children[len(n.Children)-2].render(program)

	n.Children[3].(*CompoundStmt).belongsToSwitch = true
	body, _ := n.Children[len(n.Children)-1].render(program)

	out := fmt.Sprintf("switch %s {\n%s\n}\n", condition, body)

	return out, ""
}

func (n *SwitchStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
