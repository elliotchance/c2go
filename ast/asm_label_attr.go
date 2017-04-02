package ast

type AsmLabelAttr struct {
	Address      string
	Position     string
	FunctionName string
}

func parseAsmLabelAttr(line string) AsmLabelAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> \"(?P<function>.+)\"",
		line,
	)

	return AsmLabelAttr{
		Address: groups["address"],
		Position: groups["position"],
		FunctionName: groups["function"],
	}
}
