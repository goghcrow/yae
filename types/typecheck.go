package types

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/util"
)

func Infer(expr ast.Expr, env *Env) (ty *Type, err error) {
	defer util.Recover(&err)
	return Check(expr, env), err
}

func Check(expr ast.Expr, env *Env) *Type {
	switch e := expr.(type) {
	case *ast.StrExpr:
		return Str
	case *ast.NumExpr:
		return Num
	case *ast.TimeExpr:
		// time 字面量会被 desugar 成 strtotime, 这里留着测试场景
		return Time
	case *ast.BoolExpr:
		return Bool
	case *ast.ListExpr:
		sz := len(e.Elems)
		if sz == 0 {
			ty := List(Bottom)
			e.Type = ty // attach ast
			return ty
		}

		elTy := Check(e.Elems[0], env)
		for i := 1; i < sz; i++ {
			ty := Check(e.Elems[i], env)
			typeAssert(elTy, ty, expr)
		}
		ty := List(elTy)
		e.Type = ty // attach ast
		return ty
	case *ast.MapExpr:
		sz := len(e.Pairs)
		if sz == 0 {
			ty := Map(Bottom, Bottom)
			e.Type = ty // attach ast
			return ty
		}

		kTy := Check(e.Pairs[0].Key, env)
		util.Assert(kTy.IsPrimitive(), "invalid type of map's key: %s", kTy)
		vTy := Check(e.Pairs[0].Val, env)
		for i := 1; i < sz; i++ {
			ty := Check(e.Pairs[i].Key, env)
			typeAssert(kTy, ty, expr)
			ty = Check(e.Pairs[i].Val, env)
			typeAssert(vTy, ty, expr)
		}
		ty := Map(kTy, vTy)
		e.Type = ty // attach ast
		return ty
	case *ast.ObjExpr:
		sz := len(e.Fields)
		if sz == 0 {
			ty := Obj([]Field{})
			e.Type = ty // attach ast
			return ty
		}

		fs := make([]Field, sz)
		for i, f := range e.Fields {
			fs[i] = Field{f.Name, Check(f.Val, env)}
		}

		ty := Obj(fs)
		e.Type = ty // attach ast
		return ty
	case *ast.IdentExpr:
		id := e.Name
		util.Assert(!lexer.Reserved(id), "%s reserved", id)
		ty, ok := env.Get(id)
		util.Assert(ok, "undefined %s", id)
		return ty
	case *ast.CallExpr:
		callee := e.Callee

		argSz := len(e.Args)
		args := make([]*Type, argSz)
		for i := 0; i < argSz; i++ {
			args[i] = Check(e.Args[i], env)
		}

		var fun *FunTy

		// util.Assert(callee.Type == ast.IDENT, "invalid callable %s in %s", callee, expr)
		if ident, ok := callee.(*ast.IdentExpr); ok {
			fn := ident.Name
			fun = resolveFun(env, e, fn, args)
		} else {
			f := Check(callee, env)
			util.Assert(f.Kind == KFun, "non callable of `%s` in `%s`", callee, e)
			fun = inferFun(f.Fun(), args)
			util.Assert(fun != nil, "args `%s` mismatch fun `%s`", args, f)
		}

		paramSz := len(fun.Param)
		arityAssert(paramSz, argSz, callee)

		for i := 0; i < paramSz; i++ {
			paramTy := fun.Param[i]
			argTy := args[i]
			typeAssert(paramTy, argTy, expr)
		}

		e.CalleeType = fun.Ty() // attach ast
		return fun.Return
	case *ast.SubscriptExpr:
		varTy := Check(e.Var, env)

		switch varTy.Kind {
		case KList:
			idxTy := Check(e.Idx, env)
			typeAssert(idxTy, Num, expr)
			elTy := varTy.List().El
			e.VarType = varTy // attach ast
			return elTy
		case KMap:
			idxTy := Check(e.Idx, env)
			keyTy := varTy.Map().Key
			valTy := varTy.Map().Val
			typeAssert(idxTy, keyTy, expr)
			e.VarType = varTy // attach ast
			return valTy
		default:
			util.Assert(false,
				"type mismatched, expect `list | map` actual `%s` in `%s`", varTy, e.Var)
			return nil
		}
	case *ast.MemberExpr:
		objTy := Check(e.Obj, env)
		util.Assert(objTy.Kind == KObj,
			"type mismatched, expect `%s` actual `%s` in `%s`", KObj, objTy, expr)
		obj := objTy.Obj()
		fName := e.Field.Name
		f, ok := obj.GetField(fName)
		util.Assert(ok, "undefined filed `%s` of `%s` in `%s`", fName, objTy, expr)
		e.ObjType = objTy          // attach ast
		e.Index = obj.Index[fName] // attach ast, obj index
		return f.Val
	//case *ast.IfExpr:
	//	// IF 已经 desugar 成 lazyFun 了, 这里已经没用了
	//	condKind := Check(e.Cond, env)
	//	typeAssert(condKind, Bool, expr)
	//	thenKind := Check(e.Then, env)
	//	elseKind := Check(e.Else, env)
	//	typeAssert(thenKind, elseKind, expr)
	//	return thenKind
	default:
		util.Unreachable()
		return nil
	}
}

