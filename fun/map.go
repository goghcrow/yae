package fun

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoUnusedGlobalVariable,GoSnakeCaseUsage
var (
	// LEN_MAP len :: forall k v. (map[k, v] -> num)
	LEN_MAP = func() *val.Val {
		K := types.Slot("k")
		V := types.Slot("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(LEN, []*types.Kind{mapKV}, types.Num),
			func(args ...*val.Val) *val.Val {
				return val.Num(float64(len(args[0].Map().V)))
			},
		)
	}()
	// ISSET_MAP_K isset :: forall k v. (map[k, v] -> k -> bool)
	ISSET_MAP_K = func() *val.Val {
		K := types.Slot("k")
		V := types.Slot("v")
		mapKV := types.Map(K, V)
		return val.Fun(
			types.Fun(ISSET, []*types.Kind{mapKV, K}, types.Bool),
			func(args ...*val.Val) *val.Val {
				_, ok := args[0].Map().V[args[1].Key()]
				return val.Bool(ok)
			},
		)
	}()
)
