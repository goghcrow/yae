package ast

import (
	"fmt"
	"github.com/goghcrow/yae/token"
	"github.com/goghcrow/yae/util"
	"strconv"
	"strings"
)

func Dot(e *Expr, graphName string) string {
	return newVisualizer().Dot(e, graphName)
}

var newNode = func() func() label {
	n := 0
	return func() label {
		n++
		return label("N" + strconv.Itoa(n))
	}
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

func (v *visualizer) Dot(e *Expr, graphName string) string {
	v.buf.WriteString(fmt.Sprintf("graph %q {\n", graphName))
	v.buf.WriteString(fmt.Sprintf("\tlabel=%q\n", e.String()))
	v.dot(e)
	v.buf.WriteString("}")
	return v.buf.String()
}

func (v *visualizer) dot(e *Expr) label {
	switch e.Type {
	case LITERAL:
		lit := e.Literal()
		return v.node(fmt.Sprintf("%s %s", lit.LitType, lit.Text))
	case IDENT:
		return v.node(e.Ident().Name)
	case LIST:
		l := e.List()
		n := v.node("<list>")
		v.connect(n, v.node(token.LEFT_BRACKET))
		isFst := true
		for _, el := range l.Elems {
			if isFst {
				isFst = false
			} else {
				v.connect(n, v.node(token.COMMA))
			}
			v.connect(n, v.dot(el))
		}
		v.connect(n, v.node(token.RIGHT_BRACKET))
		return n
	case MAP:
		m := e.Map()
		n := v.node("<map>")
		v.connect(n, v.node(token.LEFT_BRACKET))
		isFst := true
		for _, pair := range m.Pairs {
			if isFst {
				isFst = false
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
	case OBJ:
		o := e.Obj()
		n := v.node("<obj>")
		v.connect(n, v.node(token.LEFT_BRACE))
		isFst := true
		for _, f := range o.Fields {
			if isFst {
				isFst = false
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
	case UNARY:
		u := e.Unary()
		n := v.node("<unary>")
		if u.Prefix {
			v.connect(n, v.dot(u.LHS))
			v.connect(n, v.node(u.Name))
		} else {
			v.connect(n, v.node(u.Name))
			v.connect(n, v.dot(u.LHS))
		}
		return n
	case BINARY:
		b := e.Binary()
		n := v.node(b.Name)
		v.connect(n, v.dot(b.LHS))
		v.connect(n, v.dot(b.RHS))
		return n
	case TENARY:
		t := e.Tenary()
		n := v.node(t.Name)
		v.connect(n, v.dot(t.Left))
		v.connect(n, v.dot(t.Mid))
		v.connect(n, v.dot(t.Right))
		return n
	case CALL:
		c := e.Call()
		n := v.node("<call>")
		v.connect(n, v.dot(c.Callee))
		v.connect(n, v.node(token.LEFT_PAREN))
		isFst := true
		for _, arg := range c.Args {
			if isFst {
				isFst = false
			} else {
				v.connect(n, v.node(token.COMMA))
			}
			v.connect(n, v.dot(arg))
		}
		v.connect(n, v.node(token.RIGHT_PAREN))
		return n
	case SUBSCRIPT:
		s := e.Subscript()
		n := v.node("<subscript>")
		v.connect(n, v.dot(s.Var))
		v.connect(n, v.dot(s.Idx))
		return n
	case MEMBER:
		m := e.Member()
		n := v.node(token.DOT)
		v.connect(n, v.dot(m.Obj))
		v.connect(n, v.node(m.Field.Name))
		return n
	case IF:
		iff := e.If()
		n := v.node("<if>")
		v.connect(n, v.dot(iff.Cond))
		v.connect(n, v.dot(iff.Then))
		v.connect(n, v.dot(iff.Else))
		return n
	case GROUP:
		return v.dot(e.Group().SubExpr)
	default:
		util.Unreachable()
		return ""
	}
}
