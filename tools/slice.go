package tools

func DeleteAt[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice // Index out of bounds, return original slice
	}
	return append(slice[:index], slice[index+1:]...)
}
