package ast

import (
	"fmt"
	"strings"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
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

func (n *CallExpr) render(program *program.Program) (string, string) {
	children := n.Children
	functionName, _ := renderExpression(program, children[0])
	functionDef := getFunctionDefinition(functionName)

	if functionDef == nil {
		panic(fmt.Sprintf("unknown function: %s", functionName))
	}

	if functionDef.Substitution != "" {
		parts := strings.Split(functionDef.Substitution, ".")
		program.AddImport(strings.Join(parts[:len(parts)-1], "."))

		parts2 := strings.Split(functionDef.Substitution, "/")
		functionName = parts2[len(parts2)-1]
	}

	args := []string{}
	i := 0
	for _, arg := range children[1:] {
		e, eType := renderExpression(program, arg)

		if i > len(functionDef.ArgumentTypes)-1 {
			// This means the argument is one of the varargs
			// so we don't know what type it needs to be
			// cast to.
			args = append(args, e)
		} else {
			args = append(args, types.Cast(program, e, eType, functionDef.ArgumentTypes[i]))
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
