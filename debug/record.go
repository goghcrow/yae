package debug

import (
	"fmt"

	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

type Val struct {
	v   *val.Val
	col int // start with 1
}

type Record struct {
	vs []Val
}

func NewRecord() *Record {
	return &Record{[]Val{}}
}

func (r *Record) Clear() {
	r.vs = []Val{}
}

func (r *Record) Rec(v *val.Val, col int) {
	for _, it := range r.vs {
		if it.col == col {
			r.Rec(v, col+1)
			return
		}
	}
	r.vs = append(r.vs, Val{v, col})
}

func (r *Record) String() string {
	xs := make([]string, len(r.vs))
	for i, v := range r.vs {
		xs[i] = fmt.Sprintf("%s in col %d", v.v, v.col)
	}
	return util.JoinStr(xs, ", ", "[", "]")
}

func (r *Record) Render(src string) string {
	return newRender(src, r).render()
}
