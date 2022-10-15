package example

import (
	"github.com/goghcrow/yae"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"testing"
)

func BenchmarkVM(b *testing.B) {
	expr := yae.NewExpr().UseBytecodeCompiler()

	// typeEnv:=struct {N int `yae:"n"`}{}
	typeEnv := types.NewEnv()
	typeEnv.Put("n", types.Num)

	closure, err := expr.Compile("if(false, 1, if(true, 2+3/100, 4+2))+n", typeEnv)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// valEnv := map[string]interface{}{"n": 1,}

		valEnv := val.NewEnv()
		valEnv.Put("n", val.Num(1))

		_, err = closure(valEnv)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkClosure(b *testing.B) {
	expr := yae.NewExpr().UseClosureCompiler()
	closure, err := expr.Compile("if(false, 1, if(true, 2+3/100, 4+2))+n", struct {
		//Lst []string `yae:"lst"`
		N int `yae:"n"`
	}{})
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = closure(map[string]interface{}{
			//"lst": []string{"hello", "world"},
			"n": 1,
		})
		if err != nil {
			panic(err)
		}
	}
}