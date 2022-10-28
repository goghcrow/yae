package ast

import (
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
	"strconv"
)

func Str(s string) *StrExpr {
	v, err := strconv.Unquote(s)
	util.Assert(err == nil, "invalid string literal: %s", s)
	return &StrExpr{s, v}
}
func Num(s string) *NumExpr {
	f, err := parseNum(s)
	util.Assert(err == nil, "invalid num literal %s", s)
	return &NumExpr{s, f}
}
func Time(s string) *TimeExpr {
	ts := timelib.Strtotime(s[1 : len(s)-1]) // attach ast
	// util.Assert(ts != 0, "invalid time literal: %s", lit.Text)
	return &TimeExpr{s, ts}
}

func Var(name string) *IdentExpr                             { return &IdentExpr{name} }
func True() *BoolExpr                                        { return &BoolExpr{token.TRUE, true} }
func False() *BoolExpr                                       { return &BoolExpr{token.FALSE, false} }
func List(elems []Expr) *ListExpr                            { return &ListExpr{elems, nil} }
func Map(pairs []Pair) *MapExpr                              { return &MapExpr{pairs, nil} }
func Obj(fields []Field) *ObjExpr                            { return &ObjExpr{fields, nil} }
func Call(callee Expr, args []Expr) *CallExpr                { return &CallExpr{callee, args, nil, "", -1} }
func Subscript(varExpr Expr, expr Expr) *SubscriptExpr       { return &SubscriptExpr{varExpr, expr, nil} }
func Member(obj Expr, field *IdentExpr) *MemberExpr          { return &MemberExpr{obj, field, nil, -1} } // FieldSelection
func Group(sub Expr) *GroupExpr                              { return &GroupExpr{sub} }
func Unary(name string, expr Expr, prefix bool) *UnaryExpr   { return &UnaryExpr{name, expr, prefix} }
func Tenary(name string, l Expr, m Expr, r Expr) *TenaryExpr { return &TenaryExpr{name, l, m, r} }
func Binary(name string, fixity oper.Fixity, lhs Expr, rhs Expr) *BinaryExpr {
	return &BinaryExpr{name, fixity, lhs, rhs}
}

//func If(cond, then, els Expr) *IfExpr { return &IfExpr{cond, then, els} }
