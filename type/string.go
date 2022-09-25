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
	case TList:
		l := k.List()
		if l.El == Hole {
			return "list"
		} else {
			return fmt.Sprintf("list[%s]", l.El)
		}
	case TMap:
		m := k.Map()
		if m.Key == Hole && m.Val == Hole {
			return "map"
		} else {
			return fmt.Sprintf("map[%s, %s]", m.Key, m.Val)
		}
	case TObj:
		fs := k.Obj().Fields
		buf := "{"
		for name, kind := range fs {
			buf += fmt.Sprintf("%s: %s, ", name, kind)
		}
		return buf + "}"
	case TFun:
		f := k.Fun()
		pt := make([]string, len(f.Param))
		for i, param := range f.Param {
			pt[i] = param.String()
		}
		return fmt.Sprintf("%s(%s) %s", f.Name, strings.Join(pt, ","), f.Return)
	case THold:
		return ""
	case TTop:
		return "⊤"
	case TBottom:
		return "⊥"
	default:
		util.Unreachable()
	}
	return ""
}
