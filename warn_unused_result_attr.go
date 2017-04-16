package main

type WarnUnusedResultAttr struct {
	Address  string
	Position string
	Children []interface{}
}

func parseWarnUnusedResultAttr(line string) *WarnUnusedResultAttr {
	groups := groupsFromRegex(`<(?P<position>.*)> warn_unused_result`, line)

	return &WarnUnusedResultAttr{
		Address:  groups["address"],
		Position: groups["position"],
		Children: []interface{}{},
	}
}
