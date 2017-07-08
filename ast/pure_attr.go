package ast

type PureAttr struct {
	Address   string
	Position  string
	Implicit  bool
	Inherited bool
	Children  []Node
}

func parsePureAttr(line string) *PureAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		(?P<implicit> Implicit)?`,
		line,
	)

	return &PureAttr{
		Address:   groups["address"],
		Position:  groups["position"],
		Implicit:  len(groups["implicit"]) > 0,
		Inherited: len(groups["inherited"]) > 0,
		Children:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *PureAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
