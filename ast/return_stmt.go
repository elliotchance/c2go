package ast

type ReturnStmt struct {
	Address  string
	Position string
}

func parseReturnStmt(line string) ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return ReturnStmt{
		Address: groups["address"],
		Position: groups["position"],
	}
}
