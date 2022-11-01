package closure

import (
	"github.com/goghcrow/yae/debug"
	"time"

	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/parser/loc"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

func Compile(expr ast.Expr, env *val.Env) compiler.Closure {
	return compile(expr, env, false)
}

// DebugCompile 📢 调试模式
func DebugCompile(expr ast.Expr, env1 *val.Env) compiler.Closure {
	closure := compile(expr, env1, true)
	return func(env *val.Env) *val.Val {
		env.Dgb.(*debug.Record).Clear()
		return closure(env)
	}
}

func compile(expr ast.Expr, env1 *val.Env, dbg bool) compiler.Closure {
	closure := compile0(expr, env1, dbg)
	if dbg {
		return wrapForDebug(expr, closure)
	}
	return closure
}

func wrapForDebug(expr ast.Expr, cl compiler.Closure) compiler.Closure {
	// 调试模式编译过程会通过 golang 闭包将 compiler.Closure 与 col 绑定
	recordVal := func(col loc.DbgCol, cl compiler.Closure) compiler.Closure {
		return func(env *val.Env) *val.Val {
			v := cl(env)
			if rcd, ok := env.Dgb.(*debug.Record); ok {
				rcd.Rec(v, int(col)+1)
			}
			return v
		}
	}

	switch e := expr.(type) {
	case *ast.StrExpr, *ast.NumExpr, *ast.TimeExpr, *ast.BoolExpr,
		*ast.ListExpr, *ast.MapExpr, *ast.ObjExpr:
		return cl
	case *ast.IdentExpr:
		return recordVal(loc.DbgCol(e.Col), cl)
	case *ast.CallExpr:
		return recordVal(e.DbgCol, cl)
	case *ast.SubscriptExpr:
		return recordVal(e.DbgCol, cl)
	case *ast.MemberExpr:
		return recordVal(e.DbgCol, cl)
	default:
		util.Unreachable()
		return nil
	}
}

// compile0 编译成闭包
// 论文的思路 http://www.iro.umontreal.ca/~feeley/papers/FeeleyLapalmeCL87.pdf
// implement compilers for embedded languages
// 注意区分: env1 是编译期环境, env 是运行时环境
// 能在编译期完成的, 尽可能在编译期计算完
// 调用 Compile 之前必须先对 ast 调用 types.Check:
// 1. 处理字解析
// 2. 繁饰 list/map/obj 类型, 简化 Compile 代码
// 3. call.callee resolve
// 4. 且, Compile 中不检查错误, 假设 types.Check 已经全部检查
func compile0(expr ast.Expr, env1 *val.Env, dbg bool) compiler.Closure {
	switch e := expr.(type) {
	case *ast.StrExpr:
		s := val.Str(e.Val)
		return func(env *val.Env) *val.Val { return s }

	case *ast.NumExpr:
		n := val.Num(e.Val)
		return func(env *val.Env) *val.Val { return n }

	case *ast.TimeExpr:
		t := val.Time(time.Unix(e.Val, 0))
		return func(env *val.Env) *val.Val { return t }

	case *ast.BoolExpr:
		if e.Val {
			return func(env *val.Env) *val.Val { return val.True }
		} else {
			return func(env *val.Env) *val.Val { return val.False }
		}

	case *ast.ListExpr:
		els := e.Elems
		sz := len(els)
		if sz == 0 {
			return func(env *val.Env) *val.Val {
				// 注意空列表类型 list[nothing]
				ty := types.List(types.Bottom).List()
				return val.List(ty, 0)
			}
		}

		ty := e.Type.(*types.Type).List()
		cs := make([]compiler.Closure, sz)
		for i, el := range els {
			cs[i] = compile(el, env1, dbg)
		}

		return func(env *val.Env) *val.Val {
			l := val.List(ty, sz).List()
			for i, cl := range cs {
				l.V[i] = cl(env)
			}
			return l.Vl()
		}

	case *ast.MapExpr:
		sz := len(e.Pairs)
		if sz == 0 {
			return func(env *val.Env) *val.Val {
				// 注意空 map 类型 map[nothing, nothing]
				ty := types.Map(types.Bottom, types.Bottom).Map()
				return val.Map(ty)
			}
		}

		ty := e.Type.(*types.Type).Map()
		// 保持字面量声明的执行顺序
		cs := make([]struct{ k, v compiler.Closure }, sz)
		for i, pair := range e.Pairs {
			cs[i] = struct{ k, v compiler.Closure }{compile(pair.Key, env1, dbg), compile(pair.Val, env1, dbg)}
		}

		return func(env *val.Env) *val.Val {
			m := val.Map(ty).Map()
			for _, cl := range cs {
				k := cl.k(env).Key()
				v := cl.v(env)
				m.V[k] = v
			}
			return m.Vl()
		}

	case *ast.ObjExpr:
		sz := len(e.Fields)
		if sz == 0 {
			return func(env *val.Env) *val.Val {
				ty := types.Obj([]types.Field{}).Obj()
				return val.Obj(ty)
			}
		}
		// 保持字面量声明的执行顺序
		cs := make([]compiler.Closure, sz)
		for i, f := range e.Fields {
			cs[i] = compile(f.Val, env1, dbg)
		}
		ty := e.Type.(*types.Type).Obj()

		return func(env *val.Env) *val.Val {
			m := val.Obj(ty).Obj()
			for i, cl := range cs {
				m.V[i] = cl(env)
			}
			return m.Vl()
		}

	case *ast.IdentExpr:
		// 如果从性能角度考虑, 所有运行时的符号查找其实都可以从 map 换成 array
		// 1. 需要在编译期把符号 resolve 成数组下标
		// 2. 把运行时环境从 map 展开成数组
		id := e.Name
		return func(env *val.Env) *val.Val { return env.MustGet(id) }

	case *ast.CallExpr:
		if e.Resolved == "" {
			return dynamicDispatch(env1, e, dbg)
		} else {
			return staticDispatch(env1, e, dbg)
		}

	case *ast.SubscriptExpr:
		// 也可以 desugar 成 build-in-fun
		vac := compile(e.Var, env1, dbg)
		idxc := compile(e.Idx, env1, dbg)

		return func(env *val.Env) *val.Val {
			x := vac(env)
			switch x.Type.Kind {
			case types.KList:
				idx := int(idxc(env).Num().V)
				lst := x.List().V
				util.Assert(idx < len(lst), "out of range %d of %s", idx, x)
				return lst[idx]
			case types.KMap:
				k := idxc(env)
				v, ok := x.Map().Get(k)
				// 如果引入 null 、nil 会让类型检查复杂以及做不到 null 安全
				// 可以加一个返回 Maybe 的 get map 函数
				// 或者用 if(isset(m, k), m[k], default)
				util.Assert(ok, "undefined key %s of %s", k, x)
				return v
			default:
				util.Unreachable()
				return nil
			}
		}
	case *ast.MemberExpr:
		// 也可以 desugar 成 build-in-fun
		obj := compile(e.Obj, env1, dbg)
		idx := e.Index
		return func(env *val.Env) *val.Val {
			return obj(env).Obj().V[idx]
		}

	//case *ast.IfExpr:
	//	// IF 已经 desugar 成 lazyFun 了, 这里已经没用了
	//	cond := compile(e.Cond, env1)
	//	then := compile(e.Then, env1)
	//	els := compile(e.Else, env1)
	//
	//	// if 分支是 lazy 的 (短路)
	//	return func(env *val.Env) *val.Val {
	//		if cond(env).Bool().V {
	//			return then(env)
	//		} else {
	//			return els(env)
	//		}
	//	}

	default:
		util.Unreachable()
		return nil
	}
}

// 函数在编译期 resolve, 通过 golang 闭包的 upval 传递给运行时
func staticDispatch(env1 *val.Env, call *ast.CallExpr, dbg bool) compiler.Closure {
	var fun *val.FunVal
	if call.Index < 0 {
		fun = env1.MustGetMonoFun(call.Resolved)
	} else {
		fun = env1.MustGetPolyFuns(call.Resolved)[call.Index]
	}
	argc, cs := compileArgs(env1, call, dbg)
	return makeCallClosure(fun, argc, cs)
}

func dynamicDispatch(env1 *val.Env, call *ast.CallExpr, dbg bool) compiler.Closure {
	cc := compile(call.Callee, env1, dbg)

	argc, cs := compileArgs(env1, call, dbg)
	return func(env *val.Env) *val.Val {
		fun := cc(env).Fun()
		return makeCallClosure(fun, argc, cs)(env)
	}
}

func compileArgs(env1 *val.Env, call *ast.CallExpr, dbg bool) (int, []compiler.Closure) {
	argc := len(call.Args)
	cs := make([]compiler.Closure, argc)
	for i, arg := range call.Args {
		cs[i] = compile(arg, env1, dbg)
	}
	return argc, cs
}

func makeCallClosure(fun *val.FunVal, argc int, cs []compiler.Closure) func(env *val.Env) *val.Val {
	return func(env *val.Env) *val.Val {
		lazy := fun.Lazy
		params := fun.Type.Fun().Param
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

func thunkify(cl compiler.Closure, env *val.Env, retK *types.Type) *val.Val {
	fk := types.Fun("thunk", []*types.Type{}, retK)
	return val.Fun(fk, func(v ...*val.Val) *val.Val {
		return cl(env)
	})
}
