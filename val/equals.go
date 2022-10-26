package val

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
)

func Equals(x, y *Val) bool {
	if x == y {
		return true
	}
	if x == nil || y == nil {
		return false
	}
	if !types.Equals(x.Type, y.Type) {
		return false
	}

	switch x.Type.Kind {
	case types.KNum:
		return NumEQ(x.Num(), y.Num())
	case types.KBool:
		return x.Bool().V == y.Bool().V
	case types.KStr:
		return x.Str().V == y.Str().V
	case types.KTime:
		return x.Time().V.Equal(y.Time().V)
	case types.KList:
		return equalsList(x.List(), y.List())
	case types.KMap:
		return equalsMap(x.Map(), y.Map())
	case types.KObj:
		return equalsObj(x.Obj(), y.Obj())
	case types.KFun:
		// panic("fun is not comparable")
		return x == y
	case types.KMaybe:
		return equalsMaybe(x.Maybe(), y.Maybe())
	default:
		util.Unreachable()
		return false
	}
}

func equalsList(x, y *ListVal) bool {
	if len(x.V) != len(y.V) {
		return false
	}
	for i := 0; i < len(x.V); i++ {
		if !Equals(x.V[i], y.V[i]) {
			return false
		}
	}
	return true
}

func equalsMap(x, y *MapVal) bool {
	if len(x.V) != len(y.V) {
		return false
	}
	for k, v1 := range x.V {
		v2, ok := y.V[k]
		if !ok || !Equals(v1, v2) {
			return false
		}
	}
	return true
}

func equalsObj(x, y *ObjVal) bool {
	if len(x.V) != len(y.V) {
		return false
	}

	xk := x.Obj().Type.Obj()
	for i, v1 := range x.V {
		// 前置判断过类型相等, v2 一定存在
		v2, _ := y.Get(xk.Fields[i].Name)
		if !Equals(v1, v2) {
			return false
		}
	}
	return true
}

func equalsMaybe(x, y *MaybeVal) bool {
	if !types.Equals(x.Type.Maybe().Elem, y.Type.Maybe().Elem) {
		return false
	}
	if x.V == nil && y.V == nil {
		return true
	}
	if x.V == nil || y.V == nil {
		return false
	}
	return Equals(x.V, y.V)
}
