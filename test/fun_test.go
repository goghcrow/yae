package test

import (
	"github.com/goghcrow/yae/closure"
	"github.com/goghcrow/yae/compiler"
	"github.com/goghcrow/yae/fun"
	"github.com/goghcrow/yae/lexer"
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/trans"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"github.com/goghcrow/yae/vm"
	"testing"
)

func TestMaybe(t *testing.T) {

	// 把允许为空值或者为 nil 的变量或者字段声明称 maybe 类型, 通过 get + 默认值 获取值

	makeEnv := func() (typedEnv *types.Env, compileEvalEnv *val.Env) {
		maybeInt := types.Maybe(types.Num).Maybe()
		nothing := val.Maybe(maybeInt, nil)
		just := val.Maybe(maybeInt, val.Num(42))

		objWithMaybeField := types.Obj([]types.Field{
			{"nothing", maybeInt.Kd()},
			{"just", maybeInt.Kd()},
		}).Obj()

		tyEnv := types.NewEnv()
		vlEnv := val.NewEnv()
		tyEnv.Put("nothing", maybeInt.Kd())
		tyEnv.Put("just", maybeInt.Kd())
		tyEnv.Put("obj", objWithMaybeField.Kd())

		vlEnv.Put("nothing", nothing)
		vlEnv.Put("just", just)

		obj := val.Obj(objWithMaybeField).Obj()
		obj.Put("nothing", nothing)
		obj.Put("just", just)
		vlEnv.Put("obj", obj.Vl())

		return tyEnv, vlEnv
	}

	tests := []struct {
		expr     string
		expected *val.Val
	}{
		{"get(nothing, 100)", val.Num(100)},
		{"get(just, 100)", val.Num(42)},
		{"get(obj.nothing, 100)", val.Num(100)},
		{"get(obj.just, 100)", val.Num(42)},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			tyEnv, vlEnv := makeEnv()
			actual := eval(tt.expr, closure.Compile, tyEnv, vlEnv)
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})

		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			tyEnv, vlEnv := makeEnv()
			actual := eval(tt.expr, vm.Compile, tyEnv, vlEnv)
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})
	}
}

func TestIf(t *testing.T) {
	tests := []struct {
		expr     string
		expected *val.Val
	}{
		{"if(true, print(1), print(2))", val.Num(1)},  // print 1
		{"if(false, print(1), print(2))", val.Num(2)}, // print 2
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			actual := eval(tt.expr, closure.Compile, types.NewEnv(), val.NewEnv())
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})

		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			actual := eval(tt.expr, vm.Compile, types.NewEnv(), val.NewEnv())
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})
	}
}

