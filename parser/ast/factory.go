package ast

import (
	"strconv"

	"github.com/goghcrow/yae/parser/loc"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/token"
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/util"
)

func Str(s string, loc loc.Loc) *StrExpr {
	v, err := strconv.Unquote(s)
	util.Assert(err == nil, "invalid string literal: %s", s)
	return &StrExpr{loc, s, v}
}
func Num(s string, loc loc.Loc) *NumExpr {
	f, err := parseNum(s)
	util.Assert(err == nil, "invalid num literal %s", s)
	return &NumExpr{loc, s, f}
}
func Time(s string, loc loc.Loc) *TimeExpr {
	ts := timelib.Strtotime(s[1 : len(s)-1]) // attach ast
	// util.Assert(ts != 0, "invalid time literal: %s", lit.Text)
	return &TimeExpr{loc, s, ts}
}
func Var(name string, loc loc.Loc) *IdentExpr  { return &IdentExpr{loc, name} }
func True(loc loc.Loc) *BoolExpr               { return &BoolExpr{loc, token.TRUE, true} }
func False(loc loc.Loc) *BoolExpr              { return &BoolExpr{loc, token.FALSE, false} }
func List(elems []Expr, loc loc.Loc) *ListExpr { return &ListExpr{loc, elems, nil} }
func Map(pairs []Pair, loc loc.Loc) *MapExpr   { return &MapExpr{loc, pairs, nil} }
func Obj(fields []Field, loc loc.Loc) *ObjExpr { return &ObjExpr{loc, fields, nil} }
func Group(sub Expr, loc loc.Loc) *GroupExpr   { return &GroupExpr{loc, sub} }
func Call(callee Expr, args []Expr, col loc.DBGCol, loc loc.Loc) *CallExpr {
	return &CallExpr{Loc: loc, DBGCol: col, Callee: callee, Args: args, Index: -1}
}
func Subscript(varExpr Expr, expr Expr, col loc.DBGCol, loc loc.Loc) *SubscriptExpr {
	return &SubscriptExpr{Loc: loc, DBGCol: col, Var: varExpr, Idx: expr}
}
func Member(obj Expr, field *IdentExpr, col loc.DBGCol, loc loc.Loc) *MemberExpr {
	return &MemberExpr{Loc: loc, DBGCol: col, Obj: obj, Field: field, Index: -1}
}
func Unary(name *IdentExpr, expr Expr, prefix bool, loc loc.Loc) *UnaryExpr {
	return &UnaryExpr{Loc: loc, IdentExpr: name, LHS: expr, Prefix: prefix}
}
func Tenary(name *IdentExpr, l Expr, m Expr, r Expr, loc loc.Loc) *TenaryExpr {
	return &TenaryExpr{Loc: loc, IdentExpr: name, Left: l, Mid: m, Right: r}
}
func Binary(name *IdentExpr, fixity oper.Fixity, lhs Expr, rhs Expr, loc loc.Loc) *BinaryExpr {
	return &BinaryExpr{Loc: loc, IdentExpr: name, Fixity: fixity, LHS: lhs, RHS: rhs}
}

//func If(cond, then, els Expr, loc loc.Loc) *IfExpr { return &IfExpr{loc, cond, then, els} }
