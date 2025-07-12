package utils

func DeleteSliceElement[T comparable](slice []T, element T) []T {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func DeleteSliceElementAtIndex[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice // Index out of bounds, return original slice
	}
	return append(slice[:index], slice[index+1:]...)
}
