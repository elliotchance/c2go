package ast

// FunctionNoProtoType is a function type without parameters.
//
// Example:
//    int (*)()
type FunctionNoProtoType struct {
	Addr        Address
	Type        string
	CallingConv string
	ChildNodes  []Node
}

func parseFunctionNoProtoType(line string) *FunctionNoProtoType {
	groups := groupsFromRegex(
		"'(?P<type>.*?)' (?P<calling_conv>.*)",
		line,
	)

	return &FunctionNoProtoType{
		Addr:        ParseAddress(groups["address"]),
		Type:        groups["type"],
		CallingConv: groups["calling_conv"],
		ChildNodes:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *FunctionNoProtoType) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *FunctionNoProtoType) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *FunctionNoProtoType) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *FunctionNoProtoType) Position() Position {
	return Position{}
}
