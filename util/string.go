package util

import "strings"

func JoinStr(xs []string, sep string) string {
	return strings.Join(xs, sep)
}

// JoinStrEx Strings.Join 带前后缀版本
func JoinStrEx(xs []string, sep, start, end string) string {
	switch len(xs) {
	case 0:
		return start + end
	case 1:
		return start + xs[0] + end
	}
	n := len(start) + len(sep)*(len(xs)-1) + len(end)
	for i := 0; i < len(xs); i++ {
		n += len(xs[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(start)
	b.WriteString(xs[0])
	for _, s := range xs[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	b.WriteString(end)
	return b.String()
}
