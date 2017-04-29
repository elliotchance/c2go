package ast

import (
	"bytes"

	"github.com/elliotchance/c2go/program"
)

type CompoundStmt struct {
	Address  string
	Position string
	Children []Node

	// TODO: remove this
	BelongsToSwitch bool
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Address:         groups["address"],
		Position:        groups["position"],
		Children:        []Node{},
		BelongsToSwitch: false,
	}
}

func (n *CompoundStmt) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})

	for i, c := range n.Children {
		// A switch statement in C usually has break statements. These break
		// statements (if they are not enclosed by a scope) will be children of
		// the CompoundStmt that is directly owned by the SwitchStmt. Since the
		// behavior of Go switches are different to that of C we have to be
		// careful to translate this correctly.
		if _, ok := c.(*BreakStmt); n.BelongsToSwitch && ok {
			// Ignore the break statement at the end of the case.
			continue
		}

		// On the other hand if there is not a break statement at the end of the
		// case we need to make sure it falls-through correctly.
		if n.BelongsToSwitch && i > 0 {
			if _, ok := n.Children[i-1].(*CaseStmt); ok {
				printLine(out, "fallthrough", program.Indent)
			}
		}

		src, _ := renderExpression(program, c)
		printLine(out, src, program.Indent)
	}

	return out.String(), ""
}

func (n *CompoundStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
