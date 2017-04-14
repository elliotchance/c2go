package main

type RecordType struct {
	Address  string
	Type     string
	Children []interface{}
}

func parseRecordType(line string) *RecordType {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return &RecordType{
		Address:  groups["address"],
		Type:     groups["type"],
		Children: []interface{}{},
	}
}
