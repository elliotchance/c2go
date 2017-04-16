package ast

type BuiltinType struct {
	Address  string
	Type     string
	Children []interface{}
}

func parseBuiltinType(line string) *BuiltinType {
	groups := groupsFromRegex(
		"'(?P<type>.*?)'",
		line,
	)

	return &BuiltinType{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []interface{}{},
	}
}

func (n *BuiltinType) render(ast *Ast) (string, string) {
	return "", ""
}
