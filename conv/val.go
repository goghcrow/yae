package conv

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"reflect"
	"time"
)

func ValOf(v interface{}) (vl *val.Val, err error) {
	return valOfRV(reflect.ValueOf(v))
}

func valOfRV(rv reflect.Value) (vl *val.Val, err error) {
	defer util.Recover(&err)
	return valOf(rv, 0), nil
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
		kd := typeOf(rt, lv)
		return val.List(kd.List(), l)
	}

	vl := valOf(rv.Index(0), lv+1)
	kd := types.List(vl.Kind).List()
	lst := val.List(kd, l).List()
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
		kd := typeOf(rt, lv)
		return val.Map(kd.Map()).Map().Vl()
	}

	kVal := valOf(keys[0], lv+1)
	vVal := valOf(rv.MapIndex(keys[0]), lv+1)
	kd := types.Map(kVal.Kind, vVal.Kind)
	m := val.Map(kd.Map()).Map()
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
		kd := typeOf(rt, lv)
		return val.Obj(kd.Obj()).Obj().Vl()
	}

	vs := make([]*val.Val, 0)
	ks := make([]types.Field, 0)
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.IsExported() {
			name := ft.Tag.Get(tagName)
			if name == "" {
				name = ft.Name
			}
			fv := rv.Field(i)
			vl := valOf(fv, lv+1)

			vs = append(vs, vl)
			ks = append(ks, types.Field{Name: name, Val: vl.Kind})
		}
	}

	kd := types.Obj(ks)
	obj := val.Obj(kd.Obj()).Obj()
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
	if !types.Equals(expected.Kind, actual.Kind) {
		panic(fmt.Errorf("expect %s (%s) actual %s (%s)", expected.Kind, expected, actual.Kind, actual))
	}
}
