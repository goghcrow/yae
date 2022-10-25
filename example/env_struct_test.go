package example

import (
	"github.com/goghcrow/yae"
	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
	"time"
)

func TestStructEnv(t *testing.T) {
	type Entity struct {
		Id   int    `yae:"ID"`
		Name string `yae:"姓名"`
	}
	type Ctx struct {
		Ok  bool      `yae:"布尔"`
		N   int       `yae:"数字"`
		T   time.Time `yae:"时间"`
		Lst []*Entity `yae:"列表"`
		Obj *Entity   `yae:"对象"`
	}

	typeEnv := conv.MustTypeEnvOf(Ctx{
		Lst: []*Entity{},
		Obj: &Entity{},
	})
	expr := yae.NewExpr().EnableDebug(os.Stderr)
	closure, err := expr.Compile("if(布尔, 列表[0].姓名.len() + 数字, 0)", typeEnv)
	if err != nil {
		panic(err)
	}

	{
		valEnv := conv.MustValEnvOf(&Ctx{
			Ok: true,
			N:  42,
			T:  time.Now(),
			Lst: []*Entity{
				{Id: 42, Name: "晓"},
			},
			Obj: &Entity{Id: 42, Name: "晓"},
		})
		v, err := closure(valEnv)
		if err != nil {
			panic(err)
		}
		if !val.Equals(v, val.Num(43)) {
			t.Errorf("expect 43 actual %s", v)
		}
	}
	{
		valEnv := conv.MustValEnvOf(&Ctx{
			Ok: true,
			N:  100,
			T:  time.Now(),
			Lst: []*Entity{
				{Id: 42, Name: "晓"},
			},
			Obj: &Entity{Id: 42, Name: "晓"},
		})
		v, err := closure(valEnv)
		if err != nil {
			panic(err)
		}
		if !val.Equals(v, val.Num(101)) {
			t.Errorf("expect 101 actual %s", v)
		}
	}
}
