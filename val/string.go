package val

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"strconv"
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
		return fmt.Sprintf("%s", v.Map().V)
	case types.TObj:
		return fmt.Sprintf("%s", v.Obj().V)
	case types.TFun:
		return v.Fun().Kind.String()
	default:
		util.Unreachable()
	}
	return ""
}
