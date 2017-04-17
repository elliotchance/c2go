package ast

import "fmt"

type PredefinedExpr struct {
	Address  string
	Position string
	Type     string
	Name     string
	Lvalue   bool
	Children []Node
}

func parsePredefinedExpr(line string) *PredefinedExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' lvalue (?P<name>.*)",
		line,
	)

	return &PredefinedExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Name:     groups["name"],
		Lvalue:   true,
		Children: []Node{},
	}
}

func (n *PredefinedExpr) render(ast *Ast) (string, string) {
	if n.Name == "__PRETTY_FUNCTION__" {
		// FIXME
		return "\"void print_number(int *)\"", "const char*"
	}

	if n.Name == "__func__" {
		// FIXME
		src := fmt.Sprintf("\"%s\"", "print_number")
		return src, "const char*"
	}

	panic(fmt.Sprintf("renderExpression: unknown PredefinedExpr: %s", n.Name))
}

func (n *PredefinedExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
