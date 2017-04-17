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
	functionName, _ := renderExpression(ast, children[0])
	functionDef := getFunctionDefinition(functionName)

	if functionDef == nil {
		panic(fmt.Sprintf("unknown function: %s", functionName))
	}

	if functionDef.Substitution != "" {
		parts := strings.Split(functionDef.Substitution, ".")
		ast.addImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(functionDef.Substitution, "/")
		functionName = parts2[len(parts2)-1]
	}

	args := []string{}
	i := 0
	for _, arg := range children[1:] {
		e, eType := renderExpression(ast, arg)

		if i > len(functionDef.ArgumentTypes)-1 {
			// This means the argument is one of the varargs
			// so we don't know what type it needs to be
			// cast to.
			args = append(args, e)
		} else {
			args = append(args, cast(ast, e, eType, functionDef.ArgumentTypes[i]))
		}

		i++
	}

	parts := []string{}

	for _, v := range args {
		parts = append(parts, v)
	}

	src := fmt.Sprintf("%s(%s)", functionName, strings.Join(parts, ", "))
	return src, functionDef.ReturnType
}

func (n *CallExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
