package ast

type ConditionalOperator struct {
	Address  string
	Position string
	Type string
	Children []interface{}
}

func parseConditionalOperator(line string) *ConditionalOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)> '(?P<type>.*?)'`,
		line,
	)

	return &ConditionalOperator{
		Address: groups["address"],
		Position: groups["position"],
		Type: groups["type"],
		Children: []interface{}{},
	}
}
