package types

import "github.com/goghcrow/yae/util"

func (o *ObjTy) MustGetField(name string) *Field {
	f, ok := o.GetField(name)
	util.Assert(ok, "missing field %s", name)
	return f
}

func (o *ObjTy) GetField(name string) (*Field, bool) {
	i, ok := o.Index[name]
	if !ok {
		return nil, false
	}
	return &o.Fields[i], true
}
