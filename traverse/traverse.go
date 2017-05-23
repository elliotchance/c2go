// Package traverse contains functions for traversing the clang AST.
package traverse

import (
	"reflect"

	"github.com/elliotchance/c2go/ast"
)

// GetAllNodesOfType returns all of the nodes of the tree that match the type
// provided. The type should be a pointer to an object in the ast package.
//
// The nodes returned may reference each other and there is no guarenteed order
// in which the nodes are returned.
func GetAllNodesOfType(root ast.Node, t reflect.Type) []ast.Node {
	nodes := []ast.Node{}

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
	case *ast.AlignedAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.AlwaysInlineAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ArraySubscriptExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.AsmLabelAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.AvailabilityAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.BinaryOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.BreakStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.BuiltinType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.CallExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.CaseStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.CharacterLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.CompoundStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ConditionalOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ConstAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ConstantArrayType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ContinueStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.CompoundAssignOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.CStyleCastExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.DeclRefExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.DeclStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.DefaultStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.DeprecatedAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.DoStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ElaboratedType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.Enum:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.EnumConstantDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.EnumDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.EnumType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.FieldDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.FloatingLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.FormatAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.FunctionDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.FunctionProtoType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ForStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.GotoStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.IfStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ImplicitCastExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ImplicitValueInitExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.IncompleteArrayType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.InitListExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.IntegerLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.LabelStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.MallocAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.MaxFieldAlignmentAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.MemberExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ModeAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.NoInlineAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.NoThrowAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.NonNullAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.OffsetOfExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.PackedAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ParenExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ParenType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ParmVarDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.PointerType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.PredefinedExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.PureAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.QualType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.Record:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.RecordDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.RecordType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.RestrictAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ReturnStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.ReturnsTwiceAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.StringLiteral:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.SwitchStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.TranslationUnitDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.TransparentUnionAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.Typedef:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.TypedefDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.TypedefType:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.UnaryExprOrTypeTraitExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.UnaryOperator:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.VAArgExpr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.VarDecl:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.WarnUnusedResultAttr:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	case *ast.WhileStmt:
		for _, c := range n.Children {
			nodes = append(nodes, GetAllNodesOfType(c, t)...)
		}
	default:
		panic(n)
	}

	return nodes
}
