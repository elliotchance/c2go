package main

type BinaryOperator struct {
	Address  string
	Position string
	Type     string
	Operator string
	Children []interface{}
}

func parseBinaryOperator(line string) *BinaryOperator {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' '(?P<operator>.*?)'",
		line,
	)

	return &BinaryOperator{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Operator: groups["operator"],
		Children: []interface{}{},
	}
}