// desugar 会把所有操作符都转换成函数调用, 这里会统一处理操作符和函数
//
//goland:noinspection SpellCheckingInspection
func resolveFun(env *Env, call *ast.CallExpr, fnName string, args []*Type) *FunTy {
	// 1. 首先尝试 resolve mono fn
	monofnTy, mono := Fun(fnName, args, Bottom /*返回类型无所谓*/).Fun().Lookup()
	util.Assert(mono, "unexpected")
	f, ok := env.GetMonoFun(monofnTy)
	if ok {
		util.Assert(f.Kind == KFun, "non callable of %s in %s", monofnTy, call)
		call.Resolved = monofnTy // attach ast, 标记 callee 在环境中的 key
		monofun := f.Fun()
		return monofun
	}

	// 2. 然后依次尝试 poly fn
	// 先按 `函数名+参数个数` 查找重载的函数列表(包括泛型函数)
	polyfnTy, _ := Fun(fnName, args, TyVar("α")).Fun().Lookup()
	fks, ok := env.GetPolyFuns(polyfnTy)
	util.Assert(ok, "func `%s` has no overload func for params`%s`", fnName, Tuple(args))

	// 然后在重载函数列表中依次查找
	// 因为不支持子类型, 所以也没有最适合的规则, 找到匹配为止
	// 并实例化函数 poly ~~> mono
	for i, f := range fks {
		util.Assert(f.Kind == KFun, "non callable of %s in %s", fnName, call)
		monof := inferFun(f, args)
		if monof == nil {
			continue
		}
		call.Resolved = polyfnTy // attach ast, 标记 callee 在环境中的 key
		call.Index = i           // attach ast, 以及在泛型函数表的中的位置
		return monof
	}
	util.Assert(false, "func `%s` has no overload func for params`%s`", fnName, Tuple(args))
	return nil
}

//goland:noinspection SpellCheckingInspection
func inferFun(f *FunTy, args []*Type) *FunTy {
	// 1. 构造 psuido fun
	sx := make([]*Type, len(args))
	for i := 0; i < len(args); i++ {
		sx[i] = TyVar("s")
	}
	s := Tuple(sx)
	t := TyVar("t")
	psuidoFun := Fun(f.Name, []*Type{s}, t)

	// 2. 需要被 infer 的 fun
	fun := Fun(f.Name, []*Type{Tuple(f.Param)}, f.Return)

	// 3. 在环境 m 中 unify
	m := map[string]*Type{}
	unifyFun := Unify(psuidoFun, fun, m)
	if unifyFun == nil {
		return nil
	}

	// 4. 替换得到参数类型
	targ := Tuple(args)
	targ1 := applySubst(s, m)
	targ2 := Unify(targ1, targ, m)
	if targ2 == nil || targ2.Kind != kTuple {
		return nil
	}

	// 5. 替换得到返回类型
	// 返回值必须是具体类型
	tresult := applySubst(t, m)
	if !slotFree(tresult) {
		return nil
	}

	return Fun(f.Name, targ2.Tuple().Val, tresult).Fun()
}

func slotFree(ty *Type) bool {
	switch ty.Kind {
	case KNum, KStr, KBool, KTime, KBot, KTop:
		return true
	case KTyVar:
		return false
	case KList:
		return slotFree(ty.List().El)
	case KMap:
		return slotFree(ty.Map().Key) && slotFree(ty.Map().Val)
	case kTuple:
		for _, vk := range ty.Tuple().Val {
			if !slotFree(vk) {
				return false
			}
		}
		return true
	case KObj:
		for _, f := range ty.Obj().Fields {
			if !slotFree(f.Val) {
				return false
			}
		}
		return true
	case KFun:
		for _, param := range ty.Fun().Param {
			if !slotFree(param) {
				return false
			}
		}
		return slotFree(ty.Fun().Return)
	case KMaybe:
		return slotFree(ty.Maybe().Elem)
	default:
		util.Unreachable()
		return false
	}
}

func arityAssert(expect, actual int, f ast.Expr) {
	util.Assert(expect == actual, "arity mismatch, expect %d actual %d in `%s`", expect, actual, f)
}

func typeAssert(expect, actual *Type, f ast.Expr) {
	eq := Equals(expect, actual)
	util.Assert(eq, "type mismatched, expect `%s` actual `%s` in `%s`", expect, actual, f)
}
