package example

import (
	expr "github.com/goghcrow/yae"
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
	"time"
)

func TestMapEnv(t *testing.T) {
	ctx := map[string]interface{}{
		"ok": false,
		"n":  0,
		"t":  time.Time{},
		"lst": []*struct {
			Id   int
			Name string
		}{},
		"obj": &struct {
			Id   int
			Name string
		}{},
	}

	typeEnv, err := conv.TypeEnvOf(ctx)
	if err != nil {
		panic(err)
	}
	expr := expr.NewExpr().EnableDebug(os.Stderr)
	closure, err := expr.Compile("if(ok, lst[0].Name.len() + n, 0)", typeEnv)
	if err != nil {
		panic(err)
	}

	{
		valEnv, err := conv.ValEnvOf(map[string]interface{}{
			"ok": true,
			"n":  42,
			"t":  time.Now(),
			"lst": []*struct {
				Id   int
				Name string
			}{
				{
					Id:   100,
					Name: "晓",
				},
			},
			"obj": &struct {
				Id   int
				Name string
			}{
				Id:   100,
				Name: "晓",
			},
		})
		if err != nil {
			panic(err)
		}
		v, err := closure(valEnv)
		if err != nil {
			panic(err)
		}
		if !val.Equals(v, val.Num(43)) {
			t.Errorf("expect 43 actual %s", v)
		}
	}

	{
		valEnv, err := conv.ValEnvOf(map[string]interface{}{
			"ok": true,
			"n":  100,
			"t":  time.Now(),
			"lst": []*struct {
				Id   int
				Name string
			}{
				{
					Id:   42,
					Name: "晓",
				},
			},
			"obj": &struct {
				Id   int
				Name string
			}{
				Id:   42,
				Name: "晓",
			},
		})
		if err != nil {
			panic(err)
		}
		v, err := closure(valEnv)
		if err != nil {
			panic(err)
		}
		if !val.Equals(v, val.Num(101)) {
			t.Errorf("expect 101 actual %s", v)
		}
	}
}