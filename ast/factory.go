package ast

import (
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
)

func Ident(name string) *IdentExpr                           { return &IdentExpr{name} }
func Str(s string, v string) *StrExpr                        { return &StrExpr{s, v} }
func Num(s string, v float64) *NumExpr                       { return &NumExpr{s, v} }
func Time(s string, v int64) *TimeExpr                       { return &TimeExpr{s, v} }
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
