package ast

import (
	"bytes"

	"github.com/elliotchance/c2go/program"
)

type TranslationUnitDecl struct {
	Address  string
	Children []Node
}

func parseTranslationUnitDecl(line string) *TranslationUnitDecl {
	groups := groupsFromRegex("", line)

	return &TranslationUnitDecl{
		Address:  groups["address"],
		Children: []Node{},
	}
}

func (n *TranslationUnitDecl) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})
	for _, c := range n.Children {
		src, _ := renderExpression(program, c)
		printLine(out, src, program.Indent)
	}

	return out.String(), ""
}

func (n *TranslationUnitDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
