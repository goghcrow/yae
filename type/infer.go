package types

import (
	"fmt"
	"github.com/goghcrow/yae/util"
)

func Infer(f *FunKind, args []*Kind) *FunKind {
	fun := Fun(f.Name, []*Kind{Tuple(f.Param)}, f.Return)
	targ := Tuple(args)

	sx := make([]*Kind, len(args))
	for i := 0; i < len(args); i++ {
		sx[i] = Slot("s")
	}
	s := Tuple(sx)
	t := Slot("t")
	psuido := Fun(f.Name, []*Kind{s}, t)

	m := map[string]*Kind{}
	fmt.Println("_________")
	tfn1 := Unify(psuido, fun, m)
	fmt.Println(m)
	if tfn1 == nil {
		return nil
	}

	fmt.Println("_________")
	fmt.Println(tfn1)

	fmt.Println("_________")
	fmt.Println(s)
	targ1 := Subst(s, m)
	fmt.Println(targ1)
	fmt.Println("_________")

	targ2 := Unify(targ1, targ, m)
	fmt.Println(m)
	if targ2 == nil || targ2.Type != TTuple {
		return nil
	}
	fmt.Println(targ2)
	fmt.Println("_________")

	tresult := Subst(t, m)
	if !slotFree(tresult) {
		return nil
	}
	fmt.Println(tresult)

	return Fun(f.Name, targ2.Tuple().Val, tresult).Fun()
}

func slotFree(k *Kind) bool {
	switch k.Type {
	case TNum:
		return true
	case TStr:
		return true
	case TBool:
		return true
	case TTime:
		return true
	case TList:
		return slotFree(k.List().El)
	case TMap:
		return slotFree(k.Map().Key) && slotFree(k.Map().Val)
	case TTuple:
		for _, vk := range k.Tuple().Val {
			if !slotFree(vk) {
				return false
			}
		}
		return true
	case TObj:
		for _, fk := range k.Obj().Fields {
			if !slotFree(fk) {
				return false
			}
		}
		return true
	case TFun:
		for _, param := range k.Fun().Param {
			if !slotFree(param) {
				return false
			}
		}
		return slotFree(k.Fun().Return)
	case TSlot:
		return false
	case TTop:
		return true
	case TBottom:
		return true
	default:
		util.Unreachable()
	}
	return false
}
