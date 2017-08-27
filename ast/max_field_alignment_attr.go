package ast

import "github.com/elliotchance/c2go/util"

type MaxFieldAlignmentAttr struct {
	Addr     Address
	Position string
	Size     int
	Children []Node
}

func parseMaxFieldAlignmentAttr(line string) *MaxFieldAlignmentAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> Implicit (?P<size>\d*)`,
		line,
	)

	return &MaxFieldAlignmentAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Size:     util.Atoi(groups["size"]),
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *MaxFieldAlignmentAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *MaxFieldAlignmentAttr) Address() Address {
	return n.Addr
}
