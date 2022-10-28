package val

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
)

// Env for compile and runtime
type Env struct {
	parent *Env
	ctx    map[string]*Val
	fnTbl  map[string]interface{} // *FunVal | []*FunVal
}

func NewEnv() *Env {
	return &Env{nil, map[string]*Val{}, map[string]interface{}{}}
}

func (e *Env) Inherit(parent *Env) *Env {
	util.Assert(e.parent == nil, "env.parent != nil")
	e.parent = parent
	return e
}

func (e *Env) Derive() *Env {
	return &Env{e, map[string]*Val{}, map[string]interface{}{}}
}

func (e *Env) MustGet(name string) *Val {
	v, ok := e.Get(name)
	util.Assert(ok, "missing `%s` in env", name)
	return v
}

func (e *Env) Get(name string) (*Val, bool) {
	v, ok := e.ctx[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return v, ok
}

func (e *Env) Put(name string, val *Val) {
	// 注意这里只修改当前环境, 不修改继承
	// 如果是 scope 语义, 需要先 env:=findDefEnv(name) 然后 env.ctx[name]=val
	e.ctx[name] = val
}

func (e *Env) ForEach(f func(string, *Val)) {
	for k, v := range e.ctx {
		f(k, v)
	}
}

func (e *Env) RegisterFun(f *Val) {
	util.Assert(f.Type.Kind == types.KFun, "expect FunVal actual %s", f)
	lookup, mono := f.Type.Fun().Lookup()
	if mono {
		e.fnTbl[lookup] = f.Fun()
	} else {
		tbl, ok := e.fnTbl[lookup].([]*FunVal)
		if !ok {
			tbl = []*FunVal{}
		}
		tbl = append(tbl, f.Fun())
		e.fnTbl[lookup] = tbl
	}
}

func (e *Env) GetMonoFun(name string) (*FunVal, bool) {
	fk, ok := e.fnTbl[name].(*FunVal)
	if ok {
		return fk, true
	} else if e.parent != nil {
		return e.parent.GetMonoFun(name)
	} else {
		return nil, false
	}
}

func (e *Env) MustGetMonoFun(name string) *FunVal {
	fv, ok := e.GetMonoFun(name)
	if !ok {
		panic(fmt.Errorf("%s not defined", name))
	}
	return fv
}

func (e *Env) GetPolyFuns(name string) ([]*FunVal, bool) {
	fk, ok := e.fnTbl[name].([]*FunVal)
	if ok {
		return fk, true
	} else if e.parent != nil {
		return e.parent.GetPolyFuns(name)
	} else {
		return nil, false
	}
}

func (e *Env) MustGetPolyFuns(name string) []*FunVal {
	fv, ok := e.GetPolyFuns(name)
	if !ok {
		panic(fmt.Errorf("%s not defined", name))
	}
	return fv
}
