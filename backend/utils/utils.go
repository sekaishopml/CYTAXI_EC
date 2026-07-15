package utils

func Coalesce[T any](vals ...T) T {
	var zero T
	for _, v := range vals {
		if any(v) != any(zero) {
			return v
		}
	}
	return zero
}

func Ptr[T any](v T) *T {
	return &v
}
