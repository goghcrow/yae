package val

import (
	"fmt"
	types "github.com/goghcrow/yae/type"
	"github.com/goghcrow/yae/util"
	"strconv"
)

func (m *MapVal) Get(k *Val) (*Val, bool) {
	v, ok := m.V[k.Key()]
	return v, ok
}

func (m *MapVal) Put(k, v *Val) *MapVal {
	util.Assert(types.Equals(m.Kind.Map().Key, k.Kind),
		"type mismatched, expect `%s` actual `%s`", m.Kind.Map().Key, k)
	m.V[k.Key()] = v
	return m
}

// Key 只有 Primitive(bool,num,str,time) 允许作为 map 的 key
type Key struct {
	tag types.Type
	val string
}

func (k Key) String() string {
	return k.val
}

func (v *Val) Key() Key {
	switch v.Kind.Type {
	case types.TBool:
		return Key{v.Kind.Type, strconv.FormatBool(v.Bool().V)}
	case types.TNum:
		return Key{v.Kind.Type, fmt.Sprintf("%f", v.Num().V)}
	case types.TStr:
		return Key{v.Kind.Type, v.Str().V}
	case types.TTime:
		return Key{v.Kind.Type, v.Time().V.String()}
	default:
		panic(fmt.Errorf("invalid map key type: %s", v.Kind))
	}
}
