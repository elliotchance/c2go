package ast

type RecordType struct {
	Address string
	Type    string
}

func parseRecordType(line string) RecordType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return RecordType{
		Address: groups["address"],
		Type: groups["type"],
	}
}
