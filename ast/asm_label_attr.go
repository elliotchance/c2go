package ast

type AsmLabelAttr struct {
	Address      string
	Position     string
	FunctionName string
	Children     []interface{}
}

func parseAsmLabelAttr(line string) *AsmLabelAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> \"(?P<function>.+)\"",
		line,
	)

	return &AsmLabelAttr{
		Address:      groups["address"],
		Position:     groups["position"],
		FunctionName: groups["function"],
		Children:     []interface{}{},
	}
}

func (n *AsmLabelAttr) render(ast *Ast) (string, string) {
	return "", ""
}
