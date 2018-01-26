package ast

// Field struct
type Field struct {
	Addr       Address
	String1    string
	String2    string
	ChildNodes []Node
}

func parseField(line string) *Field {
	groups := groupsFromRegex(
		`'(?P<string1>.*?)' '(?P<string2>.*?)'`,
		line,
	)

	return &Field{
		Addr:       ParseAddress(groups["address"]),
		String1:    groups["string1"],
		String2:    groups["string2"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *Field) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *Field) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *Field) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *Field) Position() Position {
	return Position{}
}
