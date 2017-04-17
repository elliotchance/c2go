package ast

type FunctionProtoType struct {
	Address  string
	Type     string
	Kind     string
	Children []Node
}

func parseFunctionProtoType(line string) *FunctionProtoType {
	groups := groupsFromRegex(
		"'(?P<type>.*)' (?P<kind>.*)",
		line,
	)

	return &FunctionProtoType{
		Address:  groups["address"],
		Type:     groups["type"],
		Kind:     groups["kind"],
		Children: []Node{},
	}
}

func (n *FunctionProtoType) render(ast *Ast) (string, string) {
	return "", ""
}

func (n *FunctionProtoType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
