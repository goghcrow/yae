package facade

import (
	"fmt"
	"github.com/goghcrow/yae/env"
	"github.com/goghcrow/yae/env0"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"os"
	"testing"
)

func TestTmp(t *testing.T) {
	fmt.Println(util.ParseNum("0x12"))

	//input := "a - 2 + 3 * 4 == 11"
	//input := `"hello" + " World!\n"`
	//input := `"我" + " 是 " + " 谁 !\n🚑" + `
	//input := `"我" + " 是 " + " 谁 !\n🚑" + "Sp\u00e9cification"`
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

	// todo 做一张真值表测试, 做一个内置副作用函数看是否短路
	//input := "true and false"
	//input := "true and true"
	//input := "true && false"
	//input := "true && true"
	//input := "true || false"
	//input := "true or false"
	//input := "false or true"
	//input := "if true || false then 1 else 2"
	//input := "if 1+1>2 then 1 else 2"
	//input := "if (if true then false else false end) then 1 else 2"
	//input := "1+2+3" //todo 测试 infixn
	//input := "3.14 > 3"
	//input := "today() == today()"
	//input := "today()"
	//input := "today(800)" // 今天 8 点
	//input := "0x12"
	//input := "max(1,2)==2"
	//input := `now() - today()`
	//input := `(
	//			time("2006-01-02 15:04:05", "2022-09-28 02:00:00") - today()
	//	) / 3600`
	//input := "if(false,1,2) + 1"
	//input := "if false then 1 else 2"
	//input := `if(true,1,"2")`
	//input := `if(true,1,"2")`
	//input := `if(true,"1",2)`
	//input := `[1]==["1"]`
	//input := `[1]==[1]`
	//input := `[1]==[2]`

	// 右结合
	//input := `true ? 1 : true ? 2 : 3` // 1
	//input := `false ? 1 : true ? 2 : 3` // 2
	//input := `false ? 1 : false ? 2 : 3` // 3

	//input := `len([])` // 3
	//input := `len([1])`   // 3
	//input := `len([1,2])` // 3

	//input := `[true,false,true][1] == false`

	// !!!!!!!!!!! TODO ast 加通用字段, 来保存 已经推导好的值类型

	//input := `[1,2,3, ]` // 尾部可以多余逗号
	input := `[1:"1",2:"2", ]` // 尾部可以多余逗号, map必须 kv 类型一致, map[k,v], k 支持 num, str, time, bool
	// 对象的的key 只能是 ident, val 类型可以不一样

	// todo 改成这样
	// todo 对象不能有重复的 key ，typecheck 检查
	// todo map 也检查下
	// 一个类型一样, 一个类型可以不一样
	//input := `{id:42, id:"xiao"}.id`
	//input := `[:]`
	//input := `["id":"x", "name":"xiao"]["id"]`

	// todo
	// a.b
	// a[1]

	// todo contains 、in 之类也做成自定义中缀不结合操作符

	// todo 测试空数组的类型信息,
	// todo 或者 不支持空数组字面量, 把 unit 类型都删了

	expr := NewExpr().EnableDebug(os.Stderr)
	// 注册函数和加上下文变量
	//rt.Put("", val.List())
	//expr.RegisterFun(val.Fun(types.Fun("", []*types.Kind{}, types.Str), func(v ...*val.Val) *val.Val {
	//	return val.Str("Hello")
	//}).Fun())

	//expr.RegisterTransformer()

	env0 := env0.NewEnv()
	env0.Put("a", types.Num)
	env0.Put("终极答案", types.Num)

	// todo, 1. 这里的 env0 可以换成 struct + tag 反射来定义??
	// todo 2. closure 开始需要加入 根据 env0 检查 env1 类型!!!
	closure, err := expr.Compile(input, env0)
	if err != nil {
		panic(err)
	}
	env1 := env.NewEnv()
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

	//expr := parser.Parser("if(false, f((1+2)/5^4), \"hello world!\".lenth())")
	//fmt.Println(expr)
	//fmt.Println(types.Map(types.Str, types.List(types.Num)))
	//
	//m := val.Map(types.Str, types.List(types.Num)).Map()
	//
	//lst := val.List(types.Num).List()
	//lst.Add(val.Num(42))
	//m.Put(val.Str("Hello"), lst.Vl())
	//
	//v, _ := m.Get(val.Str("Hello"))
	//fmt.Println(v.List().V[0].Num().V)
	//
	//expr = parser.Parser("lst[1](a.b.c.d(1,2,3))")
	//fmt.Println(trans.Desugar(expr))
}
