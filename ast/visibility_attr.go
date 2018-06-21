package ast

// VisibilityAttr contains information for a VisibilityAttr AST line.
type VisibilityAttr struct {
	Addr        Address
	Pos         Position
	ChildNodes  []Node
	IsDefault   bool
	IsInherited bool
	IsHidden    bool
}

func parseVisibilityAttr(line string) *VisibilityAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		(?P<default> Default)?
		(?P<hidden> Hidden)?
		`,
		line,
	)

	return &VisibilityAttr{
		Addr:        ParseAddress(groups["address"]),
		Pos:         NewPositionFromString(groups["position"]),
		ChildNodes:  []Node{},
		IsDefault:   len(groups["default"]) > 0,
		IsInherited: len(groups["inherited"]) > 0,
		IsHidden:    len(groups["hidden"]) > 0,
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *VisibilityAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *VisibilityAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *VisibilityAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *VisibilityAttr) Position() Position {
	return n.Pos
}
