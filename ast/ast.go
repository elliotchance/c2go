// Package ast parses the clang AST output into AST structures.
package ast

import (
	"regexp"
	"strconv"
	"strings"
)

// Node represents any node in the AST.
type Node interface {
	Address() Address
	Children() []Node
	AddChild(node Node)
}

// Address contains the memory address (originally outputted as a hexadecimal
// string) from the clang AST. The address are not predictable between run and
// are only useful for identifying nodes in a single AST.
//
// The Address is used like a primary key when storing the tree as a flat
// structure.
type Address uint64

// ParseAddress returns the integer representation of the hexadecimal address
// (like 0x7f8a1d8ccfd0). If the address cannot be parsed, 0 is returned.
func ParseAddress(address string) Address {
	addr, _ := strconv.ParseUint(address, 0, 64)

	return Address(addr)
}

// Position returns the position of the node in the original file. If the
// position cannot be determined an empty string is returned.
func Position(node Node) string {
	switch n := node.(type) {
	case *AlignedAttr:
		return n.Position
	case *AlwaysInlineAttr:
		return n.Position
	case *ArraySubscriptExpr:
		return n.Position
	case *AsmLabelAttr:
		return n.Position
	case *AvailabilityAttr:
		return n.Position
	case *BinaryOperator:
		return n.Position
	case *BreakStmt:
		return n.Position
	case *BuiltinType:
		return ""
	case *CallExpr:
		return n.Position
	case *CaseStmt:
		return n.Position
	case *CharacterLiteral:
		return n.Position
	case *CompoundStmt:
		return n.Position
	case *ConditionalOperator:
		return n.Position
	case *ConstAttr:
		return n.Position
	case *ConstantArrayType:
		return ""
	case *ContinueStmt:
		return n.Position
	case *CompoundAssignOperator:
		return n.Position
	case *CStyleCastExpr:
		return n.Position
	case *DeclRefExpr:
		return n.Position
	case *DeclStmt:
		return n.Position
	case *DefaultStmt:
		return n.Position
	case *DeprecatedAttr:
		return n.Position
	case *DoStmt:
		return n.Position
	case *ElaboratedType:
		return ""
	case *Enum:
		return ""
	case *EnumConstantDecl:
		return n.Position
	case *EnumDecl:
		return n.Position
	case *EnumType:
		return ""
	case *FieldDecl:
		return n.Position
	case *FloatingLiteral:
		return n.Position
	case *FormatAttr:
		return n.Position
	case *FunctionDecl:
		return n.Position
	case *FunctionProtoType:
		return ""
	case *ForStmt:
		return n.Position
	case *GotoStmt:
		return n.Position
	case *IfStmt:
		return n.Position
	case *ImplicitCastExpr:
		return n.Position
	case *ImplicitValueInitExpr:
		return n.Position
	case *IncompleteArrayType:
		return ""
	case *IndirectFieldDecl:
		return n.Position
	case *InitListExpr:
		return n.Position
	case *IntegerLiteral:
		return n.Position
	case *LabelStmt:
		return n.Position
	case *MallocAttr:
		return n.Position
	case *MaxFieldAlignmentAttr:
		return n.Position
	case *MemberExpr:
		return n.Position
	case *ModeAttr:
		return n.Position
	case *NoInlineAttr:
		return n.Position
	case *NoThrowAttr:
		return n.Position
	case *NonNullAttr:
		return n.Position
	case *OffsetOfExpr:
		return n.Position
	case *PackedAttr:
		return n.Position
	case *ParenExpr:
		return n.Position
	case *ParenType:
		return ""
	case *ParmVarDecl:
		return n.Position
	case *PointerType:
		return ""
	case *PredefinedExpr:
		return n.Position
	case *PureAttr:
		return n.Position
	case *QualType:
		return ""
	case *Record:
		return ""
	case *RecordDecl:
		return n.Position
	case *RecordType:
		return ""
	case *RestrictAttr:
		return n.Position
	case *ReturnStmt:
		return n.Position
	case *ReturnsTwiceAttr:
		return n.Position
	case *StringLiteral:
		return n.Position
	case *SwitchStmt:
		return n.Position
	case *TranslationUnitDecl:
		return ""
	case *TransparentUnionAttr:
		return n.Position
	case *Typedef:
		return ""
	case *TypedefDecl:
		return n.Position
	case *TypedefType:
		return ""
	case *UnaryExprOrTypeTraitExpr:
		return n.Position
	case *UnaryOperator:
		return n.Position
	case *UnusedAttr:
		return n.Position
	case *VAArgExpr:
		return n.Position
	case *VarDecl:
		return n.Position
	case *WarnUnusedResultAttr:
		return n.Position
	case *WeakAttr:
		return n.Position
	case *WhileStmt:
		return n.Position
	default:
		panic(n)
	}
}

