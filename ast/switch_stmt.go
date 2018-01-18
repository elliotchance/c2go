package ast

// SwitchStmt is node represent 'switch'
type SwitchStmt struct {
	Addr       Address
	Pos        Position
	ChildNodes []Node
}

func parseSwitchStmt(line string) *SwitchStmt {
	groups := groupsFromRegex(`<(?P<position>.*)>`, line)

	return &SwitchStmt{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *SwitchStmt) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *SwitchStmt) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *SwitchStmt) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *SwitchStmt) Position() Position {
	return n.Pos
}
