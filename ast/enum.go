package ast

type Enum struct {
	Address string
	Name    string
}

func parseEnum(line string) Enum {
	groups := groupsFromRegex(
		"'(?P<name>.*)'",
		line,
	)

	return Enum{
		Address: groups["address"],
		Name: groups["name"],
	}
}
