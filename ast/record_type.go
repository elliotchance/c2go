package ast

type RecordType struct {
	Address  string
	Type     string
	Children []Node
}

func parseRecordType(line string) *RecordType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &RecordType{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []Node{},
	}
}

func (n *RecordType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
