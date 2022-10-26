package val

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
)

func (l *ListVal) Set(i int, v *Val) *ListVal {
	util.Assert(types.Equals(l.Type.List().El, v.Type),
		"type mismatched, expect %s actual %s", l.Type.List().El, v.Type)
	l.V[i] = v
	return l
}

func (l *ListVal) Add(vs ...*Val) *ListVal {
	for _, v := range vs {
		util.Assert(types.Equals(l.Type.List().El, v.Type),
			"type mismatched, expect %s actual %s", l.Type.List().El, v.Type)
		l.V = append(l.V, v)
	}
	return l
}
