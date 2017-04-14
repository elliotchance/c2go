package main

type IfStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseIfStmt(line string) *IfStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &IfStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
