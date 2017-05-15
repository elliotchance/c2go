package ast

type GotoStmt struct {
	Address   string
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
		Address:   groups["address"],
		Position:  groups["position"],
		Name:      groups["name"],
		Position2: groups["position2"],
		Children:  []Node{},
	}
}

func (n *GotoStmt) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
