package ast

import (
	"fmt"

	"github.com/goghcrow/yae/loc"
	"github.com/goghcrow/yae/oper"
)

type Expr interface {
	fmt.Stringer
	loc.Location
}

type (
	StrExpr struct {
		loc.Loc
		Text string
		// 👇🏻 for typecheck and compile
		Val string
	}
	NumExpr struct {
		loc.Loc
		Text string
		// 👇🏻 for typecheck and compile
		Val float64
	}
	TimeExpr struct {
		loc.Loc
		Text string
		// 👇🏻 for typecheck and compile
		Val int64
	}
	BoolExpr struct {
		loc.Loc
		Text string
		// 👇🏻 for typecheck and compile
		Val bool
	}
	ListExpr struct { // lit
		loc.Loc
		Elems []Expr
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	Pair struct {
		Key, Val Expr
	}
	MapExpr struct { // lit
		loc.Loc
		Pairs []Pair
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	Field struct {
		Name string
		Val  Expr
	}
	ObjExpr struct { // lit
		loc.Loc
		Fields []Field // 不用 map 是因为要保持声明顺序
		// 👇🏻 for typecheck and compile
		Type interface{} // *types.Type
	}
	IdentExpr struct {
		loc.Loc
		Name string
	}
	CallExpr struct {
		loc.Loc
		loc.DbgCol // for desugar and debug
		Callee     Expr
		Args       []Expr
		// 👇🏻 for typecheck and compile
		CalleeType interface{} // *types.Type
		Resolved   string
		Index      int
	}
	SubscriptExpr struct {
		loc.Loc
		loc.DbgCol // for desugar and debug
		Var        Expr
		Idx        Expr
		// 👇🏻 for typecheck and compile
		VarType interface{} // *types.Type
	}
	MemberExpr struct { // FieldSelection
		loc.Loc
		loc.DbgCol // for desugar and debug
		Obj        Expr
		Field      *IdentExpr
		// 👇🏻 for typecheck and compile
		ObjType interface{} // *types.Type
		Index   int
	}
)

// 👇🏻 会被 desugar 处理
type (
	UnaryExpr struct {
		loc.Loc
		*IdentExpr // loc for desugar and debug
		LHS        Expr
		Prefix     bool
	}
	BinaryExpr struct {
		loc.Loc
		*IdentExpr // loc for desugar and debug
		oper.Fixity
		LHS Expr
		RHS Expr
	}
	TenaryExpr struct {
		loc.Loc
		*IdentExpr // loc for desugar and debug
		Left       Expr
		Mid        Expr
		Right      Expr
	}
	GroupExpr struct { // 仅用于 String(), 会被 Desugar 会去掉
		loc.Loc
		SubExpr Expr
	}
)

//type IfExpr struct { loc.Loc; Cond Expr;Then Expr;Else Expr }
