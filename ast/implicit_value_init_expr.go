package ast

type ImplicitValueInitExpr struct {
	Addr     Address
	Position string
	Type1    string
	Type2    string
	Children []Node
}

func parseImplicitValueInitExpr(line string) *ImplicitValueInitExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type1>.*?)'(:'(?P<type2>.*)')?",
		line,
	)

	return &ImplicitValueInitExpr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Type1:    groups["type1"],
		Type2:    groups["type2"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ImplicitValueInitExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ImplicitValueInitExpr) Address() Address {
	return n.Addr
}
