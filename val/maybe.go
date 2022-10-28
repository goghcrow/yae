package val

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
)

func (v *MaybeVal) GetOrDefault(defVal *Val) *Val {
	if v.V == nil {
		if false {
			ty := v.Type.Maybe().Elem
			util.Assert(types.Equals(ty, defVal.Type),
				"expect %s actual %s %s", ty, defVal.Type, defVal)
		}
		return defVal
	} else {
		return v.V
	}
}
