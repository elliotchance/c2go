package ast

type MemberExpr struct {
	Addr       Address
	Position   string
	Type       string
	Type2      string
	Name       string
	IsLvalue   bool
	IsBitfield bool
	Address2   string
	IsPointer  bool
	Children   []Node
}

func parseMemberExpr(line string) *MemberExpr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		(?P<type2>:'.*?')?
		 lvalue
		(?P<bitfield> bitfield)?
		 (?P<pointer>[->.]+)
		(?P<name>\w+)
		 (?P<address2>[0-9a-fx]+)`,
		line,
	)

	type2 := groups["type2"]
	if type2 != "" {
		type2 = type2[2 : len(type2)-1]
	}

	return &MemberExpr{
		Addr:       ParseAddress(groups["address"]),
		Position:   groups["position"],
		Type:       groups["type"],
		Type2:      type2,
		IsPointer:  groups["pointer"] == "->",
		Name:       groups["name"],
		IsLvalue:   true,
		IsBitfield: len(groups["bitfield"]) > 0,
		Address2:   groups["address2"],
		Children:   []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *MemberExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// GetDeclRefExpr gets DeclRefExpr from MemberExpr, or nil if there is no
// DeclRefExpr
func (n *MemberExpr) GetDeclRefExpr() *DeclRefExpr {
	for _, child := range n.Children {
		res, ok := child.(*DeclRefExpr)
		if ok {
			return res
		}

		cast, ok := child.(*ImplicitCastExpr)
		if ok {
			res, ok = cast.Children[0].(*DeclRefExpr)
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
