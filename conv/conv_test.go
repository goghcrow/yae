package conv

import (
	"encoding/json"
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/val"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func setObj(obj *val.ObjVal, m map[string]*val.Val) {
	for n, v := range m {
		if !obj.Put(n, v) {
			panic(fmt.Errorf("field %s not found in %s", n, obj.Kind))
		}
	}
}

func TestConv(t *testing.T) {
	tests := []struct {
		name         string
		v            interface{}
		expectedType *types.Kind
		expectedVal  *val.Val
	}{
		{
			name:         "nil",
			v:            nil,
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name:         "int/interface",
			v:            42,
			expectedType: types.Num,
			expectedVal:  val.Num(42),
		},
		{
			name: "int/ptr",
			v: func() interface{} {
				i := 42
				return &i
			}(),
			expectedType: types.Num,
			expectedVal:  val.Num(42),
		},
		{
			name: "int/ptr",
			v: func() *int {
				i := 42
				return &i
			}(),
			expectedType: types.Num,
			expectedVal:  val.Num(42),
		},
		{
			name: "int/ptrptr",
			v: func() **int {
				i := 42
				pi := &i
				return &pi
			}(),
			expectedType: types.Num,
			expectedVal:  val.Num(42),
		},
		{
			name:         "[]any",
			v:            []interface{}{},
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name:         "[]int",
			v:            []int{1},
			expectedType: types.List(types.Num),
			expectedVal: func() *val.Val {
				lst := val.List(types.List(types.Num).List(), 1).List()
				lst.Set(0, val.Num(1))
				return lst.Vl()
			}(),
		},
		{
			name:         "[]any/int",
			v:            []interface{}{1},
			expectedType: types.List(types.Num),
			expectedVal: func() *val.Val {
				lst := val.List(types.List(types.Num).List(), 1).List()
				lst.Set(0, val.Num(1))
				return lst.Vl()
			}(),
		},
		{
			name:         "[]any/num",
			v:            []interface{}{42, 3.14},
			expectedType: types.List(types.Num),
			expectedVal: func() *val.Val {
				lst := val.List(types.List(types.Num).List(), 2).List()
				lst.Set(0, val.Num(42))
				lst.Set(1, val.Num(3.14))
				return lst.Vl()
			}(),
		},
		{
			name: "[]any/num",
			v: func() interface{} {
				a := 42
				b := 3.14
				pb := &b
				c := []interface{}{&a, &pb}
				return &c
			}(),
			expectedType: types.List(types.Num),
			expectedVal: func() *val.Val {
				lst := val.List(types.List(types.Num).List(), 2).List()
				lst.Set(0, val.Num(42))
				lst.Set(1, val.Num(3.14))
				return lst.Vl()
			}(),
		},
		{
			name:         "[]any/any",
			v:            []interface{}{42, 3.14, ""},
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name:         "[]any/any",
			v:            []interface{}{nil},
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name:         "map[string]any",
			v:            map[string]interface{}{},
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name:         "map[string]int/empty",
			v:            map[string]int{},
			expectedType: types.Map(types.Str, types.Num),
			expectedVal: func() *val.Val {
				return val.Map(types.Map(types.Str, types.Num).Map())
			}(),
		},
		{
			name:         "map[string]any",
			v:            map[string]interface{}{"a": 1, "b": 2},
			expectedType: types.Map(types.Str, types.Num),
			expectedVal: func() *val.Val {
				m := val.Map(types.Map(types.Str, types.Num).Map()).Map()
				m.Put(val.Str("a"), val.Num(1))
				m.Put(val.Str("b"), val.Num(2))
				return m.Vl()
			}(),
		},
		{
			name:         "map[string]any",
			v:            map[string]interface{}{"a": 1, "b": "2"},
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name:         "map[int]num",
			v:            map[int]interface{}{1: 1, 2: 2},
			expectedType: types.Map(types.Num, types.Num),
			expectedVal: func() *val.Val {
				m := val.Map(types.Map(types.Num, types.Num).Map()).Map()
				m.Put(val.Num(1), val.Num(1))
				m.Put(val.Num(2), val.Num(2))
				return m.Vl()
			}(),
		},
		{
			name:         "map[int64]num",
			v:            map[int64]interface{}{1: 1, 2: 2},
			expectedType: types.Map(types.Num, types.Num),
			expectedVal: func() *val.Val {
				m := val.Map(types.Map(types.Num, types.Num).Map()).Map()
				m.Put(val.Num(1), val.Num(1))
				m.Put(val.Num(2), val.Num(2))
				return m.Vl()
			}(),
		},
		{
			name:         "map[int64]any",
			v:            map[int64]interface{}{1: 1, 2: "2"},
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name:         "map[string]any",
			v:            map[string]interface{}{"a": 1, "b": 2},
			expectedType: types.Map(types.Str, types.Num),
			expectedVal: func() *val.Val {
				m := val.Map(types.Map(types.Str, types.Num).Map()).Map()
				m.Put(val.Str("a"), val.Num(1))
				m.Put(val.Str("b"), val.Num(2))
				return m.Vl()
			}(),
		},
		{
			name:         "map[string]any",
			v:            map[string]interface{}{"a": 1, "b": "2"},
			expectedType: nil,
			expectedVal:  nil,
		},
		{
			name: "struct/nested/nil",
			v: func() interface{} {
				return &struct {
					Nested *struct {
						A int
					}
				}{
					Nested: &struct {
						A int
					}{42},
				}
			}(),
			expectedType: types.Obj([]types.Field{
				{"Nested", types.Obj([]types.Field{
					{"A", types.Num},
				})},
			}),
			expectedVal: func() *val.Val {
				nestedT := types.Obj([]types.Field{
					{"A", types.Num},
				}).Obj()
				obj := val.Obj(types.Obj([]types.Field{
					{"Nested", nestedT.Kd()},
				}).Obj()).Obj()
				nested := val.Obj(nestedT).Obj()
				nested.Put("A", val.Num(42))
				obj.Put("Nested", nested.Vl())
				return obj.Vl()
			}(),
		},
		{
			name: "struct/nested/nil",
			v: func() interface{} {
				return &struct {
					Nested *struct {
						A int
					}
				}{
					Nested: nil,
				}
			}(),
			expectedType: types.Obj([]types.Field{
				{"Nested", types.Obj([]types.Field{
					{"A", types.Num},
				})},
			}),
			expectedVal: nil,
		},
		{
			name: "struct",
			v: &struct {
				A int
				B interface{}
				C *int
				D *interface{}
				E interface{}
				F *interface{}
			}{
				A: 42,
				B: 42,
				C: func() *int {
					i := 42
					return &i
				}(),
				D: func() *interface{} {
					var i interface{} = 42
					return &i
				}(),
				E: struct {
					A int
				}{42},
				F: func() *interface{} {
					var i interface{} = struct {
						A int
					}{42}
					return &i
				}(),
			},
			expectedType: func() *types.Kind {
				return types.Obj([]types.Field{
					{"A", types.Num},
					{"B", types.Num},
					{"C", types.Num},
					{"D", types.Num},
					{"E", types.Obj([]types.Field{
						{"A", types.Num},
					})},
					{"F", types.Obj([]types.Field{
						{"A", types.Num},
					})},
				})
			}(),
			expectedVal: func() *val.Val {
				obj := val.Obj(types.Obj([]types.Field{
					{"A", types.Num},
					{"B", types.Num},
					{"C", types.Num},
					{"D", types.Num},
					{"E", types.Obj([]types.Field{
						{"A", types.Num},
					})},
					{"F", types.Obj([]types.Field{
						{"A", types.Num},
					})},
				}).Obj()).Obj()

				ef := val.Obj(types.Obj([]types.Field{
					{"A", types.Num},
				}).Obj()).Obj()
				setObj(ef, map[string]*val.Val{
					"A": val.Num(42),
				})
				setObj(obj, map[string]*val.Val{
					"A": val.Num(42),
					"B": val.Num(42),
					"C": val.Num(42),
					"D": val.Num(42),
					"E": ef.Vl(),
					"F": ef.Vl(),
				})
				return obj.Vl()
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			k, err := TypeOf(tt.v)
			if err != nil {
				if tt.expectedType != nil {
					t.Errorf("[typeof] expect %s actual error `%s`", tt.expectedType, err)
				}
			} else {
				if tt.expectedType == nil {
					t.Errorf("[valof] expect %s actual %s", tt.expectedType, k)
				} else if !types.Equals(tt.expectedType, k) {
					t.Errorf("[valof] expect %s actual %s", tt.expectedType, k)
				}
			}

			v, err := ValOf(tt.v)
			if err != nil {
				if tt.expectedVal != nil {
					t.Errorf("expect %s actual error `%s`", tt.expectedVal, err)
				}
			} else {
				if tt.expectedVal == nil {
					t.Errorf("expect %s actual %s", tt.expectedVal, k)
				} else if !val.Equals(tt.expectedVal, v) {
					t.Errorf("expect %s actual %s", tt.expectedVal, k)
				}
			}
		})
	}
}

