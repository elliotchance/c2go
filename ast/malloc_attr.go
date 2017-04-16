package ast

type MallocAttr struct {
	Address  string
	Position string
	Children []interface{}
}

func parseMallocAttr(line string) *MallocAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &MallocAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *MallocAttr) render(ast *Ast) (string, string) {
	return "", ""
}
