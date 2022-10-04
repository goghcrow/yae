package ast

import (
	"github.com/goghcrow/yae/oper"
	"unsafe"
)

type Expr struct {
	Type NodeType
}

type LiteralExpr struct {
	Expr
	LitType
	Val string
}
type ListExpr struct { // lit
	Expr
	Elems []*Expr
	// 👇🏻 for typecheck and compile
	Kind interface{} // *types.Kind
}
type Pair struct{ Key, Val *Expr }
type MapExpr struct { // lit
	Expr
	Pairs []Pair
	// 👇🏻 for typecheck and compile
	Kind interface{} //*types.Kind
}
type ObjExpr struct { // lit
	Expr
	Fields map[string]*Expr
	// 👇🏻 for typecheck and compile
	Kind interface{} //*types.Kind
}
type IdentExpr struct {
	Expr
	Name string
}
type UnaryExpr struct {
	Expr
	Name   string
	LHS    *Expr
	Prefix bool
}
type BinaryExpr struct {
	Expr
	Name string
	oper.Fixity
	LHS *Expr
	RHS *Expr
}
type TenaryExpr struct {
	Expr
	Name  string
	Left  *Expr
	Mid   *Expr
	Right *Expr
}
type IfExpr struct {
	Expr
	Cond *Expr
	Then *Expr
	Else *Expr
}
type CallExpr struct {
	Expr
	Callee *Expr
	Args   []*Expr
	// 👇🏻 for typecheck and compile
	Resolved string
	Index    int
}
type SubscriptExpr struct {
	Expr
	Var *Expr
	Idx *Expr
}
type MemberExpr struct {
	Expr
	Obj   *Expr
	Field *IdentExpr
}

func (e *Expr) Ident() *IdentExpr         { return (*IdentExpr)(unsafe.Pointer(e)) }
func (e *Expr) Literal() *LiteralExpr     { return (*LiteralExpr)(unsafe.Pointer(e)) }
func (e *Expr) List() *ListExpr           { return (*ListExpr)(unsafe.Pointer(e)) }
func (e *Expr) Map() *MapExpr             { return (*MapExpr)(unsafe.Pointer(e)) }
func (e *Expr) Obj() *ObjExpr             { return (*ObjExpr)(unsafe.Pointer(e)) }
func (e *Expr) Unary() *UnaryExpr         { return (*UnaryExpr)(unsafe.Pointer(e)) }
func (e *Expr) Binary() *BinaryExpr       { return (*BinaryExpr)(unsafe.Pointer(e)) }
func (e *Expr) Tenary() *TenaryExpr       { return (*TenaryExpr)(unsafe.Pointer(e)) }
func (e *Expr) If() *IfExpr               { return (*IfExpr)(unsafe.Pointer(e)) }
func (e *Expr) Call() *CallExpr           { return (*CallExpr)(unsafe.Pointer(e)) }
func (e *Expr) Subscript() *SubscriptExpr { return (*SubscriptExpr)(unsafe.Pointer(e)) }
func (e *Expr) Member() *MemberExpr       { return (*MemberExpr)(unsafe.Pointer(e)) }
