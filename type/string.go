package types

import (
	"fmt"
	"github.com/goghcrow/yae/util"
	"strings"
)

func (k *Kind) String() string {
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
		buf := "("
		isFirst := true
		for _, kind := range val {
			if isFirst {
				buf += kind.String()
				isFirst = false
			} else {
				buf += fmt.Sprintf(", %s", kind)
			}
		}
		return buf + ")"
	case TList:
		l := k.List()
		return fmt.Sprintf("list[%s]", l.El)
	case TMap:
		m := k.Map()
		return fmt.Sprintf("map[%s, %s]", m.Key, m.Val)
	case TObj:
		fs := k.Obj().Fields
		buf := "{"
		isFirst := true
		for name, kind := range fs {
			if isFirst {
				buf += fmt.Sprintf("%s: %s", name, kind)
				isFirst = false
			} else {
				buf += fmt.Sprintf(", %s: %s", name, kind)
			}
		}
		return buf + "}"
	case TFun:
		f := k.Fun()
		pt := make([]string, len(f.Param))
		for i, param := range f.Param {
			pt[i] = param.String()
		}
		return fmt.Sprintf("%s(%s) %s", f.Name, strings.Join(pt, ","), f.Return)
	case TSlot:
		return k.Slot().Name
	case TTop:
		return "⊤"
	case TBottom:
		return "⊥"
	default:
		util.Unreachable()
	}
	return ""
}
