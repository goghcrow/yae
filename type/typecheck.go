package types

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/util"
	"strconv"
	"unsafe"
)

func Infer(env *Env, expr *ast.Expr) (kind *Kind, err error) {
	defer util.Recover(&err)
	return Check(env, expr), err
}

func Check(env *Env, expr *ast.Expr) *Kind {
	switch expr.Type {
	case ast.LITERAL:
		lit := expr.Literal()
		var err error
		switch lit.LitType {
		case ast.LIT_STR:
			lit.Val, err = strconv.Unquote(lit.Text)
			util.Assert(err == nil, "invalid string literal: %s", lit.Text)
			return Str
		case ast.LIT_TIME:
			// time 字面量会被 desugar 成 strtotime, 这里留着测试场景
			ts := util.Strtotime(lit.Text[1 : len(lit.Text)-1])
			util.Assert(ts != 0, "invalid time literal: %s", lit.Text)
			lit.Val = ts
			return Time
		case ast.LIT_NUM:
			lit.Val, err = util.ParseNum(lit.Text)
			util.Assert(err == nil, "invalid num literal %s", lit.Text)
			return Num
		case ast.LIT_TRUE:
			lit.Val = true
			return Bool
		case ast.LIT_FALSE:
			lit.Val = false
			return Bool
		default:
			util.Unreachable()
			return nil
		}

	case ast.IDENT:
		id := expr.Ident().Name
		util.Assert(!lexer.Reserved(id), "%s reserved", id)
		kind, ok := env.Get(id)
		util.Assert(ok, "undefined %s", id)
		return kind

	case ast.LIST:
		lst := expr.List()
		sz := len(lst.Elems)
		if sz == 0 {
			return List(Bottom)
		}

		elKind := Check(env, lst.Elems[0])
		for i := 1; i < sz; i++ {
			kind := Check(env, lst.Elems[i])
			typeAssert(elKind, kind, expr)
		}
		kind := List(elKind)
		lst.Kind = kind // ast 附加类型
		return kind

	case ast.MAP:
		m := expr.Map()
		sz := len(m.Pairs)
		if sz == 0 {
			return Map(Bottom, Bottom)
		}

		kKind := Check(env, m.Pairs[0].Key)
		util.Assert(kKind.IsPrimitive(), "invalid type of map's key: %s", kKind)
		vKind := Check(env, m.Pairs[0].Val)
		for i := 1; i < sz; i++ {
			kind := Check(env, m.Pairs[i].Key)
			typeAssert(kKind, kind, expr)
			kind = Check(env, m.Pairs[i].Val)
			typeAssert(vKind, kind, expr)
		}
		kind := Map(kKind, vKind)
		m.Kind = kind // ast 附加类型
		return kind

	case ast.OBJ:
		obj := expr.Obj()
		sz := len(obj.Fields)
		if sz == 0 {
			return Obj([]Field{})
		}

		fs := make([]Field, sz)
		for i, f := range obj.Fields {
			fs[i] = Field{f.Name, Check(env, f.Val)}
		}

		kind := Obj(fs)
		obj.Kind = kind // ast 附加类型
		return kind

	case ast.SUBSCRIPT:
		sub := expr.Subscript()
		varKind := Check(env, sub.Var)

		if varKind.Type == TList {
			idxKind := Check(env, sub.Idx)
			typeAssert(idxKind, Num, expr)
			elKind := varKind.List().El
			return elKind
		}
		if varKind.Type == TMap {
			idxKind := Check(env, sub.Idx)
			keyKind := varKind.Map().Key
			valKind := varKind.Map().Val
			typeAssert(idxKind, keyKind, expr)
			return valKind
		}
		util.Assert(false,
			"type mismatched, expect `list | map` actual `%s` in `%s`", varKind, sub.Var)
		return nil

	case ast.MEMBER:
		mem := expr.Member()
		objKind := Check(env, mem.Obj)
		util.Assert(objKind.Type == TObj,
			"type mismatched, expect `%s` actual `%s` in `%s`", TObj, objKind, expr)
		obj := objKind.Obj()
		fName := mem.Field.Name
		f, ok := obj.GetField(fName)
		util.Assert(ok, "undefined filed `%s` of `%s` in `%s`", fName, objKind, expr)
		mem.Index = obj.Index[fName] // attach obj index
		return f.Val

	case ast.CALL:
		call := expr.Call()
		callee := call.Callee

		argSz := len(call.Args)
		args := make([]*Kind, argSz)
		for i := 0; i < argSz; i++ {
			args[i] = Check(env, call.Args[i])
		}

		var fun *FunKind

		// util.Assert(callee.Type == ast.IDENT, "invalid callable %s in %s", callee, expr)
		if callee.Type == ast.IDENT {
			fn := callee.Ident().Name
			fun = resolveFun(env, call, fn, args)
		} else {
			f := Check(env, callee)
			util.Assert(f.Type == TFun, "non callable of `%s` in `%s`", callee, call)
			fun = inferFun(f.Fun(), args)
			util.Assert(fun != nil, "args `%s` mismatch fun `%s`", args, f)
		}

		paramSz := len(fun.Param)
		arityAssert(paramSz, argSz, callee)

		for i := 0; i < paramSz; i++ {
			paramKind := fun.Param[i]
			argKind := args[i]
			typeAssert(paramKind, argKind, expr)
		}

		return fun.Return

	// IF 已经 desugar 成 lazyFun 了, 这里已经没用了
	case ast.IF:
		iff := expr.If()
		condKind := Check(env, iff.Cond)
		typeAssert(condKind, Bool, expr)
		thenKind := Check(env, iff.Then)
		elseKind := Check(env, iff.Else)
		typeAssert(thenKind, elseKind, expr)
		return thenKind

	default:
		util.Unreachable()
		return nil
	}
}

