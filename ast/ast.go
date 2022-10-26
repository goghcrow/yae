package ast

import (
	"fmt"
	"github.com/goghcrow/yae/oper"
)

type Expr interface {
	fmt.Stringer
	_exprMakeIDEHappy()
}

type (
	StrExpr struct {
		Text string
		// 👇🏻 for typecheck and compile
		Val string
	}
	NumExpr struct {
		Text string
		// 👇🏻 for typecheck and compile
		Val float64
	}
	TimeExpr struct {
		Text string
		// 👇🏻 for typecheck and compile
		Val int64
	}
	BoolExpr struct {
		Text string
		// 👇🏻 for typecheck and compile
		Val bool
	}
	ListExpr struct { // lit
		Elems []Expr
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	Pair struct {
		Key, Val Expr
	}
	MapExpr struct { // lit
		Pairs []Pair
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	Field struct {
		Name string
		Val  Expr
	}
	ObjExpr struct { // lit
		Fields []Field // 不用 map 是因为要保持声明顺序
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	IdentExpr struct {
		Name string
	}
	CallExpr struct {
		Callee Expr
		Args   []Expr
		// 👇🏻 for typecheck and compile
		CalleeType interface{} // *types.Type
		Resolved   string
		Index      int
	}
	SubscriptExpr struct {
		Var Expr
		Idx Expr
		// 👇🏻 for typecheck and compile
		VarType interface{} // *types.Type
	}
	MemberExpr struct {
		Obj   Expr
		Field *IdentExpr
		// 👇🏻 for typecheck and compile
		ObjType interface{} // *types.Type
		Index   int
	}
)

// 👇🏻 会被 desugar 处理
type (
	UnaryExpr struct {
		Name   string
		LHS    Expr
		Prefix bool
	}
	BinaryExpr struct {
		Name string
		oper.Fixity
		LHS Expr
		RHS Expr
	}
	TenaryExpr struct {
		Name  string
		Left  Expr
		Mid   Expr
		Right Expr
	}
	GroupExpr struct { // 仅用于 String(), 会被 Desugar 会去掉
		SubExpr Expr
	}
)

func (_ *StrExpr) _exprMakeIDEHappy()       {}
func (_ *NumExpr) _exprMakeIDEHappy()       {}
func (_ *TimeExpr) _exprMakeIDEHappy()      {}
func (_ *BoolExpr) _exprMakeIDEHappy()      {}
func (_ *ListExpr) _exprMakeIDEHappy()      {}
func (_ *MapExpr) _exprMakeIDEHappy()       {}
func (_ *ObjExpr) _exprMakeIDEHappy()       {}
func (_ *IdentExpr) _exprMakeIDEHappy()     {}
func (_ *UnaryExpr) _exprMakeIDEHappy()     {}
func (_ *BinaryExpr) _exprMakeIDEHappy()    {}
func (_ *TenaryExpr) _exprMakeIDEHappy()    {}
func (_ *CallExpr) _exprMakeIDEHappy()      {}
func (_ *SubscriptExpr) _exprMakeIDEHappy() {}
func (_ *MemberExpr) _exprMakeIDEHappy()    {}
func (_ *GroupExpr) _exprMakeIDEHappy()     {}

//type IfExpr struct { Cond Expr;Then Expr;Else Expr }
//func (_ *IfExpr) _exprMakeIDEHappy()        {}
