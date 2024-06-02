package utils

func ToPtrIfNotEmpty[T comparable](x T) *T {
	var empty T
	if x == empty {
		return nil
	}
	return &x
}
