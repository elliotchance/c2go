package ast

// ParamCommandComment is a type of comment
type ParamCommandComment struct {
	Addr       Address
	Pos        Position
	Other      string
	ChildNodes []Node
}

func parseParamCommandComment(line string) *ParamCommandComment {
	groups := groupsFromRegex(
		`<(?P<position>.*)> (?P<other>.*)`,
		line,
	)

	return &ParamCommandComment{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Other:      groups["other"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParamCommandComment) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ParamCommandComment) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ParamCommandComment) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ParamCommandComment) Position() Position {
	return n.Pos
}