func TestFun(t *testing.T) {
	tests := []struct {
		expr     string
		expected *val.Val
	}{
		{`false && false`, val.False},
		{`false && true`, val.False},
		{`true && false`, val.False},
		{`true && true`, val.True},

		{`false || false`, val.False},
		{`false || true`, val.True},
		{`true || false`, val.True},
		{`true || true`, val.True},

		{`!false`, val.True},
		{`!true`, val.False},

		{`true == true`, val.True},
		{`true != true`, val.False},
		{`false == false`, val.True},
		{`false != false`, val.False},

		{`if(true, 1, 2)`, val.Num(1)},
		{`if(false, 1, 2)`, val.Num(2)},

		{`1 == 1`, val.True},
		{`1 != 1`, val.False},
		{`3.14 == 3.14`, val.True},
		{`3.14 != 3.14`, val.False},
		{`"" == ""`, val.True},
		{`"" != ""`, val.False},
		{`'today' == 'today'`, val.True},
		{`'today' != 'today'`, val.False},
		{`[] == []`, val.True},
		{`[] != []`, val.False},
		{`[1,2,3] == [1,2,3]`, val.True},
		{`[1,2,3] != [1,2,3]`, val.False},
		{`[:] == [:]`, val.True},
		{`[:] != [:]`, val.False},
		{`[42:"42"] == [42:"42"]`, val.True},
		{`[42:"42"] != [42:"42"]`, val.False},
		{`["":42] == ["":42]`, val.True},
		{`["":42] != ["":42]`, val.False},
		{`[{}] == [{}]`, val.True},
		{`[{}] != [{}]`, val.False},
		{`[{id:42, name:"晓"}] == [{id:42, name:"晓"}]`, val.True},
		{`[{id:42, name:"晓"}] != [{id:42, name:"晓"}]`, val.False},
		{`[
			{
				id: 42, 
				name: "晓", 
				list: [1,2,3], 
				map:[
					"obj": {
						id:42
					}
				]
			}
		]
		== 
		[
			{
				id: 42, 
				name: "晓", 
				list: [1,2,3], 
				map:[
					"obj": {
						id:42
					}
				]
			}
		]`, val.True},
		{`[{id:42, name:"晓", list:[1,2,3], map:["obj":{id:42}]}] != [{id:42, name:"晓", list:[1,2,3], map:["obj":{id:42}]}]`, val.False},

		{`2 > 1`, val.True},
		{`1 > 2`, val.False},
		{`2 < 1`, val.False},
		{`1 < 2`, val.True},
		{`2 >= 1`, val.True},
		{`1 >= 2`, val.False},
		{`2 <= 1`, val.False},
		{`1 <= 2`, val.True},
		{`42 >= 42`, val.True},
		{`42 >= 42`, val.True},
		{`42 <= 42`, val.True},
		{`42 <= 42`, val.True},

		{`'today' == 'today'`, val.True},
		{`'today' != 'today'`, val.False},
		{`'today 09:00' > 'today 08:00'`, val.True},
		{`'today 08:00' < 'today 09:00'`, val.True},
		{`'today 09:00' < 'today 08:00'`, val.False},
		{`'today 08:00' > 'today 09:00'`, val.False},
		{`'today 09:00' >= 'today 08:00'`, val.True},
		{`'today 08:00' <= 'today 09:00'`, val.True},
		{`'today 09:00' <= 'today 08:00'`, val.False},
		{`'today 08:00' >= 'today 09:00'`, val.False},
		{`'today 09:00' >= 'today 09:00'`, val.True},
		{`'today 09:00' <= 'today 09:00'`, val.True},
		{`'today 09:00' <= 'today 09:00'`, val.True},
		{`'today 09:00' >= 'today 09:00'`, val.True},

		{`len("Hello")`, val.Num(5)},
		{`len("晓")`, val.Num(1)},
		{`len("1a晓😁")`, val.Num(4)},

		{`len([])`, val.Num(0)},
		{`len([1])`, val.Num(1)},
		{`len([1,2])`, val.Num(2)},
		{`len([:])`, val.Num(0)},
		{`len([1:1])`, val.Num(1)},
		{`len([1:1,1:2])`, val.Num(1)},
		{`len([1:1,2:2])`, val.Num(2)},
		{`len(["a":0,"a":0])`, val.Num(1)},
		{`len(["a":0,"b":0])`, val.Num(2)},

		{`"Hello" + " " + "World!"`, val.Str("Hello World!")},
		{`"Hello" + " World!\n" == "Hello World!\n"`, val.True},
		{`"我" + " 是 " + " 谁 !\n🚑" + "\u6653"`, val.Str("我 是  谁 !\n🚑晓")},

		{"-42 == -(42)", val.True},
		{"+42 == +(42)", val.True},
		{"1.1 + 2.2 == 3.3", val.True},
		{"2.2 - 1.1 == 1.1", val.True},
		{"3.3 * 4 == 13.2", val.True},
		{"100 / 50 == 2", val.True},
		{"100 % 3 == 1", val.True},
		{"2 ^ 3 == 8", val.True},
		{"2 ^ 3 ^ 2 == 2 ^ 9", val.True}, // 右结合
		{"1 - 2 + 3 * 4 == 11", val.True},
		{"1 - (-1 + 3) * 4 == -7", val.True},
		{"(1 - -1 + 3) * 4 == 20", val.True},

		{"max(1,2) == 2", val.True},
		{"max([1,2,3]) == 3", val.True},
		{"min(1,2) == 1", val.True},
		{"min([1,2,3]) == 1", val.True},
		{`abs(-1)`, val.Num(1)},
		{`round(1.4)`, val.Num(1)},
		{`round(1.5)`, val.Num(2)},
		{`floor(1.9)`, val.Num(1)},
		{`ceil(1.1)`, val.Num(2)},

		{"print(1)", val.Num(1)},

		{`true ? 1 : true ? 2 : 3`, val.Num(1)},   // 右结合
		{`false ? 1 : true ? 2 : 3`, val.Num(2)},  // 右结合
		{`false ? 1 : false ? 2 : 3`, val.Num(3)}, // 右结合

		{`[1,2,3, ] == [1,2,3]`, val.True},             // 尾部可以多余逗号
		{`[1:"1",2:"2", ] == [1:"1",2:"2"]`, val.True}, // 尾部可以多余逗号

		// map必须 kv 类型一致, map[k,v], k 支持 num, str, time, bool
		// 对象的的key 只能是 ident, val 类型可以不一样
		//  对象不能有重复的 key ，typecheck 检查

		{`[true,false,true][0] == true`, val.True},
		{`[true,false,true][1] == false`, val.True},
		{`[true,false,true][2] == true`, val.True},

		{`{id:42, name:"晓"}.id == 42`, val.True},
		{`{id:42, name:"晓"}.name == "晓"`, val.True},
		{`[1:2, 3:4][1] == 2`, val.True},
		{`[1:2, 3:4][3] == 4`, val.True},
		{`["id":"x", "name":"xiao"]["id"] == "x"`, val.True},

		{`if(1 + 2 > 3, 'today 08:00', 'today 09:00') == 'today 09:00'`, val.True},

		{`string(42)`, val.Str("42")},
		{`string("s")`, val.Str("s")},
		{`string("[1,2]")`, val.Str("[1,2]")},

		{`match("^\\d+$", "123")`, val.True},
		{`match("^\\d+$", "123a")`, val.False},

		{`isset([1:""], 0)`, val.False},
		{`isset([1:""], 1)`, val.True},

		{`[1,2,3,n].get(2, 42) == 3`, val.True},
		{`[1,2,3,n].get(3, 42) == 42`, val.True},
		{`[1,2,3,n].get(4, 42) == 42`, val.True},

		{`["id":1,"nil":n].get("id", 42) == 1`, val.True},
		{`["id":1,"nil":n].get("nil", 42) == 42`, val.True},
		{`["id":1,"nil":n].get("not_exist", 42) == 42`, val.True},

		{`union([1,2,3], [2,3,4]) == [1,2,3,4]`, val.True},
		{`union([1,2,2,3], [2,3,3,4]) == [1,2,3,4]`, val.True},
		{`intersect([1,2,3], [2,3,4]) == [2,3] `, val.True},
		{`intersect([1,2,2,3], [2,3,3,4]) == [2,3]`, val.True},
		{`diff([1,2,3], [2,3,4]) == [1]`, val.True},
		{`diff([1,2,2,3], [2,3,3,4]) == [1]`, val.True},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			typeEnv := types.NewEnv()
			valEnv := val.NewEnv()
			typeEnv.Put("n", types.Num)
			valEnv.Put("n", nil)
			actual := eval(tt.expr, closure.Compile, typeEnv, valEnv)
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect %s actual %s", expected, actual)
			}
		})

		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			typeEnv := types.NewEnv()
			valEnv := val.NewEnv()
			typeEnv.Put("n", types.Num)
			valEnv.Put("n", nil)
			actual := eval(tt.expr, vm.Compile, typeEnv, valEnv)
			expected := tt.expected
			if !val.Equals(expected, actual) {
				t.Errorf("expect %s actual %s", expected, actual)
			}
		})
	}
}

