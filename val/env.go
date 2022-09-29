package val

import (
	types "github.com/goghcrow/yae/type"
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
	e.ctx[name] = val
}

func (e *Env) RegisterFun(f *Val) {
	util.Assert(f.Kind.Type == types.TFun, "need to register FunVal, get %s", f)
	lookup, mono := f.Kind.Fun().Lookup()
	if mono {
		e.Put(lookup, f)
	} else {
		fs := &[]*FunVal{f.Fun()}
		arr, ok := e.Get(lookup)
		if ok {
			fs = (*[]*FunVal)(unsafe.Pointer(arr))
			*fs = append(*fs, f.Fun())
		}
		e.Put(lookup, (*Val)(unsafe.Pointer(fs)))
	}
}
