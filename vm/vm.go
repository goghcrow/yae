package vm

import (
	"github.com/goghcrow/yae/val"
)

// https://www.zhihu.com/question/57754882/answer/154549716
// https://en.wikipedia.org/wiki/Threaded_code
type threading func(*VM) *val.Val

type VM struct {
	pc int
	*stack
	*bytecode
	env    *val.Env
	interp threading
}

func NewVM() *VM {
	return &VM{
		stack:  newStack(),
		interp: switchThreading,
		//interp: callThreading,
	}
}

// Interp 字节码解释器 switch threading
func (v *VM) Interp(b *bytecode, env *val.Env) *val.Val {
	v.stack = newStack()
	v.bytecode = b
	v.env = env
	return v.interp(v)
}

func (v *VM) save() *VM {
	return &VM{
		pc:       v.pc,
		stack:    v.stack,
		bytecode: v.bytecode,
	}
}

func (v *VM) reset(save *VM) {
	v.pc = save.pc
	v.stack = save.stack
	v.bytecode = save.bytecode
}

func (v *VM) doCall0(body *bytecode, env *val.Env) *val.Val {
	v.pc = 0
	v.stack = newStack()
	v.bytecode = body
	// 表达式没有局部作用域, 这里不需要 env.Derive() 子环境
	return v.interp(v)
}

// for thunk args, 无参, 也没 call convention, 直接返回
func (v *VM) call0(body *bytecode, env *val.Env) *val.Val {
	// 只有 thunk 会产生新的函数, 这里简化处理就不做 frame 了
	save := v.save()
	ret := v.doCall0(body, env)
	v.reset(save)
	// v.Push(ret)
	return ret
}