func TestTypeError(t *testing.T) {
	for _, expr := range []string{
		`if(true,1,"2")`,
		`if(true,1,"2")`,
		`if(true,"1",2)`,
		`[1]==["1"]`,
	} {
		k := typeError(expr, t)
		if k != nil {
			t.Errorf("%s  expect type error actual `%s`", expr, k)
		}
	}
}

func typeError(s string, t *testing.T) (k *types.Kind) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("%s => %v", s, r)
			k = nil
		}
	}()

	return infer(s)
}

var typecheckEnv = types.NewEnv()
var compileEnv = val.NewEnv()

func init() {
	initEnv(typecheckEnv, compileEnv)
}

func eval(s string, compile compiler.Compiler, typedEnv *types.Env, compileEvalEnv *val.Env) *val.Val {
	toks := lexer.NewLexer(oper.BuildIn()).Lex(s)
	ast := parser.NewParser(oper.BuildIn()).Parse(toks)
	ast = trans.Desugar(ast)

	_ = types.Check(ast, typedEnv.Inherit(typecheckEnv))
	valuedEnv := compileEvalEnv.Inherit(compileEnv)
	compiled := compile(ast, valuedEnv)

	runtimeEnv := val.NewEnv()
	runtimeEnv = runtimeEnv.Inherit(valuedEnv)
	return compiled(runtimeEnv)
}

func infer(s string) *types.Kind {
	toks := lexer.NewLexer(oper.BuildIn()).Lex(s)
	ast := parser.NewParser(oper.BuildIn()).Parse(toks)
	ast = trans.Desugar(ast)
	return types.Check(ast, typecheckEnv)
}

func initEnv(typecheckEnv *types.Env, compileEnv *val.Env) {
	for _, f := range fun.BuildIn() {
		typecheckEnv.RegisterFun(f.Kind)
		compileEnv.RegisterFun(f)
	}
}
