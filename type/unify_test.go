package types

import (
	"testing"
)

func TestUnify(t *testing.T) {
	// A :: ((α → β) × [γ]) → [β]
	// B :: ((γ → δ) × [γ]) → ε
	β := Slot("β")
	A := Fun("",
		[]*Kind{
			Obj([]Field{
				{"a", Fun("", []*Kind{Slot("α")}, β)},
				{"b", List(Slot("γ"))},
			}),
		},
		List(β),
	)
	γ := Slot("γ")
	B := Fun("",
		[]*Kind{
			Obj([]Field{
				{"a", Fun("", []*Kind{γ}, Slot("δ"))},
				{"b", List(γ)},
			}),
		},
		Slot("ε"),
	)
	m := map[string]*Kind{}
	t.Log(Unify(A, B, m))
	t.Log(m)
}

func TestUnify1(tt *testing.T) {
	a := Slot("a")
	fun := Fun("", []*Kind{a}, a)
	targ := Slot("str")

	s := Slot("s")
	t := Slot("t")
	psuido := Fun("", []*Kind{s}, t)

	m := map[string]*Kind{}
	tfn1 := Unify(psuido, fun, m)
	tt.Log(tfn1)

	// s 换成 a
	tt.Log("_________")
	tt.Log(s)
	targ1 := applySubst(s, m)
	tt.Log(targ1)
	tt.Log("_________")
	// a 换成 str
	targ2 := Unify(targ1, targ, m) // 带入
	tt.Log(targ2)
	tt.Log("_________")

	tresult := applySubst(t, m) // m 已经有返回值
	tt.Log(tresult)
}

func TestRecursive(t *testing.T) {
	lt := List(nil)
	lt.List().El = lt
	t.Log(lt)

	Unify(lt, lt, map[string]*Kind{})
}
