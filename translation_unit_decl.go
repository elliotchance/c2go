package main

type TranslationUnitDecl struct {
	Address  string
	Children []interface{}
}

func parseTranslationUnitDecl(line string) *TranslationUnitDecl {
	groups := groupsFromRegex("", line)

	return &TranslationUnitDecl{
		Address:  groups["address"],
		Children: []interface{}{},
	}
}
