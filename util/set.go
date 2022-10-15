package util

import (
	"reflect"
)

type void *struct{}

var null void = &struct{}{}

type IntSet map[int]void
type StrSet map[string]void
type PtrSet map[uintptr]void
type PtrPtrSet map[uintptr]map[uintptr]void // PtrPtrSet 注意 uintptr 并不持有引用

func (s IntSet) IsEmpty() bool                 { return len(s) == 0 }
func (s StrSet) IsEmpty() bool                 { return len(s) == 0 }
func (s PtrSet) IsEmpty() bool                 { return len(s) == 0 }
func (s IntSet) Add(i int)                     { s[i] = null }
func (s StrSet) Add(str string)                { s[str] = null }
func (s PtrSet) Add(ptr interface{})           { s[ptrOf(ptr)] = null }
func (s IntSet) Contains(i int) bool           { return s[i] == null }
func (s StrSet) Contains(str string) bool      { return s[str] == null }
func (s PtrSet) Contains(ptr interface{}) bool { return s[ptrOf(ptr)] == null }

func (p PtrPtrSet) Add(ptr1, ptr2 interface{}) {
	if p[ptrOf(ptr1)] == nil {
		p[ptrOf(ptr1)] = map[uintptr]void{ptrOf(ptr2): null}
	} else {
		p[ptrOf(ptr1)][ptrOf(ptr2)] = null
	}
}
func (p PtrPtrSet) Contains(ptr1, ptr2 interface{}) bool {
	if p[ptrOf(ptr1)] == nil {
		return false
	} else {
		return p[ptrOf(ptr1)][ptrOf(ptr2)] == null
	}
}

func ptrOf(ptr interface{}) uintptr { return reflect.ValueOf(ptr).Pointer() }
