package ast

// ArrayFiller is type of array filler
type ArrayFiller struct {
	ChildNodes []Node
}

func parseArrayFiller(line string) *ArrayFiller {
	return &ArrayFiller{
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ArrayFiller) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. For an ArrayFilter this will
// always be zero. See the documentation for the Address type for more
// information.
func (n *ArrayFiller) Address() Address {
	return 0
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *ArrayFiller) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *ArrayFiller) Position() Position {
	return Position{}
}
