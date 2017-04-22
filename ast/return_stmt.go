package ast

import (
	"bytes"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

type ReturnStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseReturnStmt(line string) *ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *ReturnStmt) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})
	r := "return"

	if len(n.Children) > 0 && program.FunctionName != "main" {
		re, reType := renderExpression(program, n.Children[0])
		r = "return " + types.Cast(program, re, reType, "int")
	}

	printLine(out, r, program.Indent)

	return out.String(), ""
}

func (n *ReturnStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
