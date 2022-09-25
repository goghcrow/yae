package ast

import (
	"github.com/goghcrow/yae/token"
	"unsafe"
)

type NodeType int

const (
	LITERAL NodeType = iota
	IDENT
	LIST

	// MAP
	// OBJ

	UNARY
	BINARY
	TENARY
	IF
	CALL
	SUBSCRIPT
	MEMBER
	//BEGIN
)

type LitType int

//goland:noinspection GoSnakeCaseUsage
const (
	LIT_NULL LitType = iota
	LIT_STR
	LIT_NUM
	LIT_TRUE
	LIT_FALSE
)

type Expr struct {
	Type NodeType
}

type LiteralExpr struct {
	Expr
	LitType
	Val string
}
type ListExpr struct {
	Expr
	Elems []*Expr
}
type IdentExpr struct {
	Expr
	Name string
}
type Operator struct {
	Expr
	Name string
}
type UnaryExpr struct {
	Expr
	token.Type
	LHS    *Expr
	Prefix bool
}
type BinaryExpr struct {
	Expr
	token.Type
	LHS *Expr
	RHS *Expr
}
type TenaryExpr struct {
	Expr
	token.Type
	Left  *Expr
	Mid   *Expr
	Right *Expr
}
type IfExpr struct {
	Expr
	Cond *Expr
	Else *Expr
	Then *Expr
}
type CallExpr struct {
	Expr
	Callee   *Expr
	Args     []*Expr
	Resolved string
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
type BeginExpr struct {
	Expr
	Exprs []*Expr
}

func (e *Expr) Ident() *IdentExpr         { return (*IdentExpr)(unsafe.Pointer(e)) }
func (e *Expr) Literal() *LiteralExpr     { return (*LiteralExpr)(unsafe.Pointer(e)) }
func (e *Expr) List() *ListExpr           { return (*ListExpr)(unsafe.Pointer(e)) }
func (e *Expr) Unary() *UnaryExpr         { return (*UnaryExpr)(unsafe.Pointer(e)) }
func (e *Expr) Binary() *BinaryExpr       { return (*BinaryExpr)(unsafe.Pointer(e)) }
func (e *Expr) Tenary() *TenaryExpr       { return (*TenaryExpr)(unsafe.Pointer(e)) }
func (e *Expr) If() *IfExpr               { return (*IfExpr)(unsafe.Pointer(e)) }
func (e *Expr) Call() *CallExpr           { return (*CallExpr)(unsafe.Pointer(e)) }
func (e *Expr) Subscript() *SubscriptExpr { return (*SubscriptExpr)(unsafe.Pointer(e)) }
func (e *Expr) Member() *MemberExpr       { return (*MemberExpr)(unsafe.Pointer(e)) }
