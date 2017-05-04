package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
)

type CaseStmt struct {
	Address  string
	Position string
	Children []Node
}

func parseCaseStmt(line string) *CaseStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &CaseStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *CaseStmt) render(program *program.Program) (string, string) {
	// panic(n.Children[2])

	// if n == 0 {
	// 	out += "fallthrough\n"
	// })
	c, _ := n.Children[0].render(program)
	c += ":"

	for _, s := range n.Children[1:] {
		if s != nil {
			line, _ := s.render(program)
			c += "\n" + line
		}
	}

	return fmt.Sprintf("case %s", c), ""
}

func (n *CaseStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
