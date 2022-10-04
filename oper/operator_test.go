package oper

import (
	"math/rand"
	"testing"
	"time"
)

func IsOperator(s string) bool {
	return opReg.MatchString(s)
}

func Test_IsOp(t *testing.T) {
	bytes := []rune(operators)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		rand.Shuffle(len(bytes), func(i, j int) { bytes[i], bytes[j] = bytes[j], bytes[i] })
		op := string(bytes)
		t.Log(op)
		if !IsOperator(op) {
			t.Errorf("!IsOperator(%s)", op)
		}
	}

	for _, tt := range []struct {
		s  string
		is bool
	}{
		{":!#$%^&*+./<=>?@\\ˆ|~-", true},

		{"^.^", true},
		{"^o^", false},
		{"a", false},
		{"a+", false},
		{"<-", true},
		{"->", true},
		{"<->", true},
	} {
		t.Run(tt.s, func(t *testing.T) {
			if IsOperator(tt.s) != tt.is {
				t.Errorf("IsOperator(%s) != %t", tt.s, tt.is)
			}
		})
	}
}
