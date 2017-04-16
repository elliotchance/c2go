package ast

type FormatAttr struct {
	Address      string
	Position     string
	Implicit     bool
	Inherited    bool
	FunctionName string
	Unknown1     int
	Unknown2     int
	Children     []interface{}
}

func parseFormatAttr(line string) *FormatAttr {
	groups := groupsFromRegex(
		`<(?P<position>.*)>
		(?P<implicit> Implicit)?
		(?P<inherited> Inherited)?
		 (?P<function>\w+)
		 (?P<unknown1>\d+)
		 (?P<unknown2>\d+)`,
		line,
	)

	return &FormatAttr{
		Address:      groups["address"],
		Position:     groups["position"],
		Implicit:     len(groups["implicit"]) > 0,
		Inherited:    len(groups["inherited"]) > 0,
		FunctionName: groups["function"],
		Unknown1:     atoi(groups["unknown1"]),
		Unknown2:     atoi(groups["unknown2"]),
		Children:     []interface{}{},
	}
}

func (n *FormatAttr) render(ast *Ast) (string, string) {
	return "", ""
}
