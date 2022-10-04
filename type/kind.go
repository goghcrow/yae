package types

import (
	"unsafe"
)

// bottom 只在空列表和空 map 使用, top 没有使用
// 做成 subtype 的复杂度对于表达式没必要

//goland:noinspection GoUnusedGlobalVariable
var (
	Top    = &Kind{TTop}    // ⊤  any,universal
	Bottom = &Kind{TBottom} // ⊥  ∅,never,nothing
	//Unit   = &Kind{} // null :: void

	Str  = &Kind{TStr}
	Num  = &Kind{TNum}
	Bool = &Kind{TBool}
	Time = &Kind{TTime}
)

type Kind struct {
	Type
}

type SlotKind struct {
	Kind
	Name string
}
type BoolKind struct {
	Kind
}
type NumKind struct {
	Kind
}
type StrKind struct {
	Kind
}
type TimeKind struct {
	Kind
}
type TupleKind struct {
	Kind
	Val []*Kind
}
type ListKind struct {
	Kind
	El *Kind
}
type MapKind struct {
	Kind
	Key *Kind
	Val *Kind
}
type ObjKind struct {
	Kind
	Fields map[string]*Kind
}
type FunKind struct {
	Kind
	Name   string
	Param  []*Kind
	Return *Kind
}

func (k *Kind) Slot() *SlotKind   { return (*SlotKind)(unsafe.Pointer(k)) }
func (k *Kind) Bool() *BoolKind   { return (*BoolKind)(unsafe.Pointer(k)) }
func (k *Kind) Num() *NumKind     { return (*NumKind)(unsafe.Pointer(k)) }
func (k *Kind) Str() *StrKind     { return (*StrKind)(unsafe.Pointer(k)) }
func (k *Kind) Time() *TimeKind   { return (*TimeKind)(unsafe.Pointer(k)) }
func (k *Kind) Tuple() *TupleKind { return (*TupleKind)(unsafe.Pointer(k)) }
func (k *Kind) List() *ListKind   { return (*ListKind)(unsafe.Pointer(k)) }
func (k *Kind) Map() *MapKind     { return (*MapKind)(unsafe.Pointer(k)) }
func (k *Kind) Obj() *ObjKind     { return (*ObjKind)(unsafe.Pointer(k)) }
func (k *Kind) Fun() *FunKind     { return (*FunKind)(unsafe.Pointer(k)) }

func (k *BoolKind) Kd() *Kind  { return &k.Kind }
func (k *NumKind) Kd() *Kind   { return &k.Kind }
func (k *StrKind) Kd() *Kind   { return &k.Kind }
func (k *TimeKind) Kd() *Kind  { return &k.Kind }
func (k *TupleKind) Kd() *Kind { return &k.Kind }
func (k *ListKind) Kd() *Kind  { return &k.Kind }
func (k *MapKind) Kd() *Kind   { return &k.Kind }
func (k *ObjKind) Kd() *Kind   { return &k.Kind }
func (k *FunKind) Kd() *Kind   { return &k.Kind }
