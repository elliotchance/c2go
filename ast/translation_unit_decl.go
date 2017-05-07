package ast

type TranslationUnitDecl struct {
	Address  string
	Children []Node
}

func parseTranslationUnitDecl(line string) *TranslationUnitDecl {
	groups := groupsFromRegex("", line)

	return &TranslationUnitDecl{
		Address:  groups["address"],
		Children: []Node{},
	}
}

func (n *TranslationUnitDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
