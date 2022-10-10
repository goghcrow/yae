package fun

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// STRING_ANY string :: forall a. (a -> str)
	STRING_ANY = func() *val.Val {
		a := types.Slot("a")
		return val.Fun(
			types.Fun(STRING, []*types.Kind{a}, types.Str),
			func(args ...*val.Val) *val.Val {
				return val.Str(stringify(args[0]))
			},
		)
	}()
)
