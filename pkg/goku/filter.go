// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package goku

import "strings"

// LimitDefault is the default limit for pagination when not specified.
var LimitDefault int32 = 10

// PageDefault is the default page number for pagination when not specified.
var PageDefault int32 = 1

// Filter holds search and pagination parameters for queries.
type Filter struct {
	Search string
	page   int
	limit  int
	raw    *raw
}

// raw holds raw 32-bit pagination parameters.
type raw struct {
	page  int32
	limit int32
}

// NewFilter creates a new Filter instance with the given search query, page, and limit.
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

// FulltextSearch converts the search string into a format suitable for full-text search.
//
// It removes leading/trailing whitespace and converts the search string into a boolean mode search query.
// For example, "apple banana" becomes "+apple* +banana*" (each word gets a "+" prefix and "*" suffix to ensure it appears in results).
func (f *Filter) FulltextSeach() string {
	str := strings.TrimSpace(f.Search)
	if str == "" {
		return ""
	}

	words := strings.Fields(str)

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

// PageOrDefault returns the page number, or PageDefault if the page is less than or equal to 0.
func (f *Filter) PageOrDefault() int {
	if f.page > 0 {
		return f.page
	}
	return int(PageDefault)
}

// LimitOrDefault returns the limit, or LimitDefault if the limit is less than or equal to 0.
func (f *Filter) LimitOrDefault() int {
	if f.limit > 0 {
		return f.limit
	}
	return int(LimitDefault)
}

// Offset calculates the query offset based on the current page and limit.
func (f *Filter) Offset() int {
	page := f.PageOrDefault()
	if page == 1 {
		return 0
	}
	return (page - 1) * f.LimitOrDefault()
}

// Raw returns the underlying raw 32-bit pagination structure.
func (f *Filter) Raw() *raw {
	return f.raw
}

// PageOrDefault returns the 32-bit page number, or PageDefault if invalid.
func (i *raw) PageOrDefault() int32 {
	if i.page > 0 {
		return i.page
	}
	return PageDefault
}

// LimitOrDefault returns the 32-bit limit, or LimitDefault if invalid.
func (i *raw) LimitOrDefault() int32 {
	if i.limit > 0 {
		return i.limit
	}
	return LimitDefault
}

// Offset calculates the 32-bit query offset based on the page and limit.
func (i *raw) Offset() int32 {
	page := i.PageOrDefault()
	if page == 1 {
		return 0
	}
	return (page - 1) * i.LimitOrDefault()
}
