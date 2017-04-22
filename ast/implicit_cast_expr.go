package ast

import (
	"github.com/elliotchance/c2go/program"
)

type ImplicitCastExpr struct {
	Address  string
	Position string
	Type     string
	Kind     string
	Children []Node
}

func parseImplicitCastExpr(line string) *ImplicitCastExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' <(?P<kind>.*)>",
		line,
	)

	return &ImplicitCastExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

func (n *ImplicitCastExpr) render(program *program.Program) (string, string) {
	return renderExpression(program, n.Children[0])
}

func (n *ImplicitCastExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
