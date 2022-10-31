package types

import (
	"strconv"

	"github.com/goghcrow/yae/util"
)

// TyVar 新建 Type Variable
// 📢 每次调用都生成全局唯一类型变量
var TyVar = func() func(name string) *Type {
	n := 0
	return func(name string) *Type {
		n++
		t := TypeVariable{Type{KTyVar}, name + strconv.Itoa(n)}
		return &t.Type
	}
}()

func Tuple(val []*Type) *Type {
	t := TupleTy{Type{kTuple}, val}
	return &t.Type
}

func List(el *Type) *Type {
	t := ListTy{Type{KList}, el}
	return &t.Type
}

func Map(k, v *Type) *Type {
	util.Assert(keyable(k),
		"invalid type of map's key: %s", k)
	m := MapTy{Type{KMap}, k, v}
	return &m.Type
}

func keyable(ty *Type) bool { return ty.IsPrimitive() || ty.Kind == KTyVar || ty.Kind == KBot }

func Obj(fields []Field) *Type {
	t := ObjTy{Type{KObj}, fields, nil}
	t.Index = make(map[string]int, len(fields))
	for i, f := range fields {
		j, ok := t.Index[f.Name]
		util.Assert(!ok, "duplicated field %s in %d and %d", f.Name, i, j)
		t.Index[f.Name] = i
	}
	return &t.Type
}

func Fun(name string, param []*Type, ret *Type) *Type {
	t := FunTy{Type{KFun}, name, param, ret}
	return &t.Type
}

func Maybe(elem *Type) *Type {
	t := MaybeTy{Type{KMaybe}, elem}
	return &t.Type
}
