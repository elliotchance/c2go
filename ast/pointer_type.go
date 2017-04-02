package ast

type PointerType struct {
	Address string
	Type    string
}

func parsePointerType(line string) PointerType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return PointerType{
		Address: groups["address"],
		Type: groups["type"],
	}
}
