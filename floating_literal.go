package main

import "fmt"

type FloatingLiteral struct {
	Address  string
	Position string
	Type     string
	Value    float64
	Children []interface{}
}

func parseFloatingLiteral(line string) *FloatingLiteral {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<value>.+)",
		line,
	)

	return &FloatingLiteral{
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    atof(groups["value"]),
		Children: []interface{}{},
	}
}

func (n *FloatingLiteral) Render() []string {
	return []string{fmt.Sprintf("%f", n.Value), "double"}
}
