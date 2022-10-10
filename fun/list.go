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
	// GET_LIST_NUM_ANY get :: forall a. (list[a] -> num -> a -> a)
	GET_LIST_NUM_ANY = func() *val.Val {
		a := types.Slot("a")
		listA := types.List(a)
		return val.Fun(
			types.Fun(GET, []*types.Kind{listA, types.Num, a}, a),
			func(args ...*val.Val) *val.Val {
				lst := args[0].List().V
				idx := int(args[1].Num().V)
				defVl := args[2]
				if idx >= len(lst) {
					return defVl
				}
				el := lst[idx]
				if el == nil {
					return defVl
				}
				return el
			},
		)
	}()

	// todo 集合运算
)
