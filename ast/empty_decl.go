package ast

// EmptyDecl - element of AST
type EmptyDecl struct {
	Addr       Address
	Pos        Position
	Position2  Position
	ChildNodes []Node
}

func parseEmptyDecl(line string) *EmptyDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		( (?P<position2>.*))?`,
		line,
	)

	return &EmptyDecl{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  NewPositionFromString(groups["position2"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *EmptyDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *EmptyDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *EmptyDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *EmptyDecl) Position() Position {
	return n.Pos
}
