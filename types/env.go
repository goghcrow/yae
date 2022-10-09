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
	util.Assert(slotFree(val), "expect unslot")
	// 注意这里只修改当前环境, 不修改继承
	// 如果是 scope 语义, 需要先 env:=findDefEnv(name) 然后 env.ctx[name]=val
	e.ctx[name] = val
}

func (e *Env) ForEach(f func(string, *Kind)) {
	for k, v := range e.ctx {
		f(k, v)
	}
}

func (e *Env) RegisterFun(f *Kind) {
	util.Assert(f.Type == TFun, "expect FunKind actual %s", f)
	lookup, mono := f.Fun().Lookup()
	if mono {
		e.Put(lookup, f)
	} else {
		hackPtr, ok := e.Get(lookup)
		if !ok {
			hackPtr = (*Kind)(unsafe.Pointer(&[]*FunKind{}))
			// e.Put(lookup, hackPtr)
			e.ctx[lookup] = hackPtr
		}

		fnTblPtr := (*[]*FunKind)(unsafe.Pointer(hackPtr))
		*fnTblPtr = append(*fnTblPtr, f.Fun())
	}
}
