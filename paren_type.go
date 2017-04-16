package main

type ParenType struct {
	Address  string
	Type     string
	Sugar    bool
	Children []interface{}
}

func parseParenType(line string) *ParenType {
	groups := groupsFromRegex(`'(?P<type>.*?)' sugar`, line)

	return &ParenType{
		Address:  groups["address"],
		Type:     groups["type"],
		Sugar:    true,
		Children: []interface{}{},
	}
}
