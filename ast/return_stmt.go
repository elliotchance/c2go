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

func (n *ReturnStmt) render(p *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})
	r := "return"

	if len(n.Children) > 0 && p.FunctionName != "main" {
		re, reType := renderExpression(p, n.Children[0])
		funcDef := program.GetFunctionDefinition(p.FunctionName)
		r = "return " + types.Cast(p, re, reType, funcDef.ReturnType)
	}

	printLine(out, r, p.Indent)

	return out.String(), ""
}

func (n *ReturnStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
