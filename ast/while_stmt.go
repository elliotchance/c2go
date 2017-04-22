package ast

import (
	"bytes"
	"fmt"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
)

type WhileStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseWhileStmt(line string) *WhileStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &WhileStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *WhileStmt) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})
	// TODO: The first child of a WhileStmt appears to always be null.
	// Are there any cases where it is used?
	children := n.Children[1:]

	e, eType := renderExpression(program, children[0])
	printLine(out, fmt.Sprintf("for %s {", types.Cast(program, e, eType, "bool")), program.Indent)

	body, _ := renderExpression(program, children[1])
	printLine(out, body, program.Indent+1)

	printLine(out, "}", program.Indent)

	return out.String(), ""
}

func (n *WhileStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
