package ast

type RestrictAttr struct {
	Address  string
	Position string
	Name     string
}

func parseRestrictAttr(line string) RestrictAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<name>.+)",
		line,
	)

	return RestrictAttr{
		Address: groups["address"],
		Position: groups["position"],
		Name: groups["name"],
	}
}
