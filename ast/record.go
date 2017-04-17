package ast

type Record struct {
	Address  string
	Type     string
	Children []Node
}

func parseRecord(line string) *Record {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &Record{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *Record) render(ast *Ast) (string, string) {
	return "", ""
}

func (n *Record) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
