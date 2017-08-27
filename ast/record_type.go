package ast

type RecordType struct {
	Addr     Address
	Type     string
	Children []Node
}

func parseRecordType(line string) *RecordType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &RecordType{
		Addr:     ParseAddress(groups["address"]),
		Type:     groups["type"],
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *RecordType) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
