// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package meta

import (
	resPkg "github.com/dedyf5/resik/pkg/response"
	httpUtil "github.com/dedyf5/resik/utils/http"
	numUtil "github.com/dedyf5/resik/utils/numbers"
)

func MetaSetup(total uint64, limit, currentPage int) (meta *Meta, err *resPkg.Status) {
	limitRes, err := numUtil.SafeConvert[int32](limit)
	if err != nil {
		return nil, err
	}

	pageRes, err := MetaPageSetup(total, limit, currentPage)
	if err != nil {
		return nil, err
	}

	return &Meta{Total: total, Limit: limitRes, Page: pageRes}, nil
}

func MetaPageSetup(total uint64, limit, currentPage int) (metaPage *MetaPage, err *resPkg.Status) {
	page := httpUtil.Page(total, limit, currentPage)

	first, err := numUtil.SafeConvert[int32](page.First)
	if err != nil {
		return nil, err
	}

	var previous *int32 = nil
	if page.Previous != nil {
		prevRes, err := numUtil.SafeConvert[int32](*page.Previous)
		if err != nil {
			return nil, err
		}
		previous = &prevRes
	}

	current, err := numUtil.SafeConvert[int32](page.Current)
	if err != nil {
		return nil, err
	}

	var next *int32 = nil
	if page.Next != nil {
		nextRes, err := numUtil.SafeConvert[int32](*page.Next)
		if err != nil {
			return nil, err
		}
		next = &nextRes
	}

	last, err := numUtil.SafeConvert[int32](page.Last)
	if err != nil {
		return nil, err
	}

	return &MetaPage{
		First:    first,
		Previous: previous,
		Current:  current,
		Next:     next,
		Last:     last,
	}, nil
}
