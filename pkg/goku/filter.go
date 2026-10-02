// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

import "strings"

var LimitDefault int32 = 10
var PageDefault int32 = 1

type Filter struct {
	Search string
	page   int
	limit  int
	raw    *raw
}

type raw struct {
	page  int32
	limit int32
}

func NewFilter(search string, page, limit int32) *Filter {
	return &Filter{
		Search: search,
		page:   int(page),
		limit:  int(limit),
		raw: &raw{
			page:  page,
			limit: limit,
		},
	}
}

// FulltextSeach converts the search string into a format suitable for full-text search.
//
// It removes leading/trailing whitespace and converts the search string into a boolean mode search query.
// For example, "apple banana" becomes "+apple* +banana*" (each word gets a "+" prefix and "*" suffix to ensure it appears in results).
func (f *Filter) FulltextSeach() string {
	str := strings.TrimSpace(f.Search)
	if str == "" {
		return ""
	}

	words := strings.Fields(str)
	if len(words) == 0 {
		return ""
	}

	var formattedWords []string
	for _, word := range words {
		cleanWord := strings.Map(func(r rune) rune {
			if strings.ContainsRune("+-*~<>()@\"", r) {
				return -1
			}
			return r
		}, word)

		if cleanWord != "" {
			formattedWords = append(formattedWords, "+"+cleanWord+"*")
		}
	}

	return strings.Join(formattedWords, " ")
}

func (f *Filter) PageOrDefault() int {
	if f.page > 0 {
		return f.page
	}
	return int(PageDefault)
}

func (f *Filter) LimitOrDefault() int {
	if f.limit > 0 {
		return f.limit
	}
	return int(LimitDefault)
}

func (f *Filter) Offset() int {
	page := f.PageOrDefault()
	if page == 1 {
		return 0
	}
	return (page - 1) * f.LimitOrDefault()
}

func (f *Filter) Raw() *raw {
	return f.raw
}

func (i *raw) PageOrDefault() int32 {
	if i.page > 0 {
		return i.page
	}
	return PageDefault
}

func (i *raw) LimitOrDefault() int32 {
	if i.limit > 0 {
		return i.limit
	}
	return LimitDefault
}

func (i *raw) Offset() int32 {
	page := i.PageOrDefault()
	if page == 1 {
		return 0
	}
	return (page - 1) * i.LimitOrDefault()
}
