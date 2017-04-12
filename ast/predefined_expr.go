package ast

type PredefinedExpr struct {
	Address  string
	Position string
	Type     string
	Name     string
	Lvalue   bool
	Children []interface{}
}

func parsePredefinedExpr(line string) *PredefinedExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' lvalue (?P<name>.*)",
		line,
	)

	return &PredefinedExpr{
		Address: groups["address"],
		Position: groups["position"],
		Type: groups["type"],
		Name: groups["name"],
		Lvalue: true,
		Children: []interface{}{},
	}
}
