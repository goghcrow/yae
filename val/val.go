package val

import (
	types "github.com/goghcrow/yae/type"
	"math"
	"time"
	"unsafe"
)

var (
	//True = (*Val)(unsafe.Pointer(&BoolVal{Val{types.Bool}, true}))
	//False = (*Val)(unsafe.Pointer(&BoolVal{Val{types.Bool}, false}))
	True  = &(&BoolVal{Val{types.Bool}, true}).Val
	False = &(&BoolVal{Val{types.Bool}, false}).Val
)

func init() {

	// 注意不能这么写, 应该写成 &BoolVal
	//t := BoolVal{Val{types.Bool}, true}
	//True = &t.Val
	//t = BoolVal{Val{types.Bool}, false}
	//False = &t.Val

}

type Val struct {
	Kind *types.Kind
}

type BoolVal struct {
	Val
	V bool
}
type NumVal struct {
	Val
	V float64
}
type StrVal struct {
	Val
	V string
}
type TimeVal struct {
	Val
	V time.Time
}
type ListVal struct {
	Val
	V []*Val
}
type MapVal struct {
	Val
	V map[Key]*Val
}
type ObjVal struct {
	Val
	V map[string]*Val
}
type IFun func(...*Val) *Val
type FunVal struct {
	Val
	V    IFun
	Lazy bool // 惰性求值 for and or 等短路操作符/函数, 实参会被包成 thunk
}

func (n *NumVal) IsInt() bool { return n.V == math.Trunc(n.V) }
func (n *NumVal) Int() int64  { return int64(n.V) }

func (v *Val) Bool() *BoolVal { return (*BoolVal)(unsafe.Pointer(v)) }
func (v *Val) Num() *NumVal   { return (*NumVal)(unsafe.Pointer(v)) }
func (v *Val) Str() *StrVal   { return (*StrVal)(unsafe.Pointer(v)) }
func (v *Val) Time() *TimeVal { return (*TimeVal)(unsafe.Pointer(v)) }
func (v *Val) List() *ListVal { return (*ListVal)(unsafe.Pointer(v)) }
func (v *Val) Map() *MapVal   { return (*MapVal)(unsafe.Pointer(v)) }
func (v *Val) Obj() *ObjVal   { return (*ObjVal)(unsafe.Pointer(v)) }
func (v *Val) Fun() *FunVal   { return (*FunVal)(unsafe.Pointer(v)) }

func (v *BoolVal) Vl() *Val { return &v.Val }
func (v *NumVal) Vl() *Val  { return &v.Val }
func (v *StrVal) Vl() *Val  { return &v.Val }
func (v *TimeVal) Vl() *Val { return &v.Val }
func (v *ListVal) Vl() *Val { return &v.Val }
func (v *MapVal) Vl() *Val  { return &v.Val }
func (v *ObjVal) Vl() *Val  { return &v.Val }
func (v *FunVal) Vl() *Val  { return &v.Val }
