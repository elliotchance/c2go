package ast

type IfStmt struct {
	Address  string
	Position string
}

func parseIfStmt(line string) IfStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return IfStmt{
		Address: groups["address"],
		Position: groups["position"],
	}
}
