package ast

type Enum struct {
	Address  string
	Name     string
	Children []interface{}
}

func parseEnum(line string) *Enum {
	groups := groupsFromRegex(
		"'(?P<name>.*)'",
		line,
	)

	return &Enum{
		Address:  groups["address"],
		Name:     groups["name"],
		Children: []interface{}{},
	}
}

func (n *Enum) render(ast *Ast) (string, string) {
	return "", ""
}
