package fun

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoUnusedGlobalVariable,GoSnakeCaseUsage
var (
	// LEN_MAP == :: forall a. (map[k, v] -> num)
	LEN_MAP = func() *val.Val {
		K := types.Slot("k")
		V := types.Slot("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun("len", []*types.Kind{mapKV}, types.Num),
			func(args ...*val.Val) *val.Val {
				return val.Num(float64(len(args[0].Map().V)))
			},
		)
	}()
)
