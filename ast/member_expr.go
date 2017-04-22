package ast

import (
	"fmt"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"
)

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

func (n *MemberExpr) render(program *program.Program) (string, string) {
	children := n.Children

	lhs, lhsType := renderExpression(program, children[0])
	lhsResolvedType := types.ResolveType(program, lhsType)
	rhs := n.Name
	rhsType := ""

	// FIXME: This is just a hack
	if util.InStrings(lhsResolvedType, []string{"darwin.Float2", "darwin.Double2"}) {
		rhs = util.GetExportedName(rhs)
		rhsType = "int"
	}

	src := fmt.Sprintf("%s.%s", lhs, rhs)
	return src, rhsType
}

func (n *MemberExpr) AddChild(node Node) {
	n.Children = append(n.Children, node)
}
