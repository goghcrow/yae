package val

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"strconv"
	"strings"
)

func (v *Val) String() string {
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
		return fmt.Sprintf("%s", v.List().V)
	case types.TMap:
		return stringifyMap(v.Map())
	case types.TObj:
		return stringifyObj(v.Obj())
	case types.TFun:
		return v.Fun().Kind.String()
	default:
		util.Unreachable()
		return ""
	}
}

func stringifyMap(m *MapVal) string {
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
		buf.WriteString(v.String())
	}
	buf.WriteString("]")
	return buf.String()
}

func stringifyObj(o *ObjVal) string {
	if len(o.V) == 0 {
		return "{}"
	}
	buf := &strings.Builder{}
	buf.WriteString("{")
	isFst := true
	for name, val := range o.V {
		if isFst {
			isFst = false
		} else {
			buf.WriteString(", ")
		}
		buf.WriteString(name)
		buf.WriteString(": ")
		buf.WriteString(val.String())
	}
	buf.WriteString("}")
	return buf.String()
}
