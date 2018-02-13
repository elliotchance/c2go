package ast

// GotoStmt is node represent 'goto'
type GotoStmt struct {
	Addr       Address
	Pos        Position
	Name       string
	Position2  string
	ChildNodes []Node
}

func parseGotoStmt(line string) *GotoStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<name>.*)' (?P<position2>.*)",
		line,
	)

	return &GotoStmt{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Name:       groups["name"],
		Position2:  groups["position2"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *GotoStmt) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *GotoStmt) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *GotoStmt) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *GotoStmt) Position() Position {
	return n.Pos
}
