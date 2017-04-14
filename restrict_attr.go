package main

type RestrictAttr struct {
	Address  string
	Position string
	Name     string
	Children []interface{}
}

func parseRestrictAttr(line string) *RestrictAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return &RestrictAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Name:     groups["name"],
		Children: []interface{}{},
	}
}
