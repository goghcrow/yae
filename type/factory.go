package types

import "github.com/goghcrow/yae/util"

func List(el *Kind) *Kind {
	t := ListKind{Kind{TList}, el}
	return &t.Kind
}

func Map(k, v *Kind) *Kind {
	util.Assert(k.IsPrimitive(), "invalid type of map's key: %s", k)
	m := MapKind{Kind{TMap}, k, v}
	return &m.Kind
}

func Obj(fields map[string]*Kind) *Kind {
	t := ObjKind{Kind{TObj}, fields}
	return &t.Kind
}

func Fun(name string, param []*Kind, ret *Kind) *Kind {
	t := FunKind{Kind{TFun}, name, param, ret}
	return &t.Kind
}
