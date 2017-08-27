package ast

type DeclRefExpr struct {
	Addr     Address
	Position string
	Type     string
	Lvalue   bool
	For      string
	Address2 string
	Name     string
	Type2    string
	Children []Node
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
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type:     groups["type"],
		Lvalue:   len(groups["lvalue"]) > 0,
		For:      groups["for"],
		Address2: groups["address2"],
		Name:     groups["name"],
		Type2:    groups["type2"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DeclRefExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *DeclRefExpr) Address() Address {
	return n.Addr
}
