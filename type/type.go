package types

import "github.com/goghcrow/yae/util"

type Type int

// 这里其实可以 func (t *Kind) IsXXX() 然后私有化 Type 字段
// 但是使用的地方就不能用 switch, 得写一堆 if k.isXXX(), so 直接用吧, 别改就行了
const (
	TTop Type = iota
	TBottom
	TSlot // type variable

	_Primitive_
	TNum
	TStr
	TBool
	TTime

	_Composite_ // seq
	TTuple      // 内部使用 for 类型推导
	TList
	TMap
	TObj
	TFun
)

func (t Type) IsPrimitive() bool { return t > _Primitive_ && t < _Composite_ }
func (t Type) IsComposite() bool { return t > _Composite_ }

func (t Type) String() string {
	switch t {
	case TNum:
		return "num"
	case TStr:
		return "str"
	case TBool:
		return "bool"
	case TTime:
		return "time"
	case TTuple:
		return "tuple"
	case TList:
		return "list"
	case TMap:
		return "map"
	case TObj:
		return "obj"
	case TFun:
		return "fun"
	case TSlot:
		return "slot"
	case TTop:
		return "⊤"
	case TBottom:
		return "⊥"
	default:
		util.Unreachable()
		return ""
	}
}
