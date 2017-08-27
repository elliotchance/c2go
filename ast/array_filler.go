package ast

type ArrayFiller struct {
	Children []Node
}

func parseArrayFiller(line string) *ArrayFiller {
	return &ArrayFiller{
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *ArrayFiller) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. For an ArrayFilter this will
// always be zero. See the documentation for the Address type for more
// information.
func (n *ArrayFiller) Address() Address {
	return 0
}
