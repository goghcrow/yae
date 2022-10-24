package val

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"sort"
	"strconv"
)

func (v *Val) String() string {
	return stringify(v, util.PtrSet{})
}

func stringify(v *Val, inProcess util.PtrSet) string {
	if inProcess.Contains(v) {
		return fmt.Sprintf("recursive-val %s@%p", v.Kind, v)
	} else {
		inProcess.Add(v)
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
		l := v.List()
		xs := make([]string, len(l.V))
		for i, v2 := range l.V {
			xs[i] = stringify(v2, inProcess)
		}
		return util.JoinStrEx(xs, ", ", "[", "]")
	case types.TMap:
		m := v.Map()
		if len(m.V) == 0 {
			return "[:]"
		}
		ord := make([]Key, 0, len(m.V))
		for k, _ := range m.V {
			ord = append(ord, k)
		}
		sort.SliceStable(ord, func(i, j int) bool { return ord[i].val < ord[j].val })
		xs := make([]string, 0, len(m.V))
		for _, k := range ord {
			xs = append(xs, fmt.Sprintf("%s: %s", k, stringify(m.V[k], inProcess)))
		}
		return util.JoinStrEx(xs, ", ", "[", "]")
	case types.TObj:
		o := v.Obj()
		fs := o.Kind.Obj().Fields
		ord := make([]int, len(o.V))
		for i, _ := range o.V {
			ord[i] = i
		}
		sort.SliceStable(ord, func(i, j int) bool { return fs[i].Name < fs[j].Name })
		xs := make([]string, len(o.V))
		for j, i := range ord {
			xs[j] = fmt.Sprintf("%s: %s", fs[i].Name, stringify(o.V[i], inProcess))
		}
		return util.JoinStrEx(xs, ", ", "{", "}")
	case types.TFun:
		return fmt.Sprintf("%s#%p", v.Fun().Kind.String(), v)
	case types.TMaybe:
		mb := v.Maybe()
		ks := mb.Kind.Maybe().Elem.String()
		if mb.V == nil {
			return fmt.Sprintf("Nothing#%s()", ks)
		} else {
			return fmt.Sprintf("Just#%s(%s)", ks, stringify(v.Maybe().V, inProcess))
		}
	default:
		util.Unreachable()
		return ""
	}
}
