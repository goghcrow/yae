package conv

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"reflect"
	"strings"
	"time"
)

const maxLevel = 100

var typeOfTime = reflect.TypeOf(time.Time{})

var tagName = "yae"
var tagMaybe = "maybe"

func MustTypeOf(v interface{}) *types.Kind {
	k, err := TypeOf(v)
	if err != nil {
		panic(err)
	}
	return k
}

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
	return typeOf(rv.Type(), 0), err
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
		fs := make([]types.Field, 0)
		for i := 0; i < rt.NumField(); i++ {
			ft := rt.Field(i)
			if ft.IsExported() {
				name, maybe := parseTag(ft)
				fk := typeOf(ft.Type, lv+1)
				if maybe {
					fk = types.Maybe(fk)
				}
				fs = append(fs, types.Field{Name: name, Val: fk})
			}
		}
		return types.Obj(fs)
	default:
		panic(fmt.Errorf("val: TypeOf(nil %v)", rt))
	}
}

func parseTag(r reflect.StructField) (name string, maybe bool) {
	name = r.Name
	xs := strings.Split(r.Tag.Get(tagName), ",")
	if len(xs) > 0 {
		fst := strings.TrimSpace(xs[0])
		if fst != "" {
			name = fst
		}
	}
	if len(xs) > 1 {
		maybe = strings.ToLower(strings.TrimSpace(xs[1])) == tagMaybe
	}
	return
}
