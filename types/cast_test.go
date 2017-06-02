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

func ToJson(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func TestCast(t *testing.T) {
	p := program.NewProgram()

	type args struct {
		expr     string
		fromType string
		toType   string
	}
	tests := []struct {
		args args
		want goast.Expr
	}{
		// Casting to the same type is not needed.
		{args{"1", "int", "int"}, util.NewStringLit("1")},
		{args{"2.3", "float", "float"}, util.NewStringLit("2.3")},

		// Casting between numeric types.
		{args{"1", "int", "float"}, util.NewCallExpr("float32", util.NewStringLit("1"))},
		{args{"1", "int", "double"}, util.NewCallExpr("float64", util.NewStringLit("1"))},
		{args{"1", "int", "__uint16_t"}, util.NewCallExpr("uint16", util.NewStringLit("1"))},

		// Casting to bool
		{args{"1", "int", "bool"}, util.NewBinaryExpr(util.NewStringLit("1"), token.NEQ, util.NewStringLit("0"))},

		// Casting from bool. This is a special case becuase C int and bool
		// values are very commonly used interchangably.
		{args{"1", "bool", "int"}, util.NewCallExpr("noarch.BoolToInt", util.NewStringLit("1"))},

		// String types
		// {args{"foo", "[3]char", "const char*"}, "1 != 0"},

		{args{"false", "_Bool", "bool"}, util.NewStringLit("false")},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.args)

		t.Run(name, func(t *testing.T) {
			e := util.NewStringLit(tt.args.expr)
			got, err := CastExpr(p, e, tt.args.fromType, tt.args.toType)

			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cast()%s\n", util.ShowDiff(ToJson(got), ToJson(tt.want)))
			}
		})
	}
}
