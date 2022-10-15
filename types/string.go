package types

import (
	"fmt"
	"github.com/goghcrow/yae/util"
)

func (k *Kind) String() string {
	return stringify(k, 0)
}

func stringify(k *Kind, lv int) string {
	if lv > 42 {
		// 可以用 set 精确检查, 这里简单处理
		return "*recursive?*"
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
			xs[i] = stringify(kind, lv+1)
		}
		return util.JoinStr(xs, ", ", "(", ")")
	case TList:
		l := k.List()
		return fmt.Sprintf("list[%s]", stringify(l.El, lv+1))
	case TMap:
		m := k.Map()
		return fmt.Sprintf("map[%s, %s]", stringify(m.Key, lv+1), stringify(m.Val, lv+1))
	case TObj:
		fs := k.Obj().Fields
		xs := make([]string, len(fs))
		for i, f := range fs {
			xs[i] = fmt.Sprintf("%s: %s", f.Name, stringify(f.Val, lv+1))
		}
		return util.JoinStr(xs, ", ", "{", "}")
	case TFun:
		f := k.Fun()
		xs := make([]string, len(f.Param))
		for i, p := range f.Param {
			xs[i] = stringify(p, lv+1)
		}
		pre := "func " + f.Name + "("
		post := ") " + stringify(f.Return, lv+1)
		return util.JoinStr(xs, ", ", pre, post)
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
