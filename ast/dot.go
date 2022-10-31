package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
)

func Dot(e Expr, graphName string) string {
	return newVisualizer().Dot(e, graphName)
}

var newNode = func() func() label {
	n := 0
	return func() label { n++; return label("N" + strconv.Itoa(n)) }
}()

type visualizer struct {
	buf *strings.Builder
}

func newVisualizer() *visualizer {
	return &visualizer{
		buf: &strings.Builder{},
	}
}

type label string

func (v *visualizer) node(label string) label {
	n := newNode()
	v.buf.WriteString(fmt.Sprintf("\t%s [label=%q] ;\n", n, label))
	return n
}

func (v *visualizer) connect(parent, child label) {
	v.buf.WriteString(fmt.Sprintf("\t%s -- %s ;\n", parent, child))
}

func (v *visualizer) Dot(e Expr, graphName string) string {
	v.buf.WriteString(fmt.Sprintf("graph %q {\n", graphName))
	v.buf.WriteString(fmt.Sprintf("\tlabel=%q\n", e.String()))
	v.dot(e)
	v.buf.WriteString("}")
	return v.buf.String()
}

func (v *visualizer) dot(expr Expr) label {
	switch e := expr.(type) {
	case *StrExpr:
		return v.node(fmt.Sprintf("LIT_STR %s", e.Text))
	case *NumExpr:
		return v.node(fmt.Sprintf("LIT_NUM %s", e.Text))
	case *TimeExpr:
		return v.node(fmt.Sprintf("LIT_TIME %s", e.Text))
	case *BoolExpr:
		return v.node(fmt.Sprintf("LIT_BOOL %s", e.Text))
	case *ListExpr:
		n := v.node("<list>")
		v.connect(n, v.node(token.LEFT_BRACKET))
		fst := true
		for _, el := range e.Elems {
			if fst {
				fst = false
			} else {
				v.connect(n, v.node(token.COMMA))
			}
			v.connect(n, v.dot(el))
		}
		v.connect(n, v.node(token.RIGHT_BRACKET))
		return n
	case *MapExpr:
		n := v.node("<map>")
		v.connect(n, v.node(token.LEFT_BRACKET))
		fst := true
		for _, pair := range e.Pairs {
			if fst {
				fst = false
			} else {
				v.connect(n, v.node(token.COMMA))
			}
			p := v.node(token.COLON)
			v.connect(p, v.dot(pair.Key))
			v.connect(p, v.dot(pair.Val))
			v.connect(n, p)
		}
		v.connect(n, v.node(token.RIGHT_BRACKET))
		return n
	case *ObjExpr:
		n := v.node("<obj>")
		v.connect(n, v.node(token.LEFT_BRACE))
		fst := true
		for _, f := range e.Fields {
			if fst {
				fst = false
			} else {
				v.connect(n, v.node(token.COMMA))
			}
			n1 := v.node(token.COLON)
			v.connect(n1, v.node(f.Name))
			v.connect(n1, v.dot(f.Val))
			v.connect(n, n1)
		}
		v.connect(n, v.node(token.RIGHT_BRACE))
		return n
	case *IdentExpr:
		return v.node(e.Name)
	case *UnaryExpr:
		n := v.node("<unary>")
		if e.Prefix {
			v.connect(n, v.dot(e.LHS))
			v.connect(n, v.node(e.Name))
		} else {
			v.connect(n, v.node(e.Name))
			v.connect(n, v.dot(e.LHS))
		}
		return n
	case *BinaryExpr:
		n := v.node(e.Name)
		v.connect(n, v.dot(e.LHS))
		v.connect(n, v.dot(e.RHS))
		return n
	case *TenaryExpr:
		n := v.node(e.Name)
		v.connect(n, v.dot(e.Left))
		v.connect(n, v.dot(e.Mid))
		v.connect(n, v.dot(e.Right))
		return n
	case *CallExpr:
		n := v.node("<call>")
		v.connect(n, v.dot(e.Callee))
		v.connect(n, v.node(token.LEFT_PAREN))
		fst := true
		for _, arg := range e.Args {
			if fst {
				fst = false
			} else {
				v.connect(n, v.node(token.COMMA))
			}
			v.connect(n, v.dot(arg))
		}
		v.connect(n, v.node(token.RIGHT_PAREN))
		return n
	case *SubscriptExpr:
		n := v.node("<subscript>")
		v.connect(n, v.dot(e.Var))
		v.connect(n, v.dot(e.Idx))
		return n
	case *MemberExpr:
		n := v.node(token.DOT)
		v.connect(n, v.dot(e.Obj))
		v.connect(n, v.node(e.Field.Name))
		return n
	case *GroupExpr:
		return v.dot(e.SubExpr)
	//case *IfExpr:
	//	n := v.node("<if>")
	//	v.connect(n, v.dot(e.Cond))
	//	v.connect(n, v.dot(e.Then))
	//	v.connect(n, v.dot(e.Else))
	//	return n
	default:
		util.Unreachable()
		return ""
	}
}
