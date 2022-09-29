package types

import (
	"fmt"
)

// Lookup 根据参数重命名实现函数重载
func (f *FunKind) Lookup() (string, bool) {
	if slotFree(f.Kd()) {
		// 单态函数直接根据去除返回值的签名来查找
		return fmt.Sprintf("λ<%s %s>", f.Name, f.Param), true
	} else {
		// 多态函数根据名称+参数个数来查找, ∀ universal quantification
		return fmt.Sprintf("λ<∀ %s %d>", f.Name, len(f.Param)), false
	}
}
