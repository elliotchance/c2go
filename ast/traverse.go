package ast

import (
	"reflect"
)

// GetAllNodesOfType returns all of the nodes of the tree that match the type
// provided. The type should be a pointer to an object in the ast package.
//
// The nodes returned may reference each other and there is no guarenteed order
// in which the nodes are returned.
func GetAllNodesOfType(root Node, t reflect.Type) []Node {
	nodes := []Node{}

	if reflect.TypeOf(root) == t {
		nodes = append(nodes, root)
	}

	// I know that below looks like a lot of duplicate code. However, the ast
	// package will evolve so that the Children attribute only appears on a
	// handful of nodes. And each of those attributes will be listed below.
	//
	// I would also like to minimise the amount of places we need to have these
	// comprehensive switch statements. Code that exists now or in the future
	// should try to traverse package instead.
	switch n := root.(type) {
	case *AlignedAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *AlwaysInlineAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ArraySubscriptExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *AsmLabelAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *AvailabilityAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *BinaryOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *BreakStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *BuiltinType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *CallExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *CaseStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *CharacterLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *CompoundStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ConditionalOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ConstAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ConstantArrayType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ContinueStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *CompoundAssignOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *CStyleCastExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *DeclRefExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *DeclStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *DefaultStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *DeprecatedAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *DoStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ElaboratedType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *Enum:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *EnumConstantDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *EnumDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *EnumType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *FieldDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *FloatingLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *FormatAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *FunctionDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *FunctionProtoType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ForStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *GotoStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *IfStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ImplicitCastExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ImplicitValueInitExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *IncompleteArrayType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *InitListExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *IntegerLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *LabelStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *MallocAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *MaxFieldAlignmentAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *MemberExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ModeAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *NoInlineAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *NoThrowAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *NonNullAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *OffsetOfExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *PackedAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ParenExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ParenType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ParmVarDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *PointerType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *PredefinedExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *PureAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *QualType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *Record:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *RecordDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *RecordType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *RestrictAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ReturnStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ReturnsTwiceAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *StringLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *SwitchStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *TranslationUnitDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *TransparentUnionAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *Typedef:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *TypedefDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *TypedefType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *UnaryExprOrTypeTraitExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *UnaryOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *VAArgExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *VarDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *WeakAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *WarnUnusedResultAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *WhileStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	default:
		panic(n)
	}

	return nodes
}
