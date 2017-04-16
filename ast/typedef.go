package ast

type Typedef struct {
	Address  string
	Type     string
	Children []interface{}
}

func parseTypedef(line string) *Typedef {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &Typedef{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []interface{}{},
	}
}

func (n *Typedef) render(ast *Ast) (string, string) {
	return "", ""
}
