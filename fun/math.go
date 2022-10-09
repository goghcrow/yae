package fun

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"math"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// 目前重载比较简单, 不定参数+重载的分派会变复杂, 暂时可以用重载 n 个参数来解决, 会比较啰嗦, e.g. MAX_NUM_NUM_NUM

	// MAX_NUM_NUM max :: num -> num -> num
	MAX_NUM_NUM = val.Fun(
		types.Fun(MAX, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Max(args[0].Num().V, args[1].Num().V))
		},
	)
	// MAX_LIST max :: list[num] -> num
	MAX_LIST = val.Fun(
		types.Fun(MAX, []*types.Kind{types.List(types.Num)}, types.Num),
		func(args ...*val.Val) *val.Val {
			lst := args[0].List().V
			if len(lst) == 0 {
				return val.Num(0)
			}
			max := lst[0].Num().V
			for i := 1; i < len(lst); i++ {
				max = math.Max(max, lst[i].Num().V)
			}
			return val.Num(max)
		},
	)
	// MIN_NUM_NUM min :: num -> num -> num
	MIN_NUM_NUM = val.Fun(
		types.Fun(MIN, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Min(args[0].Num().V, args[1].Num().V))
		},
	)
	// MIN_LIST min :: list[num] -> num
	MIN_LIST = val.Fun(
		types.Fun(MIN, []*types.Kind{types.List(types.Num)}, types.Num),
		func(args ...*val.Val) *val.Val {
			lst := args[0].List().V
			if len(lst) == 0 {
				return val.Num(0)
			}
			min := lst[0].Num().V
			for i := 1; i < len(lst); i++ {
				min = math.Min(min, lst[i].Num().V)
			}
			return val.Num(min)
		},
	)

	// ABS_NUM abs :: num -> num
	ABS_NUM = val.Fun(
		types.Fun(ABS, []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Abs(args[0].Num().V))
		},
	)
	// CEIL_NUM ceil :: num -> num
	CEIL_NUM = val.Fun(
		types.Fun(CEIL, []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Ceil(args[0].Num().V))
		},
	)
	// FLOOR_NUM floor :: num -> num
	FLOOR_NUM = val.Fun(
		types.Fun(FLOOR, []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Floor(args[0].Num().V))
		},
	)
	// ROUND_NUM round :: num -> num
	ROUND_NUM = val.Fun(
		types.Fun(ROUND, []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Round(args[0].Num().V))
		},
	)
)
