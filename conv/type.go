package conv

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"reflect"
	"time"
)

const maxLevel = 100

var typeOfTime = reflect.TypeOf(time.Time{})

var tagName = "yae"

func TypeOf(v interface{}) (k *types.Kind, err error) {
	return typeOfRV(reflect.ValueOf(v))
}

func typeOfRV(rv reflect.Value) (k *types.Kind, err error) {
	vl, err := valOfRV(rv)
	if err == nil {
		return vl.Kind, nil
	}
	err = nil
	defer util.Recover(&err)
	return typeOf(rv.Type(), 0), nil
}

func typeOf(rt reflect.Type, lv int) *types.Kind {
	if lv > maxLevel {
		panic("max nested depth exceeded")
	}

	for rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	if rt == typeOfTime {
		return types.Time
	}

	// 通过 type 无法获取 interface 类型
	// if rt.Kind() == reflect.Interface { }

	switch rt.Kind() {
	case reflect.Bool:
		return types.Bool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return types.Num
	case reflect.String:
		return types.Str
	case reflect.Array, reflect.Slice:
		el := typeOf(rt.Elem(), lv+1)
		return types.List(el)
	case reflect.Map:
		key := typeOf(rt.Key(), lv+1)
		el := typeOf(rt.Elem(), lv+1)
		return types.Map(key, el)
	case reflect.Struct:
		fs := make(map[string]*types.Kind, rt.NumField())
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.IsExported() {
				name := f.Tag.Get(tagName)
				if name == "" {
					name = f.Name
				}
				fk := typeOf(f.Type, lv+1)
				fs[name] = fk
			}
		}
		return types.Obj(fs)
	default:
		panic(fmt.Errorf("val: TypeOf(nil %v)", rt))
	}
}
