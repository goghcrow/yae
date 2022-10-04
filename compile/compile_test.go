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
	typEnv.Put("var_str", types.Str)
	typEnv.Put("var_num", types.Num)
	typEnv.Put("var_list_num", types.List(types.Num))

	valEnv := val.NewEnv()
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
		name string
		expr *Expr
		kind *types.Kind
		val  *val.Val
	}{
		{
			name: "str",
			expr: LitStr(`"Hello World!"`),
			kind: types.Str,
			val:  val.Str("Hello World!"),
		},
		{
			name: "str/raw",
			expr: LitStr("`raw str`"),
			kind: types.Str,
			val:  val.Str("raw str"),
		},
		{
			name: "str/escape",
			expr: LitStr(`"abc晓\r\n\t"`),
			kind: types.Str,
			val:  val.Str("abc晓\r\n\t"),
		},
		{
			name: "num/zero",
			expr: LitNum("0"),
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			name: "num/int",
			expr: LitNum("42"),
			kind: types.Num,
			val:  num,
		},
		{
			name: "num/float",
			expr: LitNum("123.456"),
			kind: types.Num,
			val:  val.Num(123.456),
		},
		{
			name: "num/e",
			expr: LitNum("123.456e-78"),
			kind: types.Num,
			val:  val.Num(123.456e-78),
		},
		{
			name: "num/neg",
			expr: LitNum("-123.456E-78"),
			kind: types.Num,
			val:  val.Num(-123.456e-78),
		},
		{
			name: "num/bin",
			expr: LitNum("0b10101"),
			kind: types.Num,
			val:  val.Num(0b10101),
		},
		{
			name: "num/oct",
			expr: LitNum("0o1234567"),
			kind: types.Num,
			val:  val.Num(0o1234567),
		},
		{
			name: "num/hex",
			expr: LitNum("0x123456789abcdef"),
			kind: types.Num,
			val:  val.Num(0x123456789abcdef),
		},
		{
			name: "bool/true",
			expr: LitTrue(),
			kind: types.Bool,
			val:  val.True,
		},
		{
			name: "bool/false",
			expr: LitFalse(),
			kind: types.Bool,
			val:  val.False,
		},
		{
			name: "var/num",
			expr: Ident("var_num"),
			kind: types.Num,
			val:  num,
		},
		{
			name: "var/str",
			expr: Ident("var_str"),
			kind: types.Str,
			val:  str,
		},
		{
			name: "var/list[num]",
			expr: Ident("var_list_num"),
			kind: types.List(types.Num),
			val:  list.Vl(),
		},
		{
			name: "list[num]",
			expr: List([]*Expr{LitNum("1"), LitNum("2")}),
			kind: types.List(types.Num),
			val:  val.List(types.List(types.Num).List(), 0).List().Add(val.Num(1), val.Num(2)).Vl(),
		},
		{
			name: "map[num, str]",
			expr: Map([]Pair{{Key: LitNum("1"), Val: LitStr("`1`")}}),
			kind: types.Map(types.Num, types.Str),
			val:  val.Map(types.Map(types.Num, types.Str).Map()).Map().Put(val.Num(1), val.Str("1")).Vl(),
		},
		{
			name: "obj{id:num, name:str}",
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
			name: "subscript/list",
			expr: Subscript(
				List([]*Expr{LitNum("1"), LitNum("2")}),
				LitNum("0"),
			),
			kind: types.Num,
			val:  val.Num(1),
		},
		{
			name: "subscript/map",
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
			name: "member",
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
			name: "call/mono",
			expr: Call(Ident("mono_itoa"), []*Expr{LitNum("42")}),
			kind: types.Str,
			val:  val.Str("42"),
		},
		{
			name: "call/poly/id/num",
			expr: Call(Ident("poly_id"), []*Expr{LitNum("42")}),
			kind: types.Num,
			val:  val.Num(42),
		},
		{
			name: "call/poly/id/str",
			expr: Call(Ident("poly_id"), []*Expr{LitStr(`"晓"`)}),
			kind: types.Str,
			val:  val.Str("晓"),
		},
		{
			name: "call/poly/len/emptyList",
			expr: Call(Ident("poly_len"), []*Expr{
				List([]*Expr{}),
			}),
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			name: "call/poly/len/list",
			expr: Call(Ident("poly_len"), []*Expr{
				List([]*Expr{LitNum("1"), LitNum("2")}),
			}),
			kind: types.Num,
			val:  val.Num(2),
		},
		{
			name: "call/poly/len/emptyMap",
			expr: Call(Ident("poly_len"), []*Expr{
				Map([]Pair{}),
			}),
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			name: "call/poly/len/map",
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
			name: "call/dynamic",
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
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
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
