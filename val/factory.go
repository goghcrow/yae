package val

import (
	types "github.com/goghcrow/yae/type"
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

func List(kind *types.ListKind, cap int) *Val {
	v := ListVal{Val{kind.Kd()}, make([]*Val, cap)}
	return &v.Val
}

func Map(kind *types.MapKind) *Val {
	v := MapVal{Val{kind.Kd()}, make(map[Key]*Val)}
	return &v.Val
}

func Obj(kind *types.ObjKind) *Val {
	v := ObjVal{Val{kind.Kd()}, make(map[string]*Val)}
	return &v.Val
}

func Fun(kind *types.Kind, f IFun) *Val {
	util.Assert(kind.Type == types.TFun, "expect Fun actual %s", kind)
	fv := FunVal{Val{kind}, f, false}
	return &fv.Val
}

func LazyFun(kind *types.Kind, f IFun) *Val {
	util.Assert(kind.Type == types.TFun, "expect Fun actual %s", kind)
	fv := FunVal{Val{kind}, f, true}
	return &fv.Val
}
