package main

import (
	"strings"
	"fmt"
)

type CallExpr struct {
	Address  string
	Position string
	Type     string
	Children []interface{}
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
		Children: []interface{}{},
	}
}

func (n *CallExpr) Render() []string {
	children := n.Children
	func_name := renderExpression(children[0])[0]

	func_def := getFunctionDefinition(func_name)

	if func_def.Substitution != "" {
		parts := strings.Split(func_def.Substitution, ".")
		addImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(func_def.Substitution, "/")
		func_name = parts2[len(parts2)-1]
	}

	args := []string{}
	i := 0
	for _, arg := range children[1:] {
		e := renderExpression(arg)

		if i > len(func_def.ArgumentTypes)-1 {
			// This means the argument is one of the varargs
			// so we don't know what type it needs to be
			// cast to.
			args = append(args, e[0])
		} else {
			args = append(args, cast(e[0], e[1], func_def.ArgumentTypes[i]))
		}

		i += 1
	}

	parts := []string{}

	for _, v := range args {
		parts = append(parts, v)
	}

	return []string{
		fmt.Sprintf("%s(%s)", func_name, strings.Join(parts, ", ")),
		func_def.ReturnType}
}
