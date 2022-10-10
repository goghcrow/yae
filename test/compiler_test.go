package test

import (
	"github.com/goghcrow/yae/closure"
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"github.com/goghcrow/yae/vm"
	"strconv"
	"testing"
)

func TestCompile(t *testing.T) {
	obj, _ := conv.ValOf(struct {
		Id   int    `yae:"id"`
		Name string `yae:"name"`
	}{42, "晓"})

	typEnv := conv.MustTypeEnvOf(struct {
		Str  string `yae:"var_str"`
		Num  int    `yae:"var_num"`
		List []int  `yae:"var_list_num"`
	}{}).Inherit(typecheckEnv)

	valEnv := conv.MustValEnvOf(map[string]interface{}{
		"var_str":      "Hello",
		"var_num":      42,
		"var_list_num": []int{1, 2, 3},
	}).Inherit(compileEnv)

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
		expr string
		kind *types.Kind
		val  *val.Val
	}{
		{
			name: "str",
			expr: `"Hello World!"`,
			kind: types.Str,
			val:  val.Str("Hello World!"),
		},
		{
			name: "str/raw",
			expr: "`raw str`",
			kind: types.Str,
			val:  val.Str("raw str"),
		},
		{
			name: "str/escape",
			expr: `"abc晓\r\n\t"`,
			kind: types.Str,
			val:  val.Str("abc晓\r\n\t"),
		},
		{
			name: "num/zero",
			expr: "0",
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			name: "num/int",
			expr: "42",
			kind: types.Num,
			val:  conv.MustValOf(42),
		},
		{
			name: "num/float",
			expr: "123.456",
			kind: types.Num,
			val:  val.Num(123.456),
		},
		{
			name: "num/e",
			expr: "123.456e-78",
			kind: types.Num,
			val:  val.Num(123.456e-78),
		},
		{
			name: "num/neg",
			expr: "-123.456E-78",
			kind: types.Num,
			val:  val.Num(-123.456e-78),
		},
		{
			name: "num/bin",
			expr: "0b10101",
			kind: types.Num,
			val:  val.Num(0b10101),
		},
		{
			name: "num/oct",
			expr: "0o1234567",
			kind: types.Num,
			val:  val.Num(0o1234567),
		},
		{
			name: "num/hex",
			expr: "0x123456789abcdef",
			kind: types.Num,
			val:  val.Num(0x123456789abcdef),
		},
		{
			name: "bool/true",
			expr: "true",
			kind: types.Bool,
			val:  val.True,
		},
		{
			name: "bool/false",
			expr: "false",
			kind: types.Bool,
			val:  val.False,
		},
		{
			name: "var/num",
			expr: "var_num",
			kind: types.Num,
			val:  conv.MustValOf(42),
		},
		{
			name: "var/str",
			expr: "var_str",
			kind: types.Str,
			val:  conv.MustValOf("Hello"),
		},
		{
			name: "var/list[num]",
			expr: "var_list_num",
			kind: types.List(types.Num),
			val:  conv.MustValOf([]int{1, 2, 3}),
		},
		{
			name: "list[num]",
			expr: "[1, 2]",
			kind: types.List(types.Num),
			val:  conv.MustValOf([]int{1, 2}),
		},
		{
			name: "map[num, str]",
			expr: "[1:`1`]",
			kind: types.Map(types.Num, types.Str),
			val:  conv.MustValOf(map[int]string{1: "1"}),
		},
		{
			name: "obj{id:num, name:str}",
			expr: "{id:42, name:`晓`}",
			kind: types.Obj([]types.Field{
				{"id", types.Num},
				{"name", types.Str},
			}),
			val: obj,
		},
		{
			name: "subscript/list",
			expr: "[1,2][0]",
			kind: types.Num,
			val:  val.Num(1),
		},
		{
			name: "subscript/map",
			expr: "[1:`1`][1]",
			kind: types.Str,
			val:  val.Str("1"),
		},
		{
			name: "member",
			expr: "{id:42,name:`晓`}.id",
			kind: types.Num,
			val:  val.Num(42),
		},
		{
			name: "call/mono",
			expr: "mono_itoa(42)",
			kind: types.Str,
			val:  val.Str("42"),
		},
		{
			name: "call/poly/id/num",
			expr: "poly_id(42)",
			kind: types.Num,
			val:  val.Num(42),
		},
		{
			name: "call/poly/id/str",
			expr: "poly_id(`晓`)",
			kind: types.Str,
			val:  val.Str("晓"),
		},
		{
			name: "call/poly/len/emptyList",
			expr: "poly_len([])",
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			name: "call/poly/len/list",
			expr: "poly_len([1,2])",
			kind: types.Num,
			val:  val.Num(2),
		},
		{
			name: "call/poly/len/emptyMap",
			expr: "poly_len([:])",
			kind: types.Num,
			val:  val.Num(0),
		},
		{
			name: "call/poly/len/map",
			expr: "poly_len([1:`1`])",
			kind: types.Num,
			val:  val.Num(1),
		},
		// 动态分派
		{
			name: "call/dynamic",
			expr: "[mono_itoa][0](1)",
			kind: types.Str,
			val:  val.Str("1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注意 typeCheck 会修改 ast 的上附加的类型信息
			parsed := parse(tt.expr)

			t.Run("typecheck", func(t *testing.T) {
				actual := types.Check(parsed, typEnv)
				expected := tt.kind
				if !types.Equals(expected, actual) {
					t.Errorf("expect `%s` actual `%s` in `%s`", expected, actual, tt.expr)
				}
			})

			t.Run("closure.Compile", func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("%v", r)
					}
				}()

				compiled := closure.Compile(parsed, valEnv)
				actual := compiled(valEnv)
				expected := tt.val
				if !val.Equals(expected, actual) {
					t.Errorf("expect %s actual %s in `%s`", expected, actual, tt.expr)
				}
			})
			t.Run("bytecode.Compile", func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("%v", r)
					}
				}()

				compiled := vm.Compile(parsed, valEnv)
				actual := compiled(valEnv)
				expected := tt.val
				if !val.Equals(expected, actual) {
					t.Errorf("expect %s actual %s in `%s`", expected, actual, tt.expr)
				}
			})

		})
	}
}
