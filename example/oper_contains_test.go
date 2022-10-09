package example

import (
	"github.com/goghcrow/yae"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
)

func TestRegisterContains(t *testing.T) {
	expr := yae.NewExpr().EnableDebug(os.Stderr)

	// 添加一个自定义操作符, 同时添加对应的函数
	// contains :: forall a.(list[a] -> a -> bool)

	expr.RegisterOperator(oper.Operator{
		Type:   "contains",
		BP:     oper.BP_TERM,
		Fixity: oper.INFIX_N,
	})
	expr.RegisterFun(func() *val.Val {
		a := types.Slot("a")
		return val.Fun(
			types.Fun("contains", []*types.Kind{types.List(a), a}, types.Bool),
			func(v ...*val.Val) *val.Val {
				for _, el := range v[0].List().V {
					if val.Equals(el, v[1]) {
						return val.True
					}
				}
				return val.False
			},
		)
	}())

	closure, err := expr.Compile(`if (lst contains 42, 142, 100)`, map[string]interface{}{
		"lst": []int{},
	})
	if err != nil {
		panic(err)
	}

	v, err := closure(map[string]interface{}{
		"lst": []int{1, 2, 3},
	})
	if err != nil {
		panic(err)
	}
	t.Log(v)

	v, err = closure(map[string]interface{}{
		"lst": []int{1, 2, 42},
	})
	if err != nil {
		panic(err)
	}
	t.Log(v)
}
