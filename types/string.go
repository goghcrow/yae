package types

import (
	"fmt"
	"github.com/goghcrow/yae/util"
	"strings"
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
		buf := &strings.Builder{}
		buf.WriteString("(")
		isFst := true
		for _, kind := range val {
			if isFst {
				isFst = false
			} else {
				buf.WriteString(", ")
			}
			buf.WriteString(stringify(kind, lv+1))
		}
		buf.WriteString(")")
		return buf.String()
	case TList:
		l := k.List()
		return fmt.Sprintf("list[%s]", stringify(l.El, lv+1))
	case TMap:
		m := k.Map()
		return fmt.Sprintf("map[%s, %s]", stringify(m.Key, lv+1), stringify(m.Val, lv+1))
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
			buf.WriteString(stringify(f.Val, lv+1))
		}
		buf.WriteString("}")
		return buf.String()
	case TFun:
		f := k.Fun()
		buf := &strings.Builder{}
		buf.WriteString("func ")
		buf.WriteString(f.Name)
		buf.WriteString("(")
		isFst := true
		for _, param := range f.Param {
			if isFst {
				isFst = false
			} else {
				buf.WriteString(", ")
			}
			buf.WriteString(stringify(param, lv+1))
		}
		buf.WriteString(") ")
		buf.WriteString(stringify(f.Return, lv+1))
		return buf.String()
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
