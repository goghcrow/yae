package types

import (
	"unsafe"
)

// 不支持递归类型
// 做成 subtype 的复杂度对于表达式没必要
// bottom 只在空列表和空 map 使用, top 和 unit 暂时没用

//goland:noinspection GoUnusedGlobalVariable
var (
	Top    = &Type{KTop} // ⊤  any,universal
	Bottom = &Type{KBot} // ⊥  ∅,never,nothing
	Unit   = Tuple([]*Type{})

	Str  = &Type{KStr}
	Num  = &Type{KNum}
	Bool = &Type{KBool}
	Time = &Type{KTime}
)

type Type struct {
	Kind
}

type (
	TypeVariable struct {
		Type
		Name string
	}
	BoolTy  struct{ Type }
	NumTy   struct{ Type }
	StrTy   struct{ Type }
	TimeTy  struct{ Type }
	TupleTy struct {
		Type
		Val []*Type
	}
	ListTy struct {
		Type
		El *Type
	}
	MapTy struct {
		Type
		Key *Type
		Val *Type
	}
	Field struct {
		Name string
		Val  *Type
	}
	ObjTy struct {
		Type
		Fields []Field
		Index  map[string]int
	}
	FunTy struct {
		Type
		Name   string
		Param  []*Type
		Return *Type
	}
	MaybeTy struct {
		Type
		Elem *Type
	}
)

func (t *Type) TyVar() *TypeVariable { return (*TypeVariable)(unsafe.Pointer(t)) }
func (t *Type) Bool() *BoolTy        { return (*BoolTy)(unsafe.Pointer(t)) }
func (t *Type) Num() *NumTy          { return (*NumTy)(unsafe.Pointer(t)) }
func (t *Type) Str() *StrTy          { return (*StrTy)(unsafe.Pointer(t)) }
func (t *Type) Time() *TimeTy        { return (*TimeTy)(unsafe.Pointer(t)) }
func (t *Type) Tuple() *TupleTy      { return (*TupleTy)(unsafe.Pointer(t)) }
func (t *Type) List() *ListTy        { return (*ListTy)(unsafe.Pointer(t)) }
func (t *Type) Map() *MapTy          { return (*MapTy)(unsafe.Pointer(t)) }
func (t *Type) Obj() *ObjTy          { return (*ObjTy)(unsafe.Pointer(t)) }
func (t *Type) Fun() *FunTy          { return (*FunTy)(unsafe.Pointer(t)) }
func (t *Type) Maybe() *MaybeTy      { return (*MaybeTy)(unsafe.Pointer(t)) }

func (t *BoolTy) Ty() *Type  { return &t.Type }
func (t *NumTy) Ty() *Type   { return &t.Type }
func (t *StrTy) Ty() *Type   { return &t.Type }
func (t *TimeTy) Ty() *Type  { return &t.Type }
func (t *TupleTy) Ty() *Type { return &t.Type }
func (t *ListTy) Ty() *Type  { return &t.Type }
func (t *MapTy) Ty() *Type   { return &t.Type }
func (t *ObjTy) Ty() *Type   { return &t.Type }
func (t *FunTy) Ty() *Type   { return &t.Type }
func (t *MaybeTy) Ty() *Type { return &t.Type }
