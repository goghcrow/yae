package types

import (
	"testing"
)

func arrow(param *Type, ret *Type) *Type {
	return Fun("", []*Type{param}, ret)
}
func product(ks ...*Type) *Type {
	return Tuple(ks)
}

//goland:noinspection NonAsciiCharacters
func TestUnify(t *testing.T) {

	t.Run("", func(t *testing.T) {
		// A :: ((α → β) × [γ]) → [β]
		// B :: ((γ → δ) × [γ]) → ε
		β := TyVar("β")
		A := arrow(
			product(arrow(TyVar("α"), β), List(TyVar("γ"))),
			List(β),
		)

		γ := TyVar("γ")
		B := arrow(
			product(arrow(γ, TyVar("δ")), List(γ)),
			TyVar("ε"),
		)
		m := map[string]*Type{}
		u := Unify(A, B, m)
		t.Log(u, m)
		// m [α2:γ4 β1:δ5 γ3:γ4 ε6:list[δ5]]

		assert(len(m) == 4)
		assert(Equals(applySubst(A, m), applySubst(B, m)))
	})

	t.Run("a = a\t{ a = a }\t{}\tSucceeds. (tautology)", func(t *testing.T) {
		m := map[string]*Type{}
		a := Num
		u := Unify(a, a, m)
		t.Log(m, u)
		assert(Equals(u, a))
		assert(len(m) == 0)
	})

	t.Run("a = b\t{ a = b }\t⊥\ta and b do not match", func(t *testing.T) {
		m := map[string]*Type{}
		a := Num
		b := Str
		u := Unify(a, b, m)
		t.Log(m, u)
		assert(u == nil)
		assert(len(m) == 0)
	})

	t.Run("X = X\t{ x = x }\t{}", func(t *testing.T) {
		m := map[string]*Type{}
		X := TyVar("X")
		u := Unify(X, X, m)
		t.Log(m, u)
		assert(Equals(u, X))
		assert(len(m) == 0)
	})

	t.Run("a = X\t{ a = x }\t{ x ↦ a }\tx is unified with the constant a", func(t *testing.T) {
		m := map[string]*Type{}
		a := Num
		X := TyVar("X")
		u := Unify(a, X, m)
		t.Log(m, u)
		assert(Equals(u, a))
		assert(len(m) == 1)
		assert(Equals(m[X.TyVar().Name], a))
	})

	t.Run("X = Y\t{ x = y }\t{ x ↦ y }", func(t *testing.T) {
		m := map[string]*Type{}
		X := TyVar("X")
		Y := TyVar("Y")
		u := Unify(X, Y, m)
		t.Log(m, u)
		assert(len(m) == 1)
		assert(Equals(m[X.TyVar().Name], Y))
	})

	t.Run("f(a,X) = f(a,b)\t{ f(a,x) = f(a,b) }\t{ x ↦ b }\tfunction and constant symbols match, x is unified with the constant b", func(t *testing.T) {
		m := map[string]*Type{}
		a := Num
		b := Str
		X := TyVar("X")
		fax := arrow(a, X)
		fab := arrow(a, b)
		u := Unify(fax, fab, m)
		t.Log(m, u)
		assert(Equals(u, fab))
		assert(len(m) == 1)
		assert(Equals(m[X.TyVar().Name], b))
	})

	t.Run("f(X) = f(Y)\t{ f(x) = f(y) }\t{ x ↦ y }\tx and y are aliased", func(t *testing.T) {
		m := map[string]*Type{}
		X := TyVar("X")
		Y := TyVar("Y")
		fx := arrow(X, X)
		fy := arrow(Y, Y)
		u := Unify(fx, fy, m)

		assert(len(m) == 1)
		assert(Equals(m[X.TyVar().Name], Y))

		fx1 := applySubst(fx, m)
		fy1 := applySubst(fy, m)
		t.Log(m, u, fx1, fy1)
		assert(Equals(fx1, fy))
		assert(Equals(fy, fy))
	})
}

func TestInfer1(tt *testing.T) {
	// id :: a -> a
	a := TyVar("a")
	fun := arrow(a, a)

	// psuido :: s -> t
	s := TyVar("s")
	t := TyVar("t")
	psuido := arrow(s, t)

	m := map[string]*Type{}
	tfn1 := Unify(psuido, fun, m)

	tt.Log(tfn1, m)
	assert(len(m) == 2)
	assert(Equals(m[s.TyVar().Name], a))
	assert(Equals(m[t.TyVar().Name], a))

	//--------------------------------

	// s 换成 a
	targ1 := applySubst(s, m)

	// a 换成 str
	targ := Str
	targ2 := Unify(targ1, targ, m)

	tt.Log(targ1, "=>", targ2)
	tt.Log(m) // m 已经有返回值

	tresult := applySubst(t, m)
	tt.Log(tresult)

	assert(Equals(tresult, Str))
}

func TestRecursive(t *testing.T) {
	lt := List(nil)
	lt.List().El = lt
	t.Log(lt)

	func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()
		Unify(lt, lt, map[string]*Type{})
		t.Fail()
	}()
}

func assert(cond bool) {
	if !cond {
		panic("")
	}
}
