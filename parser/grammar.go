package parser

import (
	"github.com/goghcrow/yae/token"
)

type grammar struct {
	lbps    [token.SIZE]token.BP
	prefixs [token.SIZE]nud
	infixs  [token.SIZE]led
}

// 前缀操作符
func (g *grammar) prefix(t token.Type, f nud) {
	g.prefixs[t] = f
	g.lbps[t] = 0
}

// 不结合中缀操作符
func (g *grammar) infix(t token.Type, f led) {
	g.infixs[t] = f
	g.lbps[t] = t.Bp()
}

// 右结合中缀操作符
func (g *grammar) infixRight(t token.Type, f led) {
	g.infixs[t] = f
	g.lbps[t] = t.Bp()
}

// 左结合中缀操作符
func (g *grammar) infixLeft(t token.Type, f led) {
	g.infixs[t] = f
	g.lbps[t] = t.Bp()
}

// 后缀操作符（可以看成中缀操作符木有右边操作数）
func (g *grammar) postfix(t token.Type, f led) {
	g.infixs[t] = f
	g.lbps[t] = t.Bp()
}
