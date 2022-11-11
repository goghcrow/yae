package ast

import (
	"strconv"

	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/pos"
	"github.com/goghcrow/yae/parser/token"
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/util"
)

func Str(s string, p pos.Pos) *StrExpr {
	v, err := strconv.Unquote(s)
	util.Assert(err == nil, "invalid string literal: %s", s)
	return &StrExpr{p, s, v}
}
func Num(s string, p pos.Pos) *NumExpr {
	f, err := parseNum(s)
	util.Assert(err == nil, "invalid num literal %s", s)
	return &NumExpr{p, s, f}
}
func Time(s string, p pos.Pos) *TimeExpr {
	ts := timelib.Strtotime(s[1 : len(s)-1]) // attach ast
	// util.Assert(ts != 0, "invalid time literal: %s", lit.Text)
	return &TimeExpr{p, s, ts}
}
func Var(name string, p pos.Pos) *IdentExpr  { return &IdentExpr{p, name} }
func True(p pos.Pos) *BoolExpr               { return &BoolExpr{p, token.TRUE, true} }
func False(p pos.Pos) *BoolExpr              { return &BoolExpr{p, token.FALSE, false} }
func List(elems []Expr, p pos.Pos) *ListExpr { return &ListExpr{p, elems, nil} }
func Map(pairs []Pair, p pos.Pos) *MapExpr   { return &MapExpr{p, pairs, nil} }
func Obj(fields []Field, p pos.Pos) *ObjExpr { return &ObjExpr{p, fields, nil} }
func Group(sub Expr, p pos.Pos) *GroupExpr   { return &GroupExpr{p, sub} }
func Call(callee Expr, args []Expr, col pos.DBGCol, p pos.Pos) *CallExpr {
	return &CallExpr{Pos: p, DBGCol: col, Callee: callee, Args: args, Index: -1}
}
func Subscript(varExpr Expr, expr Expr, col pos.DBGCol, p pos.Pos) *SubscriptExpr {
	return &SubscriptExpr{Pos: p, DBGCol: col, Var: varExpr, Idx: expr}
}
func Member(obj Expr, field *IdentExpr, col pos.DBGCol, p pos.Pos) *MemberExpr {
	return &MemberExpr{Pos: p, DBGCol: col, Obj: obj, Field: field, Index: -1}
}
func Unary(name *IdentExpr, expr Expr, prefix bool, p pos.Pos) *UnaryExpr {
	return &UnaryExpr{Pos: p, IdentExpr: name, LHS: expr, Prefix: prefix}
}
func Tenary(name *IdentExpr, l Expr, m Expr, r Expr, p pos.Pos) *TenaryExpr {
	return &TenaryExpr{Pos: p, IdentExpr: name, Left: l, Mid: m, Right: r}
}
func Binary(name *IdentExpr, fixity oper.Fixity, lhs Expr, rhs Expr, p pos.Pos) *BinaryExpr {
	return &BinaryExpr{Pos: p, IdentExpr: name, Fixity: fixity, LHS: lhs, RHS: rhs}
}

//func If(cond, then, els Expr, p pos.Pos) *IfExpr { return &IfExpr{p, cond, then, els} }
