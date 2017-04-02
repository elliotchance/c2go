package ast

type NoThrowAttr struct {
	Address  string
	Position string
}

func parseNoThrowAttr(line string) NoThrowAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>",
		line,
	)

	return NoThrowAttr{
		Address: groups["address"],
		Position: groups["position"],
	}
}
