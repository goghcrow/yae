package fun

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"strconv"
	"strings"
)

func stringify(v *val.Val) string {
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
		return v.Str().V
	case types.TTime:
		return v.Time().V.String()
	case types.TList:
		return stringifyList(v.List())
	case types.TMap:
		return stringifyMap(v.Map())
	case types.TObj:
		return stringifyObj(v.Obj())
	case types.TFun:
		return "#fun"
	default:
		util.Unreachable()
		return ""
	}
}

func stringifyList(l *val.ListVal) string {
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
		buf.WriteString(stringify(v))
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyMap(m *val.MapVal) string {
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
		buf.WriteString(stringify(v))
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyObj(v *val.ObjVal) string {
	if len(v.V) == 0 {
		return "{}"
	}

	buf := &strings.Builder{}
	buf.WriteString("{")
	isFst := true
	fs := v.Kind.Obj().Fields
	for i, vl := range v.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(fs[i].Name)
		buf.WriteString(": ")
		buf.WriteString(stringify(vl))
	}
	buf.WriteString("}")
	return buf.String()
}
