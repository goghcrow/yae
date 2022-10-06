package types

// IsSubtype k1 <: k2
//func Subtype(k1, k2 *Kind) bool {}

func Equals(x, y *Kind) bool {
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
		return Equals(m1.Key, m2.Key) && Equals(m1.Val, m2.Val)
	}

	if x.Type == TTuple {
		return equalsTuple(x.Tuple(), y.Tuple())
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

	for _, xf := range x.Fields {
		yf, ok := y.GetField(xf.Name)
		if !ok || !Equals(xf.Val, yf.Val) {
			return false
		}
	}
	return true
}

func equalsTuple(x, y *TupleKind) bool {
	xt := x.Tuple()
	yt := y.Tuple()
	if len(xt.Val) != len(yt.Val) {
		return false
	}
	for i := range xt.Val {
		if !Equals(xt.Val[i], yt.Val[i]) {
			return false
		}
	}
	return true
}

func equalsFun(x, y *FunKind) bool {
	if len(x.Param) != len(y.Param) {
		return false
	}
	for i := range x.Param {
		if !Equals(x.Param[i], y.Param[i]) {
			return false
		}
	}
	return Equals(x.Return, y.Return)
}
