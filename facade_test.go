package yae

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"testing"
	"time"
)

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
			expr:     `(2 + 3) == 2. +(3)`, // .和+之间必须有空格, 因为支持自定义操作符,会优先匹配完整操作符
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

		{
			name:     "-",
			expr:     `-42`,
			ctx:      nil,
			expected: val.Num(-42),
		},
		{
			name:     "-",
			expr:     `+42`,
			ctx:      nil,
			expected: val.Num(42),
		},
		{
			name:     "-",
			expr:     `1-1`,
			ctx:      nil,
			expected: val.Num(0),
		},
		{
			name:     "-",
			expr:     `1--1`,
			ctx:      nil,
			expected: val.Num(2),
		},
		{
			name:     "-",
			expr:     `--1`,
			ctx:      nil,
			expected: val.Num(1),
		},
		{
			name:     "-",
			expr:     `+-1`,
			ctx:      nil,
			expected: val.Num(-1),
		},
		{
			name:     "-",
			expr:     `-+1`,
			ctx:      nil,
			expected: val.Num(-1),
		},
		{
			name:     "list/empty",
			expr:     `[]`,
			ctx:      nil,
			expected: val.List(types.List(types.Bottom).List(), 0),
		},
		{
			name:     "obj/empty",
			expr:     `{}`,
			ctx:      nil,
			expected: val.Obj(types.Obj([]types.Field{}).Obj()),
		},
		{
			name:     "map/empty",
			expr:     `[:]`,
			ctx:      nil,
			expected: val.Map(types.Map(types.Bottom, types.Bottom).Map()),
		},
		{
			name:     "lit/string",
			expr:     `"\"中文\nabc\n123\u00e9"`,
			ctx:      nil,
			expected: val.Str("\"中文\nabc\n123\u00e9"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			r, err := Eval(tt.expr, tt.ctx)
			if err != nil {
				t.Errorf("[%s] expect %s but error %s", tt.expr, tt.expected, err)
			} else {
				if !val.Equals(tt.expected, r) {
					t.Errorf("[%s] expect %s actual %s", tt.expr, tt.expected, r)
				}
			}
		})
	}
}
