package test

import (
	"testing"
)

func TestDesugar(t *testing.T) {
	tests := []struct {
		expr     string
		expected string
	}{
		{
			expr:     "a.b.c(1,2)",
			expected: "c(a.b, 1, 2)",
		},
		{
			expr:     `"Hello".len()`,
			expected: `len("Hello")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%v", r)
				}
			}()
			actual := parse(tt.expr).String()
			expected := tt.expected
			if expected != actual {
				t.Errorf("expect `%s` actual `%s`", expected, actual)
			}
		})
	}
}
