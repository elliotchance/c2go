package ast

// DeclRefExpr is expression.
type DeclRefExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Type1      string
	Lvalue     bool
	For        string
	Address2   string
	Name       string
	Type2      string
	Type3      string
	ChildNodes []Node
}

func parseDeclRefExpr(line string) *DeclRefExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'(:'(?P<type1>.*?)')?
		.*?
		(?P<lvalue> lvalue)?
		 (?P<for>\w+)
		 (?P<address2>[0-9a-fx]+)
		 '(?P<name>.*?)'
		 '(?P<type2>.*?)'(:'(?P<type3>.*?)')?
		`,
		line,
	)

	return &DeclRefExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Type1:      groups["type1"],
		Lvalue:     len(groups["lvalue"]) > 0,
		For:        groups["for"],
		Address2:   groups["address2"],
		Name:       groups["name"],
		Type2:      groups["type2"],
		Type3:      groups["type3"],
		ChildNodes: []Node{},
	}
}

// FunctionDeclRefExpr - value of DeclRefExpr.For for function
var FunctionDeclRefExpr = "Function"

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DeclRefExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *DeclRefExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *DeclRefExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *DeclRefExpr) Position() Position {
	return n.Pos
}
