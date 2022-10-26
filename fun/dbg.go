package fun

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// PRINT_ANY print :: forall a. (a -> a)
	PRINT_ANY = func() *val.Val {
		a := types.TyVar("a")
		return val.Fun(
			types.Fun(PRINT, []*types.Type{a}, a),
			func(args ...*val.Val) *val.Val {
				fmt.Println(args[0])
				return args[0]
			},
		)
	}()
)