func TestConvValOf(t *testing.T) {
	lstNum := val.List(types.List(types.Num).List(), 0).List()
	lstNum.Add(val.Num(42))

	lstStr := val.List(types.List(types.Str).List(), 0).List()
	lstStr.Add(val.Str("Hello"))

	lstLstNum := val.List(types.List(types.List(types.Num)).List(), 0).List()
	lstLstNum.Add(lstNum.Vl())

	tests := []struct {
		name     string //make ide happy
		v        interface{}
		expected *val.Val
	}{
		{"bool/true", true, val.True},
		{"bool/false", false, val.False},
		{"num/int", 42, val.Num(42)},
		{"num/float", 3.14, val.Num(3.14)},
		{"str", "Hello", val.Str("Hello")},
		{"array/int", [1]int{42}, lstNum.Vl()},
		{"array/str", [1]string{"Hello"}, lstStr.Vl()},
		{"array/slice/num", [1][1]int{{42}}, lstLstNum.Vl()},
		{"slice/int", []int{42}, lstNum.Vl()},
		{"slice/str", []string{"Hello"}, lstStr.Vl()},
		{"slice/slice/num", [][]int{{42}}, lstLstNum.Vl()},
		{
			name: "map/str/num",
			v:    map[string]int64{"id": 42},
			expected: func() *val.Val {
				mapStrNum := val.Map(types.Map(types.Str, types.Num).Map()).Map()
				mapStrNum.Put(val.Str("id"), val.Num(42))
				return mapStrNum.Vl()
			}(),
		},
		{
			name:     "time",
			v:        time.Unix(1, 0),
			expected: val.Time(time.Unix(1, 0)),
		},
		{
			name: "time/ptr",
			v: func() *time.Time {
				t := time.Unix(1, 0)
				return &t
			}(),
			expected: val.Time(time.Unix(1, 0)),
		},
		{
			name: "struct",
			v: struct {
				Id   int64  `yae:"id"`
				Name string `yae:"name"`
			}{42, "晓"},
			expected: func() *val.Val {
				obj := val.Obj(types.Obj([]types.Field{
					{"id", types.Num},
					{"name", types.Str},
				}).Obj()).Obj()
				obj.Put("id", val.Num(42))
				obj.Put("name", val.Str("晓"))
				return obj.Vl()
			}()},
		{
			name: "composite",
			v: []struct {
				Id     int64  `yae:"id"`
				Name   string `yae:"name"`
				Nested struct {
					Props map[string][]string `yae:"props"`
				} `yae:"nested"`
			}{
				{
					Id:   42,
					Name: "晓",
					Nested: struct {
						Props map[string][]string `yae:"props"`
					}{
						Props: map[string][]string{
							"set": {"a", "b", "c"},
						},
					},
				},
			},
			expected: func() *val.Val {
				typeProps := types.Map(types.Str, types.List(types.Str)).Map()
				typeNested := types.Obj([]types.Field{
					{"props", typeProps.Kd()},
				}).Obj()
				typeObj := types.Obj([]types.Field{
					{"id", types.Num},
					{"name", types.Str},
					{"nested", typeNested.Kd()},
				})
				// list[{id: num, name: str, nested: {props: map[str, list[str]]}}]
				typeRes := types.List(typeObj).List()

				m := val.Map(typeProps).Map()
				m.Put(val.Str("set"), func() *val.Val {
					lstStr := val.List(types.List(types.Str).List(), 0).List()
					lstStr.Add(val.Str("a"))
					lstStr.Add(val.Str("b"))
					lstStr.Add(val.Str("c"))
					return lstStr.Vl()
				}())

				nestedVal := val.Obj(typeNested).Obj()
				nestedVal.Put("props", m.Vl())

				o := val.Obj(typeObj.Obj()).Obj()
				setObj(o, map[string]*val.Val{
					"id":     val.Num(42),
					"name":   val.Str("晓"),
					"nested": nestedVal.Vl(),
				})

				res := val.List(typeRes, 0).List()
				res.Add(o.Vl())
				return res.Vl()
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			vl, err := ValOf(tt.v)
			if err != nil {
				t.Errorf("error %s", err.Error())
			}
			if !val.Equals(vl, tt.expected) {
				t.Errorf("%s, expect %s actual %s", pretty(tt.v), tt.expected, vl)
			}
		})
	}
}

