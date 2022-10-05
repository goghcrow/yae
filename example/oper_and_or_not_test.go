package example

import (
	expr "github.com/goghcrow/yae"
	"github.com/goghcrow/yae/oper"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
)

// 添加自定义操作符, 同时添加对应的函数

const (
	NOT = "not"
	AND = "and"
	OR  = "or"
)

var ops = []oper.Operator{
	{OR, oper.BP_LOGIC_OR, oper.INFIX_L},
	{AND, oper.BP_LOGIC_AND, oper.INFIX_L},
	{NOT, oper.BP_PREFIX, oper.PREFIX},
}

//goland:noinspection GoSnakeCaseUsage
var (
	//AND_BOOL_BOOL and :: bool -> bool -> bool
	AND_BOOL_BOOL = val.LazyFun(
		types.Fun(AND, []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			thunk1 := args[0].Fun()
			if thunk1.V().Bool().V {
				thunk2 := args[1].Fun()
				return thunk2.V().Bool().Vl()
			} else {
				return val.False
			}
		},
	)
	//OR_BOOL_BOOL or :: bool -> bool -> bool
	OR_BOOL_BOOL = val.LazyFun(
		types.Fun(OR, []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			thunk1 := args[0].Fun()
			if thunk1.V().Bool().V {
				return val.True
			} else {
				thunk2 := args[1].Fun()
				return thunk2.V().Bool().Vl()
			}
		},
	)
	// NOT_BOOL not :: bool -> bool
	NOT_BOOL = val.Fun(
		types.Fun(NOT, []*types.Kind{types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!args[0].Bool().V)
		},
	)
)

func TestRegisterAndOrNot(t *testing.T) {
	expr := expr.NewExpr().EnableDebug(os.Stderr)
	expr.RegisterOperator(ops...)
	expr.RegisterFun(AND_BOOL_BOOL, OR_BOOL_BOOL, NOT_BOOL)

	// x > y && !(x <= z) || z == 42
	cachedClosure, err := expr.Compile(`X > Y and not (X <= Z) or Z == 42`, struct {
		X int
		Y int
		Z int
	}{})
	if err != nil {
		panic(err)
	}

	{
		v, err := cachedClosure(map[string]interface{}{
			"X": 100,
			"Y": 20,
			"Z": 30,
		})
		if err != nil {
			panic(err)
		}
		if v != val.True {
			t.Fail()
		}
	}

	{
		v, err := cachedClosure(struct {
			X int
			Y int
			Z int
		}{
			X: 1,
			Y: 1,
			Z: 42,
		})
		if err != nil {
			panic(err)
		}
		if v != val.True {
			t.Fail()
		}
	}

	eval := func(s string) *val.Val {
		closure, err := expr.Compile(s, nil)
		if err != nil {
			panic(err)
		}
		v, err := closure(nil)
		if err != nil {
			panic(err)
		}
		return v
	}

	tests := []struct {
		expr     string
		expected *val.Val
	}{
		{`false and false`, val.False},
		{`false and true`, val.False},
		{`true and false`, val.False},
		{`true and true`, val.True},

		{`false or false`, val.False},
		{`false or true`, val.True},
		{`true or false`, val.True},
		{`true or true`, val.True},

		{`not false`, val.True},
		{`not true`, val.False},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			actual := eval(tt.expr)
			if !val.Equals(tt.expected, actual) {
				t.Errorf("expect %s actual %s", tt.expected, actual)
			}
		})
	}
}