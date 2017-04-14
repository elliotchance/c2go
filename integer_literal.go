package main

import "strconv"

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
		Address:  groups["address"],
		Position: groups["position"],
		Type:     groups["type"],
		Value:    atoi(groups["value"]),
		Children: []interface{}{},
	}
}

func (n *IntegerLiteral) Render() []string {
	literal := n.Value

	// FIXME
	//if str(literal)[-1] == 'L':
	//    literal = '%s(%s)' % (resolveType('long'), literal[:-1])

	return []string{strconv.FormatInt(int64(literal), 10), "int"}
}
