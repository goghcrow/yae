// Generated Code; DO NOT EDIT.

package vm

import (
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"math"
	"time"
	"unicode/utf8"
	"unsafe"
)

type OpcodeHandler func(vm *VM)

var instructions [_END_]OpcodeHandler

const limit = 1024

func callThreading(v *VM) *val.Val {
	for i := 0; i < limit; i++ {
		op := opcode(v.code[v.pc])
		v.pc += 1
		if op == OP_RETURN {
			return v.Pop()
		}
		instructions[op](v)
	}
	panic("over exec limit")
}

func init() {
	instructions[OP_CONST] = OP_CONST_Handler
	instructions[OP_LOAD] = OP_LOAD_Handler
	instructions[OP_ADD_NUM] = OP_ADD_NUM_Handler
	instructions[OP_ADD_NUM_NUM] = OP_ADD_NUM_NUM_Handler
	instructions[OP_ADD_STR_STR] = OP_ADD_STR_STR_Handler
	instructions[OP_SUB_NUM] = OP_SUB_NUM_Handler
	instructions[OP_SUB_NUM_NUM] = OP_SUB_NUM_NUM_Handler
	instructions[OP_SUB_TIME_TIME] = OP_SUB_TIME_TIME_Handler
	instructions[OP_MUL_NUM_NUM] = OP_MUL_NUM_NUM_Handler
	instructions[OP_DIV_NUM_NUM] = OP_DIV_NUM_NUM_Handler
	instructions[OP_MOD_NUM_NUM] = OP_MOD_NUM_NUM_Handler
	instructions[OP_EXP_NUM_NUM] = OP_EXP_NUM_NUM_Handler
	instructions[OP_MIN_NUM_NUM] = OP_MIN_NUM_NUM_Handler
	instructions[OP_MAX_NUM_NUM] = OP_MAX_NUM_NUM_Handler
	instructions[OP_ABS_NUM] = OP_ABS_NUM_Handler
	instructions[OP_CEIL_NUM] = OP_CEIL_NUM_Handler
	instructions[OP_FLOOR_NUM] = OP_FLOOR_NUM_Handler
	instructions[OP_ROUND_NUM] = OP_ROUND_NUM_Handler
	instructions[OP_EQ_NUM_NUM] = OP_EQ_NUM_NUM_Handler
	instructions[OP_EQ_BOOL_BOOL] = OP_EQ_BOOL_BOOL_Handler
	instructions[OP_EQ_STR_STR] = OP_EQ_STR_STR_Handler
	instructions[OP_EQ_TIME_TIME] = OP_EQ_TIME_TIME_Handler
	instructions[OP_EQ_LIST_LIST] = OP_EQ_LIST_LIST_Handler
	instructions[OP_EQ_MAP_MAP] = OP_EQ_MAP_MAP_Handler
	instructions[OP_NE_NUM_NUM] = OP_NE_NUM_NUM_Handler
	instructions[OP_NE_BOOL_BOOL] = OP_NE_BOOL_BOOL_Handler
	instructions[OP_NE_STR_STR] = OP_NE_STR_STR_Handler
	instructions[OP_NE_TIME_TIME] = OP_NE_TIME_TIME_Handler
	instructions[OP_NE_LIST_LIST] = OP_NE_LIST_LIST_Handler
	instructions[OP_NE_MAP_MAP] = OP_NE_MAP_MAP_Handler
	instructions[OP_LT_NUM_NUM] = OP_LT_NUM_NUM_Handler
	instructions[OP_LT_TIME_TIME] = OP_LT_TIME_TIME_Handler
	instructions[OP_LE_NUM_NUM] = OP_LE_NUM_NUM_Handler
	instructions[OP_LE_TIME_TIME] = OP_LE_TIME_TIME_Handler
	instructions[OP_GT_NUM_NUM] = OP_GT_NUM_NUM_Handler
	instructions[OP_GT_TIME_TIME] = OP_GT_TIME_TIME_Handler
	instructions[OP_GE_NUM_NUM] = OP_GE_NUM_NUM_Handler
	instructions[OP_GE_TIME_TIME] = OP_GE_TIME_TIME_Handler
	instructions[OP_JUMP] = OP_JUMP_Handler
	instructions[OP_IF_TRUE] = OP_IF_TRUE_Handler
	instructions[OP_LOGICAL_NOT] = OP_LOGICAL_NOT_Handler
	instructions[OP_NEW_LIST] = OP_NEW_LIST_Handler
	instructions[OP_NEW_MAP] = OP_NEW_MAP_Handler
	instructions[OP_NEW_OBJ] = OP_NEW_OBJ_Handler
	instructions[OP_LIST_LOAD] = OP_LIST_LOAD_Handler
	instructions[OP_MAP_LOAD] = OP_MAP_LOAD_Handler
	instructions[OP_OBJ_LOAD] = OP_OBJ_LOAD_Handler
	instructions[OP_LEN_STR] = OP_LEN_STR_Handler
	instructions[OP_LEN_LIST] = OP_LEN_LIST_Handler
	instructions[OP_LEN_MAP] = OP_LEN_MAP_Handler
	instructions[OP_STRTOTIME_STR] = OP_STRTOTIME_STR_Handler
	instructions[OP_INVOKE_STATIC] = OP_INVOKE_STATIC_Handler
	instructions[OP_INVOKE_STATIC_LAZY] = OP_INVOKE_STATIC_LAZY_Handler
	instructions[OP_INVOKE_DYNAMIC] = OP_INVOKE_DYNAMIC_Handler

	instructions[OP_NOP] = OP_NOP_Handler
}

