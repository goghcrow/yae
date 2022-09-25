package env

import (
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

// Env for compile and runtime
type Env struct {
	parent *Env
	ctx    map[string]*val.Val
}

func NewEnv() *Env {
	return &Env{nil, map[string]*val.Val{}}
}

func (e *Env) Inherit(parent *Env) *Env {
	util.Assert(e.parent == nil, "env.parent != nil")
	e.parent = parent
	return e
}

func (e *Env) Derive() *Env {
	return &Env{e, map[string]*val.Val{}}
}

func (e *Env) Get(name string) (*val.Val, bool) {
	v, ok := e.ctx[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return v, ok
}

func (e *Env) Put(name string, val *val.Val) {
	e.ctx[name] = val
}

func (e *Env) RegisterFun(v *val.FunVal) {
	e.Put(v.Kind.Fun().OverloadName(), v.Vl())
}
