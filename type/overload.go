package types

import (
	"fmt"
	"strings"
)

// OverloadName 参数重载
// 移除参数 map list 的类型参数
// 对象可以再考虑下是 structural type 还是 nominal type
func (f *FunKind) OverloadName() string {
	pt := make([]string, len(f.Param))
	for i, param := range f.Param {
		pt[i] = param.Type.String()
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(pt, ","))
}
