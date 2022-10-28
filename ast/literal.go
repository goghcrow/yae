package ast

import (
	"errors"
	"strconv"
	"strings"
)

func parseNum(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return n, nil
	}

	if strings.HasPrefix(s, "0x") {
		n, err := strconv.ParseInt(s[2:], 16, 64)
		if err == nil {
			return float64(n), nil
		}
	}
	if strings.HasPrefix(s, "0b") {
		n, err := strconv.ParseInt(s[2:], 2, 64)
		if err == nil {
			return float64(n), nil
		}
	}
	if strings.HasPrefix(s, "0o") {
		n, err := strconv.ParseInt(s[2:], 8, 64)
		if err == nil {
			return float64(n), nil
		}
	}

	return 0, errors.New("invalid num: " + s)
}
