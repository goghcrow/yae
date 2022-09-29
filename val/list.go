package val

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
)

func (l *ListVal) Set(i int, v *Val) *ListVal {
	util.Assert(types.Equals(l.Kind.List().El, v.Kind),
		"invalid type, expect `%s` actual `%s`", l.Kind.List().El, v.Kind)
	l.V[i] = v
	return l
}

func (l *ListVal) Add(vs ...*Val) *ListVal {
	for _, v := range vs {
		util.Assert(types.Equals(l.Kind.List().El, v.Kind),
			"invalid type, expect `%s` actual `%s`", l.Kind.List().El, v.Kind)
		l.V = append(l.V, v)
	}
	return l
}
