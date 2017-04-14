package main

type ImplicitCastExpr struct {
	Address  string
	Position string
	Type     string
	Kind     string
	Children []interface{}
}

func parseImplicitCastExpr(line string) *ImplicitCastExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*)' <(?P<kind>.*)>",
		line,
	)

	return &ImplicitCastExpr{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []interface{}{},
	}
}
