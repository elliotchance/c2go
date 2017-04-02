package ast

type AvailabilityAttr struct {
	Address     string
	Position    string
	OS          string
	Version     string
	Unknown1    int
	Unknown2    int
	Unavailable bool
	Message1    string
	Message2    string
}

func parseAvailabilityAttr(line string) AvailabilityAttr {
	groups := groupsFromRegex(
		"<(?P<position>.*)> (?P<os>\\w+) (?P<version>[\\d.]+) (?P<unknown1>[\\d.]+) (?P<unknown2>[\\d.]+)(?P<unavalable> Unavailable)? \"(?P<message1>.*?)\" \"(?P<message2>.*?)\"",
		line,
	)

	return AvailabilityAttr{
		Address: groups["address"],
		Position: groups["position"],
		OS: groups["os"],
		Version: groups["version"],
		Unknown1: atoi(groups["unknown1"]),
		Unknown2: atoi(groups["unknown2"]),
		Unavailable: len(groups["unavalable"]) > 0,
		Message1: removeQuotes(groups["message1"]),
		Message2: removeQuotes(groups["message2"]),
	}
}
