package ast

import "github.com/elliotchance/c2go/util"

type MaxFieldAlignmentAttr struct {
	Addr       Address
	Position   string
	Size       int
	ChildNodes []Node
}

func parseMaxFieldAlignmentAttr(line string) *MaxFieldAlignmentAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> Implicit (?P<size>\d*)`,
		line,
	)

	return &MaxFieldAlignmentAttr{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Size:       util.Atoi(groups["size"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *MaxFieldAlignmentAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *MaxFieldAlignmentAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *MaxFieldAlignmentAttr) Children() []Node {
	return n.ChildNodes
}
