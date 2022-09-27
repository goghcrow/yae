package env0

import (
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"unsafe"
)

// Env for typeChecker
type Env struct {
	parent *Env
	ctx    map[string]*types.Kind
}

func NewEnv() *Env {
	return &Env{nil, map[string]*types.Kind{}}
}

func (e *Env) Inherit(parent *Env) *Env {
	util.Assert(e.parent == nil, "env.parent != nil")
	e.parent = parent
	return e
}

func (e *Env) Derive() *Env {
	return &Env{e, map[string]*types.Kind{}}
}

func (e *Env) Get(name string) (*types.Kind, bool) {
	v, ok := e.ctx[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return v, ok
}

func (e *Env) Put(name string, val *types.Kind) {
	e.ctx[name] = val
}

func (e *Env) ForEach(f func(string, *types.Kind)) {
	for k, v := range e.ctx {
		f(k, v)
	}
}

func (e *Env) RegisterFun(f *types.FunKind) {
	lookup, mono := f.Lookup()
	if mono {
		e.Put(lookup, f.Kd())
	} else {
		fs := &[]*types.FunKind{f}
		arr, ok := e.Get(lookup)
		if ok {
			fs = (*[]*types.FunKind)(unsafe.Pointer(arr))
			*fs = append(*fs, f)
		}
		e.Put(lookup, (*types.Kind)(unsafe.Pointer(fs)))
	}
}
