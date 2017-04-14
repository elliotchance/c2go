package main

type DeclStmt struct {
	Address  string
	Position string
	Children []interface{}
}

func parseDeclStmt(line string) *DeclStmt {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return &DeclStmt{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
