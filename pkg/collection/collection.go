// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package collection provides utility functions for manipulating slices
// and collections.
package collection

// Map transforms a slice of type T into a slice of type U by applying
// the provided mapping function to each element.
func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// Unique extracts unique values from a slice of type T using a key extraction
// function, returning a slice of unique keys of type K while preserving order.
func Unique[T any, K comparable](items []T, f func(T) K) []K {
	keys := make(map[K]bool, len(items))
	var list []K
	for _, item := range items {
		val := f(item)
		if !keys[val] {
			keys[val] = true
			list = append(list, val)
		}
	}
	return list
}
