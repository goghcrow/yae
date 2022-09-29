package compile

import (
	. "github.com/goghcrow/yae/ast"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"strconv"
	"testing"
)

func TestCompile(t *testing.T) {
	str := val.Str("Hello")
	num := val.Num(42)
	list := val.List(types.List(types.Num).List(), 0).List().
		Add(val.Num(1), val.Num(2), val.Num(3))
	obj := val.Obj(types.Obj(map[string]*types.Kind{
		"id":   types.Num,
		"name": types.Str,
	}).Obj()).Obj()
	obj.V["id"] = val.Num(42)
	obj.V["name"] = val.Str("晓")

	typEnv := types.NewEnv()
	valEnv := val.NewEnv()

	typEnv.Put("var_str", types.Str)
	typEnv.Put("var_num", types.Num)
	typEnv.Put("var_list_num", types.List(types.Num))

	valEnv.Put("var_str", str)
	valEnv.Put("var_num", num)
	valEnv.Put("var_list_num", list.Vl())

	{
		// num -> str
		fun := types.Fun("mono_itoa", []*types.Kind{types.Num}, types.Str)
		typEnv.RegisterFun(fun)
		valEnv.RegisterFun(val.Fun(fun, func(v ...*val.Val) *val.Val {
			return val.Str(strconv.Itoa(int(v[0].Num().V)))
		}))

		// 测试动态分派
		lookup, _ := fun.Fun().Lookup()
		monoFunKind, _ := typEnv.Get(lookup)
		typEnv.Put("mono_itoa", monoFunKind)

		monoFunVal, _ := valEnv.Get(lookup)
		valEnv.Put("mono_itoa", monoFunVal)
	}

	{
		// id :: forall a.(a -> a)
		a := types.Slot("a")
		fk := types.Fun("poly_id", []*types.Kind{a}, a)
		typEnv.RegisterFun(fk)
		valEnv.RegisterFun(val.Fun(fk, func(v ...*val.Val) *val.Val {
			return v[0]
		}))
	}

	{
		// len :: forall a.(list a -> num)
		a := types.Slot("a")
		fk := types.Fun("poly_len", []*types.Kind{types.List(a)}, types.Num)
		typEnv.RegisterFun(fk)
		valEnv.RegisterFun(val.Fun(fk, func(v ...*val.Val) *val.Val {
			return val.Num(float64(len(v[0].List().V)))
		}))
	}

	{
		// len :: forall a b.(map a b -> num)
		k := types.Slot("k")
		v := types.Slot("v")
		fk := types.Fun("poly_len", []*types.Kind{types.Map(k, v)}, types.Num)
		typEnv.RegisterFun(fk)
		valEnv.RegisterFun(val.Fun(fk, func(v ...*val.Val) *val.Val {
			return val.Num(float64(len(v[0].Map().V)))
		}))
	}

	tests := []struct {
		expr *Expr
		kind *types.Kind
		val  *val.Val
	}{
		{
			expr: LitStr("`raw str`"),
			kind: types.Str,
			val:  val.Str("raw str"),
		},
		{
			expr: LitStr(`"abc晓\r\n\t"`),
			kind: types.Str,
			val:  val.Str("abc晓\r\n\t"),
		},
		{
			expr: LitNum("0"),
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			expr: LitNum("42"),
			kind: types.Num,
			val:  num,
		},
		{
			expr: LitNum("123.456"),
			kind: types.Num,
			val:  val.Num(123.456),
		},
		{
			expr: LitNum("123.456e-78"),
			kind: types.Num,
			val:  val.Num(123.456e-78),
		},
		{
			expr: LitNum("-123.456E-78"),
			kind: types.Num,
			val:  val.Num(-123.456e-78),
		},
		{
			expr: LitNum("0b10101"),
			kind: types.Num,
			val:  val.Num(0b10101),
		},
		{
			expr: LitNum("0o1234567"),
			kind: types.Num,
			val:  val.Num(0o1234567),
		},
		{
			expr: LitNum("0x123456789abcdef"),
			kind: types.Num,
			val:  val.Num(0x123456789abcdef),
		},
		{
			expr: LitTrue(),
			kind: types.Bool,
			val:  val.True,
		},
		{
			expr: LitFalse(),
			kind: types.Bool,
			val:  val.False,
		},
		{
			expr: Ident("var_num"),
			kind: types.Num,
			val:  num,
		},
		{
			expr: Ident("var_str"),
			kind: types.Str,
			val:  str,
		},
		{
			expr: Ident("var_list_num"),
			kind: types.List(types.Num),
			val:  list.Vl(),
		},

		{
			expr: List([]*Expr{LitNum("1"), LitNum("2")}),
			kind: types.List(types.Num),
			val:  val.List(types.List(types.Num).List(), 0).List().Add(val.Num(1), val.Num(2)).Vl(),
		},
		{
			expr: Map([]Pair{{Key: LitNum("1"), Val: LitStr("`1`")}}),
			kind: types.Map(types.Num, types.Str),
			val:  val.Map(types.Map(types.Num, types.Str).Map()).Map().Put(val.Num(1), val.Str("1")).Vl(),
		},
		{
			expr: Obj(map[string]*Expr{
				"id":   LitNum("42"),
				"name": LitStr(`"晓"`),
			}),
			kind: types.Obj(map[string]*types.Kind{
				"id":   types.Num,
				"name": types.Str,
			}),
			val: obj.Vl(),
		},
		{
			expr: Subscript(
				List([]*Expr{LitNum("1"), LitNum("2")}),
				LitNum("0"),
			),
			kind: types.Num,
			val:  val.Num(1),
		},
		{
			expr: Subscript(
				Map([]Pair{
					{Key: LitNum("1"), Val: LitStr("`1`")},
				}),
				LitNum("1"),
			),
			kind: types.Str,
			val:  val.Str("1"),
		},
		{
			expr: Member(
				Obj(map[string]*Expr{
					"id":   LitNum("42"),
					"name": LitStr(`"晓"`),
				}),
				Ident("id").Ident(),
			),
			kind: types.Num,
			val:  val.Num(42),
		},
		{
			expr: Call(Ident("mono_itoa"), []*Expr{LitNum("42")}),
			kind: types.Str,
			val:  val.Str("42"),
		},
		{
			expr: Call(Ident("poly_id"), []*Expr{LitNum("42")}),
			kind: types.Num,
			val:  val.Num(42),
		},
		{
			expr: Call(Ident("poly_id"), []*Expr{LitStr(`"晓"`)}),
			kind: types.Str,
			val:  val.Str("晓"),
		},
		{
			expr: Call(Ident("poly_len"), []*Expr{
				List([]*Expr{}),
			}),
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			expr: Call(Ident("poly_len"), []*Expr{
				List([]*Expr{LitNum("1"), LitNum("2")}),
			}),
			kind: types.Num,
			val:  val.Num(2),
		},
		{
			expr: Call(Ident("poly_len"), []*Expr{
				Map([]Pair{}),
			}),
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			expr: Call(Ident("poly_len"), []*Expr{
				Map([]Pair{
					{Key: LitNum("1"), Val: LitStr("`1`")},
				}),
			}),
			kind: types.Num,
			val:  val.Num(1),
		},
		// 动态分派
		{
			// [poly_id][0](1),
			expr: Call(
				Subscript(
					List([]*Expr{
						Ident("mono_itoa"),
					}),
					LitNum("0"),
				),
				[]*Expr{
					LitNum("1"),
				},
			),
			kind: types.Str,
			val:  val.Str("1"),
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			{
				// 注意 typeCheck 会修改 ast 的上附加的类型信息
				actual := types.TypeCheck(typEnv, tt.expr)
				expected := tt.kind
				if !types.Equals(expected, actual) {
					t.Errorf("expect `%s` actual `%s` in `%s`", expected, actual, tt.expr)
				}
			}
			{
				closure := Compile(valEnv, tt.expr)
				actual := closure(valEnv)
				expected := tt.val
				if !val.Equals(expected, actual) {
					t.Errorf("expect %s actual %s in `%s`", expected, actual, tt.expr)
				}
			}
		})
	}
}
