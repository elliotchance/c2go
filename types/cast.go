package types

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	goast "go/ast"

	"strconv"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"
)

// GetArrayTypeAndSize returns the size and type of a fixed array. If the type
// is not an array with a fixed size then the type return will be an empty
// string, and the size will be -1.
func GetArrayTypeAndSize(s string) (string, int) {
	match := regexp.MustCompile(`(.*) \[(\d+)\]`).FindStringSubmatch(s)
	if len(match) > 0 {
		return match[1], util.Atoi(match[2])
	}

	return "", -1
}

func CastExpr(p *program.Program, expr ast.Expr, fromType, toType string) (ast.Expr, error) {
	// Let's assume that anything can be converted to a void pointer.
	if toType == "void *" {
		return expr, nil
	}

	fromType, err := ResolveType(p, fromType)
	if err != nil {
		return expr, err
	}

	toType, err = ResolveType(p, toType)
	if err != nil {
		return expr, err
	}

	if fromType == "null" && toType == "[][]byte" {
		return util.NewNil(), nil
	}

	if fromType == "null" && toType == "float64" {
		return util.NewStringLit("0.0"), nil
	}

	if fromType == "null" && toType == "bool" {
		return goast.NewIdent("false"), nil
	}

	// FIXME: This is a hack to avoid casting in some situations.
	if fromType == "" || toType == "" {
		return expr, nil
	}

	if fromType == "null" && toType == "[]byte" {
		return util.NewNil(), nil
	}

	// This if for linux.
	if fromType == "*_IO_FILE" && toType == "*noarch.File" {
		return expr, nil
	}

	if fromType == toType {
		return expr, nil
	}

	// Compatible integer types
	types := []string{
		// Integer types
		"byte",
		"int", "int16", "int32", "int64",
		"uint16", "uint32", "uint64",

		// Floating-point types.
		"float32", "float64",

		// Known aliases
		"__uint16_t", "size_t",

		// Darwin specific
		"__darwin_ct_rune_t", "darwin.Darwin_ct_rune_t",
	}
	for _, v := range types {
		if fromType == v && toType == "bool" {
			return &goast.BinaryExpr{
				X:  expr,
				Op: token.NEQ,
				Y: &goast.BasicLit{
					Kind:  token.STRING,
					Value: "0",
				},
			}, nil
		}
	}

	// In the forms of:
	// - `string` -> `[]byte`
	// - `string` -> `char *[13]`
	match1 := regexp.MustCompile(`\[\]byte`).FindStringSubmatch(toType)
	match2 := regexp.MustCompile(`char \*\[(\d+)\]`).FindStringSubmatch(toType)
	if fromType == "string" && (len(match1) > 0 || len(match2) > 0) {
		// Construct a byte array from "first":
		//
		//     var str []byte = []byte{'f','i','r','s','t'}

		value := &goast.CompositeLit{
			Type: &goast.ArrayType{
				Elt: goast.NewIdent("byte"),
			},
			Elts: []goast.Expr{},
		}

		strValue, err := strconv.Unquote(expr.(*goast.BasicLit).Value)
		if err != nil {
			panic(fmt.Sprintf("Failed to Unquote %s\n", expr.(*goast.BasicLit).Value))
		}

		for _, c := range []byte(strValue) {
			value.Elts = append(value.Elts, &goast.BasicLit{
				Kind:  token.CHAR,
				Value: fmt.Sprintf("%q", c),
			})
		}

		value.Elts = append(value.Elts, &goast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		})

		return value, nil
	}

	// In the forms of:
	// - `[7]byte` -> `string`
	// - `char *[12]` -> `string`
	match1 = regexp.MustCompile(`\[(\d+)\]byte`).FindStringSubmatch(fromType)
	match2 = regexp.MustCompile(`char \*\[(\d+)\]`).FindStringSubmatch(fromType)
	if (len(match1) > 0 || len(match2) > 0) && toType == "string" {
		size := 0
		if len(match1) > 0 {
			size = util.Atoi(match1[1])
		} else {
			size = util.Atoi(match2[1])
		}

		// The following code builds this:
		//
		//     string(expr[:size - 1])
		//
		return util.NewCallExpr(
			"string",
			&goast.SliceExpr{
				X: expr,
				High: &goast.BasicLit{
					Kind:  token.INT,
					Value: strconv.Itoa(size - 1),
				},
			},
		), nil
	}

	// Anything that is a pointer can be compared to nil
	if fromType[0] == '*' && toType == "bool" {
		return &goast.BinaryExpr{
			X:  expr,
			Op: token.NEQ,
			Y:  util.NewNil(),
		}, nil
	}

	if fromType == "[]byte" && toType == "bool" {
		return util.NewUnaryExpr(
			token.NOT, util.NewCallExpr("noarch.CStringIsNull", expr),
		), nil
	}

	if fromType == "int" && toType == "*int" {
		return util.NewNil(), nil
	}
	if fromType == "int" && toType == "*byte" {
		return &goast.BasicLit{
			Kind:  token.STRING,
			Value: `""`,
		}, nil
	}

	if fromType == "_Bool" && toType == "bool" {
		return expr, nil
	}

	if util.InStrings(fromType, types) && util.InStrings(toType, types) {
		return util.NewCallExpr(toType, expr), nil
	}

	p.AddImport("github.com/elliotchance/c2go/noarch")

	leftName := fromType
	rightName := toType

	if strings.Index(leftName, ".") != -1 {
		parts := strings.Split(leftName, ".")
		leftName = parts[len(parts)-1]
	}
	if strings.Index(rightName, ".") != -1 {
		parts := strings.Split(rightName, ".")
		rightName = parts[len(parts)-1]
	}

	functionName := fmt.Sprintf("noarch.%sTo%s",
		util.GetExportedName(leftName), util.GetExportedName(rightName))

	return util.NewCallExpr(functionName, expr), nil
}

func IsNullExpr(n goast.Expr) bool {
	if p1, ok := n.(*goast.ParenExpr); ok {
		if p2, ok := p1.X.(*goast.BasicLit); ok && p2.Value == "0" {
			return true
		}
	}

	return false
}
