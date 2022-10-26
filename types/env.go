package types

import (
	"github.com/goghcrow/yae/util"
)

// Env for typeChecker
type Env struct {
	parent *Env
	ctx    map[string]*Type
	fnTbl  map[string]interface{} // *FunTy | []*FunTy
}

func NewEnv() *Env {
	return &Env{nil, map[string]*Type{}, map[string]interface{}{}}
}

func (e *Env) Inherit(parent *Env) *Env {
	util.Assert(e.parent == nil, "env.parent != nil")
	e.parent = parent
	return e
}

func (e *Env) Derive() *Env {
	return &Env{e, map[string]*Type{}, map[string]interface{}{}}
}

func (e *Env) Get(name string) (*Type, bool) {
	v, ok := e.ctx[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return v, ok
}

func (e *Env) Put(name string, val *Type) {
	util.Assert(slotFree(val), "expect mono type actual %s", val)
	// 注意这里只修改当前环境, 不修改继承
	// 如果是 scope 语义, 需要先 env:=findDefEnv(name) 然后 env.ctx[name]=val
	e.ctx[name] = val
}

func (e *Env) ForEach(f func(string, *Type)) {
	for k, v := range e.ctx {
		f(k, v)
	}
}

func (e *Env) RegisterFun(f *Type) {
	util.Assert(f.Kind == KFun, "expect FunTy actual %s", f)
	lookup, mono := f.Fun().Lookup()
	if mono {
		util.Assert(slotFree(f), "expect mono type actual %s", f)
		e.fnTbl[lookup] = f.Fun()
	} else {
		tbl, ok := e.fnTbl[lookup].([]*FunTy)
		if !ok {
			tbl = []*FunTy{}
		}
		tbl = append(tbl, f.Fun())
		e.fnTbl[lookup] = tbl
	}
}

func (e *Env) GetMonoFun(name string) (*FunTy, bool) {
	fk, ok := e.fnTbl[name].(*FunTy)
	if ok {
		return fk, true
	} else if e.parent != nil {
		return e.parent.GetMonoFun(name)
	} else {
		return nil, false
	}
}

func (e *Env) GetPolyFuns(name string) ([]*FunTy, bool) {
	fk, ok := e.fnTbl[name].([]*FunTy)
	if ok {
		return fk, true
	} else if e.parent != nil {
		return e.parent.GetPolyFuns(name)
	} else {
		return nil, false
	}
}
