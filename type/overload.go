package types

import (
	"fmt"
	"strings"
)

// func (f *FunKind) mangler(m) string {}

// Lookup 根据参数重命名实现函数重载
func (f *FunKind) Lookup() (string, bool) {
	if slotFree(f.Kd()) {
		pt := make([]string, len(f.Param))
		for i, param := range f.Param {
			// pt[i] = param.Type.String()
			pt[i] = param.String()
		}
		return fmt.Sprintf("%s(%s)", f.Name, strings.Join(pt, ",")), true
	} else {
		return fmt.Sprintf("**%s**(%d)", f.Name, len(f.Param)), false
	}
}
