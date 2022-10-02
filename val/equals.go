package val

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"math"
)

func Equals(x, y *Val) bool {
	if x == nil || y == nil {
		return false
	}
	if !types.Equals(x.Kind, y.Kind) {
		return false
	}

	switch x.Kind.Type {
	case types.TNum:
		return equalsNum(x.Num(), y.Num())
	case types.TBool:
		return x.Bool().V == y.Bool().V
	case types.TStr:
		return x.Str().V == y.Str().V
	case types.TTime:
		return x.Time().V.Equal(y.Time().V)
	case types.TList:
		return equalsList(x.List(), y.List())
	case types.TMap:
		return equalsMap(x.Map(), y.Map())
	case types.TObj:
		return equalsObj(x.Obj(), y.Obj())
	case types.TFun:
		panic("fun is not comparable")
	default:
		util.Unreachable()
	}
	return false
}

const epsilon = 1e-9

func equalsNum(x, y *NumVal) bool {
	return math.Abs(x.V-y.V) < epsilon
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
	for k, v1 := range x.V {
		v2, ok := y.V[k]
		if !ok || !Equals(v1, v2) {
			return false
		}
	}
	return true
}
