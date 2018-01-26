package ast

// AsmLabelAttr is a type of attribute for assembler label
type AsmLabelAttr struct {
	Addr         Address
	Pos          Position
	Inherited    bool
	FunctionName string
	ChildNodes   []Node
}

func parseAsmLabelAttr(line string) *AsmLabelAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<inherited> Inherited)?
		 "(?P<function>.+)"`,
		line,
	)

	return &AsmLabelAttr{
		Addr:         ParseAddress(groups["address"]),
		Pos:          NewPositionFromString(groups["position"]),
		Inherited:    len(groups["inherited"]) > 0,
		FunctionName: groups["function"],
		ChildNodes:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *AsmLabelAttr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *AsmLabelAttr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *AsmLabelAttr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *AsmLabelAttr) Position() Position {
	return n.Pos
}
