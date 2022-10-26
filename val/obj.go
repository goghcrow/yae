package val

func (o *ObjVal) Get(field string) (*Val, bool) {
	i, ok := o.Type.Obj().Index[field]
	if !ok {
		return nil, false
	}
	return o.V[i], true
}

func (o *ObjVal) Put(field string, val *Val) bool {
	i, ok := o.Type.Obj().Index[field]
	if !ok {
		return false
	}

	o.V[i] = val
	return true
}
