package ast

type IntegerLiteral struct {
	Address  string
	Position string
	Type     string
	Value    int
	Children []interface{}
}

func parseIntegerLiteral(line string) *IntegerLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>\\d+)",
		line,
	)

	return &IntegerLiteral{
		Address: groups["address"],
		Position: groups["position"],
		Type: groups["type"],
		Value: atoi(groups["value"]),
		Children: []interface{}{},
	}
}
