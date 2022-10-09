package val

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"strconv"
	"strings"
)

func (v *Val) String() string {
	return stringify(v, 0)
}

func stringify(v *Val, lv int) string {
	if lv > 42 {
		// 可以用 set 精确检查, 这里简单处理
		return "*recursive?*"
	}

	switch v.Kind.Type {
	case types.TNum:
		n := v.Num()
		if n.IsInt() {
			return fmt.Sprintf("%d", n.Int())
		} else {
			return fmt.Sprintf("%f", n.V)
		}
	case types.TBool:
		return strconv.FormatBool(v.Bool().V)
	case types.TStr:
		return fmt.Sprintf("%q", v.Str().V)
	case types.TTime:
		return v.Time().V.String()
	case types.TList:
		return stringifyList(v.List(), lv)
	case types.TMap:
		return stringifyMap(v.Map(), lv)
	case types.TObj:
		return stringifyObj(v.Obj(), lv)
	case types.TFun:
		return v.Fun().Kind.String()
	default:
		util.Unreachable()
		return ""
	}
}

func stringifyList(l *ListVal, lv int) string {
	if len(l.V) == 0 {
		return "[]"
	}
	buf := &strings.Builder{}
	buf.WriteString("[")
	isFst := true
	for _, v := range l.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(stringify(v, lv+1))
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyMap(m *MapVal, lv int) string {
	if len(m.V) == 0 {
		return "[:]"
	}
	buf := &strings.Builder{}
	buf.WriteString("[")
	isFst := true
	for k, v := range m.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(k.String())
		buf.WriteString(": ")
		buf.WriteString(stringify(v, lv+1))
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyObj(o *ObjVal, lv int) string {
	if len(o.V) == 0 {
		return "{}"
	}
	buf := &strings.Builder{}
	buf.WriteString("{")
	isFst := true
	fs := o.Kind.Obj().Fields
	for i, val := range o.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(fs[i].Name)
		buf.WriteString(": ")
		buf.WriteString(stringify(val, lv+1))
	}
	buf.WriteString("}")
	return buf.String()
}
