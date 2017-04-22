package ast

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/elliotchance/c2go/program"
)

type Node interface {
	render(program *program.Program) (string, string)
	AddChild(node Node)
}

func Parse(line string) Node {

	nodeName := strings.SplitN(line, " ", 2)[0]

	switch nodeName {
	case "AlwaysInlineAttr":
		return parseAlwaysInlineAttr(line)
	case "ArraySubscriptExpr":
		return parseArraySubscriptExpr(line)
	case "AsmLabelAttr":
		return parseAsmLabelAttr(line)
	case "AvailabilityAttr":
		return parseAvailabilityAttr(line)
	case "BinaryOperator":
		return parseBinaryOperator(line)
	case "BreakStmt":
		return parseBreakStmt(line)
	case "BuiltinType":
		return parseBuiltinType(line)
	case "CallExpr":
		return parseCallExpr(line)
	case "CharacterLiteral":
		return parseCharacterLiteral(line)
	case "CompoundStmt":
		return parseCompoundStmt(line)
	case "ConditionalOperator":
		return parseConditionalOperator(line)
	case "ConstAttr":
		return parseConstAttr(line)
	case "ConstantArrayType":
		return parseConstantArrayType(line)
	case "CStyleCastExpr":
		return parseCStyleCastExpr(line)
	case "DeclRefExpr":
		return parseDeclRefExpr(line)
	case "DeclStmt":
		return parseDeclStmt(line)
	case "DeprecatedAttr":
		return parseDeprecatedAttr(line)
	case "DoStmt":
		return parseDoStmt(line)
	case "ElaboratedType":
		return parseElaboratedType(line)
	case "Enum":
		return parseEnum(line)
	case "EnumConstantDecl":
		return parseEnumConstantDecl(line)
	case "EnumDecl":
		return parseEnumDecl(line)
	case "EnumType":
		return parseEnumType(line)
	case "FieldDecl":
		return parseFieldDecl(line)
	case "FloatingLiteral":
		return parseFloatingLiteral(line)
	case "FormatAttr":
		return parseFormatAttr(line)
	case "FunctionDecl":
		return parseFunctionDecl(line)
	case "FunctionProtoType":
		return parseFunctionProtoType(line)
	case "ForStmt":
		return parseForStmt(line)
	case "IfStmt":
		return parseIfStmt(line)
	case "ImplicitCastExpr":
		return parseImplicitCastExpr(line)
	case "IntegerLiteral":
		return parseIntegerLiteral(line)
	case "MallocAttr":
		return parseMallocAttr(line)
	case "MemberExpr":
		return parseMemberExpr(line)
	case "ModeAttr":
		return parseModeAttr(line)
	case "NoThrowAttr":
		return parseNoThrowAttr(line)
	case "NonNullAttr":
		return parseNonNullAttr(line)
	case "ParenExpr":
		return parseParenExpr(line)
	case "ParenType":
		return parseParenType(line)
	case "ParmVarDecl":
		return parseParmVarDecl(line)
	case "PointerType":
		return parsePointerType(line)
	case "PredefinedExpr":
		return parsePredefinedExpr(line)
	case "PureAttr":
		return parsePureAttr(line)
	case "QualType":
		return parseQualType(line)
	case "Record":
		return parseRecord(line)
	case "RecordDecl":
		return parseRecordDecl(line)
	case "RecordType":
		return parseRecordType(line)
	case "RestrictAttr":
		return parseRestrictAttr(line)
	case "ReturnStmt":
		return parseReturnStmt(line)
	case "StringLiteral":
		return parseStringLiteral(line)
	case "TranslationUnitDecl":
		return parseTranslationUnitDecl(line)
	case "Typedef":
		return parseTypedef(line)
	case "TypedefDecl":
		return parseTypedefDecl(line)
	case "TypedefType":
		return parseTypedefType(line)
	case "UnaryOperator":
		return parseUnaryOperator(line)
	case "VarDecl":
		return parseVarDecl(line)
	case "WarnUnusedResultAttr":
		return parseWarnUnusedResultAttr(line)
	case "WhileStmt":
		return parseWhileStmt(line)
	case "NullStmt":
		return nil
	default:
		panic("unknown node type: '" + line + "'")
	}
}

func groupsFromRegex(rx, line string) map[string]string {
	// We remove tabs and newlines from the regex. This is purely cosmetic,
	// as the regex input can be quite long and it's nice for the caller to
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

func printLine(out *bytes.Buffer, line string, indent int) {
	out.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat("\t", indent), line))
}

func Render(program *program.Program, node Node) string {
	src, _ := node.render(program)
	return src
}

func renderExpression(program *program.Program, node Node) (string, string) {
	if node == nil {
		return "", "unknown54"
	}

	return node.render(program)
}

func getFunctionParams(f *FunctionDecl) []*ParmVarDecl {
	r := []*ParmVarDecl{}
	for _, n := range f.Children {
		if v, ok := n.(*ParmVarDecl); ok {
			r = append(r, v)
		}
	}

	return r
}

func getFunctionReturnType(f string) string {
	// The type of the function will be the complete prototype, like:
	//
	//     __inline_isfinitef(float) int
	//
	// will have a type of:
	//
	//     int (float)
	//
	// The arguments will handle themselves, we only care about the
	// return type ('int' in this case)
	returnType := strings.TrimSpace(strings.Split(f, "(")[0])

	if returnType == "" {
		panic(fmt.Sprintf("unable to extract the return type from: %s", f))
	}

	return returnType
}
