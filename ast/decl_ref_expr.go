package ast

type DeclRefExpr struct {
	Addr       Address
	Pos        string
	Type       string
	Lvalue     bool
	For        string
	Address2   string
	Name       string
	Type2      string
	ChildNodes []Node
}

func parseDeclRefExpr(line string) *DeclRefExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		.*?
		(?P<lvalue> lvalue)?
		 (?P<for>\w+)
		 (?P<address2>[0-9a-fx]+)
		 '(?P<name>.*?)'
		 '(?P<type2>.*?)'`,
		line,
	)

	return &DeclRefExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        groups["position"],
		Type:       groups["type"],
		Lvalue:     len(groups["lvalue"]) > 0,
		For:        groups["for"],
		Address2:   groups["address2"],
		Name:       groups["name"],
		Type2:      groups["type2"],
		ChildNodes: []Node{},
	}
}

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
	return NewPositionFromString(n.Pos)
}
