package conv

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"reflect"
)

func TypeEnvOf(v interface{}) (*types.Env, error) {
	if v == nil {
		return types.NewEnv(), nil
	}

	rv, ok := reflectMap(v)
	if ok {
		return typeEnvOfMap(rv)
	}

	kind, err := TypeOf(v)
	if err != nil {
		return nil, err
	}
	if kind.Type != types.TObj {
		return nil, fmt.Errorf("expect struct type actual %s", reflect.TypeOf(v))
	}
	env := types.NewEnv()
	for name, k := range kind.Obj().Fields {
		env.Put(name, k)
	}
	return env, nil
}

func typeEnvOfMap(rv reflect.Value) (*types.Env, error) {
	keys := rv.MapKeys()
	env := types.NewEnv()
	for i := 0; i < len(keys); i++ {
		name := keys[i].String()
		typ, err := typeOfRV(rv.MapIndex(keys[i]))
		if err != nil {
			return nil, err
		}
		env.Put(name, typ)
	}
	return env, nil
}

func reflectMap(v interface{}) (reflect.Value, bool) {
	rv := reflect.ValueOf(v)
	if isNil(rv) {
		return rv, false
	}
	rt := rv.Type()
	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		rt = rv.Type()
	}
	if rt.Kind() != reflect.Map || rt.Key().Kind() != reflect.String {
		return rv, false
	}
	return rv, true
}
