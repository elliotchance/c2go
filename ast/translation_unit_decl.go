package ast

type TranslationUnitDecl struct {
	Address string
}

func parseTranslationUnitDecl(line string) TranslationUnitDecl {
	groups := groupsFromRegex("", line)

	return TranslationUnitDecl{
		Address: groups["address"],
	}
}
