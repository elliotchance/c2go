package ast

// ElaboratedType is elaborated type
type ElaboratedType struct {
	Addr       Address
	Type       string
	Tags       string
	ChildNodes []Node
}

func parseElaboratedType(line string) *ElaboratedType {
	groups := groupsFromRegex(
		"'(?P<type>.*?)' (?P<tags>.+)",
		line,
	)

	return &ElaboratedType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		Tags:       groups["tags"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ElaboratedType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *ElaboratedType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ElaboratedType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ElaboratedType) Position() Position {
	return Position{}
}
