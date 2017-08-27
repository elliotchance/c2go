package ast

type TranslationUnitDecl struct {
	Addr     Address
	Children []Node
}

func parseTranslationUnitDecl(line string) *TranslationUnitDecl {
	groups := groupsFromRegex("", line)

	return &TranslationUnitDecl{
		Addr:     ParseAddress(groups["address"]),
		Children: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *TranslationUnitDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *TranslationUnitDecl) Address() Address {
	return n.Addr
}
