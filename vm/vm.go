package vm

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"math"
	"unsafe"
)

type VM struct {
	*stack
}

func NewVM() *VM {
	return &VM{
		stack: newStack(),
	}
}

func (v *VM) Run() {}

//https://www.zhihu.com/question/57754882/answer/154549716

// Interp 字节码解释器
func (v *VM) Interp(b *bytecode, env *val.Env) *val.Val {
	v.stack = newStack()
	return v.interp(b, env)
}

func (v *VM) interp(b *bytecode, env *val.Env) *val.Val {
	for {
		opcode := op(b.code[v.pc])
		v.pc += 1

		switch opcode {

		case OP_CONST:
			c, w := b.readConst(v.pc)
			v.pc += w
			v.Push(c.(*val.Val))

		case OP_LOAD:
			c, w := b.readConst(v.pc)
			v.pc += w
			vl, _ := env.Get(c.(string))
			v.Push(vl)

		case OP_ADD_NUM:
			// nothing to do

		case OP_ADD_NUM_NUM:
			rhs := v.Pop().Num().V
			lhs := v.Pop().Num().V
			v.Push(val.Num(lhs + rhs))

		case OP_SUB_NUM:
			n := v.Pop().Num().V
			v.Push(val.Num(-n))

		case OP_SUB_NUM_NUM:
			rhs := v.Pop().Num().V
			lhs := v.Pop().Num().V
			v.Push(val.Num(lhs - rhs))

		case OP_MUL_NUM_NUM:
			rhs := v.Pop().Num().V
			lhs := v.Pop().Num().V
			v.Push(val.Num(lhs * rhs))

		case OP_DIV_NUM_NUM:
			rhs := v.Pop().Num().V
			lhs := v.Pop().Num().V
			v.Push(val.Num(lhs / rhs))

		case OP_MOD_NUM_NUM:
			rhs := v.Pop().Num().V
			lhs := v.Pop().Num().V
			v.Push(val.Num(float64(int64(lhs) % int64(rhs))))

		case OP_EXP_NUM_NUM:
			rhs := v.Pop().Num().V
			lhs := v.Pop().Num().V
			v.Push(val.Num(math.Pow(lhs, rhs)))

		case OP_EQ_NUM_NUM:
			rhs := v.Pop().Num()
			lhs := v.Pop().Num()
			v.Push(val.Bool(val.NumEQ(lhs, rhs)))

		case OP_EQ_BOOL_BOOL:
			rhs := v.Pop().Str().V
			lhs := v.Pop().Str().V
			v.Push(val.Bool(lhs == rhs))

		case OP_EQ_STR_STR:
			rhs := v.Pop().Bool().V
			lhs := v.Pop().Bool().V
			v.Push(val.Bool(lhs == rhs))

		case OP_NE_NUM_NUM:
			rhs := v.Pop().Num()
			lhs := v.Pop().Num()
			v.Push(val.Bool(val.NumNE(lhs, rhs)))

		case OP_NE_BOOL_BOOL:
			rhs := v.Pop().Str().V
			lhs := v.Pop().Str().V
			v.Push(val.Bool(lhs != rhs))

		case OP_NE_STR_STR:
			rhs := v.Pop().Bool().V
			lhs := v.Pop().Bool().V
			v.Push(val.Bool(lhs != rhs))

		case OP_LT_NUM_NUM:
			rhs := v.Pop().Num()
			lhs := v.Pop().Num()
			v.Push(val.Bool(val.NumLT(lhs, rhs)))

		case OP_LT_TIME_TIME:
			rhs := v.Pop().Time().V
			lhs := v.Pop().Time().V
			v.Push(val.Bool(lhs.Before(rhs)))

		case OP_LE_NUM_NUM:
			rhs := v.Pop().Num()
			lhs := v.Pop().Num()
			v.Push(val.Bool(val.NumLE(lhs, rhs)))

		case OP_LE_TIME_TIME:
			rhs := v.Pop().Time().V
			lhs := v.Pop().Time().V
			v.Push(val.Bool(lhs.Before(rhs) || lhs.Equal(rhs)))

		case OP_GT_NUM_NUM:
			rhs := v.Pop().Num()
			lhs := v.Pop().Num()
			v.Push(val.Bool(val.NumGT(lhs, rhs)))

		case OP_GT_TIME_TIME:
			rhs := v.Pop().Time().V
			lhs := v.Pop().Time().V
			v.Push(val.Bool(lhs.After(rhs)))

		case OP_GE_NUM_NUM:
			rhs := v.Pop().Num()
			lhs := v.Pop().Num()
			v.Push(val.Bool(val.NumGE(lhs, rhs)))

		case OP_GE_TIME_TIME:
			rhs := v.Pop().Time().V
			lhs := v.Pop().Time().V
			v.Push(val.Bool(lhs.After(rhs) || lhs.Equal(rhs)))

		case OP_JUMP:
			off, _ := b.readMediumInt(v.pc)
			v.pc = off

		case OP_IF_TRUE:
			fpc, w := b.readMediumInt(v.pc)
			v.pc += w
			if !v.Pop().Bool().V {
				v.pc = fpc
			}

		case OP_LOGICAL_NOT:
			v.Push(val.Bool(!v.Pop().Bool().V))

		case OP_NEW_LIST:
			kd, w := b.readConst(v.pc)
			v.pc += w
			sz, w := b.readMediumInt(v.pc)
			v.pc += w
			lst := val.List(kd.(*types.Kind).List(), sz).List()
			for i := 0; i < sz; i++ {
				lst.V[sz-i-1] = v.Pop()
			}
			v.Push(lst.Vl())

		case OP_NEW_MAP:
			kd, w := b.readConst(v.pc)
			v.pc += w
			sz, w := b.readMediumInt(v.pc)
			v.pc += w

			m := val.Map(kd.(*types.Kind).Map()).Map()
			for i := 0; i < sz; i++ {
				vl := v.Pop()
				key := v.Pop()
				m.V[key.Key()] = vl
			}
			v.Push(m.Vl())

		case OP_NEW_OBJ:
			kd, w := b.readConst(v.pc)
			v.pc += w

			objK := kd.(*types.Kind).Obj()
			sz := len(objK.Fields)
			o := val.Obj(objK).Obj()
			for i := 0; i < sz; i++ {
				o.V[sz-1-i] = v.Pop()
			}
			v.Push(o.Vl())

		case OP_LIST_LOAD:
			idx := int(v.Pop().Num().V)
			lst := v.Pop().List().V
			util.Assert(idx < len(lst), "out of range %d of %s", idx, lst)
			v.Push(lst[idx])

		case OP_MAP_LOAD:
			key := v.Pop()
			m := v.Pop().Map()
			vl, ok := m.Get(key)
			util.Assert(ok, "undefined key %s of %s", key, m)
			v.Push(vl)

		case OP_OBJ_LOAD:
			idx, w := b.readMediumInt(v.pc)
			v.pc += w
			o := v.Pop().Obj()
			v.Push(o.V[idx])

		case OP_INVOKE_STATIC:
			fv, w := b.readConst(v.pc)
			v.pc += w
			argc := b.readUint8(v.pc)
			v.pc += 1

			f := fv.(*val.Val).Fun()
			args := make([]*val.Val, argc)
			for i := 0; i < argc; i++ {
				args[argc-1-i] = v.Pop()
			}
			v.Push(f.Call(args...))

		case OP_INVOKE_STATIC_LAZY:
			fv, w := b.readConst(v.pc)
			v.pc += w
			argc := b.readUint8(v.pc)
			v.pc += 1

			f := fv.(*val.Val).Fun()
			params := f.Kind.Fun().Param
			args := make([]*val.Val, argc)

			// 对于表达式而言, 没有局部作用域, 这里可以完全简化为递归执行字节码, 不需要常规的 call 调用
			// 进而 vm 也只需要一个栈, 不需要 vm 关联 *stack 或者 frame 的传统的做法
			for i := 0; i < argc; i++ {
				thunkPtr := v.Pop()
				thunk := (*thunkVal)(unsafe.Pointer(thunkPtr))
				body := thunk.bytecode
				thunkK := types.Fun("thunk", []*types.Kind{}, params[argc-1-i])
				args[argc-1-i] = val.Fun(thunkK, func(...*val.Val) *val.Val {
					return v.call0(body, env)
				})
			}
			v.Push(f.Call(args...))

		case OP_INVOKE_DYNAMIC:
			argc := b.readUint8(v.pc)
			v.pc += 1
			args := make([]*val.Val, argc)
			for i := 0; i < argc; i++ {
				args[argc-1-i] = v.Pop()
			}
			f := v.Pop().Fun()
			v.Push(f.Call(args...))

		case OP_RETURN:
			return v.Pop()
		case OP_NOP:

		default:
			panic(fmt.Errorf("unsupported op %d", opcode))
		}

	}
}

// for thunk args
func (v *VM) call0(body *bytecode, env *val.Env) *val.Val {
	// 只有 thunk 会产生新的函数, 这里简化处理就不做 frame 了
	callerStack := v.stack
	v.stack = newStack()
	// 表达式没有局部作用域, 这里不需要 env.Derive() 子环境
	ret := v.interp(body, env)
	v.stack = callerStack
	// v.Push(ret)
	return ret
}
