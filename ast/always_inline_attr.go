package ast

type AlwaysInlineAttr struct {
	Address  string
	Position string
}

func parseAlwaysInlineAttr(line string) AlwaysInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> always_inline",
		line,
	)

	return AlwaysInlineAttr{
		Address: groups["address"],
		Position: groups["position"],
	}
}
