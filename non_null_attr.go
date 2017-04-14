package main

type NonNullAttr struct {
	Address  string
	Position string
	Children []interface{}
}

func parseNonNullAttr(line string) *NonNullAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> 1",
		line,
	)

	return &NonNullAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
