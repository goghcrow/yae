package ast

import (
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/pos"
)

// 👇🏻 会被 desugar 处理
type (
	UnaryExpr struct {
		pos.Pos
		*IdentExpr // pos for desugar and debug
		LHS        Expr
		Prefix     bool
	}
	BinaryExpr struct {
		pos.Pos
		*IdentExpr // pos for desugar and debug
		oper.Fixity
		LHS Expr
		RHS Expr
	}
	TenaryExpr struct {
		pos.Pos
		*IdentExpr // pos for desugar and debug
		Left       Expr
		Mid        Expr
		Right      Expr
	}
	GroupExpr struct { // 仅用于 String(), 会被 Desugar 会去掉
		pos.Pos
		SubExpr Expr
	}
)

func (_ *UnaryExpr) isExpr()  {}
func (_ *BinaryExpr) isExpr() {}
func (_ *TenaryExpr) isExpr() {}
func (_ *GroupExpr) isExpr()  {}
