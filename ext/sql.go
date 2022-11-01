package ext

import (
	"fmt"

	"github.com/goghcrow/yae"
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/ext/sql"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

func CompileToSql(be BoolExpr, env1 *types.Env) func(v interface{}) (string, error) {
	expr := yae.NewExpr().
		UseBuiltIn(false).
		RegisterOperator(oper.BuiltIn()...).
		UseCompiler(sql.Compile).
		RegisterFun(sql.BuiltIn()...)

	closure := expr.CompileExpr(be.expr(), env1)

	return func(v interface{}) (s string, err error) {
		env, ok := v.(*val.Env)
		if !ok {
			env, err = conv.ValEnvOf(v)
			if err != nil {
				return "", err
			}
		}
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()
		env.ForEach(func(name string, v *val.Val) {
			ty, ok := env1.Get(name)
			util.Assert(ok, "undefined %s", name)
			util.Assert(types.Equals(ty, v.Type),
				"type mismatched, expect `%s` actual `%s` %s", ty, v.Type, v)
		})
		return closure(env).Str().V, nil
	}
}
