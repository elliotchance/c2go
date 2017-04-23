package ast

import (
	"strings"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

type NonNullAttr struct {
	Address  string
	Position string
	A        int
	B        int
	Children []Node
}

func parseNonNullAttr(line string) *NonNullAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<a> \d+)(?P<b> \d+)?`,
		line,
	)

	b := 0
	if groups["b"] != "" {
		b = util.Atoi(strings.TrimSpace(groups["b"]))
	}

	return &NonNullAttr{
		Address:  groups["address"],
		Position: groups["position"],
		A:        util.Atoi(strings.TrimSpace(groups["a"])),
		B:        b,
		Children: []Node{},
	}
}

func (n *NonNullAttr) render(program *program.Program) (string, string) {
	return "", ""
}

func (n *NonNullAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
