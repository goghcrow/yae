package test

import (
	"fmt"
	"testing"

	"github.com/goghcrow/yae"
)

func TestPowerDebug(t *testing.T) {
	type obj struct {
		num int
	}
	for _, tt := range []struct {
		s      string
		v      interface{}
		expect string
	}{
		{
			`obj.num + lst[1] > "hello".len()`,
			map[string]interface{}{
				"obj": obj{num: 42},
				"lst": []int{1, 2, 3},
			},
			`obj.num + lst[1] > "hello".len()
|  |    | |  |   |            |
|  42   44|  2   true         5
{num: 42} [1, 2, 3]`},
		{
			s: `s1 + " " + s2`,
			v: map[string]interface{}{
				"s1": "Hello\nWorld!",
				"s2": "123\n456\n789",
			},
			expect: `s1 + " " + s2
|  |     | |
|  |     | "123\n456\n789"
|  |     "Hello\nWorld! 123\n456\n789"
|  "Hello\nWorld! "
"Hello\nWorld!"`,
		},
		//{
		//	s: "变量1 + 变量2 + max(变量1, 变量2)",
		//	v: struct {
		//		a int `yae:"变量1"`
		//		b int `yae:"变量2"`
		//	}{
		//		a: 1, b: 2,
		//	},
		//	expect: "",
		//},
	} {
		t.Run(tt.s, func(t *testing.T) {
			_, actual, err := yae.Debug(tt.s, tt.v)
			if err != nil {
				t.Errorf("error: %s", err.Error())
			}
			if tt.expect != actual {
				t.Errorf("expect %s actual %s", tt.expect, actual)
			}
			fmt.Println(actual)
		})
	}
	t.Log(len("变量"))
}
