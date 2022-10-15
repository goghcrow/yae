package types

import (
	"fmt"
	"github.com/goghcrow/yae/util"
)

func (k *Kind) String() string {
	return stringify(k, util.PtrSet{})
}

func stringify(k *Kind, inProcess util.PtrSet) string {
	if inProcess.Contains(k) {
		return fmt.Sprintf("recursive-kind %s@%p", k.Type, k)
	} else {
		inProcess.Add(k)
	}

	switch k.Type {
	case TNum:
		return "num"
	case TStr:
		return "str"
	case TBool:
		return "bool"
	case TTime:
		return "time"
	case TTuple:
		val := k.Tuple().Val
		xs := make([]string, len(val))
		for i, kind := range val {
			xs[i] = stringify(kind, inProcess)
		}
		return util.JoinStrEx(xs, ", ", "(", ")")
	case TList:
		l := k.List()
		return fmt.Sprintf("list[%s]", stringify(l.El, inProcess))
	case TMap:
		m := k.Map()
		return fmt.Sprintf("map[%s, %s]", stringify(m.Key, inProcess), stringify(m.Val, inProcess))
	case TObj:
		fs := k.Obj().Fields
		xs := make([]string, len(fs))
		for i, f := range fs {
			xs[i] = fmt.Sprintf("%s: %s", f.Name, stringify(f.Val, inProcess))
		}
		return util.JoinStrEx(xs, ", ", "{", "}")
	case TFun:
		f := k.Fun()
		xs := make([]string, len(f.Param))
		for i, p := range f.Param {
			xs[i] = stringify(p, inProcess)
		}
		pre := "func " + f.Name + "("
		post := ") " + stringify(f.Return, inProcess)
		return util.JoinStrEx(xs, ", ", pre, post)
	case TSlot:
		return k.Slot().Name
	case TTop:
		return "⊤"
	case TBottom:
		return "⊥"
	default:
		util.Unreachable()
		return ""
	}
}
