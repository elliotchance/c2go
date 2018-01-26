package ast

// FunctionProtoType is function proto type
type FunctionProtoType struct {
	Addr       Address
	Type       string
	Kind       string
	ChildNodes []Node
}

func parseFunctionProtoType(line string) *FunctionProtoType {
	groups := groupsFromRegex(
		"'(?P<type>.*?)' (?P<kind>.*)",
		line,
	)

	return &FunctionProtoType{
		Addr:       ParseAddress(groups["address"]),
		Type:       groups["type"],
		Kind:       groups["kind"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FunctionProtoType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FunctionProtoType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FunctionProtoType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FunctionProtoType) Position() Position {
	return Position{}
}
