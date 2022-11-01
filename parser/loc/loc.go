package loc

import (
	"fmt"

	"github.com/goghcrow/yae/util"
)

type Location interface {
	GetLoc() Loc
}

func (l Loc) GetLoc() Loc { return l }

type Loc struct {
	Pos    int // include
	PosEnd int // exclude
	Col    int
	Line   int
}

type DbgCol int // for debug render

var UnknownCol DbgCol = -1
var Unknown = Loc{-1, -1, -1, -1}

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

func (l Loc) Merge(other Loc) Loc {
	util.Assert(other.Pos > l.Pos, "expect right loc")
	l.PosEnd = other.PosEnd
	return l
}

func (l Loc) String() string {
	return fmt.Sprintf("pos %d-%d line %d col %d", l.Pos+1, l.PosEnd+1, l.Line+1, l.Col+1)
}

func (l Loc) Span(runes []rune) string { return string(runes[l.Pos:l.PosEnd]) }
