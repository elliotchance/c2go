package transpiler

import (
	"github.com/elliotchance/c2go/ast"
	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/types"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"strings"
)

func getSizeOfCType(cType string) int {
	// Remove the signedness keyword. This does not effect the size.
	if strings.HasPrefix(cType, "signed ") {
		cType = cType[7:]
	}
	if strings.HasPrefix(cType, "unsigned ") {
		cType = cType[9:]
	}

	// FIXME: The pointer size will be different on different platforms. We
	// should find out the correct size at runtime.
	pointerSize := 4

	switch cType {
	case "char":
		return 1

	case "short":
		return 2

	case "int":
		return 4

	case "long":
		return 8

	default:
		// If we cannot determine the type we can assume it is a type of
		// pointer.
		return pointerSize
	}
}

func transpileUnaryExprOrTypeTraitExpr(n *ast.UnaryExprOrTypeTraitExpr, p *program.Program) (
	*goast.BasicLit, string, []goast.Stmt, []goast.Stmt, error) {
	return util.NewIntLit(getSizeOfCType(n.Type2)), types.ResolveType(p, n.Type1), nil, nil, nil
}
