package g

type Nothing struct{}

func Empty[T any]() T {
	var empty T
	return empty
}

func IsEmpty[C comparable](val C) bool {
	var empty C
	return val == empty
}

func Pointer[T any](val T) *T {
	return &val
}

func TransformSlice[In, Out any](slice []In, tf func(In) Out) []Out {
	result := make([]Out, len(slice))
	for i, in := range slice {
		result[i] = tf(in)
	}
	return result
}

func FilterSlice[T any](slice []T, filter func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, val := range slice {
		if filter(val) {
			result = append(result, val)
		}
	}
	return result
}

func SliceToKeys[K comparable](slice []K) map[K]Nothing {
	keys := make(map[K]Nothing, len(slice))
	for _, key := range slice {
		keys[key] = Empty[Nothing]()
	}
	return keys
}

func KeysToSlice[K comparable](keys map[K]Nothing) []K {
	slice := make([]K, 0, len(keys))
	for key := range keys {
		slice = append(slice, key)
	}
	return slice
}

func MapToSlice[T any, K comparable](mapped map[K]T) []T {
	slice := make([]T, 0, len(mapped))
	for _, val := range mapped {
		slice = append(slice, val)
	}
	return slice
}
