package ast

type CompoundAssignOperator struct {
	Address               string
	Position              string
	Type                  string
	Opcode                string
	ComputationLHSType    string
	ComputationResultType string
	Children              []Node
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
		Address:               groups["address"],
		Position:              groups["position"],
		Type:                  groups["type"],
		Opcode:                groups["opcode"],
		ComputationLHSType:    groups["clhstype"],
		ComputationResultType: groups["crestype"],
		Children:              []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *CompoundAssignOperator) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
