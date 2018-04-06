package ast

// FormatArgAttr is a type of attribute that is optionally attached to a
// function definition.
type FormatArgAttr struct {
	Addr       Address
	Pos        Position
	Arg        string
	ChildNodes []Node
}

func parseFormatArgAttr(line string) *FormatArgAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 *(?P<arg>\d+)`,
		line,
	)

	return &FormatArgAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Arg:        groups["arg"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FormatArgAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FormatArgAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FormatArgAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FormatArgAttr) Position() Position {
	return n.Pos
}
