package ast

// DeprecatedAttr is a type of attribute that is optionally attached to a variable
// or struct field definition.
type DeprecatedAttr struct {
	Addr        Address
	Pos         Position
	Message1    string
	Message2    string
	IsInherited bool
	ChildNodes  []Node
}

func parseDeprecatedAttr(line string) *DeprecatedAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<inherited> Inherited)? "(?P<message1>.*?)"(?P<message2> ".*?")?`,
		line,
	)

	return &DeprecatedAttr{
		Addr:        ParseAddress(groups["address"]),
		Pos:         NewPositionFromString(groups["position"]),
		Message1:    removeQuotes(groups["message1"]),
		Message2:    removeQuotes(groups["message2"]),
		IsInherited: len(groups["inherited"]) > 0,
		ChildNodes:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DeprecatedAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *DeprecatedAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *DeprecatedAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *DeprecatedAttr) Position() Position {
	return n.Pos
}
