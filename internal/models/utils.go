package models

func Map[T any, V any](source []T, converter func(T) V) []V {
	if source == nil {
		return []V{}
	}

	result := make([]V, len(source))
	for i, v := range source {
		result[i] = converter(v)
	}

	return result
}
