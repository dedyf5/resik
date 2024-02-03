// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

var LimitDefault uint = 10
var PageDefault uint = 1

type Filter struct {
	Search string
	Page   uint
	Limit  uint
}

func (f *Filter) PageOrDefault() uint {
	if f.Page > 0 {
		return f.Page
	}
	return PageDefault
}

func (f *Filter) LimitOrDefault() uint {
	if f.Limit > 0 {
		return f.Limit
	}
	return LimitDefault
}

func (f *Filter) Offset() uint {
	page := f.PageOrDefault()
	if page == 1 {
		return 0
	}
	return (page - 1) * f.LimitOrDefault()
}
