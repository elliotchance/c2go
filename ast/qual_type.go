package ast

type QualType struct {
	Address  string
	Type     string
	Kind     string
	Children []Node
}

func parseQualType(line string) *QualType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &QualType{
		Address:  groups["address"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

func (n *QualType) render(ast *Ast) (string, string) {
	return "", ""
}
