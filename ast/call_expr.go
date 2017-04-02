package ast

type CallExpr struct {
	Address  string
	Position string
	Type     string
}

func parseCallExpr(line string) CallExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)'",
		line,
	)

	return CallExpr{
		Address: groups["address"],
		Position: groups["position"],
		Type: groups["type"],
	}
}
