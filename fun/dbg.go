package fun

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// PRINT_ANY print :: forall a. (a -> a)
	PRINT_ANY = func() *val.Val {
		a := types.Slot("a")
		return val.Fun(
			types.Fun("print", []*types.Kind{a}, a),
			func(args ...*val.Val) *val.Val {
				fmt.Println(args[0])
				return args[0]
			},
		)
	}()
)
