package ast

// NoThrowAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type NoThrowAttr struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
	Implicit   bool
	Inherited  bool
}

func parseNoThrowAttr(line string) *NoThrowAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		(?P<implicit> Implicit)?
		`,
		line,
	)

	return &NoThrowAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
		Inherited:  len(groups["inherited"]) > 0,
		Implicit:   len(groups["implicit"]) > 0,
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *NoThrowAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *NoThrowAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *NoThrowAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *NoThrowAttr) Position() Position {
	return n.Pos
}
