// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package meta

import httpUtil "github.com/dedyf5/resik/utils/http"

func MetaSetup(total uint64, limit, currentPage int) *Meta {
	return &Meta{
		Total: total,
		Limit: int32(limit),
		Page:  MetaPageSetup(total, limit, currentPage),
	}
}

func MetaPageSetup(total uint64, limit, currentPage int) *MetaPage {
	page := httpUtil.Page(total, limit, currentPage)
	var previous *int32 = nil
	if page.Previous != nil {
		tmp := int32(*page.Previous)
		previous = &tmp
	}
	var next *int32 = nil
	if page.Next != nil {
		tmp := int32(*page.Next)
		next = &tmp
	}
	return &MetaPage{
		First:    int32(page.First),
		Previous: previous,
		Current:  int32(page.Current),
		Next:     next,
		Last:     int32(page.Last),
	}
}
