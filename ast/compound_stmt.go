package ast

type CompoundStmt struct {
	Address  string
	Position string
}

func parseCompoundStmt(line string) CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return CompoundStmt{
		Address: groups["address"],
		Position: groups["position"],
	}
}
