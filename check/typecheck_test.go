package check

import (
	"github.com/goghcrow/yae/ast"
	"github.com/goghcrow/yae/env0"
	types "github.com/goghcrow/yae/type"
	"testing"
)

func TestTypeCheck(t *testing.T) {
	e0 := env0.NewEnv()

	for _, tt := range []struct {
		expr *ast.Expr
		kind *types.Kind
	}{
		{
			expr: ast.Literal(ast.LIT_STR, "s"),
			kind: types.Str,
		},
		{
			expr: ast.Literal(ast.LIT_STR, ``),
			kind: types.Str,
		},
	} {
		actual := TypeCheck(e0, tt.expr)
		expected := tt.kind
		if !types.Equals(expected, actual) {
			t.Errorf("expect %s actual %s in %s", expected, actual, tt.expr)
		}
	}
}
