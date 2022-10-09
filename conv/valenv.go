package conv

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/val"
	"reflect"
)

func ValEnvOf(v interface{}) (*val.Env, error) {
	if v == nil {
		return val.NewEnv(), nil
	}

	rv, ok := reflectMap(v)
	if ok {
		return valEnvOfMap(rv)
	}

	vl, err := ValOf(v)
	if err != nil {
		return nil, err
	}
	if vl.Kind.Type != types.TObj {
		return nil, fmt.Errorf("expect struct type actual %s", reflect.TypeOf(v))
	}
	env := val.NewEnv()

	fs := vl.Obj().Kind.Obj().Fields
	for i, ov := range vl.Obj().V {
		env.Put(fs[i].Name, ov)
	}
	return env, nil
}

func valEnvOfMap(rv reflect.Value) (*val.Env, error) {
	keys := rv.MapKeys()
	env := val.NewEnv()
	for i := 0; i < len(keys); i++ {
		name := keys[i].String()
		typ, err := valOfRV(rv.MapIndex(keys[i]))
		if err != nil {
			return nil, err
		}
		env.Put(name, typ)
	}
	return env, nil
}
