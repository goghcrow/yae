package types

// IsSubtype k1 <: k2
//func Subtype(k1, k2 *Kind) bool {}

func Equals(x, y *Kind) bool {
	return equals(x, y, 0)
}

func equals(x, y *Kind, lv int) bool {
	if lv > 42 {
		// 可以用 set 精确检查 recursive, 这里简化处理
		return true
	}
	if x == nil || y == nil {
		return false
	}
	if x.Type != y.Type {
		return false
	}

	if x.Type == TSlot {
		// 这里可以直接比较是因为每次调用 Slot 都会生成唯一值
		return x.Slot().Name == y.Slot().Name
	}

	if x.Type == TMap {
		m1 := x.Map()
		m2 := y.Map()
		return equals(m1.Key, m2.Key, lv+1) && equals(m1.Val, m2.Val, lv+1)
	}

	if x.Type == TTuple {
		return equalsTuple(x.Tuple(), y.Tuple(), lv)
	}

	if x.Type == TList {
		return equals(x.List().El, y.List().El, lv+1)
	}

	// structural type (without sequence of fields)
	if x.Type == TObj {
		return equalsObj(x.Obj(), y.Obj(), lv)
	}

	if x.Type == TFun {
		return equalsFun(x.Fun(), y.Fun(), lv)
	}

	return true
}

func equalsObj(x *ObjKind, y *ObjKind, lv int) bool {
	if len(x.Fields) != len(y.Fields) {
		return false
	}

	for _, xf := range x.Fields {
		yf, ok := y.GetField(xf.Name)
		if !ok || !equals(xf.Val, yf.Val, lv+1) {
			return false
		}
	}
	return true
}

func equalsTuple(x, y *TupleKind, lv int) bool {
	xt := x.Tuple()
	yt := y.Tuple()
	if len(xt.Val) != len(yt.Val) {
		return false
	}
	for i := range xt.Val {
		if !equals(xt.Val[i], yt.Val[i], lv+1) {
			return false
		}
	}
	return true
}

func equalsFun(x, y *FunKind, lv int) bool {
	if len(x.Param) != len(y.Param) {
		return false
	}
	for i := range x.Param {
		if !equals(x.Param[i], y.Param[i], lv+1) {
			return false
		}
	}
	return equals(x.Return, y.Return, lv+1)
}
