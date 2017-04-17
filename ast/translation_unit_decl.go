package ast

import (
	"bytes"
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

func (n *TranslationUnitDecl) render(ast *Ast) (string, string) {
	out := bytes.NewBuffer([]byte{})
	for _, c := range n.Children {
		src, _ := renderExpression(ast, c)
		printLine(out, src, ast.indent)
	}

	return out.String(), ""
}

func (n *TranslationUnitDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
