package val

import (
	"fmt"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"strconv"
)

func (m *MapVal) Get(k *Val) (*Val, bool) {
	v, ok := m.V[k.Key()]
	return v, ok
}

func (m *MapVal) Put(k, v *Val) {
	util.Assert(types.Equals(m.Type.Map().Key, k.Type),
		"type mismatched, expect `%s` actual `%s`", m.Type.Map().Key, k)
	m.V[k.Key()] = v
}

// Key 只有 Primitive(bool,num,str,time) 允许作为 map 的 key
type Key struct {
	tag types.Kind
	val string
}

func (k Key) String() string {
	return k.val
}

func (v *Val) Key() Key {
	switch v.Type.Kind {
	case types.KBool:
		return Key{v.Type.Kind, strconv.FormatBool(v.Bool().V)}
	case types.KNum:
		if v.Num().IsInt() {
			return Key{v.Type.Kind, fmt.Sprintf("%.0f", v.Num().V)}
		} else {
			return Key{v.Type.Kind, fmt.Sprintf("%f", v.Num().V)}
		}
	case types.KStr:
		return Key{v.Type.Kind, fmt.Sprintf("%q", v.Str().V)}
	case types.KTime:
		return Key{v.Type.Kind, fmt.Sprintf("%q", v.Time().V.String())}
	default:
		panic(fmt.Errorf("invalid map key type: %s", v.Type))
	}
}
