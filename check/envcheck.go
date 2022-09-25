package check

import (
	"github.com/goghcrow/yae/env"
	"github.com/goghcrow/yae/env0"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
)

// EnvCheck 检查闭包执行的上下文, env0: 编译期环境 env1: 执行期环境
func EnvCheck(env0 *env0.Env, env1 *env.Env) {
	// 遍历 0 检查 1 是否有以及类型
	env0.ForEach(func(name string, kind *types.Kind) {
		v, ok := env1.Get(name)
		util.Assert(ok, "undefined %s", name)
		util.Assert(types.Equals(kind, v.Kind), "expect %s get %s", kind, v.Kind)
	})
}
