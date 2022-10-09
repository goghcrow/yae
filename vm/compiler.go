package vm

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"math"
	"strconv"
	"time"
	"unsafe"
)

// 字节码的本质是 ast 后序遍历的产生的线性序列

func Compile(expr *ast.Expr, env1 *val.Env) compiler.Closure {
	bytecode := NewCompile().Compile(expr, env1)
	return func(env *val.Env) *val.Val {
		return NewVM().Interp(bytecode, env)
	}
}

func NewCompile() *Compiler {
	return &Compiler{cp: &cp{}}
}

// 编译的结果不需要序列化, 保存在内存即可,
// 所以常量池可以奔放一点, 直接处理成 []interface{}
type cp struct {
	data []interface{}
}

type Compiler struct {
	*cp
}

func (c *Compiler) Compile(expr *ast.Expr, env *val.Env) *bytecode {
	b := &bytecode{cp: c.cp}
	b.compile(c, expr, env)
	b.end()
	return b
}

type bytecode struct {
	code []byte
	*cp
}

func (b *bytecode) compile(c *Compiler, expr *ast.Expr, env *val.Env) {
	switch expr.Type {

	case ast.LITERAL:
		lit := expr.Literal()
		switch lit.LitType {
		case ast.LIT_STR:
			s, err := strconv.Unquote(lit.Text)
			util.Assert(err == nil, "invalid string literal: %s", lit.Val)
			b.emitOP(OP_CONST)
			b.emitConst(val.Str(s))
		case ast.LIT_TIME:
			ts := util.Strtotime(lit.Text[1 : len(lit.Text)-1])
			b.emitOP(OP_CONST)
			b.emitConst(val.Time(time.Unix(ts, 0)))
		case ast.LIT_NUM:
			n, err := util.ParseNum(lit.Text)
			util.Assert(err == nil, "invalid num literal %s", lit.Val)
			b.emitOP(OP_CONST)
			b.emitConst(val.Num(n))
		case ast.LIT_TRUE:
			b.emitOP(OP_CONST)
			b.emitConst(val.True)
		case ast.LIT_FALSE:
			b.emitOP(OP_CONST)
			b.emitConst(val.False)
		default:
			util.Unreachable()
		}

	case ast.IDENT:
		b.emitOP(OP_LOAD)
		b.emitConst(expr.Ident().Name)

	case ast.LIST:
		lst := expr.List()
		els := lst.Elems
		for _, el := range els {
			b.compile(c, el, env)
		}
		b.emitOP(OP_NEW_LIST)
		b.emitConst(lst.Kind)
		b.emitMediumInt(len(els))

	case ast.MAP:
		m := expr.Map()
		for _, p := range m.Pairs {
			b.compile(c, p.Key, env)
			b.compile(c, p.Val, env)
		}
		b.emitOP(OP_NEW_MAP)
		b.emitConst(m.Kind)
		b.emitMediumInt(len(m.Pairs))

	case ast.OBJ:
		o := expr.Obj()
		for _, f := range o.Fields {
			b.compile(c, f.Val, env)
		}
		b.emitOP(OP_NEW_OBJ)
		b.emitConst(o.Kind)

	case ast.SUBSCRIPT:
		sub := expr.Subscript()
		b.compile(c, sub.Var, env)
		b.compile(c, sub.Idx, env)
		kd := sub.VarKind.(*types.Kind)
		switch kd.Type {
		case types.TList:
			b.emitOP(OP_LIST_LOAD)
		case types.TMap:
			b.emitOP(OP_MAP_LOAD)
		default:
			util.Unreachable()
		}

	case ast.MEMBER:
		mem := expr.Member()
		b.compile(c, mem.Obj, env)
		b.emitOP(OP_OBJ_LOAD)
		b.emitMediumInt(mem.Index)

	case ast.CALL:
		call := expr.Call()
		if call.Resolved == "" {
			b.compileInvokeDynamic(c, call, env)
		} else {
			b.compileInvokeStatic(c, call, env)
		}

	default:
		util.Unreachable()
	}
}

