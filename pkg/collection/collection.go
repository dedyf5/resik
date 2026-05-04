// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package collection

func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

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
