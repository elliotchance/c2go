package ast

// VerbatimLineComment is a type of comment
type VerbatimLineComment struct {
	Addr       Address
	Pos        Position
	Text       string
	ChildNodes []Node
}

func parseVerbatimLineComment(line string) *VerbatimLineComment {
	groups := groupsFromRegex(
		`<(?P<position>.*)> Text="(?P<text>.*)"`,
		line,
	)

	return &VerbatimLineComment{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Text:       groups["text"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *VerbatimLineComment) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *VerbatimLineComment) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *VerbatimLineComment) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *VerbatimLineComment) Position() Position {
	return n.Pos
}
