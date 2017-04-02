package ast

type ForStmt struct {
	Address  string
	Position string
}

func parseForStmt(line string) ForStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return ForStmt{
		Address: groups["address"],
		Position: groups["position"],
	}
}
