package types

type Type int

// 这里其实可以 func (t *Kind) IsXXX() 然后私有化 Type 字段
// 但是使用的地方就不能用 switch, 得写一堆 if k.isXXX(), so 直接用吧, 别改就行了
const (
	TTop Type = iota
	TBottom
	TSlot // type variable

	TNum
	TStr
	TBool
	TTime

	TTuple // 内部使用 for 类型推导
	TList
	TMap
	TObj
	TFun
)

func (t Type) IsPrimitive() bool { return t >= TNum && t <= TTime }
func (t Type) IsComposite() bool { return t >= TTuple }

func (t Type) String() string {
	return [...]string{
		"⊤", "⊥", "slot ",
		"num", "str", "bool", "time", // primitive
		"Tuple", "list", "map", "obj", "fun", // composite
	}[t]
}
