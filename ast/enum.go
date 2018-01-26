package ast

// Enum struct
type Enum struct {
	Addr       Address
	Name       string
	ChildNodes []Node
}

func parseEnum(line string) *Enum {
	groups := groupsFromRegex(
		"'(?P<name>.*?)'",
		line,
	)

	return &Enum{
		Addr:       ParseAddress(groups["address"]),
		Name:       groups["name"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *Enum) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *Enum) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *Enum) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *Enum) Position() Position {
	return Position{}
}
