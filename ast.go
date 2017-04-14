package main

import (
	"regexp"
	"strconv"
	"strings"
)

func Parse(line string) interface{} {
	nodeName := strings.SplitN(line, " ", 2)[0]
	var node interface{}

	switch nodeName {
	case "AlwaysInlineAttr":
		node = parseAlwaysInlineAttr(line)
	case "ArraySubscriptExpr":
		node = parseArraySubscriptExpr(line)
	case "AsmLabelAttr":
		node = parseAsmLabelAttr(line)
	case "AvailabilityAttr":
		node = parseAvailabilityAttr(line)
	case "BinaryOperator":
		node = parseBinaryOperator(line)
	case "BreakStmt":
		node = parseBreakStmt(line)
	case "BuiltinType":
		node = parseBuiltinType(line)
	case "CallExpr":
		node = parseCallExpr(line)
	case "CharacterLiteral":
		node = parseCharacterLiteral(line)
	case "CompoundStmt":
		node = parseCompoundStmt(line)
	case "ConditionalOperator":
		node = parseConditionalOperator(line)
	case "ConstAttr":
		node = parseConstAttr(line)
	case "ConstantArrayType":
		node = parseConstantArrayType(line)
	case "CStyleCastExpr":
		node = parseCStyleCastExpr(line)
	case "DeclRefExpr":
		node = parseDeclRefExpr(line)
	case "DeclStmt":
		node = parseDeclStmt(line)
	case "DeprecatedAttr":
		node = parseDeprecatedAttr(line)
	case "ElaboratedType":
		node = parseElaboratedType(line)
	case "Enum":
		node = parseEnum(line)
	case "EnumConstantDecl":
		node = parseEnumConstantDecl(line)
	case "EnumDecl":
		node = parseEnumDecl(line)
	case "EnumType":
		node = parseEnumType(line)
	case "FieldDecl":
		node = parseFieldDecl(line)
	case "FloatingLiteral":
		node = parseFloatingLiteral(line)
	case "FormatAttr":
		node = parseFormatAttr(line)
	case "FunctionDecl":
		node = parseFunctionDecl(line)
	case "FunctionProtoType":
		node = parseFunctionProtoType(line)
	case "ForStmt":
		node = parseForStmt(line)
	case "IfStmt":
		node = parseIfStmt(line)
	case "ImplicitCastExpr":
		node = parseImplicitCastExpr(line)
	case "IntegerLiteral":
		node = parseIntegerLiteral(line)
	case "MallocAttr":
		node = parseMallocAttr(line)
	case "MemberExpr":
		node = parseMemberExpr(line)
	case "ModeAttr":
		node = parseModeAttr(line)
	case "NoThrowAttr":
		node = parseNoThrowAttr(line)
	case "NonNullAttr":
		node = parseNonNullAttr(line)
	case "ParenExpr":
		node = parseParenExpr(line)
	case "ParmVarDecl":
		node = parseParmVarDecl(line)
	case "PointerType":
		node = parsePointerType(line)
	case "PredefinedExpr":
		node = parsePredefinedExpr(line)
	case "QualType":
		node = parseQualType(line)
	case "Record":
		node = parseRecord(line)
	case "RecordDecl":
		node = parseRecordDecl(line)
	case "RecordType":
		node = parseRecordType(line)
	case "RestrictAttr":
		node = parseRestrictAttr(line)
	case "ReturnStmt":
		node = parseReturnStmt(line)
	case "StringLiteral":
		node = parseStringLiteral(line)
	case "TranslationUnitDecl":
		node = parseTranslationUnitDecl(line)
	case "Typedef":
		node = parseTypedef(line)
	case "TypedefDecl":
		node = parseTypedefDecl(line)
	case "TypedefType":
		node = parseTypedefType(line)
	case "UnaryOperator":
		node = parseUnaryOperator(line)
	case "VarDecl":
		node = parseVarDecl(line)
	case "WhileStmt":
		node = parseWhileStmt(line)
	default:
		panic("'" + line + "'")
	}

	return node
}

func groupsFromRegex(rx, line string) map[string]string {
	// We remove tabs and newlines from the regex. This is purely cosmetic
	// as the regex input can be quite lone and its nice for the caller to
	// be able to format it in a more readable way.
	fullRegexp := "(?P<address>[0-9a-fx]+) " +
		strings.Replace(strings.Replace(rx, "\n", "", -1), "\t", "", -1)
	re := regexp.MustCompile(fullRegexp)

	match := re.FindStringSubmatch(line)
	if len(match) == 0 {
		panic("could not match regexp '" + fullRegexp +
			"' with string '" + line + "'")
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return result
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

func removeQuotes(s string) string {
	s = strings.TrimSpace(s)

	if s == `""` {
		return ""
	}
	if s == `''` {
		return ""
	}

	if len(s) >= 2 && s[0] == '"' && s[len(s) - 1] == '"' {
		return s[1:len(s) - 2]
	}
	if len(s) >= 2 && s[0] == '\'' && s[len(s) - 1] == '\'' {
		return s[1:len(s) - 1]
	}

	return s
}

func atof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}

	return f
}

func unescapeString(s string) string {
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\r", "\r", -1)
	s = strings.Replace(s, "\\t", "\t", -1)

	return s
}
