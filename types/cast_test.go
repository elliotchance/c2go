package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/elliotchance/c2go/program"
	"github.com/elliotchance/c2go/util"

	goast "go/ast"
	"go/token"
)

func toJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func TestCast(t *testing.T) {
	p := program.NewProgram()

	type args struct {
		expr     goast.Expr
		fromType string
		toType   string
	}
	tests := []struct {
		args args
		want goast.Expr
	}{
		// Casting to the same type is not needed.
		{args{util.NewIntLit(1), "int", "int"}, util.NewIntLit(1)},
		{args{util.NewFloatLit(2.3), "float", "float"}, util.NewFloatLit(2.3)},

		// Casting between numeric types.
		{args{util.NewIntLit(1), "int", "float"}, util.NewCallExpr("float32", util.NewIntLit(1))},
		{args{util.NewIntLit(1), "int", "double"}, util.NewCallExpr("float64", util.NewIntLit(1))},
		{args{util.NewIntLit(1), "int", "__uint16_t"}, util.NewCallExpr("uint16", util.NewIntLit(1))},

		// Casting to bool
		{args{util.NewIntLit(1), "int", "bool"}, util.NewBinaryExpr(util.NewIntLit(1), token.NEQ, util.NewIntLit(0), "bool", false)},

		// Casting from bool. This is a special case becuase C int and bool
		// values are very commonly used interchangably.
		{args{util.NewIntLit(1), "bool", "int"}, util.NewCallExpr("noarch.BoolToInt", util.NewIntLit(1))},

		// String types
		// {args{"foo", "[3]char", "const char*"}, "1 != 0"},

		{args{util.NewIdent("false"), "_Bool", "bool"}, util.NewIdent("false")},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.args)

		t.Run(name, func(t *testing.T) {
			got, err := CastExpr(p, tt.args.expr, tt.args.fromType, tt.args.toType)

			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cast()%s\n", util.ShowDiff(toJSON(got), toJSON(tt.want)))
			}
		})
	}
}
