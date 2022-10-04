package fun

import (
	"github.com/goghcrow/yae/oper"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// EQ_BOOL_BOOL == :: bool -> bool -> bool
	EQ_BOOL_BOOL = val.Fun(
		types.Fun(oper.EQ, []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Bool().V == args[1].Bool().V)
		},
	)
	// NE_BOOL_BOOL != :: bool -> bool -> bool
	NE_BOOL_BOOL = val.Fun(
		types.Fun(oper.NE, []*types.Kind{types.Bool, types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Bool().V != args[1].Bool().V)
		},
	)
	// EQ_NUM_NUM == :: num -> num -> bool
	EQ_NUM_NUM = val.Fun(
		types.Fun(oper.EQ, []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(val.Equals(args[0], args[1]))
		},
	)
	// NE_NUM_NUM != :: num -> num -> bool
	NE_NUM_NUM = val.Fun(
		types.Fun(oper.NE, []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!val.Equals(args[0], args[1]))
		},
	)
	// EQ_STR_STR == :: str -> str -> bool
	EQ_STR_STR = val.Fun(
		types.Fun(oper.EQ, []*types.Kind{types.Str, types.Str}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Str().V == args[1].Str().V)
		},
	)
	// NE_STR_STR != :: str -> str -> bool
	NE_STR_STR = val.Fun(
		types.Fun(oper.NE, []*types.Kind{types.Str, types.Str}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Str().V != args[1].Str().V)
		},
	)
	// EQ_TIME_TIME == :: time -> time -> bool
	EQ_TIME_TIME = val.Fun(
		types.Fun(oper.EQ, []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	// NE_TIME_TIME != :: time -> time -> bool
	NE_TIME_TIME = val.Fun(
		types.Fun(oper.NE, []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	//EQ_LIST_LIST == :: forall a. (list[a] -> list[a] -> bool)
	EQ_LIST_LIST = func() *val.Val {
		T := types.Slot("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(oper.EQ, []*types.Kind{listT, listT}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(val.Equals(args[0], args[1]))
			},
		)
	}()
	//NE_LIST_LIST != :: forall a. (list[a] -> list[a] -> bool)
	NE_LIST_LIST = func() *val.Val {
		T := types.Slot("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(oper.NE, []*types.Kind{listT, listT}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(!val.Equals(args[0], args[1]))
			},
		)
	}()
	//EQ_MAP_MAP == :: forall k v . (map[k,v] -> map[k,v] -> bool)
	EQ_MAP_MAP = func() *val.Val {
		K := types.Slot("k")
		V := types.Slot("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(oper.EQ, []*types.Kind{mapKV, mapKV}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(val.Equals(args[0], args[1]))
			},
		)
	}()
	//NE_MAP_MAP != :: forall k v . (map[k,v] -> map[k,v] -> bool)
	NE_MAP_MAP = func() *val.Val {
		K := types.Slot("k")
		V := types.Slot("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(oper.NE, []*types.Kind{mapKV, mapKV}, types.Bool),
			func(args ...*val.Val) *val.Val {
				return val.Bool(!val.Equals(args[0], args[1]))
			},
		)
	}()

	// GT_NUM_NUM > :: num -> num -> bool
	GT_NUM_NUM = val.Fun(
		types.Fun(oper.GT, []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V > args[1].Num().V && !val.Equals(args[0], args[1]))
		},
	)
	// GE_NUM_NUM >= :: num -> num -> bool
	GE_NUM_NUM = val.Fun(
		types.Fun(oper.GE, []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V > args[1].Num().V || val.Equals(args[0], args[1]))
		},
	)
	// LT_NUM_NUM < :: num -> num -> bool
	LT_NUM_NUM = val.Fun(
		types.Fun(oper.LT, []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V < args[1].Num().V && !val.Equals(args[0], args[1]))
		},
	)
	// LE_NUM_NUM <= :: num -> num -> bool
	LE_NUM_NUM = val.Fun(
		types.Fun(oper.LE, []*types.Kind{types.Num, types.Num}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Num().V < args[1].Num().V || val.Equals(args[0], args[1]))
		},
	)

	// GT_TIME_TIME > :: time -> time -> bool
	GT_TIME_TIME = val.Fun(
		types.Fun(oper.GT, []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.After(args[1].Time().V))
		},
	)
	// GE_TIME_TIME >= :: time -> time -> bool
	GE_TIME_TIME = val.Fun(
		types.Fun(oper.GE, []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.After(args[1].Time().V) || args[0].Time().V.Equal(args[1].Time().V))
		},
	)
	// LT_TIME_TIME < :: time -> time -> bool
	LT_TIME_TIME = val.Fun(
		types.Fun(oper.LT, []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Before(args[1].Time().V))
		},
	)
	// LE_TIME_TIME <= :: time -> time -> bool
	LE_TIME_TIME = val.Fun(
		types.Fun(oper.LE, []*types.Kind{types.Time, types.Time}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(args[0].Time().V.Before(args[1].Time().V) || args[0].Time().V.Equal(args[1].Time().V))
		},
	)
)
