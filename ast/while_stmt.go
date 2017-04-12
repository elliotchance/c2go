package ast

type WhileStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseWhileStmt(line string) *WhileStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &WhileStmt{
		Address: groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
