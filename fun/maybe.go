package fun

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// GET_MAYBE get :: forall a. (maybe[a] -> a -> a)
	GET_MAYBE = func() *val.Val {
		T := types.Slot("a")
		maybeT := types.Maybe(T)
		return val.Fun(
			types.Fun(GET, []*types.Kind{maybeT, T}, T),
			func(args ...*val.Val) *val.Val {
				mb := args[0].Maybe()
				if mb.V == nil {
					return args[1]
				} else {
					return mb.V
				}
			},
		)
	}()
)
