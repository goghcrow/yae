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
		buf := &strings.Builder{}
		buf.WriteString("(")
		isFst := true
		for _, kind := range val {
			if isFst {
				isFst = false
			} else {
				buf.WriteString(", ")
			}
			buf.WriteString(kind.String())
		}
		buf.WriteString(")")
		return buf.String()
	case TList:
		l := k.List()
		return fmt.Sprintf("list[%s]", l.El)
	case TMap:
		m := k.Map()
		return fmt.Sprintf("map[%s, %s]", m.Key, m.Val)
	case TObj:
		fs := k.Obj().Fields
		buf := &strings.Builder{}
		buf.WriteString("{")
		isFst := true
		for _, f := range fs {
			if isFst {
				isFst = false
			} else {
				buf.WriteString(", ")
			}
			buf.WriteString(f.Name)
			buf.WriteString(": ")
			buf.WriteString(f.Val.String())
		}
		buf.WriteString("}")
		return buf.String()
	case TFun:
		f := k.Fun()
		pt := make([]string, len(f.Param))
		for i, param := range f.Param {
			pt[i] = param.String()
		}
		return fmt.Sprintf("func %s(%s) %s", f.Name, strings.Join(pt, ","), f.Return)
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
