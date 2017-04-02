package ast

type MallocAttr struct {
	Address  string
	Position string
}

func parseMallocAttr(line string) MallocAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return MallocAttr{
		Address: groups["address"],
		Position: groups["position"],
	}
}
