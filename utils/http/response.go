// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"
	"math"

	resPkg "github.com/dedyf5/resik/pkg/response"
)

func ResponseFromStatus(status *resPkg.Status) resPkg.Response {
	if status.IsError() {
		return ResponseErrorAuto(status)
	}
	return ResponseSuccessAuto(status)
}

func ResponseSuccessAuto(status *resPkg.Status) resPkg.Response {
	res := resPkg.Response{
		Status: resPkg.ResponseStatus{
			Code:    fmt.Sprintf("%d.1", status.Code),
			Message: status.MessageOrDefault(),
		},
		Data: status.Data,
		Meta: ResponseMetaFromStatusMeta(status.Meta),
	}
	return res
}

func ResponseErrorAuto(status *resPkg.Status) resPkg.Response {
	return resPkg.Response{
		Status: resPkg.ResponseStatus{
			Code:    fmt.Sprintf("%d.1", status.Code),
			Message: status.MessageOrDefault(),
			Detail:  status.Detail,
		},
	}
}

func ResponseMetaFromStatusMeta(statusMeta *resPkg.Meta) *resPkg.ResponseMeta {
	if statusMeta == nil {
		return nil
	}
	return &resPkg.ResponseMeta{
		Total: statusMeta.Total,
		Page:  Page(statusMeta.Total, statusMeta.Limit, statusMeta.PageCurrent),
		Limit: statusMeta.Limit,
	}
}

func Page(total uint64, limit, pageCurrent int) *resPkg.ResponsePage {
	current := pageCurrent
	if pageCurrent == 0 {
		current = 1
	}
	res := resPkg.ResponsePage{
		First:    1,
		Previous: nil,
		Current:  current,
		Next:     nil,
		Last:     1,
	}
	if total == 0 {
		return &res
	}

	var previous *int = nil
	if current > 1 {
		tmp := current - 1
		previous = &tmp
	}

	last := int(math.Ceil(float64(total) / float64(limit)))
	var next *int = nil
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

	return &resPkg.ResponsePage{
		First:    1,
		Previous: previous,
		Current:  current,
		Next:     next,
		Last:     last,
	}
}

func LoggerFromStatus(status *resPkg.Status) resPkg.Log {
	return resPkg.Log{
		Response: ResponseFromStatus(status),
		Message:  status.CauseErrorMessageOrDefault(),
	}
}

func LoggerErrorAuto(status *resPkg.Status) resPkg.Log {
	return resPkg.Log{
		Response: ResponseErrorAuto(status),
		Message:  status.CauseErrorMessageOrDefault(),
	}
}
