package fun

import (
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"math"
)

// 数学计算可以先转换成 big.Int/big.Float 进行计算再转回来
// 或者数字类型, 使用 big.Int/big.Float 表示
// 大部分表达式的场景对精度要求没那么高, 简单处理成 float64

//goland:noinspection GoSnakeCaseUsage
var (
	// ADD_NUM + :: num -> num
	ADD_NUM = val.Fun(
		types.Fun(oper.PLUS, []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return args[0]
		},
	)
	// ADD_NUM_NUM + :: num -> num -> num
	ADD_NUM_NUM = val.Fun(
		types.Fun(oper.PLUS, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V + args[1].Num().V)
		},
	)

	// SUB_NUM - :: num -> num
	SUB_NUM = val.Fun(
		types.Fun(oper.SUB, []*types.Kind{types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(-args[0].Num().V)
		},
	)
	// SUB_NUM_NUM - :: num -> num -> num
	SUB_NUM_NUM = val.Fun(
		types.Fun(oper.SUB, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V - args[1].Num().V)
		},
	)
	// SUB_TIME_TIME - :: time -> time -> num
	SUB_TIME_TIME = val.Fun(
		types.Fun(oper.SUB, []*types.Kind{types.Time, types.Time}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Time().V.Sub(args[1].Time().V).Seconds())
		},
	)

	// MUL_NUM_NUM * :: num -> num -> num
	MUL_NUM_NUM = val.Fun(
		types.Fun(oper.MUL, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V * args[1].Num().V)
		},
	)
	// DIV_NUM_NUM / :: num -> num -> num
	DIV_NUM_NUM = val.Fun(
		types.Fun(oper.DIV, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(args[0].Num().V / args[1].Num().V)
		},
	)
	// MOD_NUM_NUM % :: num -> num -> num
	MOD_NUM_NUM = val.Fun(
		types.Fun(oper.MOD, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(float64(int64(args[0].Num().V) % int64(args[1].Num().V)))
		},
	)
	// EXP_NUM_NUM ^ :: num -> num -> num
	EXP_NUM_NUM = val.Fun(
		types.Fun(oper.EXP, []*types.Kind{types.Num, types.Num}, types.Num),
		func(args ...*val.Val) *val.Val {
			return val.Num(math.Pow(args[0].Num().V, args[1].Num().V))
		},
	)
)
