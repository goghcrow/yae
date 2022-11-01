package fun

import (
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// EQ_BOOL_BOOL == :: bool -> bool -> bool
	EQ_BOOL_BOOL = val.Fun(
		types.Fun(oper.EQ, []*types.Type{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Bool().V == args[1].Bool().V)
		},
	)
	// NE_BOOL_BOOL != :: bool -> bool -> bool
	NE_BOOL_BOOL = val.Fun(
		types.Fun(oper.NE, []*types.Type{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Bool().V != args[1].Bool().V)
		},
	)
	// EQ_NUM_NUM == :: num -> num -> bool
	EQ_NUM_NUM = val.Fun(
		types.Fun(oper.EQ, []*types.Type{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.NumEQ(args[0].Num(), args[1].Num()))
		},
	)
	// NE_NUM_NUM != :: num -> num -> bool
	NE_NUM_NUM = val.Fun(
		types.Fun(oper.NE, []*types.Type{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.NumNE(args[0].Num(), args[1].Num()))
		},
	)
	// EQ_STR_STR == :: str -> str -> bool
	EQ_STR_STR = val.Fun(
		types.Fun(oper.EQ, []*types.Type{types.Str, types.Str}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Str().V == args[1].Str().V)
		},
	)
	// NE_STR_STR != :: str -> str -> bool
	NE_STR_STR = val.Fun(
		types.Fun(oper.NE, []*types.Type{types.Str, types.Str}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Str().V != args[1].Str().V)
		},
	)
	// EQ_TIME_TIME == :: time -> time -> bool
	EQ_TIME_TIME = val.Fun(
		types.Fun(oper.EQ, []*types.Type{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	// NE_TIME_TIME != :: time -> time -> bool
	NE_TIME_TIME = val.Fun(
		types.Fun(oper.NE, []*types.Type{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	// EQ_LIST_LIST == :: forall a. (list[a] -> list[a] -> bool)
	EQ_LIST_LIST = func() *val.Val {
		T := types.TyVar("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(oper.EQ, []*types.Type{listT, listT}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(val.Equals(args[0], args[1]))
			},
		)
	}()
	// NE_LIST_LIST != :: forall a. (list[a] -> list[a] -> bool)
	NE_LIST_LIST = func() *val.Val {
		T := types.TyVar("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(oper.NE, []*types.Type{listT, listT}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(!val.Equals(args[0], args[1]))
			},
		)
	}()
	// EQ_MAP_MAP == :: forall k v. (map[k,v] -> map[k,v] -> bool)
	EQ_MAP_MAP = func() *val.Val {
		K := types.TyVar("k")
		V := types.TyVar("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(oper.EQ, []*types.Type{mapKV, mapKV}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(val.Equals(args[0], args[1]))
			},
		)
	}()
	// NE_MAP_MAP != :: forall k v. (map[k,v] -> map[k,v] -> bool)
	NE_MAP_MAP = func() *val.Val {
		K := types.TyVar("k")
		V := types.TyVar("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(oper.NE, []*types.Type{mapKV, mapKV}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(!val.Equals(args[0], args[1]))
			},
		)
	}()

	// GT_NUM_NUM > :: num -> num -> bool
	GT_NUM_NUM = val.Fun(
		types.Fun(oper.GT, []*types.Type{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.NumGT(args[0].Num(), args[1].Num()))
		},
	)
	// GE_NUM_NUM >= :: num -> num -> bool
	GE_NUM_NUM = val.Fun(
		types.Fun(oper.GE, []*types.Type{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.NumGE(args[0].Num(), args[1].Num()))
		},
	)
	// LT_NUM_NUM < :: num -> num -> bool
	LT_NUM_NUM = val.Fun(
		types.Fun(oper.LT, []*types.Type{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.NumLT(args[0].Num(), args[1].Num()))
		},
	)
	// LE_NUM_NUM <= :: num -> num -> bool
	LE_NUM_NUM = val.Fun(
		types.Fun(oper.LE, []*types.Type{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.NumLE(args[0].Num(), args[1].Num()))
		},
	)

	// GT_TIME_TIME > :: time -> time -> bool
	GT_TIME_TIME = val.Fun(
		types.Fun(oper.GT, []*types.Type{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.After(args[1].Time().V))
		},
	)
	// GE_TIME_TIME >= :: time -> time -> bool
	GE_TIME_TIME = val.Fun(
		types.Fun(oper.GE, []*types.Type{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.After(args[1].Time().V) || args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	// LT_TIME_TIME < :: time -> time -> bool
	LT_TIME_TIME = val.Fun(
		types.Fun(oper.LT, []*types.Type{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Before(args[1].Time().V))
		},
	)
	// LE_TIME_TIME <= :: time -> time -> bool
	LE_TIME_TIME = val.Fun(
		types.Fun(oper.LE, []*types.Type{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Before(args[1].Time().V) || args[0].Time().V.Equal(args[1].Time().V))
		},
	)
)
