package util

type void struct{}

var null void = struct{}{}

type StrSet map[string]void

func (s StrSet) Add(str string) {
	s[str] = null
}

func (s StrSet) Contains(str string) bool {
	_, ok := s[str]
	return ok
}
