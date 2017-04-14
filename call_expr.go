package main

type CallExpr struct {
	Address  string
	Position string
	Type     string
	Children []interface{}
}

func parseCallExpr(line string) *CallExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)'",
		line,
	)

	return &CallExpr{
		Address: groups["address"],
		Position: groups["position"],
		Type: groups["type"],
		Children: []interface{}{},
	}
}
