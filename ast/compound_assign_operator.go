package ast

// CompoundAssignOperator is type of compound assign operator
type CompoundAssignOperator struct {
	Addr                  Address
	Pos                   Position
	Type                  string
	Opcode                string
	ComputationLHSType    string
	ComputationResultType string
	ChildNodes            []Node
}

func parseCompoundAssignOperator(line string) *CompoundAssignOperator {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.+?)'
		 '(?P<opcode>.+?)'
		 ComputeLHSTy='(?P<clhstype>.+?)'
		 ComputeResultTy='(?P<crestype>.+?)'`,
		line,
	)

	return &CompoundAssignOperator{
		Addr:                  ParseAddress(groups["address"]),
		Pos:                   NewPositionFromString(groups["position"]),
		Type:                  groups["type"],
		Opcode:                groups["opcode"],
		ComputationLHSType:    groups["clhstype"],
		ComputationResultType: groups["crestype"],
		ChildNodes:            []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CompoundAssignOperator) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *CompoundAssignOperator) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *CompoundAssignOperator) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *CompoundAssignOperator) Position() Position {
	return n.Pos
}