//goland:noinspection GoSnakeCaseUsage
func OP_CONST_Handler(v *VM) {
	c, w := v.readConst(v.pc)
	v.pc += w
	v.Push(c.(*val.Val))
}

//goland:noinspection GoSnakeCaseUsage
func OP_LOAD_Handler(v *VM) {
	c, w := v.readConst(v.pc)
	v.pc += w
	vl, _ := v.env.Get(c.(string))
	v.Push(vl)
}

//goland:noinspection GoSnakeCaseUsage
func OP_ADD_NUM_Handler(v *VM) {
	// nothing to do
}

//goland:noinspection GoSnakeCaseUsage
func OP_ADD_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(lhs + rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_ADD_STR_STR_Handler(v *VM) {
	rhs := v.Pop().Str().V
	lhs := v.Pop().Str().V
	v.Push(val.Str(lhs + rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_SUB_NUM_Handler(v *VM) {
	v.Push(val.Num(-v.Pop().Num().V))
}

//goland:noinspection GoSnakeCaseUsage
func OP_SUB_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(lhs - rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_SUB_TIME_TIME_Handler(v *VM) {
	rhs := v.Pop().Time().V
	lhs := v.Pop().Time().V
	v.Push(val.Num(lhs.Sub(rhs).Seconds()))
}

//goland:noinspection GoSnakeCaseUsage
func OP_MUL_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(lhs * rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_DIV_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(lhs / rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_MOD_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(float64(int64(lhs) % int64(rhs))))
}

//goland:noinspection GoSnakeCaseUsage
func OP_EXP_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(math.Pow(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_MIN_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(math.Min(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_MAX_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num().V
	lhs := v.Pop().Num().V
	v.Push(val.Num(math.Max(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_ABS_NUM_Handler(v *VM) {
	v.Push(val.Num(math.Abs(v.Pop().Num().V)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_CEIL_NUM_Handler(v *VM) {
	v.Push(val.Num(math.Ceil(v.Pop().Num().V)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_FLOOR_NUM_Handler(v *VM) {
	v.Push(val.Num(math.Floor(v.Pop().Num().V)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_ROUND_NUM_Handler(v *VM) {
	v.Push(val.Num(math.Round(v.Pop().Num().V)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_EQ_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num()
	lhs := v.Pop().Num()
	v.Push(val.Bool(val.NumEQ(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_EQ_BOOL_BOOL_Handler(v *VM) {
	rhs := v.Pop().Bool().V
	lhs := v.Pop().Bool().V
	v.Push(val.Bool(lhs == rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_EQ_STR_STR_Handler(v *VM) {
	rhs := v.Pop().Str().V
	lhs := v.Pop().Str().V
	v.Push(val.Bool(lhs == rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_EQ_TIME_TIME_Handler(v *VM) {
	rhs := v.Pop().Time().V
	lhs := v.Pop().Time().V
	v.Push(val.Bool(lhs.Equal(rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_EQ_LIST_LIST_Handler(v *VM) {
	rhs := v.Pop()
	lhs := v.Pop()
	v.Push(val.Bool(val.Equals(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_EQ_MAP_MAP_Handler(v *VM) {
	rhs := v.Pop()
	lhs := v.Pop()
	v.Push(val.Bool(val.Equals(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NE_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num()
	lhs := v.Pop().Num()
	v.Push(val.Bool(val.NumNE(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NE_BOOL_BOOL_Handler(v *VM) {
	rhs := v.Pop().Bool().V
	lhs := v.Pop().Bool().V
	v.Push(val.Bool(lhs != rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NE_STR_STR_Handler(v *VM) {
	rhs := v.Pop().Str().V
	lhs := v.Pop().Str().V
	v.Push(val.Bool(lhs != rhs))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NE_TIME_TIME_Handler(v *VM) {
	rhs := v.Pop().Time().V
	lhs := v.Pop().Time().V
	v.Push(val.Bool(!lhs.Equal(rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NE_LIST_LIST_Handler(v *VM) {
	rhs := v.Pop()
	lhs := v.Pop()
	v.Push(val.Bool(!val.Equals(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NE_MAP_MAP_Handler(v *VM) {
	rhs := v.Pop()
	lhs := v.Pop()
	v.Push(val.Bool(!val.Equals(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_LT_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num()
	lhs := v.Pop().Num()
	v.Push(val.Bool(val.NumLT(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_LT_TIME_TIME_Handler(v *VM) {
	rhs := v.Pop().Time().V
	lhs := v.Pop().Time().V
	v.Push(val.Bool(lhs.Before(rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_LE_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num()
	lhs := v.Pop().Num()
	v.Push(val.Bool(val.NumLE(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_LE_TIME_TIME_Handler(v *VM) {
	rhs := v.Pop().Time().V
	lhs := v.Pop().Time().V
	v.Push(val.Bool(lhs.Before(rhs) || lhs.Equal(rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_GT_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num()
	lhs := v.Pop().Num()
	v.Push(val.Bool(val.NumGT(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_GT_TIME_TIME_Handler(v *VM) {
	rhs := v.Pop().Time().V
	lhs := v.Pop().Time().V
	v.Push(val.Bool(lhs.After(rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_GE_NUM_NUM_Handler(v *VM) {
	rhs := v.Pop().Num()
	lhs := v.Pop().Num()
	v.Push(val.Bool(val.NumGE(lhs, rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_GE_TIME_TIME_Handler(v *VM) {
	rhs := v.Pop().Time().V
	lhs := v.Pop().Time().V
	v.Push(val.Bool(lhs.After(rhs) || lhs.Equal(rhs)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_JUMP_Handler(v *VM) {
	off, _ := v.readMediumInt(v.pc)
	v.pc = off
}

//goland:noinspection GoSnakeCaseUsage
func OP_IF_TRUE_Handler(v *VM) {
	fpc, w := v.readMediumInt(v.pc)
	v.pc += w
	if !v.Pop().Bool().V {
		v.pc = fpc
	}
}

//goland:noinspection GoSnakeCaseUsage
func OP_LOGICAL_NOT_Handler(v *VM) {
	v.Push(val.Bool(!v.Pop().Bool().V))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NEW_LIST_Handler(v *VM) {
	kd, w := v.readConst(v.pc)
	v.pc += w
	sz, w := v.readMediumInt(v.pc)
	v.pc += w
	lst := val.List(kd.(*types.Kind).List(), sz).List()
	for i := 0; i < sz; i++ {
		lst.V[sz-i-1] = v.Pop()
	}
	v.Push(lst.Vl())
}

//goland:noinspection GoSnakeCaseUsage
func OP_NEW_MAP_Handler(v *VM) {
	kd, w := v.readConst(v.pc)
	v.pc += w
	sz, w := v.readMediumInt(v.pc)
	v.pc += w

	m := val.Map(kd.(*types.Kind).Map()).Map()
	for i := 0; i < sz; i++ {
		vl := v.Pop()
		key := v.Pop()
		m.V[key.Key()] = vl
	}
	v.Push(m.Vl())
}

//goland:noinspection GoSnakeCaseUsage
func OP_NEW_OBJ_Handler(v *VM) {
	kd, w := v.readConst(v.pc)
	v.pc += w

	objK := kd.(*types.Kind).Obj()
	sz := len(objK.Fields)
	o := val.Obj(objK).Obj()
	for i := 0; i < sz; i++ {
		o.V[sz-1-i] = v.Pop()
	}
	v.Push(o.Vl())
}

//goland:noinspection GoSnakeCaseUsage
func OP_LIST_LOAD_Handler(v *VM) {
	idx := int(v.Pop().Num().V)
	lst := v.Pop().List().V
	util.Assert(idx < len(lst), "out of range %d of %s", idx, lst)
	v.Push(lst[idx])
}

//goland:noinspection GoSnakeCaseUsage
func OP_MAP_LOAD_Handler(v *VM) {
	key := v.Pop()
	m := v.Pop().Map()
	vl, ok := m.Get(key)
	util.Assert(ok, "undefined key %s of %s", key, m)
	v.Push(vl)
}

//goland:noinspection GoSnakeCaseUsage
func OP_OBJ_LOAD_Handler(v *VM) {
	idx, w := v.readMediumInt(v.pc)
	v.pc += w
	o := v.Pop().Obj()
	v.Push(o.V[idx])
}

//goland:noinspection GoSnakeCaseUsage
func OP_LEN_STR_Handler(v *VM) {
	v.Push(val.Num(float64(utf8.RuneCountInString(v.Pop().Str().V))))
}

//goland:noinspection GoSnakeCaseUsage
func OP_LEN_LIST_Handler(v *VM) {
	v.Push(val.Num(float64(len(v.Pop().List().V))))
}

//goland:noinspection GoSnakeCaseUsage
func OP_LEN_MAP_Handler(v *VM) {
	v.Push(val.Num(float64(len(v.Pop().Map().V))))
}

//goland:noinspection GoSnakeCaseUsage
func OP_STRTOTIME_STR_Handler(v *VM) {
	ts := timelib.Strtotime(v.Pop().Str().V)
	v.Push(val.Time(time.Unix(ts, 0)))
}

//goland:noinspection GoSnakeCaseUsage
func OP_INVOKE_STATIC_Handler(v *VM) {
	fv, w := v.readConst(v.pc)
	v.pc += w
	argc := v.readUint8(v.pc)
	v.pc += 1

	f := fv.(*val.Val).Fun()
	args := make([]*val.Val, argc)
	for i := 0; i < argc; i++ {
		args[argc-1-i] = v.Pop()
	}
	v.Push(f.Call(args...))
}

//goland:noinspection GoSnakeCaseUsage
func OP_INVOKE_STATIC_LAZY_Handler(v *VM) {
	fv, w := v.readConst(v.pc)
	v.pc += w
	argc := v.readUint8(v.pc)
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
			return v.call0(body, v.env)
		})
	}
	v.Push(f.Call(args...))
}

//goland:noinspection GoSnakeCaseUsage
func OP_INVOKE_DYNAMIC_Handler(v *VM) {
	argc := v.readUint8(v.pc)
	v.pc += 1
	args := make([]*val.Val, argc)
	for i := 0; i < argc; i++ {
		args[argc-1-i] = v.Pop()
	}
	f := v.Pop().Fun()
	v.Push(f.Call(args...))
}

//goland:noinspection GoSnakeCaseUsage
func OP_NOP_Handler(v *VM) {}
