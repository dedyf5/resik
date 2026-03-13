// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import "math"

func ResponseMetaSetup(total int64, limit, pageCurrent int32) *ResponseMeta {
	return &ResponseMeta{
		Total: total,
		Limit: limit,
		Page:  ResponsePageSetup(total, limit, pageCurrent),
	}
}

func ResponsePageSetup(total int64, limit, pageCurrent int32) *ResponsePage {
	current := pageCurrent
	if pageCurrent == 0 {
		current = 1
	}

	res := ResponsePage{
		First:    1,
		Previous: nil,
		Current:  current,
		Next:     nil,
		Last:     1,
	}

	if total == 0 {
		return &res
	}

	var previous *int32 = nil
	if current > 1 {
		tmp := current - 1
		previous = &tmp
	}

	last := int32(math.Ceil(float64(total) / float64(limit)))

	var next *int32 = nil

	if current != last {
		tmp := current + 1
		next = &tmp
	}

	if current > last {
		previous = &last
		next = nil
	}

	if last == 0 {
		last = 1
	}

	if previous != nil {
		if *previous == current {
			previous = nil
		}
	}

	return &ResponsePage{
		First:    1,
		Previous: previous,
		Current:  current,
		Next:     next,
		Last:     last,
	}
}
