package ast

type PureAttr struct {
	Addr      Address
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
		Addr:      ParseAddress(groups["address"]),
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

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *PureAttr) Address() Address {
	return n.Addr
}
