package main

import "fmt"

type MemberExpr struct {
	Address  string
	Position string
	Type     string
	Lvalue   bool
	Name     string
	Address2 string
	Children []interface{}
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
		Children: []interface{}{},
	}
}

func (n *MemberExpr) Render() []string {
	children := n.Children

	lhs := renderExpression(children[0])
	lhs_type := resolveType(lhs[1])
	rhs := n.Name

	if inStrings(lhs_type, []string{"darwin.Float2", "darwin.Double2"}) {
		rhs = getExportedName(rhs)
	}

	return []string{
		fmt.Sprintf("%s.%s", lhs[0], rhs),
		lhs[1],
	}
}
