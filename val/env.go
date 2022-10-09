package val

import (
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"unsafe"
)

// Env for compile and runtime
type Env struct {
	parent *Env
	ctx    map[string]*Val
}

func NewEnv() *Env {
	return &Env{nil, map[string]*Val{}}
}

func (e *Env) Inherit(parent *Env) *Env {
	util.Assert(e.parent == nil, "env.parent != nil")
	e.parent = parent
	return e
}

func (e *Env) Derive() *Env {
	return &Env{e, map[string]*Val{}}
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

func (e *Env) RegisterFun(f *Val) {
	util.Assert(f.Kind.Type == types.TFun, "expect FunVal actual %s", f)
	lookup, mono := f.Kind.Fun().Lookup()
	if mono {
		e.Put(lookup, f)
	} else {
		hackPtr, ok := e.Get(lookup)
		if !ok {
			hackPtr = (*Val)(unsafe.Pointer(&[]*FunVal{}))
			e.Put(lookup, hackPtr)
		}

		fnTblPtr := (*[]*FunVal)(unsafe.Pointer(hackPtr))
		*fnTblPtr = append(*fnTblPtr, f.Fun())
	}
}
