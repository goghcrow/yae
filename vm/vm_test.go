package vm

import (
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/trans"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"testing"
)

var typecheckEnv = types.NewEnv()
var compileEnv = val.NewEnv()

func init() {
	initEnv(typecheckEnv, compileEnv)
}
func initEnv(typecheckEnv *types.Env, compileEnv *val.Env) {
	for _, f := range fun.BuildIn() {
		typecheckEnv.RegisterFun(f.Kind)
		compileEnv.RegisterFun(f)
	}
}

func TestX(t *testing.T) {
	//input := `{id:42,name:"晓", map:["lst":[1,2]]}.map["lst"][1] + 100`
	input := `if(false,1,{id:42,name:"晓", map:["lst":[1,2]]}.map["lst"][1] + 100)`
	//input := "if(false,1,2)"

	ops := oper.BuildIn()
	toks := lexer.NewLexer(ops).Lex(input)
	expr := parser.NewParser(ops).Parse(toks)
	expr = trans.Desugar(expr)

	//dot := ast.Dot(expr, "")
	//cmd := exec.Command("open", "https://dreampuf.github.io/GraphvizOnline/#"+url.PathEscape(dot))
	//_, _ = cmd.Output()

	typeEnv, _ := conv.TypeEnvOf(struct {
		Lst []string `yae:"lst"`
		N   int      `yae:"n"`
	}{})
	_ = types.Check(expr, typeEnv.Inherit(typecheckEnv))

	bytecode := NewCompile().Compile(expr, compileEnv)
	t.Log(bytecode)

	runtimeEnv, _ := conv.ValEnvOf(map[string]interface{}{
		"lst": []string{"hello", "world"},
		"n":   1,
	})
	runtimeEnv = runtimeEnv.Inherit(compileEnv)

	r := NewVM().Interp(bytecode, runtimeEnv)
	t.Log(r)

	r = NewVM().Interp(bytecode, runtimeEnv)
	t.Log(r)
}
