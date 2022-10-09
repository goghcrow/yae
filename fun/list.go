package fun

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// LEN_LIST len :: forall a. (list[a] -> num)
	LEN_LIST = func() *val.Val {
		T := types.Slot("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(LEN, []*types.Kind{listT}, types.Num),
			func(args ...*val.Val) *val.Val {
				return val.Num(float64(len(args[0].List().V)))
			},
		)
	}()

	// todo 集合运算
)
