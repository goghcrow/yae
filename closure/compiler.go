package closure

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/compiler"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"time"
	"unsafe"
)

// Compile 编译成闭包
// 论文的思路 http://www.iro.umontreal.ca/~feeley/papers/FeeleyLapalmeCL87.pdf
// implement compilers for embedded languages
// 注意区分: env1 是编译期环境, env 是运行时环境
// 能在编译期完成的, 尽可能在编译期计算完
// 调用 Compile 之前必须先对 ast 调用 types.Check:
// 1. 处理字解析
// 2. 繁饰 list/map/obj 类型, 简化 Compile 代码
// 3. call.callee resolve
// 4. 且, Compile 中不检查错误, 假设 types.Check 已经全部检查
func Compile(expr *ast.Expr, env1 *val.Env) compiler.Closure {
	switch expr.Type {
	case ast.LITERAL:
		lit := expr.Literal()
		switch lit.LitType {
		case ast.LIT_STR:
			s := val.Str(lit.Val.(string))
			return func(env *val.Env) *val.Val { return s }
		case ast.LIT_TIME:
			t := val.Time(time.Unix(lit.Val.(int64), 0))
			return func(env *val.Env) *val.Val { return t }
		case ast.LIT_NUM:
			n := val.Num(lit.Val.(float64))
			return func(env *val.Env) *val.Val { return n }
		case ast.LIT_TRUE:
			return func(env *val.Env) *val.Val { return val.True }
		case ast.LIT_FALSE:
			return func(env *val.Env) *val.Val { return val.False }
		default:
			util.Unreachable()
			return nil
		}

	case ast.IDENT:
		// 如果从性能角度考虑, 所有运行时的符号查找其实都可以从 map 换成 array
		// 1. 需要在编译期把符号 resolve 成数组下标
		// 2. 把运行时环境从 map 展开成数组
		id := expr.Ident().Name
		return func(env *val.Env) *val.Val {
			v, _ := env.Get(id)
			return v
		}

	case ast.LIST:
		lst := expr.List()
		els := lst.Elems
		sz := len(els)
		if sz == 0 {
			return func(env *val.Env) *val.Val {
				// 注意空列表类型 list[nothing]
				kind := types.List(types.Bottom).List()
				return val.List(kind, 0)
			}
		}

		kind := lst.Kind.(*types.Kind).List()
		cs := make([]compiler.Closure, sz)
		for i, el := range els {
			cs[i] = Compile(el, env1)
		}

		return func(env *val.Env) *val.Val {
			l := val.List(kind, sz).List()
			for i, c := range cs {
				l.V[i] = c(env)
			}
			return l.Vl()
		}

	case ast.MAP:
		m := expr.Map()
		sz := len(m.Pairs)
		if sz == 0 {
			return func(env *val.Env) *val.Val {
				// 注意空 map 类型 map[nothing, nothing]
				kind := types.Map(types.Bottom, types.Bottom).Map()
				return val.Map(kind)
			}
		}

		kind := m.Kind.(*types.Kind).Map()
		// 保持字面量声明的执行顺序
		cs := make([]struct{ k, v compiler.Closure }, sz)
		for i, pair := range m.Pairs {
			cs[i] = struct{ k, v compiler.Closure }{Compile(pair.Key, env1), Compile(pair.Val, env1)}
		}

		return func(env *val.Env) *val.Val {
			m := val.Map(kind).Map()
			for _, c := range cs {
				k := c.k(env).Key()
				v := c.v(env)
				m.V[k] = v
			}
			return m.Vl()
		}

	case ast.OBJ:
		obj := expr.Obj()
		sz := len(obj.Fields)
		if sz == 0 {
			return func(env *val.Env) *val.Val {
				kind := types.Obj([]types.Field{}).Obj()
				return val.Obj(kind)
			}
		}
		// 保持字面量声明的执行顺序
		cs := make([]compiler.Closure, sz)
		for i, f := range obj.Fields {
			cs[i] = Compile(f.Val, env1)
		}
		kind := obj.Kind.(*types.Kind).Obj()
		return func(env *val.Env) *val.Val {
			m := val.Obj(kind).Obj()
			for i, c := range cs {
				m.V[i] = c(env)
			}
			return m.Vl()
		}

	case ast.SUBSCRIPT:
		// 也可以 desugar 成 build-in-fun
		sub := expr.Subscript()
		vac := Compile(sub.Var, env1)
		idxc := Compile(sub.Idx, env1)

		return func(env *val.Env) *val.Val {
			x := vac(env)
			switch x.Kind.Type {
			case types.TList:
				idx := int(idxc(env).Num().V)
				lst := x.List().V
				util.Assert(idx < len(lst), "out of range %d of %s", idx, x)
				return lst[idx]
			case types.TMap:
				k := idxc(env)
				v, ok := x.Map().Get(k)
				// 如果引入 null 、nil 会让类型检查复杂以及做不到安全, 可以业务逻辑里头处理 null
				util.Assert(ok, "undefined key %s of %s", k, x)
				return v
			default:
				util.Unreachable()
				return nil
			}
		}

	case ast.MEMBER:
		// 也可以 desugar 成 build-in-fun
		mem := expr.Member()
		obj := Compile(mem.Obj, env1)
		idx := mem.Index
		return func(env *val.Env) *val.Val {
			return obj(env).Obj().V[idx]
		}

	case ast.CALL:
		call := expr.Call()
		if call.Resolved == "" {
			return dynamicDispatch(env1, call)
		} else {
			return staticDispatch(env1, call)
		}

	// IF 已经 desugar 成 lazyFun 了, 这里已经没用了
	case ast.IF:
		iff := expr.If()
		cond := Compile(iff.Cond, env1)
		then := Compile(iff.Then, env1)
		els := Compile(iff.Else, env1)

		// if 分支是 lazy 的 (短路)
		return func(env *val.Env) *val.Val {
			if cond(env).Bool().V {
				return then(env)
			} else {
				return els(env)
			}
		}

	default:
		util.Unreachable()
		return nil
	}
}

