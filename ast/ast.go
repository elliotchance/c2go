package ast

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
	case "MallocAttr":
		node = parseMallocAttr(line)
	case "NoThrowAttr":
		node = parseNoThrowAttr(line)
	case "NotNullAttr":
		node = parseNotNullAttr(line)
	case "PointerType":
		node = parsePointerType(line)
	case "QualType":
		node = parseQualType(line)
	case "Record":
		node = parseRecord(line)
	case "RecordDecl":
		node = parseRecordDecl(line)
	case "RecordType":
		node = parseRecordType(line)
	case "ReturnStmt":
		node = parseReturnStmt(line)
	case "TranslationUnitDecl":
		node = parseTranslationUnitDecl(line)
	case "Typedef":
		node = parseTypedef(line)
	case "TypedefType":
		node = parseTypedefType(line)
	case "WhileStmt":
		node = parseWhileStmt(line)
	default:
		panic(nodeName)
	}

	return node
}

func groupsFromRegex(rx, line string) map[string]string {
	fullRegexp := "(?P<address>[0-9a-fx]+) " + rx
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

	if len(s) >= 2 && s[0] == '"' && s[len(s) - 1] == '"' {
		return s[1:len(s) - 2]
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
