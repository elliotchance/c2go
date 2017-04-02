package ast

type DeclStmt struct {
	Address  string
	Position string
}

func parseDeclStmt(line string) DeclStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return DeclStmt{
		Address: groups["address"],
		Position: groups["position"],
	}
}
