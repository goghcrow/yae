package types

func (o *ObjKind) GetField(name string) (*Field, bool) {
	i, ok := o.Index[name]
	if !ok {
		return nil, false
	}
	return &o.Fields[i], true
}