func TestTypeOf(t *testing.T) {
	tt := time.Time{}
	ttptr := &tt
	ttptrptr := &ttptr

	kind, err := TypeOf(tt)
	if err != nil {
		t.Errorf(err.Error())
	}
	if kind != types.Time {
		t.Errorf("expect time, actual %s", kind)
	}

	kind, err = TypeOf(ttptr)
	if err != nil {
		t.Errorf(err.Error())
	}
	if kind != types.Time {
		t.Errorf("expect time actual %s", kind)
	}

	kind, err = TypeOf(ttptrptr)
	if err != nil {
		t.Errorf(err.Error())
	}
	if kind != types.Time {
		t.Errorf("expect time, actual %s", kind)
	}

	kind, err = TypeOf(map[string]interface{}{})
	if err == nil || kind != nil {
		t.Errorf("expect err")
	}
}

func TestPtrInterface(t *testing.T) {
	var i interface{} = 42
	k, err := TypeOf(i)
	assert(err == nil)
	assert(k == types.Num)

	k, err = TypeOf(&i)
	assert(err == nil)
	assert(k == types.Num)

	v, err := ValOf(i)
	assert(err == nil)
	assert(val.Equals(v, val.Num(42)))

	v, err = ValOf(&i)
	assert(err == nil)
	assert(val.Equals(v, val.Num(42)))
}

