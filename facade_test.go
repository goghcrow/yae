package expr

import (
	"fmt"
	"github.com/goghcrow/yae/conv"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
	"time"
)

func lit(expr string) *val.Val {
	r, err := Eval(expr, struct{}{})
	if err != nil {
		panic(err)
	}
	return r
}

func TestEval(t *testing.T) {
	ref := func(v interface{}) interface{} { return &v }
	tests := []struct {
		name     string
		expr     string
		ctx      interface{}
		expected *val.Val
	}{
		// '时间字面量' 支持的语法参见 https://github.com/goghcrow/strtotime
		// "字符串"
		{
			name:     "hhmm",
			expr:     `"距离今天零点已经过去 " + string( ('today 08:00' - 'today') / 60 ) + " 分钟"`,
			ctx:      nil,
			expected: val.Str(fmt.Sprintf("距离今天零点已经过去 %d 分钟", 60*8)),
		},
		{
			name:     "hhmm",
			expr:     `"距离今天零点已经过去 " + string( (strtotime("today" + hhmm) - 'today') / 60 ) + " 分钟"`,
			ctx:      map[string]interface{}{"hhmm": "09:00"},
			expected: val.Str(fmt.Sprintf("距离今天零点已经过去 %d 分钟", 60*9)),
		},
		{
			name:     "strtotime",
			expr:     `'first day of next month'`,
			ctx:      nil,
			expected: val.Time(time.Unix(util.Strtotime("first day of next month"), 0)),
		},
		{
			name:     "map",
			expr:     "终极答案 * 100",
			ctx:      map[string]int{"终极答案": 42},
			expected: val.Num(4200),
		},
		{
			name:     "map/interface",
			expr:     "终极答案 * 100",
			ctx:      map[string]interface{}{"终极答案": 42},
			expected: val.Num(4200),
		},
		{
			name:     "map/interface",
			expr:     "终极答案 * 100",
			ctx:      &map[string]interface{}{"终极答案": 42},
			expected: val.Num(4200),
		},
		{
			name:     "map/interface",
			expr:     "终极答案 * 100",
			ctx:      ref(&map[string]interface{}{"终极答案": 42}),
			expected: val.Num(4200),
		},
		{
			name: "struct/tag",
			expr: "终极答案 * 100",
			ctx: struct {
				Answer int `yae:"终极答案"`
			}{Answer: 42},
			expected: val.Num(4200),
		},
		{
			name:     "struct/tag",
			expr:     `if(false, (1+2)/5^4, "hello world!".len())`,
			ctx:      nil,
			expected: val.Num(12),
		},
		{
			name:     "desugar",
			expr:     `(2 + 3) == 2.+(3)`,
			ctx:      nil,
			expected: val.True,
		},
		{
			name:     "desugar",
			expr:     `strtotime("today") == "today".strtotime()`,
			ctx:      nil,
			expected: val.True,
		},
		{
			name:     "desugar",
			expr:     `'today' == strtotime("today")`,
			ctx:      nil,
			expected: val.True,
		},
		{
			name:     "desugar",
			expr:     `('today +1 day' - 'today') / 3600 == 24`,
			ctx:      nil,
			expected: val.True,
		},
		{
			name: "time",
			expr: `if(时间 < 'today 08:00', 10, 20)`,
			ctx: map[string]time.Time{
				"时间": time.Now(),
			},
			expected: func() *val.Val {
				if time.Now().Unix() < util.Strtotime("today 08:00") {
					return val.Num(10)
				} else {
					return val.Num(20)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Eval(tt.expr, tt.ctx)
			if err != nil {
				t.Errorf("[%s] expected %s but error %s", tt.expr, tt.expected, err)
			} else {
				if !val.Equals(tt.expected, r) {
					t.Errorf("[%s] expected %s actual %s", tt.expr, tt.expected, r)
				}
			}
		})
	}
}

func TestManualEnv(t *testing.T) {
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

	expr := NewExpr().EnableDebug(os.Stderr)
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
			t.Errorf("expected 43 actual %s", v)
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
			t.Errorf("expected 101 actual %s", v)
		}
	}
}

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

	typeEnv, err := conv.TypeEnvOf(Ctx{})
	if err != nil {
		panic(err)
	}
	expr := NewExpr().EnableDebug(os.Stderr)
	closure, err := expr.Compile("if(布尔, 列表[0].姓名.len() + 数字, 0)", typeEnv)
	if err != nil {
		panic(err)
	}

	{
		valEnv, err := conv.ValEnvOf(&Ctx{
			Ok: true,
			N:  42,
			T:  time.Now(),
			Lst: []*Entity{
				{Id: 42, Name: "晓"},
			},
			Obj: &Entity{Id: 42, Name: "晓"},
		})
		if err != nil {
			panic(err)
		}
		v, err := closure(valEnv)
		if err != nil {
			panic(err)
		}
		if !val.Equals(v, val.Num(43)) {
			t.Errorf("expected 43 actual %s", v)
		}
	}
	{
		valEnv, err := conv.ValEnvOf(&Ctx{
			Ok: true,
			N:  100,
			T:  time.Now(),
			Lst: []*Entity{
				{Id: 42, Name: "晓"},
			},
			Obj: &Entity{Id: 42, Name: "晓"},
		})
		if err != nil {
			panic(err)
		}
		v, err := closure(valEnv)
		if err != nil {
			panic(err)
		}
		if !val.Equals(v, val.Num(101)) {
			t.Errorf("expected 101 actual %s", v)
		}
	}
}

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
	expr := NewExpr().EnableDebug(os.Stderr)
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
			t.Errorf("expected 43 actual %s", v)
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
			t.Errorf("expected 101 actual %s", v)
		}
	}
}

func TestRegisterXXX(t *testing.T) {
	type Ctx struct {
		A      int `yae:"a"`
		Answer int `yae:"终极答案"`
	}

	expr := NewExpr().EnableDebug(os.Stderr)

	compileTimeEnv, err := conv.TypeEnvOf(Ctx{})
	if err != nil {
		panic(err)
	}
	closure, err := expr.Compile("a + 终极答案", compileTimeEnv)
	if err != nil {
		panic(err)
	}

	runtimeEnv, err := conv.ValEnvOf(Ctx{
		A:      1,
		Answer: 42,
	})
	if err != nil {
		panic(err)
	}
	v, err := closure(runtimeEnv)
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
	if v.Kind == types.Str {
		fmt.Println(v.Str().V)
	}
}
