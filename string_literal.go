package main

type StringLiteral struct {
	Address  string
	Position string
	Type     string
	Value    string
	Lvalue   bool
	Children []interface{}
}

func parseStringLiteral(line string) *StringLiteral {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*)' lvalue "(?P<value>.*)"`,
		line,
	)

	return &StringLiteral{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    unescapeString(groups["value"]),
		Lvalue:   true,
		Children: []interface{}{},
	}
}
