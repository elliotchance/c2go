package ast

type NonNullAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseNonNullAttr(line string) *NonNullAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> 1",
		line,
	)

	return &NonNullAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *NonNullAttr) render(ast *Ast) (string, string) {
	return "", ""
}

func (n *NonNullAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
