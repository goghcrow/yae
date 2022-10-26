package val

import (
	"github.com/goghcrow/yae/types"
	"math"
	"time"
	"unsafe"
)

var (
	True  = &(&BoolVal{Val{types.Bool}, true}).Val
	False = &(&BoolVal{Val{types.Bool}, false}).Val
	//True = (*Val)(unsafe.Pointer(&BoolVal{Val{types.Bool}, true}))
	//False = (*Val)(unsafe.Pointer(&BoolVal{Val{types.Bool}, false}))
)

type Val struct {
	Type *types.Type
}

type (
	BoolVal struct {
		Val
		V bool
	}
	NumVal struct {
		Val
		V float64 // 需要高精度考虑 big.Int / big.Float
	}
	StrVal struct {
		Val
		V string
	}
	TimeVal struct {
		Val
		V time.Time
	}
	ListVal struct {
		Val
		V []*Val
	}
	MapVal struct {
		Val
		V map[Key]*Val
	}
	ObjVal struct {
		Val
		V []*Val
	}
	IFun   func(...*Val) *Val
	FunVal struct {
		Val
		V IFun
		// 惰性求值 for and or 等短路操作符/函数, 实参会被包成 thunk
		Lazy bool
		// 不是闭包, 不需要引用 env
	}
	MaybeVal struct {
		Val
		V *Val
	}
)

func (v *NumVal) IsInt() bool { return v.V == math.Trunc(v.V) }
func (v *NumVal) Int() int64  { return int64(v.V) }

func (v *Val) Bool() *BoolVal   { return (*BoolVal)(unsafe.Pointer(v)) }
func (v *Val) Num() *NumVal     { return (*NumVal)(unsafe.Pointer(v)) }
func (v *Val) Str() *StrVal     { return (*StrVal)(unsafe.Pointer(v)) }
func (v *Val) Time() *TimeVal   { return (*TimeVal)(unsafe.Pointer(v)) }
func (v *Val) List() *ListVal   { return (*ListVal)(unsafe.Pointer(v)) }
func (v *Val) Map() *MapVal     { return (*MapVal)(unsafe.Pointer(v)) }
func (v *Val) Obj() *ObjVal     { return (*ObjVal)(unsafe.Pointer(v)) }
func (v *Val) Fun() *FunVal     { return (*FunVal)(unsafe.Pointer(v)) }
func (v *Val) Maybe() *MaybeVal { return (*MaybeVal)(unsafe.Pointer(v)) }

func (v *BoolVal) Vl() *Val  { return &v.Val }
func (v *NumVal) Vl() *Val   { return &v.Val }
func (v *StrVal) Vl() *Val   { return &v.Val }
func (v *TimeVal) Vl() *Val  { return &v.Val }
func (v *ListVal) Vl() *Val  { return &v.Val }
func (v *MapVal) Vl() *Val   { return &v.Val }
func (v *ObjVal) Vl() *Val   { return &v.Val }
func (v *FunVal) Vl() *Val   { return &v.Val }
func (v *MaybeVal) Vl() *Val { return &v.Val }
