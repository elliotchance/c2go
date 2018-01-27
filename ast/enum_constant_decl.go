package ast

// EnumConstantDecl is node represents a enum constant declaration.
type EnumConstantDecl struct {
	Addr       Address
	Pos        Position
	Position2  string
	Referenced bool
	Name       string
	Type       string
	ChildNodes []Node
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
		Pos:        NewPositionFromString(groups["position"]),
		Position2:  groups["position2"],
		Referenced: len(groups["referenced"]) > 0,
		Name:       groups["name"],
		Type:       groups["type"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *EnumConstantDecl) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *EnumConstantDecl) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *EnumConstantDecl) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *EnumConstantDecl) Position() Position {
	return n.Pos
}
