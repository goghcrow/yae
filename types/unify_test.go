package types

import (
	"testing"
)

func arrow(param *Kind, ret *Kind) *Kind {
	return Fun("", []*Kind{param}, ret)
}
func product(ks ...*Kind) *Kind {
	return Tuple(ks)
}

//goland:noinspection NonAsciiCharacters
func TestUnify(t *testing.T) {

	t.Run("", func(t *testing.T) {
		// A :: ((α → β) × [γ]) → [β]
		// B :: ((γ → δ) × [γ]) → ε
		β := Slot("β")
		A := arrow(
			product(arrow(Slot("α"), β), List(Slot("γ"))),
			List(β),
		)

		γ := Slot("γ")
		B := arrow(
			product(arrow(γ, Slot("δ")), List(γ)),
			Slot("ε"),
		)
		m := map[string]*Kind{}
		u := Unify(A, B, m)
		t.Log(u, m)
		// m [α2:γ4 β1:δ5 γ3:γ4 ε6:list[δ5]]

		assert(len(m) == 4)
		assert(Equals(applySubst(A, m), applySubst(B, m)))
	})

	t.Run("a = a\t{ a = a }\t{}\tSucceeds. (tautology)", func(t *testing.T) {
		m := map[string]*Kind{}
		a := Num
		u := Unify(a, a, m)
		t.Log(m, u)
		assert(Equals(u, a))
		assert(len(m) == 0)
	})

	t.Run("a = b\t{ a = b }\t⊥\ta and b do not match", func(t *testing.T) {
		m := map[string]*Kind{}
		a := Num
		b := Str
		u := Unify(a, b, m)
		t.Log(m, u)
		assert(u == nil)
		assert(len(m) == 0)
	})

	t.Run("X = X\t{ x = x }\t{}", func(t *testing.T) {
		m := map[string]*Kind{}
		X := Slot("X")
		u := Unify(X, X, m)
		t.Log(m, u)
		assert(Equals(u, X))
		assert(len(m) == 0)
	})

	t.Run("a = X\t{ a = x }\t{ x ↦ a }\tx is unified with the constant a", func(t *testing.T) {
		m := map[string]*Kind{}
		a := Num
		X := Slot("X")
		u := Unify(a, X, m)
		t.Log(m, u)
		assert(Equals(u, a))
		assert(len(m) == 1)
		assert(Equals(m[X.Slot().Name], a))
	})

	t.Run("X = Y\t{ x = y }\t{ x ↦ y }", func(t *testing.T) {
		m := map[string]*Kind{}
		X := Slot("X")
		Y := Slot("Y")
		u := Unify(X, Y, m)
		t.Log(m, u)
		assert(len(m) == 1)
		assert(Equals(m[X.Slot().Name], Y))
	})

	t.Run("f(a,X) = f(a,b)\t{ f(a,x) = f(a,b) }\t{ x ↦ b }\tfunction and constant symbols match, x is unified with the constant b", func(t *testing.T) {
		m := map[string]*Kind{}
		a := Num
		b := Str
		X := Slot("X")
		fax := arrow(a, X)
		fab := arrow(a, b)
		u := Unify(fax, fab, m)
		t.Log(m, u)
		assert(Equals(u, fab))
		assert(len(m) == 1)
		assert(Equals(m[X.Slot().Name], b))
	})

	t.Run("f(X) = f(Y)\t{ f(x) = f(y) }\t{ x ↦ y }\tx and y are aliased", func(t *testing.T) {
		m := map[string]*Kind{}
		X := Slot("X")
		Y := Slot("Y")
		fx := arrow(X, X)
		fy := arrow(Y, Y)
		u := Unify(fx, fy, m)

		assert(len(m) == 1)
		assert(Equals(m[X.Slot().Name], Y))

		fx1 := applySubst(fx, m)
		fy1 := applySubst(fy, m)
		t.Log(m, u, fx1, fy1)
		assert(Equals(fx1, fy))
		assert(Equals(fy, fy))
	})
}

func TestInfer1(tt *testing.T) {
	// id :: a -> a
	a := Slot("a")
	fun := arrow(a, a)

	// psuido :: s -> t
	s := Slot("s")
	t := Slot("t")
	psuido := arrow(s, t)

	m := map[string]*Kind{}
	tfn1 := Unify(psuido, fun, m)

	tt.Log(tfn1, m)
	assert(len(m) == 2)
	assert(Equals(m[s.Slot().Name], a))
	assert(Equals(m[t.Slot().Name], a))

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

	Unify(lt, lt, map[string]*Kind{})
}

func assert(cond bool) {
	if !cond {
		panic("")
	}
}
