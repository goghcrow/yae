package fun

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// GET_MAYBE get :: forall a. (maybe[a] -> a -> a)
	GET_MAYBE = func() *val.Val {
		T := types.TyVar("a")
		maybeT := types.Maybe(T)
		return val.Fun(
			types.Fun(GET, []*types.Type{maybeT, T}, T),
			func(args ...*val.Val) *val.Val {
				return args[0].Maybe().GetOrDefault(args[1])
			},
		)
	}()
)
