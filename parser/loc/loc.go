package loc

import (
	"fmt"

	"github.com/goghcrow/yae/util"
)

type Locatable interface {
	Location() Loc
}

type Loc struct {
	Pos    int // include
	PosEnd int // exclude
	Col    int
	Line   int
}

// Location 自己实现 Locatable 方便内嵌继承
func (l Loc) Location() Loc { return l }

var Unknown = Loc{-1, -1, -1, -1}

type DBGCol int // for debug render
var UnknownCol DBGCol = -1

// Move Cursor
func (l *Loc) Move(r rune) {
	l.Pos++
	if r == '\n' {
		l.Line++
		l.Col = 0
	} else {
		l.Col++
	}
}

func (l Loc) String() string {
	return fmt.Sprintf("pos %d-%d line %d col %d", l.Pos+1, l.PosEnd+1, l.Line+1, l.Col+1)
}

func (l Loc) Span(runes []rune) string { return string(runes[l.Pos:l.PosEnd]) }

func Range(from, to Locatable) Loc {
	l1 := from.Location()
	l2 := to.Location()
	util.Assert(l2.Pos >= l1.Pos, "expect right loc")
	l1.PosEnd = l2.PosEnd
	return l2
}
