package main

type ModeAttr struct {
	Address  string
	Position string
	Name     string
	Children []interface{}
}

func parseModeAttr(line string) *ModeAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return &ModeAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Name:     groups["name"],
		Children: []interface{}{},
	}
}
