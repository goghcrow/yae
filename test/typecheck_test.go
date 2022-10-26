package test

import (
	"github.com/goghcrow/yae/fun"
	. "github.com/goghcrow/yae/types"
	"testing"
)

func TestInfer(t *testing.T) {
	env := NewEnv()

	// id :: forall a. (a -> a)
	env.RegisterFun(func() *Type {
		a := TyVar("a")
		return Fun("id", []*Type{a}, a)
	}())

	// list :: forall a. (a -> list[a])
	env.RegisterFun(func() *Type {
		a := TyVar("a")
		return Fun("list", []*Type{a}, List(a))
	}())

	// if :: forall a. (bool -> α -> α -> α)
	env.RegisterFun(func() *Type {
		a := TyVar("a")
		return Fun("if", []*Type{Bool, a, a}, a)
	}())

	// has :: forall k v. (map[k, v] -> k -> bool)
	env.RegisterFun(func() *Type {
		k := TyVar("k")
		v := TyVar("v")
		return Fun("has", []*Type{Map(k, v), k}, Bool)
	}())

	// get :: forall k v. (map[k, v] -> k -> v)
	env.RegisterFun(func() *Type {
		k := TyVar("k")
		v := TyVar("v")
		return Fun("get", []*Type{Map(k, v), k}, v)
	}())

	env.RegisterFun(fun.STRTOTIME_STR.Fun().Type)

	tests := []struct {
		s   string
		t   *Type
		err bool
	}{
		{"id(1)", Num, false},
		{"id(true)", Bool, false},
		{`id("s")`, Str, false},
		{"id('now')", Time, false},

		{"id([])", List(Bottom), false},
		{"id([:])", Map(Bottom, Bottom), false},
		{"id({})", Obj([]Field{}), false},

		{"id([1])", List(Num), false},
		{`id([1:"s"])`, Map(Num, Str), false},
		{`id({id:1})`, Obj([]Field{
			{"id", Num},
		}), false},

		{`id([{a:[1:1]}])`, List(Obj([]Field{
			{"a", Map(Num, Num)},
		})), false},

		{`id({id:1, lst:[1], map:[1:"s"], obj:{id:1}})`, Obj([]Field{
			{"id", Num},
			{"lst", List(Num)},
			{"map", Map(Num, Str)},
			{"obj", Obj([]Field{
				{"id", Num},
			})},
		}), false},

		{"list(1)", List(Num), false},
		{"list(true)", List(Bool), false},
		{`list("s")`, List(Str), false},
		{"list('now')", List(Time), false},

		{"list([])", List(List(Bottom)), false},
		{"list([:])", List(Map(Bottom, Bottom)), false},
		{"list({})", List(Obj([]Field{})), false},

		{"list([1])", List(List(Num)), false},
		{`list([1:"s"])`, List(Map(Num, Str)), false},
		{`list({id:1})`, List(Obj([]Field{
			{"id", Num},
		})), false},

		{`list([{a:[1:1]}])`, List(List(Obj([]Field{
			{"a", Map(Num, Num)},
		}))), false},

		{`list({id:1, lst:[1], map:[1:"s"], obj:{id:1}})`, List(Obj([]Field{
			{"id", Num},
			{"lst", List(Num)},
			{"map", Map(Num, Str)},
			{"obj", Obj([]Field{
				{"id", Num},
			})},
		})), false},

		{`has(["a":42], 1)`, nil, true},
		{`has(["a":42], "")`, Bool, false},
		{`get([42:"a"], "")`, Str, true},
		{`get(["a":42], "")`, Num, false},
		{`if( has(["a":42], "a"), get(["a":42], "a"),  "1")`, Num, true},
		{`if( has(["a":42], "a"), get(["a":42], "a"),  1)`, Num, false},
		{`if( has([1:"42"], 1), get([1:"42"], 1),  "")`, Str, false},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			actual, err := Infer(parse(tt.s), env)
			if tt.err {
				if err == nil {
					t.Errorf("expect err actual %s", actual)
				}
			} else {
				if !Equals(actual, tt.t) {
					t.Errorf("expect %s actual %s", tt.t, actual)
				}
			}
		})
	}
}
