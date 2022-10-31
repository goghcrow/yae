package debug

import (
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/goghcrow/yae/util"
)

type render struct {
	src       string
	rec       *Record
	lines     []*strings.Builder
	startCols []int // startCols[i] 表示 lines[i] 的第一个非空列的下标
}

func newRender(src string, rec *Record) *render {
	util.Assert(!strings.Contains(src, "\n"), "debug mode can not contains line-break")
	return &render{
		src:       src,
		rec:       rec,
		lines:     []*strings.Builder{},
		startCols: []int{},
	}
}

func (r *render) render() string {
	r.renderAssertExpr()
	r.sortValues()
	r.renderValues()
	return r.linesToString()
}

func (r *render) renderAssertExpr() {
	// 第一行表达式
	fstLine := &strings.Builder{}
	fstLine.WriteString(r.src)
	r.lines = append(r.lines, fstLine)
	r.startCols = append(r.startCols, 0)

	// 空行
	emptyLine := &strings.Builder{}
	r.lines = append(r.lines, emptyLine)
	r.startCols = append(r.startCols, 0)
}

func (r *render) sortValues() {
	// 倒序, 按起始列从右往左渲染, 根据每行起始列判断是否有空间
	// 如果从左向右需要保存每行 end 列
	sort.SliceStable(r.rec.vs, func(i, j int) bool {
		return r.rec.vs[j].col < r.rec.vs[i].col
	})
}

var splitByLine = regexp.MustCompile("\r\n|\r|\n")

func (r *render) renderValues() {
	vals := r.rec.vs
nextValue:
	for i := 0; i < len(vals); i++ {
		val := vals[i]
		startCol := val.col

		// 跳过 token.UnknownLoc
		if startCol < 1 {
			continue
		}
		// 跳过相同列, 只渲染相等列最后的 val
		if i+1 < len(vals) && vals[i+1].col == startCol {
			continue
		}

		str := val.v.String()
		strs := splitByLine.Split(str, -1)

		// 如果渲染的文本换行, 忽略空白位置, 总是在新行渲染
		endCol := math.MaxInt
		if len(strs) == 1 {
			endCol = startCol + runeCount(str) // exclusive
		}

		// 找空行渲染, 如果非空行, 需要在垂直 startCol 方向填充 |
		for j := 1; j < len(r.lines); j++ {
			if endCol < r.startCols[j] { // 该行 [startCol:endCol] 未被占用
				r.placeString(r.lines[j], str, startCol)
				r.startCols[j] = startCol // 更新该行首个非空列位置
				continue nextValue
			} else { // 无空位置
				r.placeString(r.lines[j], "|", startCol) // 被占用的列替换成 |
				if j > 1 {                               // 第二行是空行, 不放任何值
					r.startCols[j] = startCol + 1 // + 1, 值和|之前不需要空白
				}
			}
		}

		// 新行渲染, 每行的开始位置一致 startCol
		for _, s := range strs {
			newline := &strings.Builder{}
			r.lines = append(r.lines, newline)
			r.placeString(newline, s, startCol)
			r.startCols = append(r.startCols, startCol)
		}
	}
}

func (r *render) linesToString() string {
	xs := make([]string, len(r.lines))
	for i, l := range r.lines {
		xs[i] = l.String()
	}
	return strings.Join(xs, "\n")
}

func (r *render) placeString(line *strings.Builder, str string, col int) {
	for runeCount(line.String()) < col {
		line.WriteByte(' ')
	}
	start := col - 1
	end := start + runeCount(str)
	replace(line, start, end, str)
}

func replace(buf *strings.Builder, start, end int, str string) {
	runes := []rune(buf.String())
	*buf = strings.Builder{}
	if end > len(runes) {
		buf.WriteString(string(runes[0:start]))
		buf.WriteString(str)
	} else {
		buf.WriteString(string(runes[0:start]))
		buf.WriteString(str)
		buf.WriteString(string(runes[end:]))
	}
}

func runeCount(s string) int { return utf8.RuneCountInString(s) }
