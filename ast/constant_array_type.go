package ast

type ConstantArrayType struct {
	Address string
	Type    string
	Size    int
}

func parseConstantArrayType(line string) ConstantArrayType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<size>\\d+)",
		line,
	)

	return ConstantArrayType{
		Address: groups["address"],
		Type: groups["type"],
		Size: atoi(groups["size"]),
	}
}
