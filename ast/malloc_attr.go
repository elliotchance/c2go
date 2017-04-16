package ast

type MallocAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseMallocAttr(line string) *MallocAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &MallocAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *MallocAttr) render(ast *Ast) (string, string) {
	return "", ""
}
