package requests

type Former[T any] interface {
	ToForm() T
}

func ConvertToListForm[T1 Former[T2], T2 any](i []T1) []T2 {
	var o []T2
	for _, v := range i {
		o = append(o, v.ToForm())
	}
	return o
}
