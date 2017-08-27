package ast

import (
	"strings"

	"github.com/elliotchance/c2go/util"
)

type NonNullAttr struct {
	Addr     Address
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
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		A:        util.Atoi(strings.TrimSpace(groups["a"])),
		B:        b,
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *NonNullAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
