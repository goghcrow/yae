package pos

import (
	"fmt"

	"github.com/goghcrow/yae/util"
)

type Positionable interface {
	Position() Pos
}

type Pos struct {
	Idx    int // include
	IdxEnd int // exclude
	Col    int
	Line   int
}

// Position 自己实现 Positionable 方便内嵌继承
func (p Pos) Position() Pos { return p }

var Unknown = Pos{-1, -1, -1, -1}

type DBGCol int // for debug render
var UnknownCol DBGCol = -1

// Move Cursor
func (p *Pos) Move(r rune) {
	p.Idx++
	if r == '\n' {
		p.Line++
		p.Col = 0
	} else {
		p.Col++
	}
}

func (p Pos) String() string {
	return fmt.Sprintf("pos %d-%d line %d col %d", p.Idx+1, p.IdxEnd+1, p.Line+1, p.Col+1)
}

func (p Pos) Span(runes []rune) string { return string(runes[p.Idx:p.IdxEnd]) }

func Range(from, to Positionable) Pos {
	l1 := from.Position()
	l2 := to.Position()
	util.Assert(l2.Idx >= l1.Idx, "expect right pos")
	l1.IdxEnd = l2.IdxEnd
	return l2
}
