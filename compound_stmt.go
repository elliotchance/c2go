package main

type CompoundStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseCompoundStmt(line string) *CompoundStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &CompoundStmt{
		Address: groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