// desugar 会把所有操作符都转换成函数调用, 这里会统一处理操作符和函数
//
//goland:noinspection SpellCheckingInspection
func resolveFun(env *Env, call *ast.CallExpr, fnName string, args []*Kind) *FunKind {
	// 1. 首先尝试 resolve mono fn
	monofk, mono := Fun(fnName, args, Bottom /*返回类型无所谓*/).Fun().Lookup()
	util.Assert(mono, "unexpected")
	f, ok := env.Get(monofk)
	if ok {
		util.Assert(f.Type == TFun, "non callable of %s in %s", call.Resolved, call)
		call.Resolved = monofk // ast 标记 callee 在环境中的 key
		monofun := f.Fun()
		return monofun
	}

	// 2. 然后依次尝试 poly fn
	// 这里 hack 参见: type/env.go::RegisterFun
	// 先按 `函数名+参数个数` 查找重载的函数列表(包括泛型函数)
	polyfk, _ := Fun(fnName, args, Slot("α")).Fun().Lookup()
	fs, ok := env.Get(polyfk)
	util.Assert(ok, "func `%s` has no overload for params`%s`", fnName, args)

	// 然后在重载函数列表中依次查找
	// 因为不支持子类型, 所以也没有最适合的规则, 找到匹配为止
	// 并实例化函数 poly ~~> mono
	fks := (*[]*FunKind)(unsafe.Pointer(fs))
	for i, f := range *fks {
		util.Assert(f.Type == TFun, "non callable of %s in %s", fnName, call)
		monof := inferFun(f, args)
		if monof == nil {
			continue
		}
		call.Resolved = polyfk // ast 标记 callee 在环境中的 key
		call.Index = i         // 以及在泛型函数表的中的位置
		return monof
	}
	util.Assert(false, "func `%s` has no overload for params`%s`", fnName, args)
	return nil
}

//goland:noinspection SpellCheckingInspection
func inferFun(f *FunKind, args []*Kind) *FunKind {
	// 1. 构造 psuido fun
	sx := make([]*Kind, len(args))
	for i := 0; i < len(args); i++ {
		sx[i] = Slot("s")
	}
	s := Tuple(sx)
	t := Slot("t")
	psuidoFun := Fun(f.Name, []*Kind{s}, t)

	// 2. 需要被 infer 的 fun
	fun := Fun(f.Name, []*Kind{Tuple(f.Param)}, f.Return)

	// 3. 在环境 m 中 unify
	m := map[string]*Kind{}
	unifyFun := Unify(psuidoFun, fun, m)
	if unifyFun == nil {
		return nil
	}

	// 4. 替换得到参数类型
	targ := Tuple(args)
	targ1 := applySubst(s, m)
	targ2 := Unify(targ1, targ, m)
	if targ2 == nil || targ2.Type != TTuple {
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

func slotFree(k *Kind) bool {
	switch k.Type {
	case TNum, TStr, TBool, TTime, TBottom, TTop:
		return true
	case TSlot:
		return false
	case TList:
		return slotFree(k.List().El)
	case TMap:
		return slotFree(k.Map().Key) && slotFree(k.Map().Val)
	case TTuple:
		for _, vk := range k.Tuple().Val {
			if !slotFree(vk) {
				return false
			}
		}
		return true
	case TObj:
		for _, f := range k.Obj().Fields {
			if !slotFree(f.Val) {
				return false
			}
		}
		return true
	case TFun:
		for _, param := range k.Fun().Param {
			if !slotFree(param) {
				return false
			}
		}
		return slotFree(k.Fun().Return)
	default:
		util.Unreachable()
		return false
	}
}

func arityAssert(expect, actual int, f *ast.Expr) {
	util.Assert(expect == actual, "arity mismatch, expect %d actual %d in `%s`", expect, actual, f)
}

func typeAssert(expect, actual *Kind, f *ast.Expr) {
	eq := Equals(expect, actual)
	util.Assert(eq, "type mismatched, expect `%s` actual `%s` in `%s`", expect, actual, f)
}
