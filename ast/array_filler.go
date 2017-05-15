package ast

type ArrayFiller struct {
	Children []Node
}

func parseArrayFiller(line string) *ArrayFiller {
	return &ArrayFiller{
		Children: []Node{},
	}
}

func (n *ArrayFiller) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
