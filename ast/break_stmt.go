package ast

type BreakStmt struct {
	Address  string
	Position string
}

func ParseBreakStmt(line string) BreakStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return BreakStmt{
		Address: groups["address"],
		Position: groups["position"],
	}
}
