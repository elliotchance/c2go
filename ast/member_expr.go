package ast

type MemberExpr struct {
	Address   string
	Position  string
	Type      string
	Lvalue    bool
	Name      string
	Address2  string
	IsPointer bool
	Children  []Node
}

func parseMemberExpr(line string) *MemberExpr {
	// 0x7fcc758e34a0 <col:8, col:12> 'int' lvalue ->_w 0x7fcc758d60c8
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		 '(?P<type>.*?)'
		 lvalue
		 (?P<pointer>[->.]+)
		(?P<name>\w+)
		 (?P<address2>[0-9a-fx]+)`,
		line,
	)

	return &MemberExpr{
		Address:   groups["address"],
		Position:  groups["position"],
		Type:      groups["type"],
		Lvalue:    true,
		IsPointer: groups["pointer"] == "->",
		Name:      groups["name"],
		Address2:  groups["address2"],
		Children:  []Node{},
	}
}

// AddChild adds a new child node. Child nodes can then be accessed with the
// Children attribute.
func (n *MemberExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}

// GetDeclRefExpr gets DeclRefExpr from MemberExpr, or nil if there is no DeclRefExpr
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
