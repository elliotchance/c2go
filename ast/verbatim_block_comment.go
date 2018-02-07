package ast

// VerbatimBlockComment is a type of comment
type VerbatimBlockComment struct {
	Addr       Address
	Pos        Position
	Name       string
	CloseName  string
	ChildNodes []Node
}

func parseVerbatimBlockComment(line string) *VerbatimBlockComment {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 Name="(?P<name>.*?)"
		 CloseName="(?P<close_name>.*?)"`,
		line,
	)

	return &VerbatimBlockComment{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Name:       groups["name"],
		CloseName:  groups["close_name"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *VerbatimBlockComment) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *VerbatimBlockComment) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *VerbatimBlockComment) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *VerbatimBlockComment) Position() Position {
	return n.Pos
}
