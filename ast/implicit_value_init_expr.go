package ast

type ImplicitValueInitExpr struct {
	Address  string
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
		Address:  groups["address"],
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
