package main

type TypedefType struct {
	Address string
	Type    string
	Tags    string
	Children []interface{}
}

func parseTypedefType(line string) *TypedefType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<tags>.+)",
		line,
	)

	return &TypedefType{
		Address: groups["address"],
		Type: groups["type"],
		Tags: groups["tags"],
		Children: []interface{}{},
	}
}