// Parse takes the coloured output of the clang AST command and returns a root
// node for the AST.
func Parse(line string) Node {
	// This is a special case. I'm not sure if it's a bug in the clang AST
	// dumper. It should have children.
	if line == "array filler" {
		return parseArrayFiller(line)
	}

	nodeName := strings.SplitN(line, " ", 2)[0]

	switch nodeName {
	case "AlignedAttr":
		return parseAlignedAttr(line)
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
	case "CaseStmt":
		return parseCaseStmt(line)
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
	case "ContinueStmt":
		return parseContinueStmt(line)
	case "CompoundAssignOperator":
		return parseCompoundAssignOperator(line)
	case "CStyleCastExpr":
		return parseCStyleCastExpr(line)
	case "DeclRefExpr":
		return parseDeclRefExpr(line)
	case "DeclStmt":
		return parseDeclStmt(line)
	case "DefaultStmt":
		return parseDefaultStmt(line)
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
	case "GotoStmt":
		return parseGotoStmt(line)
	case "IfStmt":
		return parseIfStmt(line)
	case "ImplicitCastExpr":
		return parseImplicitCastExpr(line)
	case "ImplicitValueInitExpr":
		return parseImplicitValueInitExpr(line)
	case "IncompleteArrayType":
		return parseIncompleteArrayType(line)
	case "IndirectFieldDecl":
		return parseIndirectFieldDecl(line)
	case "InitListExpr":
		return parseInitListExpr(line)
	case "IntegerLiteral":
		return parseIntegerLiteral(line)
	case "LabelStmt":
		return parseLabelStmt(line)
	case "MallocAttr":
		return parseMallocAttr(line)
	case "MaxFieldAlignmentAttr":
		return parseMaxFieldAlignmentAttr(line)
	case "MemberExpr":
		return parseMemberExpr(line)
	case "ModeAttr":
		return parseModeAttr(line)
	case "NoInlineAttr":
		return parseNoInlineAttr(line)
	case "NoThrowAttr":
		return parseNoThrowAttr(line)
	case "NonNullAttr":
		return parseNonNullAttr(line)
	case "OffsetOfExpr":
		return parseOffsetOfExpr(line)
	case "PackedAttr":
		return parsePackedAttr(line)
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
	case "ReturnsTwiceAttr":
		return parseReturnsTwiceAttr(line)
	case "StringLiteral":
		return parseStringLiteral(line)
	case "SwitchStmt":
		return parseSwitchStmt(line)
	case "TranslationUnitDecl":
		return parseTranslationUnitDecl(line)
	case "TransparentUnionAttr":
		return parseTransparentUnionAttr(line)
	case "Typedef":
		return parseTypedef(line)
	case "TypedefDecl":
		return parseTypedefDecl(line)
	case "TypedefType":
		return parseTypedefType(line)
	case "UnaryExprOrTypeTraitExpr":
		return parseUnaryExprOrTypeTraitExpr(line)
	case "UnaryOperator":
		return parseUnaryOperator(line)
	case "UnusedAttr":
		return parseUnusedAttr(line)
	case "VAArgExpr":
		return parseVAArgExpr(line)
	case "VarDecl":
		return parseVarDecl(line)
	case "WarnUnusedResultAttr":
		return parseWarnUnusedResultAttr(line)
	case "WeakAttr":
		return parseWeakAttr(line)
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
