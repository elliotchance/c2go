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
