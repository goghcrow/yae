package ast

import (
	"github.com/goghcrow/yae/token"
)

func Ident(name string) *Expr {
	e := IdentExpr{Expr{IDENT}, name}
	return &e.Expr
}

func Literal(typ LitType, lit string) *Expr {
	e := LiteralExpr{Expr{LITERAL}, typ, lit}
	return &e.Expr
}

func LitTrue() *Expr          { return Literal(LIT_TRUE, token.TRUE.Name()) }
func LitFalse() *Expr         { return Literal(LIT_FALSE, token.FALSE.Name()) }
func LitNum(lit string) *Expr { return Literal(LIT_NUM, lit) }
func LitStr(lit string) *Expr { return Literal(LIT_STR, lit) }

func List(elems []*Expr) *Expr {
	e := ListExpr{Expr{LIST}, elems, nil}
	return &e.Expr
}

func Map(pairs []Pair) *Expr {
	e := MapExpr{Expr{MAP}, pairs, nil}
	return &e.Expr
}

func Obj(fields map[string]*Expr) *Expr {
	e := ObjExpr{Expr{OBJ}, fields, nil}
	return &e.Expr
}

func Unary(t token.Type, expr *Expr, prefix bool) *Expr {
	e := UnaryExpr{Expr{UNARY}, t, expr, prefix}
	return &e.Expr
}

func Binary(t token.Type, lhs *Expr, rhs *Expr) *Expr {
	e := BinaryExpr{Expr{BINARY}, t, lhs, rhs}
	return &e.Expr
}

func Tenary(t token.Type, l *Expr, m *Expr, r *Expr) *Expr {
	e := TenaryExpr{Expr{TENARY}, t, l, m, r}
	return &e.Expr
}

func If(cond, then, els *Expr) *Expr {
	e := IfExpr{Expr{IF}, cond, then, els}
	return &e.Expr
}

func Call(callee *Expr, args []*Expr) *Expr {
	e := CallExpr{Expr{CALL}, callee, args, "", -1}
	return &e.Expr
}

func Subscript(varExpr *Expr, expr *Expr) *Expr {
	e := SubscriptExpr{Expr{SUBSCRIPT}, varExpr, expr}
	return &e.Expr
}

// Member FieldSelection
func Member(obj *Expr, field *IdentExpr) *Expr {
	e := MemberExpr{Expr{MEMBER}, obj, field}
	return &e.Expr
}

// Begin for multi expr, (begin ...)
func Begin(exprs []*Expr) *Expr {
	e := BeginExpr{Expr{BEGIN}, exprs}
	return &e.Expr
}
