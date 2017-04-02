package ast

type Typedef struct {
	Address string
	Type    string
}

func parseTypedef(line string) Typedef {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return Typedef{
		Address: groups["address"],
		Type: groups["type"],
	}
}
