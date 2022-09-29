package types

import (
	"github.com/goghcrow/yae/util"
	"unsafe"
)

// Env for typeChecker
type Env struct {
	parent *Env
	ctx    map[string]*Kind
}

func NewEnv() *Env {
	return &Env{nil, map[string]*Kind{}}
}

func (e *Env) Inherit(parent *Env) *Env {
	util.Assert(e.parent == nil, "env.parent != nil")
	e.parent = parent
	return e
}

func (e *Env) Derive() *Env {
	return &Env{e, map[string]*Kind{}}
}

func (e *Env) Get(name string) (*Kind, bool) {
	v, ok := e.ctx[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return v, ok
}

func (e *Env) Put(name string, val *Kind) {
	e.ctx[name] = val
}

func (e *Env) ForEach(f func(string, *Kind)) {
	for k, v := range e.ctx {
		f(k, v)
	}
}

func (e *Env) RegisterFun(f *Kind) {
	util.Assert(f.Type == TFun, "need to register FunKind, get %s", f)
	lookup, mono := f.Fun().Lookup()
	if mono {
		e.Put(lookup, f)
	} else {
		// 多态函数一组签名 hack, []*types.FunKind
		arr, ok := e.Get(lookup)
		if !ok {
			arr = (*Kind)(unsafe.Pointer(&[]*FunKind{}))
			e.Put(lookup, arr)
		}

		fs := (*[]*FunKind)(unsafe.Pointer(arr))
		// 注意这里要更新, *fs =
		*fs = append(*fs, f.Fun())
	}
}
