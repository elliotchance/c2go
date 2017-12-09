package ast

// CStyleCastExpr is expression.
type CStyleCastExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Type2      string
	Kind       string
	ChildNodes []Node
}

// CStyleCastExprNullToPointer - string of kind NullToPointer
var CStyleCastExprNullToPointer = "NullToPointer"

// CStyleCastExprToVoid - string of kind ToVoid
var CStyleCastExprToVoid = "ToVoid"

func parseCStyleCastExpr(line string) *CStyleCastExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type1>.*?)'(:'(?P<type2>.*?)')? <(?P<kind>.*)>",
		line,
	)

	return &CStyleCastExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type1"],
		Type2:      groups["type2"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CStyleCastExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CStyleCastExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CStyleCastExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CStyleCastExpr) Position() Position {
	return n.Pos
}
