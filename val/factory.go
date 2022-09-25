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

func List(el *types.Kind) *Val {
	v := ListVal{Val{types.List(el)}, make([]*Val, 0)}
	return &v.Val
}

func Map(kk, vk *types.Kind) *Val {
	v := MapVal{Val{types.Map(kk, vk)}, make(map[Key]*Val)}
	return &v.Val
}

func Obj(fs map[string]*types.Kind) *Val {
	v := ObjVal{Val{types.Obj(fs)}, make(map[string]*Val)}
	return &v.Val
}

func Fun(kind *types.Kind, f IFun) *Val {
	util.Assert(kind.Type == types.TFun, "type error")
	fv := FunVal{Val{kind}, f, false}
	return &fv.Val
}

func LazyFun(kind *types.Kind, f IFun) *Val {
	util.Assert(kind.Type == types.TFun, "type error")
	fv := FunVal{Val{kind}, f, true}
	return &fv.Val
}
