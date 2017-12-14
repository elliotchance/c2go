package ast

// ImplicitCastExpr is expression.
type ImplicitCastExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Type2      string
	Kind       string
	ChildNodes []Node
}

// ImplicitCastExprArrayToPointerDecay - constant
const ImplicitCastExprArrayToPointerDecay = "ArrayToPointerDecay"

func parseImplicitCastExpr(line string) *ImplicitCastExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		(:'(?P<type2>.*?)')?
		 <(?P<kind>.*)>`,
		line,
	)

	return &ImplicitCastExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Type2:      groups["type2"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ImplicitCastExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ImplicitCastExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ImplicitCastExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ImplicitCastExpr) Position() Position {
	return n.Pos
}
