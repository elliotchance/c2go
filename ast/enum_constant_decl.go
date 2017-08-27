package ast

type EnumConstantDecl struct {
	Addr       Address
	Position   string
	Position2  string
	Referenced bool
	Name       string
	Type       string
	Children   []Node
}

func parseEnumConstantDecl(line string) *EnumConstantDecl {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		( (?P<position2>[^ ]+))?
		( (?P<referenced>referenced))?
		 (?P<name>.+)
		 '(?P<type>.+?)'`,
		line,
	)

	return &EnumConstantDecl{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Position2:  groups["position2"],
		Referenced: len(groups["referenced"]) > 0,
		Name:       groups["name"],
		Type:       groups["type"],
		Children:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *EnumConstantDecl) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
