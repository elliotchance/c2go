package ast

// ParenExpr is expression.
type ParenExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Type2      string
	Lvalue     bool
	IsBitfield bool
	ChildNodes []Node
}

func parseParenExpr(line string) *ParenExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type1>.*?)'(:'(?P<type2>.*)')?
		(?P<lvalue> lvalue)?
		(?P<bitfield> bitfield)?
		`,
		line,
	)

	return &ParenExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type1"],
		Type2:      groups["type2"],
		Lvalue:     len(groups["lvalue"]) > 0,
		IsBitfield: len(groups["bitfield"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ParenExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ParenExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ParenExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ParenExpr) Position() Position {
	return n.Pos
}
