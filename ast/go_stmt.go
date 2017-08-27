package ast

type GotoStmt struct {
	Addr      Address
	Position  string
	Name      string
	Position2 string
	Children  []Node
}

func parseGotoStmt(line string) *GotoStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<name>.*)' (?P<position2>.*)",
		line,
	)

	return &GotoStmt{
		Addr:      ParseAddress(groups["address"]),
		Position:  groups["position"],
		Name:      groups["name"],
		Position2: groups["position2"],
		Children:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *GotoStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *GotoStmt) Address() Address {
	return n.Addr
}
