package compile

import (
	"github.com/goghcrow/yae/ast"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"strconv"
	"unsafe"
)

type Closure func(env *val.Env) *val.Val

func (c Closure) String() string { return "Closure" }

// Compile 编译成闭包
// 论文的思路 http://www.iro.umontreal.ca/~feeley/papers/FeeleyLapalmeCL87.pdf
// implement compilers for embedded languages
// 注意区分: env1 是编译期环境, env 是运行时环境
// 能在编译期完成的, 尽可能在编译期计算完
// 调用 compile 之前必须调用 types.TypeCheck 需要在 ast 中繁饰部分类型信息
func Compile(env1 *val.Env, expr *ast.Expr) Closure {
	switch expr.Type {
	case ast.LITERAL:
		lit := expr.Literal()
		switch lit.LitType {
		case ast.LIT_STR:
			unquote, _ := strconv.Unquote(lit.Val)
			s := val.Str(unquote)
			return func(env *val.Env) *val.Val { return s }
		//case ast.LIT_TIME:
		//	t := time.Unix(util.Strtotime(lit.Val[1:len(lit.Val)-1]), 0)
		//	return func(env *val.Env) *val.Val { return val.Time(t) }
		case ast.LIT_NUM:
			v, _ := util.ParseNum(lit.Val)
			n := val.Num(v)
			return func(env *val.Env) *val.Val { return n }
		case ast.LIT_TRUE:
			return func(env *val.Env) *val.Val { return val.True }
		case ast.LIT_FALSE:
			return func(env *val.Env) *val.Val { return val.False }
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
				kind := types.List(types.Bottom).List()
				return val.List(kind, 0)
			}
		}

		kind := lst.Kind.(*types.Kind).List()
		cs := make([]Closure, sz)
		for i, el := range els {
			cs[i] = Compile(env1, el)
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
				kind := types.Map(types.Bottom, types.Bottom).Map()
				return val.Map(kind)
			}
		}

		kind := m.Kind.(*types.Kind).Map()
		cs := make([]struct{ k, v Closure }, sz)
		for i, pair := range m.Pairs {
			cs[i] = struct{ k, v Closure }{Compile(env1, pair.Key), Compile(env1, pair.Val)}
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
		fs := obj.Fields
		sz := len(fs)
		if sz == 0 {
			return func(env *val.Env) *val.Val {
				kind := types.Obj(map[string]*types.Kind{}).Obj()
				return val.Obj(kind)
			}
		}

		kind := obj.Kind.(*types.Kind).Obj()
		cs := make(map[string]Closure, sz)
		for n, v := range fs {
			cs[n] = Compile(env1, v)
		}

		return func(env *val.Env) *val.Val {
			k := make(map[string]*types.Kind, sz)
			m := val.Obj(kind).Obj()
			for n, c := range cs {
				v := c(env)
				m.V[n] = v
				k[n] = v.Kind
			}
			return m.Vl()
		}

	case ast.SUBSCRIPT:
		// 也可以 desugar 成 build-in-fun
		sub := expr.Subscript()
		va := Compile(env1, sub.Var)
		idx := Compile(env1, sub.Idx)

		return func(env *val.Env) *val.Val {
			x := va(env)
			if x.Kind.Type == types.TList {
				return x.List().V[int(idx(env).Num().V)]
			} else if x.Kind.Type == types.TMap {
				v, ok := x.Map().Get(idx(env))
				// 这里如果引入 null 、nil 之类又会牵扯到子类型, 实际引入了 bottom 类型
				// 可以业务逻辑里头处理 null, 或者引入一个特殊函数(e.g. isset、has)处理
				// 不能定义成普通函数, 因为要 hack 类型检查
				util.Assert(ok, "undefined key %s of %s", idx(env), x.Map())
				return v
			}
			util.Unreachable()
			return nil
		}

	case ast.MEMBER:
		// 也可以 desugar 成 build-in-fun
		mem := expr.Member()
		obj := Compile(env1, mem.Obj)
		field := mem.Field.Name

		return func(env *val.Env) *val.Val {
			v, _ := obj(env).Obj().V[field]
			return v
		}

	case ast.CALL:
		call := expr.Call()
		if call.Resolved == "" {
			return dynamicDispatch(env1, call)
		} else {
			return staticDispatch(env1, call)
		}

	//IF 已经 desugar 成 lazyfun 了, 这里已经没用了
	//case ast.IF:
	//	iff := expr.If()
	//	cond := Compile(env1, iff.Cond)
	//	then := Compile(env1, iff.Then)
	//	els := Compile(env1, iff.Else)
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
	}
	return nil
}

// 函数在编译期 resolve, 通过 golang 闭包的 upval 传递给运行时
func staticDispatch(env1 *val.Env, call *ast.CallExpr) Closure {
	f, _ := env1.Get(call.Resolved)

	// 多态函数, 这里有点 hack 手动狗头
	if call.Index >= 0 {
		f = (*(*[]*val.FunVal)(unsafe.Pointer(f)))[call.Index].Vl()
	}
	fun := f.Fun()

	sz := len(call.Args)
	cs := make([]Closure, sz)
	for i, arg := range call.Args {
		cs[i] = Compile(env1, arg)
	}

	lazy := fun.Lazy
	retK := fun.Kind.Fun().Return

	return func(env *val.Env) *val.Val {
		args := make([]*val.Val, sz)
		if lazy {
			// 惰性求值函数参数会被包装成 thunk, 注意没有缓存
			for i := 0; i < sz; i++ {
				args[i] = thunkify(cs[i], env, retK)
			}
		} else {
			for i := 0; i < sz; i++ {
				args[i] = cs[i](env)
			}
		}
		return fun.V(args...)
	}
}

func dynamicDispatch(env1 *val.Env, call *ast.CallExpr) Closure {
	cc := Compile(env1, call.Callee)

	sz := len(call.Args)
	cs := make([]Closure, sz)
	for i, arg := range call.Args {
		cs[i] = Compile(env1, arg)
	}

	return func(env *val.Env) *val.Val {
		fun := cc(env).Fun()
		lazy := fun.Lazy
		retK := fun.Kind.Fun().Return

		args := make([]*val.Val, sz)
		if lazy {
			// 惰性求值函数参数会被包装成 thunk, 注意没有缓存
			for i := 0; i < sz; i++ {
				args[i] = thunkify(cs[i], env, retK)
			}
		} else {
			for i := 0; i < sz; i++ {
				args[i] = cs[i](env)
			}
		}
		return fun.V(args...)
	}
}

func thunkify(c Closure, env *val.Env, retK *types.Kind) *val.Val {
	fk := types.Fun("thunk", []*types.Kind{}, retK)
	return val.Fun(fk, func(v ...*val.Val) *val.Val {
		return c(env)
	})
}
