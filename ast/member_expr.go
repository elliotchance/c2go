package ast

type MemberExpr struct {
    Address  string
    Position string
    Type     string
    Lvalue   bool
    Name     string
    Address2 string
    Children []Node
}

func parseMemberExpr(line string) *MemberExpr {
    groups := groupsFromRegex(
        `<(?P<position>.*)>
		 '(?P<type>.*?)'
		 (?P<tags>.*?)
		(?P<name>\w+)
		 (?P<address2>[0-9a-fx]+)`,
        line,
    )

    return &MemberExpr{
        Address:  groups["address"],
        Position: groups["position"],
        Type:     groups["type"],
        Lvalue:   true,
        Name:     groups["name"],
        Address2: groups["address2"],
        Children: []Node{},
    }
}

func (n *MemberExpr) AddChild(node Node) {
    n.Children = append(n.Children, node)
}

func (n *MemberExpr) GetDeclRef() *DeclRefExpr {
    child := n.Children[0]
    res, ok := child.(*DeclRefExpr)
    if !ok {
        cast, ok := child.(*ImplicitCastExpr)
        if ok {
            res, _ = cast.Children[0].(*DeclRefExpr)
        }
    }

    return res
}
