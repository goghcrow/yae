package test

import (
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/trans"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/vm"
	"testing"
)

func TestVM(t *testing.T) {
	//input := `{id:42,name:"晓", map:["lst":[1,2]]}.map["lst"][1] + 100`
	//input := `if(false,1,{id:42,name:"晓", map:["lst":[1,2]]}.map["lst"][1] + 100)`
	//input := "1 + 1 > 1 && 1 < 1 || !false"
	input := "if(false, 1, if(true, 2+3, 4+2))+n"
	//input := "if(false, 1, 2)"
	//input := "(1 + 2) ^ (3 % 4) * 5 - 42 / 100"

	ops := oper.BuildIn()
	toks := lexer.NewLexer(ops).Lex(input)
	expr := parser.NewParser(ops).Parse(toks)
	expr = trans.Desugar(expr)

	//dot := ast.Dot(expr, "")
	//cmd := exec.Command("open", "https://dreampuf.github.io/GraphvizOnline/#"+url.PathEscape(dot))
	//_, _ = cmd.Output()

	typeEnv := conv.MustTypeEnvOf(struct {
		Lst []string `yae:"lst"`
		N   int      `yae:"n"`
	}{}).Inherit(typecheckEnv)
	_ = types.Check(expr, typeEnv)

	bytecode := vm.NewCompile().Compile(expr, compileEnv)
	t.Log(bytecode)

	runtimeEnv := conv.MustValEnvOf(map[string]interface{}{
		"lst": []string{"hello", "world"},
		"n":   1,
	}).Inherit(compileEnv)

	r := vm.NewVM().Interp(bytecode, runtimeEnv)
	t.Log(r)

	r = vm.NewVM().Interp(bytecode, runtimeEnv)
	t.Log(r)
}
