package tools

func ToPtr[T any](value T) *T {
	return &value
}
