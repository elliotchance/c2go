package ast

type FunctionProtoType struct {
	Address string
	Type    string
	Kind    string
	Children []interface{}
}

func parseFunctionProtoType(line string) *FunctionProtoType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &FunctionProtoType{
		Address: groups["address"],
		Type: groups["type"],
		Kind: groups["kind"],
		Children: []interface{}{},
	}
}
