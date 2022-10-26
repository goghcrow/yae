package test

import (
	"github.com/goghcrow/yae/closure"
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"github.com/goghcrow/yae/vm"
	"testing"
)

func TestMaybe(t *testing.T) {
	t.Run("tag-maybe", func(t *testing.T) {
		type s struct {
			Name string `yae:"name,maybe"`
			Age  int    `yae:",maybe"`
		}
		ty, err := conv.TypeOf(s{})
		if err != nil {
			panic(err)
		}
		assert(ty.Kind == types.KObj)

		assert(types.Equals(
			// 这里别名 name 小写开头
			ty.Obj().MustGetField("name").Val,
			types.Maybe(types.Str),
		))
		assert(types.Equals(
			// 这里没有别名 Age 大写开头
			ty.Obj().MustGetField("Age").Val,
			types.Maybe(types.Num),
		))
	})

	t.Run("nil-maybe", func(t *testing.T) {
		type s struct {
			Name    string
			Age     int
			NamePtr *string // default nil
			AgePtr  *int    // default nil
		}
		ty, err := conv.TypeOf(s{ /**/ })
		if err != nil {
			panic(err)
		}
		assert(ty.Kind == types.KObj)
		assert(types.Equals(
			ty.Obj().MustGetField("Name").Val,
			types.Str,
		))
		assert(types.Equals(
			ty.Obj().MustGetField("Age").Val,
			types.Num,
		))
		assert(types.Equals(
			ty.Obj().MustGetField("NamePtr").Val,
			types.Maybe(types.Str),
		))
		assert(types.Equals(
			ty.Obj().MustGetField("AgePtr").Val,
			types.Maybe(types.Num),
		))
	})

	// 把允许为空值或者为 nil 的变量或者字段声明称 maybe 类型, 通过 get + 默认值 获取值

	makeEnv := func() (typedEnv *types.Env, compileEvalEnv *val.Env) {
		maybeInt := types.Maybe(types.Num).Maybe()
		nothing := val.Nothing(types.Num)
		just := val.Just(types.Num, val.Num(42))

		objWithMaybeField := types.Obj([]types.Field{
			{"nothing", maybeInt.Ty()},
			{"just", maybeInt.Ty()},
		}).Obj()

		tyEnv := types.NewEnv()
		vlEnv := val.NewEnv()
		tyEnv.Put("nothing", maybeInt.Ty())
		tyEnv.Put("just", maybeInt.Ty())
		tyEnv.Put("obj", objWithMaybeField.Ty())

		vlEnv.Put("nothing", nothing)
		vlEnv.Put("just", just)

		obj := val.Obj(objWithMaybeField).Obj()
		obj.Put("nothing", nothing)
		obj.Put("just", just)
		vlEnv.Put("obj", obj.Vl())

		return tyEnv, vlEnv
	}

	tests := []struct {
		expr     string
		expected *val.Val
	}{
		{"get(nothing, 100)", val.Num(100)},
		{"get(just, 100)", val.Num(42)},
		{"get(obj.nothing, 100)", val.Num(100)},
		{"get(obj.just, 100)", val.Num(42)},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			tyEnv, vlEnv := makeEnv()
			actual := eval(tt.expr, closure.Compile, tyEnv, vlEnv)
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})

		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			tyEnv, vlEnv := makeEnv()
			actual := eval(tt.expr, vm.Compile, tyEnv, vlEnv)
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})
	}
}
