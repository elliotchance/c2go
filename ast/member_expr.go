package ast

type MemberExpr struct {
	Address  string
	Position string
	Type     string
	Lvalue   bool
	Name     string
	Address2 string
}

func parseMemberExpr(line string) MemberExpr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> '(?P<type>.*?)' (?P<tags>.*?)(?P<name>\\w+) (?P<address2>[0-9a-fx]+)",
		line,
	)

	return MemberExpr{
		Address: groups["address"],
		Position: groups["position"],
		Type: groups["type"],
		Lvalue: true,
		Name: groups["name"],
		Address2: groups["address2"],
	}
}
