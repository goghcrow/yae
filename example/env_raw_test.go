package example

import (
	expr "github.com/goghcrow/yae"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
	"time"
)

func TestRawEnv(t *testing.T) {
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

	entity := types.Obj(map[string]*types.Kind{
		"ID": types.Num,
		"姓名": types.Str,
	})
	entityLst := types.List(entity)

	typeEnv := types.NewEnv()
	typeEnv.Put("布尔", types.Bool)
	typeEnv.Put("数字", types.Num)
	typeEnv.Put("时间", types.Time)
	typeEnv.Put("列表", entityLst)
	typeEnv.Put("对象", entity)

	expr := expr.NewExpr().EnableDebug(os.Stderr)
	closure, err := expr.Compile("if(布尔, 列表[0].姓名.len() + 数字, 0)", typeEnv)
	if err != nil {
		panic(err)
	}

	{
		obj := val.Obj(entity.Obj()).Obj()
		obj.V["ID"] = val.Num(42)
		obj.V["姓名"] = val.Str("晓")
		lst := val.List(entityLst.List(), 0).List()
		lst.Add(obj.Vl())

		valEnv := val.NewEnv()
		valEnv.Put("布尔", val.True)
		valEnv.Put("数字", val.Num(42))
		valEnv.Put("时间", val.Time(time.Now()))
		valEnv.Put("列表", lst.Vl())
		valEnv.Put("对象", obj.Vl())

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
		obj := val.Obj(entity.Obj()).Obj()
		obj.V["ID"] = val.Num(42)
		obj.V["姓名"] = val.Str("晓")
		lst := val.List(entityLst.List(), 0).List()
		lst.Add(obj.Vl())

		valEnv := val.NewEnv()
		valEnv.Put("布尔", val.True)
		valEnv.Put("数字", val.Num(100))
		valEnv.Put("时间", val.Time(time.Now()))
		valEnv.Put("列表", lst.Vl())
		valEnv.Put("对象", obj.Vl())

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
