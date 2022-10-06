package val

func (v *ObjVal) GetVal(field string) (*Val, bool) {
	i, ok := v.Kind.Obj().Index[field]
	if !ok {
		return nil, false
	}
	return v.V[i], true
}

func (v *ObjVal) PutVal(field string, val *Val) bool {
	i, ok := v.Kind.Obj().Index[field]
	if !ok {
		return false
	}

	v.V[i] = val
	return true
}
