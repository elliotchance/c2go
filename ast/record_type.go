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

func (n *RecordType) render(ast *Ast) (string, string) {
	return "", ""
}
