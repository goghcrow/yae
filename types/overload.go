package types

import (
	"fmt"
)

type FunKind int

const (
	MonoFun = 1
	PolyFun = 2
)

// OverLoaded 根据参数重命名实现函数重载
func (f *FunTy) OverLoaded() (key string, fk FunKind) {
	if slotFree(f.Ty()) {
		// 单态函数直接根据去除返回值的签名来查找
		return fmt.Sprintf("λ %s %s", f.Name, Tuple(f.Param)), MonoFun
	} else {
		// for 支持 universal quantification
		// 多态函数根据名称+参数个数来查找
		return fmt.Sprintf("∀.λ %s %d", f.Name, len(f.Param)), PolyFun
	}
}
