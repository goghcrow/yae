package fun

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// LEN_LIST == :: forall a. (list[a] -> num)
	LEN_LIST = func() *val.Val {
		T := types.Slot("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun("len", []*types.Kind{listT}, types.Num),
			func(args ...*val.Val) *val.Val {
				return val.Num(float64(len(args[0].List().V)))
			},
		)
	}()
)
