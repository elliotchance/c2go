package ast

import (
	"bytes"
	"strings"

	"github.com/elliotchance/c2go/program"
)

type EnumDecl struct {
	Address   string
	Position  string
	Position2 string
	Name      string
	Children  []Node
}

func parseEnumDecl(line string) *EnumDecl {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<position2> 0x[^ ]+)?(?P<name>.*)",
		line,
	)

	return &EnumDecl{
		Address:   groups["address"],
		Position:  groups["position"],
		Position2: groups["position2"],
		Name:      strings.TrimSpace(groups["name"]),
		Children:  []Node{},
	}
}

func (n *EnumDecl) render(program *program.Program) (string, string) {
	out := bytes.NewBuffer([]byte{})

	out.WriteString("const (\n")
	for _, c := range n.Children {
		e, _ := c.render(program)
		out.WriteString(e + "\n")
	}
	out.WriteString(")\n")

	return out.String(), "unknown17"
}

func (n *EnumDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
