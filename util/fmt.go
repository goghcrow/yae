package util

import "strconv"

func FmtInt(n int64) string     { return strconv.FormatInt(n, 10) }
func FmtFloat(n float64) string { return strconv.FormatFloat(n, 'f', -1, 64) }
