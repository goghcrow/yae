package fun

import (
	"fmt"

	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
)

//goland:noinspection GoSnakeCaseUsage
var (
	// LEN_LIST len :: forall a. (list[a] -> num)
	LEN_LIST = func() *val.Val {
		T := types.TyVar("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(LEN, []*types.Type{listT}, types.Num),
			func(args ...*val.Val) *val.Val {
				return val.Num(float64(len(args[0].List().V)))
			},
		)
	}()
	// GET_LIST_NUM_ANY get :: forall a. (list[a] -> num -> a -> a)
	GET_LIST_NUM_ANY = func() *val.Val {
		a := types.TyVar("a")
		listA := types.List(a)
		return val.Fun(
			types.Fun(GET, []*types.Type{listA, types.Num, a}, a),
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
	// UNION_LIST_LIST union :: forall a. (list[a] -> list[a] -> list[a])
	UNION_LIST_LIST = func() *val.Val {
		T := types.TyVar("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(UNION, []*types.Type{listT, listT}, listT),
			func(args ...*val.Val) *val.Val {
				lhs := valSetOf(args[0].List().V)
				rhs := valSetOf(args[1].List().V)
				res := val.List(args[0].List().Type.List(), 0).List()
				res.V = union(lhs, rhs)
				fmt.Println(res)
				return res.Vl()
			},
		)
	}()
	// INTERSECT_LIST_LIST intersect :: forall a. (list[a] -> list[a] -> list[a])
	INTERSECT_LIST_LIST = func() *val.Val {
		T := types.TyVar("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(INTERSECT, []*types.Type{listT, listT}, listT),
			func(args ...*val.Val) *val.Val {
				lhs := valSetOf(args[0].List().V)
				rhs := valSetOf(args[1].List().V)
				res := val.List(args[0].List().Type.List(), 0).List()
				res.V = intersect(lhs, rhs)
				return res.Vl()
			},
		)
	}()
	// DIFF_LIST_LIST diff :: forall a. (list[a] -> list[a] -> list[a])
	DIFF_LIST_LIST = func() *val.Val {
		T := types.TyVar("a")
		listT := types.List(T)
		return val.Fun(
			types.Fun(DIFF, []*types.Type{listT, listT}, listT),
			func(args ...*val.Val) *val.Val {
				lhs := valSetOf(args[0].List().V)
				rhs := valSetOf(args[1].List().V)
				res := val.List(args[0].List().Type.List(), 0).List()
				res.V = diff(lhs, rhs)
				return res.Vl()
			},
		)
	}()
)

type valSet struct {
	m    map[string]*val.Val
	link []string
}

func valSetOf(xs []*val.Val) *valSet {
	m := map[string]*val.Val{}
	l := make([]string, 0, len(xs))
	for _, v := range xs {
		hash := v.String()
		if _, ok := m[hash]; !ok {
			m[hash] = v
			l = append(l, hash)
		}
	}
	return &valSet{m, l}
}
func union(x, y *valSet) []*val.Val {
	res := make([]*val.Val, 0)
	for _, k := range x.link {
		res = append(res, x.m[k])
	}
	for _, k := range y.link {
		if _, ok := x.m[k]; !ok {
			res = append(res, y.m[k])
		}
	}
	return res
}
func intersect(x, y *valSet) []*val.Val {
	res := make([]*val.Val, 0)
	for _, k := range x.link {
		if v, ok := y.m[k]; ok {
			res = append(res, v)
		}
	}
	return res
}
func diff(x, y *valSet) []*val.Val {
	res := make([]*val.Val, 0)
	for _, k := range x.link {
		if _, ok := y.m[k]; !ok {
			res = append(res, x.m[k])
		}
	}
	return res
}
