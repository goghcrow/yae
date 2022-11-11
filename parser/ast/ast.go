package ast

import (
	"fmt"

	"github.com/goghcrow/yae/parser/pos"
)

type Expr interface {
	isExpr() // guard method
	pos.Positionable
	fmt.Stringer
}

type (
	StrExpr struct {
		pos.Pos
		Text string
		// 👇🏻 for typecheck and compile
		Val string
	}
	NumExpr struct {
		pos.Pos
		Text string
		// 👇🏻 for typecheck and compile
		Val float64
	}
	TimeExpr struct {
		pos.Pos
		Text string
		// 👇🏻 for typecheck and compile
		Val int64
	}
	BoolExpr struct {
		pos.Pos
		Text string
		// 👇🏻 for typecheck and compile
		Val bool
	}
	ListExpr struct { // lit
		pos.Pos
		Elems []Expr
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	Pair struct {
		Key, Val Expr
	}
	MapExpr struct { // lit
		pos.Pos
		Pairs []Pair
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	Field struct {
		Name string
		Val  Expr
	}
	ObjExpr struct { // lit
		pos.Pos
		Fields []Field // 不用 map 是因为要保持声明顺序
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	IdentExpr struct {
		pos.Pos
		Name string
	}
	CallExpr struct {
		pos.Pos
		pos.DBGCol // for desugar and debug
		Callee     Expr
		Args       []Expr
		// 👇🏻 for typecheck and compile
		CalleeType interface{} // *types.Type
		Resolved   string
		Index      int
	}
	SubscriptExpr struct {
		pos.Pos
		pos.DBGCol // for desugar and debug
		Var        Expr
		Idx        Expr
		// 👇🏻 for typecheck and compile
		VarType interface{} // *types.Type
	}
	MemberExpr struct { // FieldSelection
		pos.Pos
		pos.DBGCol // for desugar and debug
		Obj        Expr
		Field      *IdentExpr
		// 👇🏻 for typecheck and compile
		ObjType interface{} // *types.Type
		Index   int
	}
)

//type IfExpr struct { pos.Pos; Cond Expr;Then Expr;Else Expr }

func (_ *StrExpr) isExpr()       {}
func (_ *NumExpr) isExpr()       {}
func (_ *TimeExpr) isExpr()      {}
func (_ *BoolExpr) isExpr()      {}
func (_ *ListExpr) isExpr()      {}
func (_ *MapExpr) isExpr()       {}
func (_ *ObjExpr) isExpr()       {}
func (_ *IdentExpr) isExpr()     {}
func (_ *CallExpr) isExpr()      {}
func (_ *SubscriptExpr) isExpr() {}
func (_ *MemberExpr) isExpr()    {}
