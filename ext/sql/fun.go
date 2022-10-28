package sql

import (
	"fmt"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

func ds(v *val.Val) string                       { return v.Str().V }
func s(format string, a ...interface{}) *val.Val { return val.Str(fmt.Sprintf(format, a...)) }

var logicalFunPrecTbl = map[*val.Val]oper.BP{
	LOGIC_AND_BOOL_BOOL: oper.BP_LOGIC_AND,
	LOGIC_OR_BOOL_BOOL:  oper.BP_LOGIC_OR,
	LOGIC_NOT_BOOL:      oper.BP_PREFIX,
}

//goland:noinspection GoSnakeCaseUsage
var (
	// LOGIC_AND_BOOL_BOOL and :: bool -> bool -> bool
	LOGIC_AND_BOOL_BOOL = val.Fun(
		types.Fun(AND, []*types.Type{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return s("%s %s %s", ds(args[0]), AND, ds(args[1]))
		},
	)
	// LOGIC_OR_BOOL_BOOL or :: bool -> bool -> bool
	LOGIC_OR_BOOL_BOOL = val.Fun(
		types.Fun(OR, []*types.Type{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return s("%s %s %s", ds(args[0]), OR, ds(args[1]))
		},
	)
	// LOGIC_NOT_BOOL not :: bool -> bool
	LOGIC_NOT_BOOL = val.Fun(
		types.Fun(NOT, []*types.Type{types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return s("%s %s", NOT, ds(args[0]))
		},
	)
)

var (
	postfixUnary = func(oper string) val.IFun {
		return func(args ...*val.Val) *val.Val {
			//return s("(%s %s)", ds(args[0]), oper)
			return s("%s %s", ds(args[0]), oper)
		}
	}
	binary = func(oper string) val.IFun {
		return func(args ...*val.Val) *val.Val {
			//return s("(%s %s %s)", ds(args[0]), oper, ds(args[1]))
			return s("%s %s %s", ds(args[0]), oper, ds(args[1]))
		}
	}
)

//goland:noinspection GoSnakeCaseUsage
var (
	EQ_BOOL_BOOL = val.Fun(types.Fun(EQ, []*types.Type{types.Bool, types.Bool}, types.Bool), binary(EQ)) // == :: bool -> bool -> bool
	NE_BOOL_BOOL = val.Fun(types.Fun(NE, []*types.Type{types.Bool, types.Bool}, types.Bool), binary(NE)) // <> :: bool -> bool -> bool

	EQ_NUM_NUM = val.Fun(types.Fun(EQ, []*types.Type{types.Num, types.Num}, types.Bool), binary(EQ)) // = :: num -> num -> bool
	NE_NUM_NUM = val.Fun(types.Fun(NE, []*types.Type{types.Num, types.Num}, types.Bool), binary(NE)) // <> :: num -> num -> bool

	EQ_STR_STR = val.Fun(types.Fun(EQ, []*types.Type{types.Str, types.Str}, types.Bool), binary(EQ)) // = :: str -> str -> bool
	NE_STR_STR = val.Fun(types.Fun(NE, []*types.Type{types.Str, types.Str}, types.Bool), binary(NE)) // <> :: str -> str -> bool

	EQ_TIME_TIME = val.Fun(types.Fun(EQ, []*types.Type{types.Time, types.Time}, types.Bool), binary(EQ)) // = :: time -> time -> bool
	NE_TIME_TIME = val.Fun(types.Fun(NE, []*types.Type{types.Time, types.Time}, types.Bool), binary(NE)) // <> :: time -> time -> bool

	GT_NUM_NUM = val.Fun(types.Fun(GT, []*types.Type{types.Num, types.Num}, types.Bool), binary(GT)) // > :: num -> num -> bool
	GE_NUM_NUM = val.Fun(types.Fun(GE, []*types.Type{types.Num, types.Num}, types.Bool), binary(GE)) // >= :: num -> num -> bool
	LT_NUM_NUM = val.Fun(types.Fun(LT, []*types.Type{types.Num, types.Num}, types.Bool), binary(LT)) // < :: num -> num -> bool
	LE_NUM_NUM = val.Fun(types.Fun(LE, []*types.Type{types.Num, types.Num}, types.Bool), binary(LE)) // <= :: num -> num -> bool

	GT_TIME_TIME = val.Fun(types.Fun(GT, []*types.Type{types.Time, types.Time}, types.Bool), binary(GT)) // > :: time -> time -> bool
	GE_TIME_TIME = val.Fun(types.Fun(GE, []*types.Type{types.Time, types.Time}, types.Bool), binary(GE)) // >= :: time -> time -> bool
	LT_TIME_TIME = val.Fun(types.Fun(LT, []*types.Type{types.Time, types.Time}, types.Bool), binary(LT)) // < :: time -> time -> bool
	LE_TIME_TIME = val.Fun(types.Fun(LE, []*types.Type{types.Time, types.Time}, types.Bool), binary(LE)) // <= :: time -> time -> bool
)

//goland:noinspection GoSnakeCaseUsage
var (
	// LIKE_STR_STR like :: str -> str -> bool
	LIKE_STR_STR = val.Fun(types.Fun(LIKE, []*types.Type{types.Str, types.Str}, types.Bool), binary(LIKE))

	// BETWEEN_NUM_NUM_NUM between :: num -> num -> num -> bool
	BETWEEN_NUM_NUM_NUM = val.Fun(
		types.Fun(BETWEEN,
			[]*types.Type{types.Num, types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return s("%s %s %s %s %s", ds(args[0]), BETWEEN, ds(args[1]), AND, ds(args[2]))
		},
	)
	// BETWEEN_TIME_TIME_TIME between :: time -> time -> time -> bool
	BETWEEN_TIME_TIME_TIME = val.Fun(
		types.Fun(BETWEEN,
			[]*types.Type{types.Time, types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return s("%s %s %s %s %s", ds(args[0]), BETWEEN, ds(args[1]), AND, ds(args[2]))
		},
	)

	// IS_NULL_A isnull :: a -> bool
	IS_NULL_A = val.Fun(func() *types.Type {
		a := types.TyVar("a")
		return types.Fun(ISNULL, []*types.Type{a}, types.Bool)
	}(), postfixUnary("IS NULL"))

	// IN_LIST in :: a -> list[a] -> bool
	IN_LIST = val.Fun(func() *types.Type {
		a := types.TyVar("a")
		return types.Fun(IN, []*types.Type{a, types.List(a)}, types.Bool)
	}(), binary(IN))
)

// todo sql func
var ()
