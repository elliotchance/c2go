package main

import (
	"fmt"
	"strings"
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
	functionName := renderExpression(children[0])[0]
	functionDef := getFunctionDefinition(functionName)

	if functionDef == nil {
		panic(fmt.Sprintf("unknown function: %s", functionName))
	}

	if functionDef.Substitution != "" {
		parts := strings.Split(functionDef.Substitution, ".")
		addImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(functionDef.Substitution, "/")
		functionName = parts2[len(parts2)-1]
	}

	args := []string{}
	i := 0
	for _, arg := range children[1:] {
		e := renderExpression(arg)

		if i > len(functionDef.ArgumentTypes)-1 {
			// This means the argument is one of the varargs
			// so we don't know what type it needs to be
			// cast to.
			args = append(args, e[0])
		} else {
			args = append(args, cast(e[0], e[1], functionDef.ArgumentTypes[i]))
		}

		i += 1
	}

	parts := []string{}

	for _, v := range args {
		parts = append(parts, v)
	}

	return []string{
		fmt.Sprintf("%s(%s)", functionName, strings.Join(parts, ", ")),
		functionDef.ReturnType}
}
