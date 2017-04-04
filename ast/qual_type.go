package ast

type QualType struct {
	Address string
	Type    string
	Kind    string
	Children []interface{}
}

func parseQualType(line string) *QualType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &QualType{
		Address: groups["address"],
		Type: groups["type"],
		Kind: groups["kind"],
		Children: []interface{}{},
	}
}
