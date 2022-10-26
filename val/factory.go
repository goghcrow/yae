package val

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"time"
)

func Bool(b bool) *Val {
	if b {
		return True
	} else {
		return False
	}
}

func Num(n float64) *Val {
	v := NumVal{Val{types.Num}, n}
	return &v.Val
}

func Str(s string) *Val {
	v := StrVal{Val{types.Str}, s}
	return &v.Val
}

func Time(t time.Time) *Val {
	v := TimeVal{Val{types.Time}, t}
	return &v.Val
}

func List(ty *types.ListTy, cap int) *Val {
	v := ListVal{Val{ty.Ty()}, make([]*Val, cap)}
	return &v.Val
}

func Map(ty *types.MapTy) *Val {
	v := MapVal{Val{ty.Ty()}, make(map[Key]*Val)}
	return &v.Val
}

func Obj(ty *types.ObjTy) *Val {
	v := ObjVal{Val{ty.Ty()}, make([]*Val, len(ty.Fields))}
	return &v.Val
}

func Fun(ty *types.Type, f IFun) *Val {
	util.Assert(ty.Kind == types.KFun, "expect Fun actual %s", ty)
	fv := FunVal{Val{ty}, f, false}
	return &fv.Val
}

func LazyFun(ty *types.Type, f IFun) *Val {
	util.Assert(ty.Kind == types.KFun, "expect Fun actual %s", ty)
	fv := FunVal{Val{ty}, f, true}
	return &fv.Val
}

func Nothing(elem *types.Type) *Val {
	mb := MaybeVal{Val{types.Maybe(elem)}, nil}
	return &mb.Val
}

func Just(elem *types.Type, v *Val) *Val {
	util.Assert(types.Equals(elem, v.Type), "expect %s actual %s", elem, v.Type)
	mb := MaybeVal{Val{types.Maybe(elem)}, v}
	return &mb.Val
}
