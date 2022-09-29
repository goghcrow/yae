package types

import (
	"github.com/goghcrow/yae/ast"
	lex "github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/util"
	"strconv"
	"unsafe"
)

// TypeCheck check & infer
func TypeCheck(env *Env, expr *ast.Expr) *Kind {
	switch expr.Type {
	case ast.LITERAL:
		lit := expr.Literal()
		switch lit.LitType {
		case ast.LIT_STR:
			_, err := strconv.Unquote(lit.Val)
			util.Assert(err == nil, "invalid string literal: %s", lit.Val)
			return Str
		case ast.LIT_NUM:
			_, err := util.ParseNum(lit.Val)
			util.Assert(err == nil, "invalid num literal %s", lit.Val)
			return Num
		case ast.LIT_TRUE:
			return Bool
		case ast.LIT_FALSE:
			return Bool
		}

	case ast.IDENT:
		id := expr.Ident().Name
		util.Assert(!lex.Reserved(id), "%s reserved", id)
		kind, ok := env.Get(id)
		util.Assert(ok, "undefined %s", id)
		return kind

	case ast.LIST:
		lst := expr.List()
		sz := len(lst.Elems)
		if sz == 0 {
			return List(Bottom)
		}

		elKind := TypeCheck(env, lst.Elems[0])
		for i := 1; i < sz; i++ {
			kind := TypeCheck(env, lst.Elems[i])
			typeAssert(elKind, kind, expr)
		}
		kind := List(elKind)
		lst.Kind = kind
		return kind

	case ast.MAP:
		m := expr.Map()
		sz := len(m.Pairs)
		if sz == 0 {
			return Map(Bottom, Bottom)
		}

		kKind := TypeCheck(env, m.Pairs[0].Key)
		util.Assert(kKind.IsPrimitive(), "invalid type of map's key: %s", kKind)
		vKind := TypeCheck(env, m.Pairs[0].Val)
		for i := 1; i < sz; i++ {
			kind := TypeCheck(env, m.Pairs[i].Key)
			typeAssert(kKind, kind, expr)
			kind = TypeCheck(env, m.Pairs[i].Val)
			typeAssert(vKind, kind, expr)
		}
		kind := Map(kKind, vKind)
		m.Kind = kind
		return kind

	case ast.OBJ:
		obj := expr.Obj()
		sz := len(obj.Fields)
		if sz == 0 {
			return Obj(map[string]*Kind{})
		}

		fs := make(map[string]*Kind, sz)
		for name, val := range obj.Fields {
			_, ok := fs[name]
			util.Assert(!ok, "duplicated field %s in %s", name, expr)
			fs[name] = TypeCheck(env, val)
		}

		kind := Obj(fs)
		obj.Kind = kind
		return kind

	case ast.SUBSCRIPT:
		sub := expr.Subscript()
		varKind := TypeCheck(env, sub.Var)

		if varKind.Type == TList {
			idxKind := TypeCheck(env, sub.Idx)
			typeAssert(idxKind, Num, expr)
			elKind := varKind.List().El
			return elKind
		}
		if varKind.Type == TMap {
			idxKind := TypeCheck(env, sub.Idx)
			keyKind := varKind.Map().Key
			valKind := varKind.Map().Val
			typeAssert(idxKind, keyKind, expr)
			return valKind
		}
		util.Assert(false,
			"type mismatched, expect `list | map` actual `%s` in `%s`", varKind, sub.Var)

	case ast.MEMBER:
		mem := expr.Member()
		objKind := TypeCheck(env, mem.Obj)
		util.Assert(objKind.Type == TObj,
			"type mismatched, expect `%s` actual `%s` in `%s`", TObj, objKind, expr)
		fn := mem.Field.Name
		fk, ok := objKind.Obj().Fields[fn]
		util.Assert(ok, "undefined filed `%s` of `%s` in `%s`", fn, objKind, expr)
		return fk

	case ast.CALL:
		call := expr.Call()
		callee := call.Callee

		argSz := len(call.Args)
		args := make([]*Kind, argSz)
		for i := 0; i < argSz; i++ {
			args[i] = TypeCheck(env, call.Args[i])
		}

		var fun *FunKind

		if callee.Type == ast.IDENT {
			util.Assert(callee.Type == ast.IDENT, "invalid callable %s in %s", callee, expr)
			fn := callee.Ident().Name
			fun = resolveFun(env, call, fn, args)
		} else {
			f := TypeCheck(env, callee)
			util.Assert(f.Type == TFun, "non callable of `%s` in `%s`", callee, call)
			fun = inferFun(f.Fun(), args)
			util.Assert(fun != nil, "args `%s` mismatch fun `%s`", args, f)
		}

		paramSz := len(fun.Param)
		arityAssert(paramSz, argSz, callee)

		for i := 0; i < paramSz; i++ {
			paramKind := fun.Param[i]
			argKind := TypeCheck(env, call.Args[i])
			typeAssert(paramKind, argKind, expr)
		}

		return fun.Return

	// IF 已经 desugar 成 lazyfun 了, 这里已经没用了
	case ast.IF:
		iff := expr.If()
		condKind := TypeCheck(env, iff.Cond)
		typeAssert(condKind, Bool, expr)
		thenKind := TypeCheck(env, iff.Then)
		elseKind := TypeCheck(env, iff.Else)
		typeAssert(thenKind, elseKind, expr)
		return thenKind

	default:
		util.Unreachable()
	}

	return nil
}

