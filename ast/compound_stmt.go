package ast

import (
	"bytes"

	"github.com/elliotchance/c2go/program"
)

type CompoundStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *CompoundStmt) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})

	for _, c := range n.Children {
		src, _ := renderExpression(program, c)
		printLine(out, src, program.Indent)
	}

	return out.String(), ""
}

func (n *CompoundStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
