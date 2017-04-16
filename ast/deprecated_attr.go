package ast

type DeprecatedAttr struct {
	Address  string
	Position string
	Message1 string
	Message2 string
	Children []interface{}
}

func parseDeprecatedAttr(line string) *DeprecatedAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> "(?P<message1>.*?)"(?P<message2> ".*?")?`,
		line,
	)

	return &DeprecatedAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Message1: removeQuotes(groups["message1"]),
		Message2: removeQuotes(groups["message2"]),
		Children: []interface{}{},
	}
}

func (n *DeprecatedAttr) render(ast *Ast) (string, string) {
	return "", ""
}
