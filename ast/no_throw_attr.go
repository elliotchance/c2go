package ast

type NoThrowAttr struct {
	Address  string
	Position string
	Children []interface{}
}

func parseNoThrowAttr(line string) *NoThrowAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &NoThrowAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}

func (n *NoThrowAttr) render(ast *Ast) (string, string) {
	return "", ""
}