// 函数在编译期 resolve, 通过 golang 闭包的 upval 传递给运行时
func staticDispatch(env1 *val.Env, call *ast.CallExpr) compiler.Closure {
	f, _ := env1.Get(call.Resolved)
	// 多态函数, 这里有点 hack, 手动狗头
	if call.Index >= 0 {
		fnTbl := *(*[]*val.FunVal)(unsafe.Pointer(f))
		f = fnTbl[call.Index].Vl()
	}

	fun := f.Fun()
	argc, cs := compileArgs(env1, call)
	return makeCallClosure(fun, argc, cs)
}

func dynamicDispatch(env1 *val.Env, call *ast.CallExpr) compiler.Closure {
	cc := Compile(call.Callee, env1)

	argc, cs := compileArgs(env1, call)
	return func(env *val.Env) *val.Val {
		fun := cc(env).Fun()
		return makeCallClosure(fun, argc, cs)(env)
	}
}

func compileArgs(env1 *val.Env, call *ast.CallExpr) (int, []compiler.Closure) {
	argc := len(call.Args)
	cs := make([]compiler.Closure, argc)
	for i, arg := range call.Args {
		cs[i] = Compile(arg, env1)
	}
	return argc, cs
}

func makeCallClosure(fun *val.FunVal, argc int, cs []compiler.Closure) func(env *val.Env) *val.Val {
	return func(env *val.Env) *val.Val {
		lazy := fun.Lazy
		params := fun.Kind.Fun().Param
		args := make([]*val.Val, argc)
		if lazy {
			// 惰性求值函数参数会被包装成 thunk, 注意没有缓存
			for i := 0; i < argc; i++ {
				args[i] = thunkify(cs[i], env, params[i])
			}
		} else {
			for i := 0; i < argc; i++ {
				args[i] = cs[i](env)
			}
		}
		return fun.Call(args...)
	}
}

func thunkify(c compiler.Closure, env *val.Env, retK *types.Kind) *val.Val {
	fk := types.Fun("thunk", []*types.Kind{}, retK)
	return val.Fun(fk, func(v ...*val.Val) *val.Val {
		return c(env)
	})
}