//goland:noinspection SpellCheckingInspection
func resolveFun(env *Env, call *ast.CallExpr, fnName string, args []*Kind) *FunKind {
	monofk, mono := Fun(fnName, args, Bottom /*返回类型无所谓*/).Fun().Lookup()
	util.Assert(mono, "")
	f, ok := env.Get(monofk)
	if ok {
		util.Assert(f.Type == TFun, "non callable of %s in %s", call.Resolved, call)
		call.Resolved = monofk
		monofun := f.Fun()
		return monofun
	}

	polyfk, _ := Fun(fnName, args, Slot("α")).Fun().Lookup()
	fs, ok := env.Get(polyfk)
	util.Assert(ok, "undefined fun %s in %s", monofk, call)

	fks := (*[]*FunKind)(unsafe.Pointer(fs))
	for i, f := range *fks {
		util.Assert(f.Type == TFun, "non callable of %s in %s", fnName, call)
		monof := inferFun(f, args)
		if monof == nil {
			continue
		}
		call.Resolved = polyfk
		call.Index = i
		return monof
	}
	util.Assert(false, "undefined fun `%s` in `%s`", monofk, call)
	return nil
}

//goland:noinspection SpellCheckingInspection
func inferFun(f *FunKind, args []*Kind) *FunKind {
	fun := Fun(f.Name, []*Kind{Tuple(f.Param)}, f.Return)
	targ := Tuple(args)

	sx := make([]*Kind, len(args))
	for i := 0; i < len(args); i++ {
		sx[i] = Slot("s")
	}
	s := Tuple(sx)
	t := Slot("t")
	psuido := Fun(f.Name, []*Kind{s}, t)

	m := map[string]*Kind{}
	//fmt.Println("_________")
	tfn1 := unify(psuido, fun, m)
	//fmt.Println(m)
	if tfn1 == nil {
		return nil
	}

	//fmt.Println("_________")
	//fmt.Println(tfn1)

	//fmt.Println("_________")
	//fmt.Println(s)
	targ1 := subst(s, m)
	//fmt.Println(targ1)
	//fmt.Println("_________")

	targ2 := unify(targ1, targ, m)
	//fmt.Println(m)
	if targ2 == nil || targ2.Type != TTuple {
		return nil
	}
	//fmt.Println(targ2)
	//fmt.Println("_________")

	tresult := subst(t, m)
	if !slotFree(tresult) {
		return nil
	}
	//fmt.Println(tresult)

	return Fun(f.Name, targ2.Tuple().Val, tresult).Fun()
}

func slotFree(k *Kind) bool {
	switch k.Type {
	case TNum:
		return true
	case TStr:
		return true
	case TBool:
		return true
	case TTime:
		return true
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
		for _, fk := range k.Obj().Fields {
			if !slotFree(fk) {
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
	case TSlot:
		return false
	case TTop:
		return true
	case TBottom:
		return true
	default:
		util.Unreachable()
	}
	return false
}

func arityAssert(expect, actual int, f *ast.Expr) {
	util.Assert(expect == actual, "arity mismatch, expect %d actual %d in `%s`", expect, actual, f)
}

func typeAssert(expect, actual *Kind, f *ast.Expr) {
	eq := Equals(expect, actual)
	util.Assert(eq, "type mismatched, expect `%s` actual `%s` in `%s`", expect, actual, f)
}
