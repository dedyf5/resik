// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

var LimitDefault int = 10
var PageDefault int = 1

type Filter struct {
	Search string
	Page   int
	Limit  int
}

func (f *Filter) PageOrDefault() int {
	if f.Page > 0 {
		return f.Page
	}
	return PageDefault
}

func (f *Filter) LimitOrDefault() int {
	if f.Limit > 0 {
		return f.Limit
	}
	return LimitDefault
}

func (f *Filter) Offset() int {
	page := f.PageOrDefault()
	if page == 1 {
		return 0
	}
	return (page - 1) * f.LimitOrDefault()
}
