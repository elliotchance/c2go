package ast

type DeprecatedAttr struct {
	Addr     Address
	Position string
	Message1 string
	Message2 string
	Children []Node
}

func parseDeprecatedAttr(line string) *DeprecatedAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)> "(?P<message1>.*?)"(?P<message2> ".*?")?`,
		line,
	)

	return &DeprecatedAttr{
		Addr:     ParseAddress(groups["address"]),
		Position: groups["position"],
		Message1: removeQuotes(groups["message1"]),
		Message2: removeQuotes(groups["message2"]),
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *DeprecatedAttr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *DeprecatedAttr) Address() Address {
	return n.Addr
}
