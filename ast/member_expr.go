package ast

import "fmt"

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

func (n *MemberExpr) render(ast *Ast) (string, string) {
	children := n.Children

	lhs, lhsType := renderExpression(ast, children[0])
	lhsResolvedType := resolveType(ast, lhsType)
	rhs := n.Name

	if inStrings(lhsResolvedType, []string{"darwin.Float2", "darwin.Double2"}) {
		rhs = getExportedName(rhs)
	}

	src := fmt.Sprintf("%s.%s", lhs, rhs)
	return src, children[0].(*DeclRefExpr).Type
}

func (n *MemberExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
