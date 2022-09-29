package facade

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
	"time"
)

func TestTmp(t *testing.T) {
	//input := "终极答案 == 42"
	//input := `"\"中文\nabc\n123\u00e9"`
	//input := "`a\nb`"
	//input := "`a\nb` + \"`\" + `HELLO`"

	//input := `1a+1`
	//input := `1`

	//input := `iff`
	//input := `.^.`
	//input := `. ^.`

	//input := `-1`

	//input := "if true || false then 1 else 2"
	//input := "if 1+1>2 then 1 else 2"
	//input := "if (if true then false else false end) then 1 else 2"
	//input := "1+2+3" //todo 测试 infixn
	//input := "today()"
	//input := "today(800)" // 今天 8 点
	//input := "0x12"
	//input := `now() - today()`
	//input := `(
	//			time("2006-01-02 15:04:05", "2022-09-28 02:00:00") - today()
	//	) / 3600`
	//input := "if(false,1,2) + 1"
	input := "if false then 1 else 2"

	// todo contains 、in 之类也做成自定义中缀不结合操作符
	// todo 测试空数组的类型信息,

	expr := NewExpr().EnableDebug(os.Stderr)
	// 注册函数和加上下文变量
	//rt.Put("", val.List())
	//expr.RegisterFun(val.Fun(types.Fun("", []*types.Kind{}, types.Str), func(v ...*val.Val) *val.Val {
	//	return val.Str("Hello")
	//}).Fun())

	//expr.RegisterTransformer()

	env0 := types.NewEnv()
	env0.Put("a", types.Num)
	env0.Put("终极答案", types.Num)

	// todo, 1. 这里的 env0 可以换成 struct + tag 反射来定义??
	// todo 2. closure 开始需要加入 根据 env0 检查 env1 类型!!!
	closure, err := expr.Compile(input, env0)
	if err != nil {
		panic(err)
	}
	env1 := val.NewEnv()
	env1.Put("a", val.Num(2))
	env1.Put("终极答案", val.Num(42))
	v, err := closure(env1)
	if err != nil {
		fmt.Println(err)
		//panic(err)
	}
	//v.Kind
	//switch v.Kind.Type {
	//case types.TBool:
	//}
	fmt.Println(v)
	if v.Kind == types.Str {
		fmt.Println(v.Str().V)
	}

	////////////////////////////////////

	//{
	//	m := val.Map(types.Map(types.Str, types.List(types.Num)).Map()).Map()
	//	lst := val.List(types.List(types.Num).List(), 0).List()
	//	lst.Add(val.Num(42))
	//	m.Put(val.Str("Hello"), lst.Vl())
	//
	//	v, _ := m.Get(val.Str("Hello"))
	//	fmt.Println(v.List().V[0].Num().V)
	//}

	//expr := Parser(lex.Lex("if(false, f((1+2)/5^4), \"hello world!\".length())"))
	//
	//fmt.Println(expr)
	//fmt.Println(types.Map(types.Str, types.List(types.Num)))
	//fmt.Println(trans.Desugar(expr))
	//
	//{
	//	actual := Parser(lex.Lex("a.b.c(1,2)"))
	//	expected := Call(
	//		Ident("c"),
	//		[]*Expr{
	//			Member(Ident("a"), Ident("b").Ident()),
	//			LitNum("1"),
	//			LitNum("2"),
	//		},
	//	)
	//	if !reflect.DeepEqual(expected, actual) {
	//		t.Errorf("expect `%s` actual `%s`", expected, actual)
	//	}
	//}
}

func TestTmp1(t *testing.T) {
	input := `if(some_time_var < today("08:00"), 10, 20)`

	expr := NewExpr()
	compileEnv := types.NewEnv()
	compileEnv.Put("some_time_var", types.Time)
	closure, err := expr.Compile(input, compileEnv)
	if err != nil {
		panic(err)
	}

	runtimeEnv := val.NewEnv()
	runtimeEnv.Put("some_time_var", val.Time(time.Now()))
	v, err := closure(runtimeEnv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s = %s\n", input, v)
}
