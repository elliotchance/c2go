package ast

type AlwaysInlineAttr struct {
	Address  string
	Position string
	Children []Node
}

func parseAlwaysInlineAttr(line string) *AlwaysInlineAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> always_inline",
		line,
	)

	return &AlwaysInlineAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []Node{},
	}
}

func (n *AlwaysInlineAttr) render(ast *Ast) (string, string) {
	return "", ""
}
