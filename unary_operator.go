package main

type UnaryOperator struct {
	Address  string
	Position string
	Type     string
	IsLvalue bool
	IsPrefix bool
	Operator string
	Children []interface{}
}

func parseUnaryOperator(line string) *UnaryOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		(?P<lvalue> lvalue)?
		(?P<prefix> prefix)?
		(?P<postfix> postfix)?
		 '(?P<operator>.*?)'`,
		line,
	)

	return &UnaryOperator{
		Address: groups["address"],
		Position: groups["position"],
		Type: groups["type"],
		IsLvalue: len(groups["lvalue"]) > 0,
		IsPrefix: len(groups["prefix"]) > 0,
		Operator: groups["operator"],
		Children: []interface{}{},
	}
}
