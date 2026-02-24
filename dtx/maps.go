package dtx

import (
	"maps"
	"slices"
)

// MapKeys returns a slice of a maps's keys, unlike maps.Keys() which returns an iterator
func MapKeys[K comparable, V any](m map[K]V) (keys []K) {
	return slices.Collect(maps.Keys(m))
}

// MapValues returns a slice of a maps's values, unlike maps.Values() which returns an iterator
func MapValues[V any, K comparable](m map[K]V) (values []V) {
	return slices.Collect(maps.Values(m))
}

func ToLookupMap[K comparable, T any](slice []T, fn func(T) K) (lookup map[K]struct{}) {
	lookup = make(map[K]struct{})
	for _, el := range slice {
		lookup[fn(el)] = struct{}{}
	}
	return lookup
}

func ToMap[T any, K comparable, V any](slice []T, fn func(T) (K, V)) (m map[K]V) {
	m = make(map[K]V)
	for _, el := range slice {
		k, v := fn(el)
		m[k] = v
	}
	return m
}
