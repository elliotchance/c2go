package ast

import (
	"fmt"
	"strings"
)

type CallExpr struct {
	Address  string
	Position string
	Type     string
	Children []Node
}

func parseCallExpr(line string) *CallExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)'",
		line,
	)

	return &CallExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *CallExpr) render(ast *Ast) (string, string) {
	children := n.Children
	func_name, _ := renderExpression(ast, children[0])

	func_def := getFunctionDefinition(func_name)

	if func_def.Substitution != "" {
		parts := strings.Split(func_def.Substitution, ".")
		ast.addImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(func_def.Substitution, "/")
		func_name = parts2[len(parts2)-1]
	}

	args := []string{}
	i := 0
	for _, arg := range children[1:] {
		e, eType := renderExpression(ast, arg)

		if i > len(func_def.ArgumentTypes)-1 {
			// This means the argument is one of the varargs
			// so we don't know what type it needs to be
			// cast to.
			args = append(args, e)
		} else {
			args = append(args, cast(ast, e, eType, func_def.ArgumentTypes[i]))
		}

		i++
	}

	parts := []string{}

	for _, v := range args {
		parts = append(parts, v)
	}

	src := fmt.Sprintf("%s(%s)", func_name, strings.Join(parts, ", "))
	return src, func_def.ReturnType
}

func (n *CallExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
