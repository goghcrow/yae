package types

type Kind int

const (
	KTop Kind = iota
	KBot
	KTyVar

	KNum
	KStr
	KBool
	KTime

	// kTuple 类型推导函数参数使用
	kTuple
	KList
	KMap
	KObj
	KFun
	KMaybe
)

func (k Kind) IsPrimitive() bool { return k >= KNum && k <= KTime }
func (k Kind) IsComposite() bool { return k >= kTuple }

func (k Kind) String() string {
	return [...]string{
		"⊤", "⊥", "typevar ",
		"num", "str", "bool", "time", // primitive
		"Tuple", "list", "map", "obj", "fun", "maybe", // composite
	}[k]
}
