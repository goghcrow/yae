package fun

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"strconv"
)

func stringify(v *val.Val) string {
	return stringify0(v, util.PtrSet{})
}

func stringify0(v *val.Val, inProcess util.PtrSet) string {
	if v.Type.Kind.IsComposite() {
		if inProcess.Contains(v) {
			return fmt.Sprintf("recursive-val %s@%p", v.Type, v)
		} else {
			inProcess.Add(v)
		}
	}

	switch v.Type.Kind {
	case types.KNum:
		n := v.Num()
		if n.IsInt() {
			return util.FmtInt(n.Int())
		} else {
			return util.FmtFloat(n.V)
		}
	case types.KBool:
		return strconv.FormatBool(v.Bool().V)
	case types.KStr:
		return v.Str().V
	case types.KTime:
		return v.Time().V.String()
	case types.KList:
		l := v.List()
		xs := make([]string, len(l.V))
		for i, v2 := range l.V {
			xs[i] = stringify0(v2, inProcess)
		}
		return util.JoinStr(xs, ", ", "[", "]")
	case types.KMap:
		m := v.Map()
		if len(m.V) == 0 {
			return "[:]"
		}
		xs := make([]string, 0, len(m.V))
		for k, v := range m.V {
			xs = append(xs, fmt.Sprintf("%s: %s", k, stringify0(v, inProcess)))
		}
		return util.JoinStr(xs, ", ", "[", "]")
	case types.KObj:
		o := v.Obj()
		fs := o.Type.Obj().Fields
		xs := make([]string, len(o.V))
		for i, v2 := range o.V {
			xs[i] = fmt.Sprintf("%s: %s", fs[i].Name, stringify0(v2, inProcess))
		}
		return util.JoinStr(xs, ", ", "{", "}")
	case types.KFun:
		return "#fun"
	case types.KMaybe:
		if v.Maybe().V == nil {
			return "Nothing()"
		} else {
			return fmt.Sprintf("Just(%s)", stringify0(v.Maybe().V, inProcess))
		}
	default:
		util.Unreachable()
		return ""
	}
}
