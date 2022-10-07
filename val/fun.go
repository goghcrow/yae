package val

func (f *FunVal) Call(args ...*Val) *Val {
	return f.V(args...)
}