func TestReflectInterfaceElem(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var i interface{} = 42

		assert(reflect.ValueOf(i).Type().Kind().String() == "int")
		assert(reflect.ValueOf(i).Int() == 42)

		// 注意这里: &i 变成 interface,而不是 int, 需要 elem
		assert(reflect.ValueOf(&i).Elem( /*deref*/ ).Interface().(int) == 42)
		assert(reflect.ValueOf(&i).Kind().String() == "ptr")
		assert(reflect.ValueOf(&i).Elem( /*deref*/ ).Kind().String() == "interface")
		assert(reflect.ValueOf(&i).Elem( /*deref*/ ).Elem( /*de_iface*/ ).Kind().String() == "int")

		assert(reflect.ValueOf(&i).Pointer() == uintptr(unsafe.Pointer(&i)))
		// t.Log(*(*int)(unsafe.Pointer(reflect.ValueOf(&i).Pointer()))) // 注意这里不是 *int
		assert((*(*interface{})(unsafe.Pointer(reflect.ValueOf(&i).Pointer()))).(int) == 42)

		assert(reflect.ValueOf(&i).Elem( /*deref*/ ).Interface().(int) == 42)
		assert(reflect.ValueOf(&i).Elem( /*deref*/ ).Elem( /*de_iface*/ ).Int() == 42)
	})

	t.Run("", func(t *testing.T) {
		v := reflect.ValueOf(struct{ a interface{} }{1})
		assert(v.Field(0).Kind().String() == "interface")
		assert(v.Field(0).Elem( /*de_iface*/ ).Kind().String() == "int")
	})

	t.Run("", func(t *testing.T) {
		i := 42
		v := reflect.ValueOf(struct{ a interface{} }{&i})
		assert(v.Field(0).Kind().String() == "interface")
		assert(v.Field(0).Elem( /*de_iface*/ ).Kind().String() == "ptr")
		assert(v.Field(0).Elem( /*de_iface*/ ).Elem().Kind().String() == "int")

		assert(*(*int)(unsafe.Pointer(v.Field(0).Elem( /*de_iface*/ ).Pointer())) == 42)
		assert(v.Field(0).Elem( /*de_iface*/ ).Elem().Int() == 42)
	})
}

func pretty(v interface{}) string {
	s, _ := json.Marshal(v)
	return string(s)
}

func assert(cond bool) {
	if !cond {
		panic(nil)
	}
}
