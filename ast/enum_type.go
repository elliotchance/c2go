package ast

type EnumType struct {
	Address string
	Name    string
}

func parseEnumType(line string) EnumType {
	groups := groupsFromRegex(
		"'(?P<name>.*)'",
		line,
	)

	return EnumType{
		Address: groups["address"],
		Name: groups["name"],
	}
}
