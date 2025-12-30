package flagx

type Flag[T any] struct {
	ptr        *T
	type_      string
	stringFunc func(p T) string
	setFunc    func(s string, p *T) error
}

func (f *Flag[T]) Set(s string) error {
	return f.setFunc(s, f.ptr)
}

func (f Flag[T]) String() string {
	if f.ptr == nil {
		return "nil"
	}
	v := *f.ptr
	return f.stringFunc(v)
}

func (f Flag[T]) Type() string {
	return f.type_
}
