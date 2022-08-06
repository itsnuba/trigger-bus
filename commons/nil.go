package commons

func NilOrValue[T any](a *T, b T) T {
	if a != nil {
		return *a
	} else {
		return b
	}
}

func SetMapIfNotNil[T any](m map[string]any, k string, v *T) {
	if v != nil {
		m[k] = v
	}
}