func (b *bytecode) compileInvokeStatic(c *Compiler, call *ast.CallExpr, env *val.Env) {
	f, _ := env.Get(call.Resolved)
	// 多态函数, 这里有点 hack, 手动狗头
	if call.Index >= 0 {
		fnTbl := *(*[]*val.FunVal)(unsafe.Pointer(f))
		f = fnTbl[call.Index].Vl()
	}

	fun := f.Fun()
	lazy := fun.Lazy
	params := fun.Kind.Fun().Param

	// 性能考虑, 加入 intrinsic, 生成特化的 opcode
	intrCBN, ok := intrinsicsCallByNeed[f]
	if ok {
		intrCBN(c, b, call.Args, env)
		return
	}

	for i, arg := range call.Args {
		if lazy {
			b.emitOP(OP_CONST)
			b.emitConst(newThunk(c.Compile(arg, env), params[i]))
		} else {
			b.compile(c, arg, env)
		}
	}

	// 性能考虑, 加入 intrinsic, 生成特化的 opcode
	intrCBV, ok := intrinsicsCallByValue[f]
	if ok {
		intrCBV(c, b, env)
		return
	}

	if lazy {
		b.emitOP(OP_INVOKE_STATIC_LAZY)
	} else {
		b.emitOP(OP_INVOKE_STATIC)
	}

	b.emitConst(f)
	b.emitUint8(len(call.Args)) // 参数最多 256
}

type thunkVal struct {
	val.Val
	*bytecode
	// 不是闭包, 不需要引用 env
}

func newThunk(b *bytecode, retK *types.Kind) *val.Val {
	fk := types.Fun("thunk", []*types.Kind{}, retK)
	v := thunkVal{val.Val{Kind: fk}, b}
	return &v.Val
}

func (b *bytecode) compileInvokeDynamic(c *Compiler, call *ast.CallExpr, env *val.Env) {
	b.compile(c, call.Callee, env)
	// dynamic 不支持 lazy
	for _, arg := range call.Args {
		b.compile(c, arg, env)
	}
	b.emitOP(OP_INVOKE_DYNAMIC)
	b.emitUint8(len(call.Args)) // 参数最多 256
}

func (b *bytecode) addConst(v interface{}) int {
	b.data = append(b.data, v)
	return len(b.data) - 1
}

func (b *bytecode) emit(ops ...byte) { b.code = append(b.code, ops...) }
func (b *bytecode) emitOP(op op)     { b.emit(byte(op)) }

func (b *bytecode) emitUint8(ui int) {
	util.Assert(ui <= math.MaxUint8, "overflow")
	b.emit(uint8(ui))
}

func (b *bytecode) readUint8(offset int) int {
	return int(b.code[offset])
}

func (b *bytecode) emitUint16(ui int) {
	util.Assert(ui <= math.MaxUint16, "overflow")
	b.emit(uint16ToByte(uint16(ui))...)
}

func (b *bytecode) readUint16(offset int) int {
	return int(byteToUInt16(b.code[offset : offset+2]))
}

func (b *bytecode) placeholderUint16() func(int) {
	offset := len(b.code)
	b.emitUint16(0)
	return func(i int) {
		util.Assert(i <= math.MaxUint16, "overflow")
		copy(b.code[offset:], uint16ToByte(uint16(i)))
	}
}

func (b *bytecode) emitConst(v interface{}) {
	b.emitUint16(b.addConst(v))
}

func (b *bytecode) readConst(offset int) (c interface{}, sz int) {
	return b.data[b.readUint16(offset)], 2
}

// list size, map pair size, object member
func (b *bytecode) emitMediumInt(i int) {
	b.emitUint16(i)
}

func (b *bytecode) readMediumInt(offset int) (n, i int) {
	return b.readUint16(offset), 2
}

func (b *bytecode) placeholderForMediumInt() func(int) {
	return b.placeholderUint16()
}

func (b *bytecode) end() {
	b.emitOP(OP_RETURN)
}
