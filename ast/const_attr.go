package ast

type ConstAttr struct {
	Address  string
	Position string
	Tags     string
}

func parseConstAttr(line string) ConstAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)>(?P<tags>.*)",
		line,
	)

	return ConstAttr{
		Address: groups["address"],
		Position: groups["position"],
		Tags: groups["tags"],
	}
}
