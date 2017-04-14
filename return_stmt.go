package main

type ReturnStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseReturnStmt(line string) *ReturnStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &ReturnStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
