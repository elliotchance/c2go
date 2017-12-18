package ast

// MemberExpr is expression.
type MemberExpr struct {
	Addr       Address
	Pos        Position
	Type       string
	Type2      string
	Name       string
	IsLvalue   bool
	IsBitfield bool
	Address2   string
	IsPointer  bool
	ChildNodes []Node
}

func parseMemberExpr(line string) *MemberExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		(?P<type2>:'.*?')?
		(?P<lvalue> lvalue)?
		(?P<bitfield> bitfield)?
		 (?P<pointer>[->.]+)
		(?P<name>\w+)?
		 (?P<address2>[0-9a-fx]+)`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	return &MemberExpr{
		Addr:       ParseAddress(groups["address"]),
		Pos:        NewPositionFromString(groups["position"]),
		Type:       groups["type"],
		Type2:      type2,
		IsPointer:  groups["pointer"] == "->",
		Name:       groups["name"],
		IsLvalue:   len(groups["lvalue"]) > 0,
		IsBitfield: len(groups["bitfield"]) > 0,
		Address2:   groups["address2"],
		ChildNodes: []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *MemberExpr) AddChild(node Node) {
	n.ChildNodes = append(n.ChildNodes, node)
}

// GetDeclRefExpr gets DeclRefExpr from MemberExpr, or nil if there is no
// DeclRefExpr
func (n *MemberExpr) GetDeclRefExpr() *DeclRefExpr {
	for _, child := range n.ChildNodes {
		res, ok := child.(*DeclRefExpr)
		if ok {
			return res
		}

		cast, ok := child.(*ImplicitCastExpr)
		if ok {
			res, ok = cast.ChildNodes[0].(*DeclRefExpr)
			if ok {
				return res
			}

		}
	}

	// There is no DeclRefExpr
	return nil
}

// Address returns the numeric address of the node. See the documentation for
// the Address type for more information.
func (n *MemberExpr) Address() Address {
	return n.Addr
}

// Children returns the child nodes. If this node does not have any children or
// this node does not support children it will always return an empty slice.
func (n *MemberExpr) Children() []Node {
	return n.ChildNodes
}

// Position returns the position in the original source code.
func (n *MemberExpr) Position() Position {
	return n.Pos
}
