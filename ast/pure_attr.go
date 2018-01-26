package ast

// PureAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type PureAttr struct {
	Addr       Address
	Pos        Position
	Implicit   bool
	Inherited  bool
	ChildNodes []Node
}

func parsePureAttr(line string) *PureAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		(?P<implicit> Implicit)?`,
		line,
	)

	return &PureAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Implicit:   len(groups["implicit"]) > 0,
		Inherited:  len(groups["inherited"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *PureAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *PureAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *PureAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *PureAttr) Position() Position {
	return n.Pos
}
