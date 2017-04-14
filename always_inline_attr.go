package main

type AlwaysInlineAttr struct {
	Address  string
	Position string
	Children []interface{}
}

func parseAlwaysInlineAttr(line string) *AlwaysInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> always_inline",
		line,
	)

	return &AlwaysInlineAttr{
		Address: groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
