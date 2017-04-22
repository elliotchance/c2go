package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
)

type ArraySubscriptExpr struct {
	Address  string
	Position string
	Type     string
	Kind     string
	Children []Node
}

func parseArraySubscriptExpr(line string) *ArraySubscriptExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<kind>.*)",
		line,
	)

	return &ArraySubscriptExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

func (n *ArraySubscriptExpr) render(program *program.Program) (string, string) {
	children := n.Children
	expression, _ := renderExpression(program, children[0])
	index, _ := renderExpression(program, children[1])
	src := fmt.Sprintf("%s[%s]", expression, index)
	return src, "unknown1"
}

func (n *ArraySubscriptExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
