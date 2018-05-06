package ast

// WeakAttr for the WeakAttr node
type WeakAttr struct {
	Addr       Address
	Pos        Position
	Inherited  bool
	ChildNodes []Node
}

func parseWeakAttr(line string) *WeakAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>(?P<inherited> Inherited)?`,
		line,
	)

	return &WeakAttr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Inherited:  len(groups["inherited"]) > 0,
		ChildNodes: []Node{},
	}
}

// AddChild method to implements Node interface
func (n *WeakAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *WeakAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *WeakAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *WeakAttr) Position() Position {
	return n.Pos
}
