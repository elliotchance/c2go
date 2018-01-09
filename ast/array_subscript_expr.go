package ast

// ArraySubscriptExpr is expression.
type ArraySubscriptExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Type2      string
	IsLvalue   bool
	ChildNodes []Node
}

func parseArraySubscriptExpr(line string) *ArraySubscriptExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*?)'(:'(?P<type2>.*?)')?
		(?P<lvalue> lvalue)?`,
		line,
	)

	return &ArraySubscriptExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Type2:      groups["type2"],
		IsLvalue:   len(groups["lvalue"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ArraySubscriptExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ArraySubscriptExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ArraySubscriptExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ArraySubscriptExpr) Position() Position {
	return n.Pos
}
