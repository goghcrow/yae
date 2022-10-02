package fun

import (
	"github.com/goghcrow/yae/token"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
)

// 逻辑操作符支持短路, 都需要声明成 lazy, 手动控制 eval

// 关于 if :: forall α. (bool -> α -> a -> α)
// if 可以通过定义关键词实现成语言结构, 也可以定义成普通 Lazy 函数(类比 scheme 中的特殊表, 或者宏)
// 如果定义成语言结构， 需要在 typecheck 和 compile 特殊处理
// 如果实现成函数, 则需要支持 ∀, 支持泛型

//goland:noinspection GoSnakeCaseUsage
var (
	// IF_BOOL_A_A if :: forall a. (bool -> α -> α -> α)
	// if 可以声明成惰性求值的泛型函数
	IF_BOOL_A_A = func() *val.Val {
		T := types.Slot("a")
		return val.LazyFun(
			types.Fun(token.IF.Name(), []*types.Kind{types.Bool, T, T}, T),
			// 注意 if 是 lazyFun, 参数都是 thunk
			func(args ...*val.Val) *val.Val {
				if args[0].Fun().V().Bool().V {
					return args[1].Fun().V()
				} else {
					return args[2].Fun().V()
				}
			},
		)
	}()
	//AND_BOOL_BOOL and :: bool -> bool -> bool
	AND_BOOL_BOOL = val.LazyFun(
		types.Fun(token.LOGIC_AND.Name(), []*types.Kind{types.Bool, types.Bool}, types.Bool),
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
		types.Fun(token.LOGIC_OR.Name(), []*types.Kind{types.Bool, types.Bool}, types.Bool),
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
		types.Fun(token.LOGIC_NOT.Name(), []*types.Kind{types.Bool}, types.Bool),
		func(args ...*val.Val) *val.Val {
			return val.Bool(!args[0].Bool().V)
		},
	)
)
