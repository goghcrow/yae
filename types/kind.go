package types

type Kind int

const (
	KTop Kind = iota
	KBot
	KTyVar

	kPrimitiveBegin
	KNum
	KStr
	KBool
	KTime

	kCompositeBegin
	kTuple // 类型推导函数参数使用
	KList
	KMap
	KObj
	KFun
	KMaybe
)

var kinds = [...]string{
	KTop:   "⊤",
	KBot:   "⊥",
	KTyVar: "typevar",

	KNum:  "num",
	KStr:  "str",
	KBool: "bool",
	KTime: "time",

	kTuple: "Tuple",
	KList:  "list",
	KMap:   "map",
	KObj:   "obj",
	KFun:   "fun",
	KMaybe: "maybe",
}

func (k Kind) IsPrimitive() bool { return k > kPrimitiveBegin && k < kCompositeBegin }
func (k Kind) IsComposite() bool { return k > kCompositeBegin }
func (k Kind) String() string    { return kinds[k] }
