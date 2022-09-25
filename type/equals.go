package types

// IsSubtype k1 <: k2
//func IsSubtype(k1, k2 *Kind) bool {}

func Equals(x, y *Kind) bool {
	if x == Hole || y == Hole {
		// 给 typecheck 开个后门, 给本地函数定义用, 函数需要自己保证类型函数
		// 这里实际上是 top 类型了
		return true
	}
	if x.Type != y.Type {
		return false
	}

	if x.Type == TMap {
		m1 := x.Map()
		m2 := y.Map()
		return Equals(m1.Key, m2.Key) && Equals(m1.Val, m2.Val)
	}

	if x.Type == TList {
		return Equals(x.List().El, y.List().El)
	}

	// structural type (without sequence of fields)
	if x.Type == TObj {
		return equalsObj(x.Obj(), y.Obj())
	}

	if x.Type == TFun {
		return equalsFun(x.Fun(), y.Fun())
	}

	return true
}

func equalsObj(x *ObjKind, y *ObjKind) bool {
	if len(x.Fields) != len(y.Fields) {
		return false
	}
	for name, f1 := range x.Fields {
		f2, ok := y.Fields[name]
		if !ok || !Equals(f1, f2) {
			return false
		}
	}
	return true
}

func equalsFun(x *FunKind, y *FunKind) bool {
	sz1 := len(x.Param)
	sz2 := len(y.Param)

	if sz1 != sz2 {
		return false
	}

	for i := 0; i < sz1; i++ {
		if Equals(x.Param[i], y.Param[i]) {
			return false
		}
	}

	return Equals(x.Return, y.Return)
}
