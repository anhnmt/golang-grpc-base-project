package utils

// ToProto converts a slice of models to a slice of proto
func ToProto[T, V any](list []*T, fn func(*T) *V) []*V {
	result := make([]*V, len(list))

	for i, v := range list {
		result[i] = fn(v)
	}

	return result
}
