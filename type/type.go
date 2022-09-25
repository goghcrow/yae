package types

import "github.com/goghcrow/yae/util"

type Type int

const (
	TTop Type = iota // seq
	TBottom
	THold

	_Primitive_
	TNum
	TStr
	TBool
	TTime

	_Composite_ // seq
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
	case TList:
		return "list"
	case TMap:
		return "map"
	case TObj:
		return "obj"
	case TFun:
		return "fun"
	case THold:
		return ""
	case TTop:
		return "⊤"
	case TBottom:
		return "⊥"
	default:
		util.Unreachable()
	}
	return ""
}
