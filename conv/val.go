package conv

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"reflect"
	"time"
)

func MustValOf(v interface{}) *val.Val {
	vl, err := ValOf(v)
	if err != nil {
		panic(err)
	}
	return vl
}

func ValOf(v interface{}) (vl *val.Val, err error) {
	return valOfRV(reflect.ValueOf(v))
}

func valOfRV(rv reflect.Value) (vl *val.Val, err error) {
	defer util.Recover(&err)
	return valOf(rv, 0), err
}

func valOf(rv reflect.Value, lv int) *val.Val {
	if lv > maxLevel {
		panic("max nested depth exceeded")
	}
	if isNil(rv) {
		panic(fmt.Errorf("val: Of(nil %v)", rv))
	}
	rt := rv.Type()
	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		rt = rv.Type()
	}
	if rt == typeOfTime {
		return val.Time(rv.Interface().(time.Time))
	}

	// 数字可能会丢失精度
	switch rv.Kind() {
	case reflect.Bool:
		return val.Bool(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Num(float64(rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Num(float64(rv.Uint()))
	case reflect.Float32, reflect.Float64:
		return val.Num(rv.Float())
	case reflect.String:
		return val.Str(rv.String())
	case reflect.Array, reflect.Slice:
		return valOfSlice(rv, lv)
	case reflect.Map:
		return valOfMap(rv, lv)
	case reflect.Struct:
		return valOfStruct(rv, lv)
	default:
		panic(fmt.Errorf("val: ValOf(nil %v)", rt))
	}
}

func valOfSlice(rv reflect.Value, lv int) *val.Val {
	rt := rv.Type()
	l := rv.Len()
	if l == 0 {
		ty := typeOf(rt, lv)
		return val.List(ty.List(), l)
	}

	vl := valOf(rv.Index(0), lv+1)
	ty := types.List(vl.Type).List()
	lst := val.List(ty, l).List()
	lst.Set(0, vl)

	for i := 1; i < l; i++ {
		tmp := valOf(rv.Index(i), lv+1)
		assertTypeEquals(vl, tmp)
		lst.Set(i, tmp)
	}
	return lst.Vl()
}

func valOfMap(rv reflect.Value, lv int) *val.Val {
	rt := rv.Type()
	keys := rv.MapKeys()
	if len(keys) == 0 {
		ty := typeOf(rt, lv)
		return val.Map(ty.Map()).Map().Vl()
	}

	kVal := valOf(keys[0], lv+1)
	vVal := valOf(rv.MapIndex(keys[0]), lv+1)
	ty := types.Map(kVal.Type, vVal.Type)
	m := val.Map(ty.Map()).Map()
	m.Put(kVal, vVal)

	for i := 1; i < len(keys); i++ {
		tmpKVal := valOf(keys[i], lv+1)
		assertTypeEquals(kVal, tmpKVal)
		tmpVVal := valOf(rv.MapIndex(keys[i]), lv+1)
		assertTypeEquals(vVal, tmpVVal)
		m.Put(tmpKVal, tmpVVal)
	}
	return m.Vl()
}

func valOfStruct(rv reflect.Value, lv int) *val.Val {
	rt := rv.Type()
	if rt.NumField() == 0 {
		ty := typeOf(rt, lv)
		return val.Obj(ty.Obj()).Obj().Vl()
	}

	vs := make([]*val.Val, 0)
	ks := make([]types.Field, 0)
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.IsExported() {
			name, maybe := parseTag(ft)
			v := rv.Field(i)
			var vl *val.Val
			if isNil(v) {
				vl = val.Nothing(typeOf(ft.Type, 0))
			} else {
				vl = valOf(v, lv+1)
				if maybe {
					vl = val.Just(vl.Type, vl)
				}
			}

			vs = append(vs, vl)
			ks = append(ks, types.Field{Name: name, Val: vl.Type})
		}
	}

	ty := types.Obj(ks).Obj()
	obj := val.Obj(ty).Obj()
	obj.V = vs
	return obj.Vl()
}

func isNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func assertTypeEquals(expected, actual *val.Val) {
	if !types.Equals(expected.Type, actual.Type) {
		panic(fmt.Errorf("expect %s (%s) actual %s (%s)", expected.Type, expected, actual.Type, actual))
	}
}
