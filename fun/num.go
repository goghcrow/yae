package fun

import (
	"github.com/goghcrow/yae/token"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"math"
)

// 数学计算可以先转换成 big.Int/big.Float 进行计算再转回来
// 或者数字类型, 使用 big.Int/big.Float 表示
// 大部分表达式的场景对精度要求没那么高, 简单处理成 float64

//goland:noinspection GoSnakeCaseUsage
var (
	// PLUS_NUM + :: num -> num
	PLUS_NUM = val.Fun(
		types.Fun(token.PLUS.Name(), []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return args[0]
		},
	)
	// PLUS_NUM_NUM + :: num -> num -> num
	PLUS_NUM_NUM = val.Fun(
		types.Fun(token.PLUS.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V + args[1].Num().V)
		},
	)

	// MINUS_NUM - :: num -> num
	MINUS_NUM = val.Fun(
		types.Fun(token.MINUS.Name(), []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(-args[0].Num().V)
		},
	)
	// MINUS_NUM_NUM - :: num -> num -> num
	MINUS_NUM_NUM = val.Fun(
		types.Fun(token.MINUS.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V - args[1].Num().V)
		},
	)
	// MINUS_TIME_TIME - :: time -> time -> num
	MINUS_TIME_TIME = val.Fun(
		types.Fun(token.MINUS.Name(), []*types.Kind{types.Time, types.Time}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Time().V.Sub(args[1].Time().V).Seconds())
		},
	)

	// MUL_NUM_NUM * :: num -> num -> num
	MUL_NUM_NUM = val.Fun(
		types.Fun(token.MUL.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V * args[1].Num().V)
		},
	)
	// DIV_NUM_NUM / :: num -> num -> num
	DIV_NUM_NUM = val.Fun(
		types.Fun(token.DIV.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V / args[1].Num().V)
		},
	)
	// MOD_NUM_NUM % :: num -> num -> num
	MOD_NUM_NUM = val.Fun(
		types.Fun(token.MOD.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(float64(int64(args[0].Num().V) % int64(args[1].Num().V)))
		},
	)
	// EXP_NUM_NUM ^ :: num -> num -> num
	EXP_NUM_NUM = val.Fun(
		types.Fun(token.EXP.Name(), []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Pow(args[0].Num().V, args[1].Num().V))
		},
	)

	// 目前重载比较简单, 不定参数+重载的分派会变复杂, 暂时可以用重载 n 个参数来解决, 会比较啰嗦, e.g. MAX_NUM_NUM_NUM

	// MAX_NUM_NUM + :: num -> num -> num
	MAX_NUM_NUM = val.Fun(
		types.Fun("max", []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Max(args[0].Num().V, args[1].Num().V))
		},
	)
	// MAX_NUM_NUM_NUM + :: num -> num -> num -> num
	MAX_NUM_NUM_NUM = val.Fun(
		types.Fun("max", []*types.Kind{types.Num, types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(
				math.Max(
					math.Max(args[0].Num().V, args[1].Num().V),
					args[2].Num().V,
				),
			)
		},
	)
	// MIN_NUM_NUM + :: num -> num -> num
	MIN_NUM_NUM = val.Fun(
		types.Fun("min", []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Min(args[0].Num().V, args[1].Num().V))
		},
	)
	// MIN_NUM_NUM_NUM + :: num -> num -> num -> num
	MIN_NUM_NUM_NUM = val.Fun(
		types.Fun("min", []*types.Kind{types.Num, types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(
				math.Min(
					math.Min(args[0].Num().V, args[1].Num().V),
					args[2].Num().V,
				),
			)
		},
	)
)
