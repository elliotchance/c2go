package ast

type AsmLabelAttr struct {
	Addr         Address
	Position     string
	Inherited    bool
	FunctionName string
	Children     []Node
}

func parseAsmLabelAttr(line string) *AsmLabelAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		 "(?P<function>.+)"`,
		line,
	)

	return &AsmLabelAttr{
		Addr:         ParseAddress(groups["address"]),
		Position:     groups["position"],
		Inherited:    len(groups["inherited"]) > 0,
		FunctionName: groups["function"],
		Children:     []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *AsmLabelAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
