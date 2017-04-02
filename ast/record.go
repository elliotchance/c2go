package ast

type Record struct {
	Address string
	Type    string
}

func parseRecord(line string) Record {
	groups := groupsFromRegex(
		"'(?P<type>.*)'",
		line,
	)

	return Record{
		Address: groups["address"],
		Type: groups["type"],
	}
}
