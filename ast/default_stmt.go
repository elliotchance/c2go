package ast

// DefaultStmt is node represent 'default'
type DefaultStmt struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
}

func parseDefaultStmt(line string) *DefaultStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &DefaultStmt{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DefaultStmt) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *DefaultStmt) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *DefaultStmt) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *DefaultStmt) Position() Position {
	return n.Pos
}
