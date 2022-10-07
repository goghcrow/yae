package vm

import (
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

const stackInit = 42
const stackGrow = 500

type stack struct {
	pc    int
	stack []*val.Val // 这里 stack slot 简单处理成 *val.Val
	sp    int
}

func newStack() *stack {
	return &stack{
		stack: make([]*val.Val, stackInit),
	}
}

func (s *stack) growStack() {
	n := make([]*val.Val, s.sp+stackGrow)
	copy(n, s.stack)
	s.stack = n
}

func (s *stack) Push(v *val.Val) {
	if s.sp == len(s.stack) {
		s.growStack()
	}
	s.stack[s.sp] = v
	s.sp++
}

func (s *stack) Empty() bool {
	return s.sp == 0
}

func (s *stack) Pop() *val.Val {
	util.Assert(!s.Empty(), "")
	v := s.stack[s.sp-1]
	s.sp--
	return v
}
