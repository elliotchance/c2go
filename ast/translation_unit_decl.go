package ast

// TranslationUnitDecl is node represents a translation unit declaration.
type TranslationUnitDecl struct {
	Addr       Address
	ChildNodes []Node
}

func parseTranslationUnitDecl(line string) *TranslationUnitDecl {
	groups := groupsFromRegex("<(?P<position>.*)> <(?P<position2>.*)>", line)

	return &TranslationUnitDecl{
		Addr:       ParseAddress(groups["address"]),
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *TranslationUnitDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *TranslationUnitDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *TranslationUnitDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *TranslationUnitDecl) Position() Position {
	return Position{}
}
