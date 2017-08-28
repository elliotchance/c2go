package ast

type ModeAttr struct {
	Addr       Address
	Pos        Position
	Name       string
	ChildNodes []Node
}

func parseModeAttr(line string) *ModeAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return &ModeAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Name:       groups["name"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ModeAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ModeAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ModeAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ModeAttr) Position() Position {
	return n.Pos
}
