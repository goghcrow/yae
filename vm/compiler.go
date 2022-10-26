package vm

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"math"
	"time"
)

// 字节码的本质是 ast 后序遍历的产生的线性序列

func Compile(expr ast.Expr, env1 *val.Env) compiler.Closure {
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

func (c *Compiler) Compile(expr ast.Expr, env *val.Env) *bytecode {
	b := &bytecode{cp: c.cp}
	b.compile(c, expr, env)
	b.end()
	return b
}

type bytecode struct {
	code []byte
	*cp
}

func (b *bytecode) compile(c *Compiler, expr ast.Expr, env *val.Env) {
	switch e := expr.(type) {
	case *ast.StrExpr:
		b.emitOP(OP_CONST)
		b.emitConst(val.Str(e.Val))

	case *ast.NumExpr:
		b.emitOP(OP_CONST)
		b.emitConst(val.Num(e.Val))

	case *ast.TimeExpr:
		b.emitOP(OP_CONST)
		b.emitConst(val.Time(time.Unix(e.Val, 0)))

	case *ast.BoolExpr:
		b.emitOP(OP_CONST)
		if e.Val {
			b.emitConst(val.True)
		} else {
			b.emitConst(val.False)
		}

	case *ast.ListExpr:
		els := e.Elems
		for _, el := range els {
			b.compile(c, el, env)
		}
		b.emitOP(OP_NEW_LIST)
		b.emitConst(e.Type)
		b.emitMediumInt(len(els))

	case *ast.MapExpr:
		for _, p := range e.Pairs {
			b.compile(c, p.Key, env)
			b.compile(c, p.Val, env)
		}
		b.emitOP(OP_NEW_MAP)
		b.emitConst(e.Type)
		b.emitMediumInt(len(e.Pairs))

	case *ast.ObjExpr:
		for _, f := range e.Fields {
			b.compile(c, f.Val, env)
		}
		b.emitOP(OP_NEW_OBJ)
		b.emitConst(e.Type)

	case *ast.IdentExpr:
		b.emitOP(OP_LOAD)
		b.emitConst(e.Name)

	case *ast.CallExpr:
		if e.Resolved == "" {
			b.compileInvokeDynamic(c, e, env)
		} else {
			b.compileInvokeStatic(c, e, env)
		}

	case *ast.SubscriptExpr:
		b.compile(c, e.Var, env)
		b.compile(c, e.Idx, env)
		ty := e.VarType.(*types.Type)
		switch ty.Kind {
		case types.KList:
			b.emitOP(OP_LIST_LOAD)
		case types.KMap:
			b.emitOP(OP_MAP_LOAD)
		default:
			util.Unreachable()
		}

	case *ast.MemberExpr:
		b.compile(c, e.Obj, env)
		b.emitOP(OP_OBJ_LOAD)
		b.emitMediumInt(e.Index)

	default:
		util.Unreachable()
	}
}

func (b *bytecode) compileInvokeStatic(c *Compiler, call *ast.CallExpr, env *val.Env) {
	var fun *val.FunVal
	if call.Index < 0 {
		fun = env.MustGetMonoFun(call.Resolved)
	} else {
		fun = env.MustGetPolyFuns(call.Resolved)[call.Index]
	}
	f := fun.Vl()

	lazy := fun.Lazy
	params := fun.Type.Fun().Param

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
	opCBV, ok := intrinsicsCallByValue[f]
	if ok {
		b.emitOP(opCBV)
		return
	}

	if lazy {
		b.emitOP(OP_CALL_BY_NEED)
	} else {
		b.emitOP(OP_CALL_BY_VALUE)
	}

	b.emitConst(f)
	b.emitUint8(len(call.Args)) // 参数最多 256
}

type thunkVal struct {
	val.Val
	*bytecode
	// 不是闭包, 不需要引用 env
}

func newThunk(b *bytecode, retK *types.Type) *val.Val {
	fk := types.Fun("thunk", []*types.Type{}, retK)
	v := thunkVal{val.Val{Type: fk}, b}
	return &v.Val
}

func (b *bytecode) compileInvokeDynamic(c *Compiler, call *ast.CallExpr, env *val.Env) {
	b.compile(c, call.Callee, env)
	// dynamic 不支持 lazy
	for _, arg := range call.Args {
		b.compile(c, arg, env)
	}
	b.emitOP(OP_DYNAMIC_CALL)
	b.emitUint8(len(call.Args)) // 参数最多 256
}

func (b *bytecode) addConst(v interface{}) int {
	b.data = append(b.data, v)
	return len(b.data) - 1
}

func (b *bytecode) emit(ops ...byte) { b.code = append(b.code, ops...) }
func (b *bytecode) emitOP(op opcode) { b.emit(byte(op)) }

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
