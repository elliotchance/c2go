package ast

// RestrictAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type RestrictAttr struct {
	Addr       Address
	Pos        Position
	Name       string
	ChildNodes []Node
}

func parseRestrictAttr(line string) *RestrictAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return &RestrictAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Name:       groups["name"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *RestrictAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *RestrictAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *RestrictAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *RestrictAttr) Position() Position {
	return n.Pos
}
