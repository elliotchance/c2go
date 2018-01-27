package ast

// CompoundStmt is node represents a compound of nodes
type CompoundStmt struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node

	// TODO: remove this
	BelongsToSwitch bool
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Addr:            ParseAddress(groups["address"]),
		Pos:             NewPositionFromString(groups["position"]),
		ChildNodes:      []Node{},
		BelongsToSwitch: false,
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CompoundStmt) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CompoundStmt) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CompoundStmt) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CompoundStmt) Position() Position {
	return n.Pos
}
