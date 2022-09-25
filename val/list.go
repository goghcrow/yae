package val

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
)

func (l *ListVal) Set(i int, v *Val) {
	util.Assert(types.Equals(l.Kind.List().El, v.Kind),
		"invalid type, expect %s get %s", l.Kind.List().El, v.Kind)
	l.V[i] = v
}

func (l *ListVal) Add(v *Val) {
	util.Assert(types.Equals(l.Kind.List().El, v.Kind),
		"invalid type, expect %s get %s", l.Kind.List().El, v.Kind)
	l.V = append(l.V, v)
}
