package ast

// BlockCommandComment is a type of comment
type BlockCommandComment struct {
	Addr       Address
	Pos        Position
	Name       string
	ChildNodes []Node
}

func parseBlockCommandComment(line string) *BlockCommandComment {
	groups := groupsFromRegex(
		`<(?P<position>.*)> Name="(?P<name>.*)"`,
		line,
	)

	return &BlockCommandComment{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Name:       groups["name"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *BlockCommandComment) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *BlockCommandComment) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *BlockCommandComment) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *BlockCommandComment) Position() Position {
	return n.Pos
}
