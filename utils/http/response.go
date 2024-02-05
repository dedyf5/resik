// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"

	"github.com/dedyf5/resik/entities/response"
	statusPkg "github.com/dedyf5/resik/pkg/status"
)

func ResponseFromStatus(status *statusPkg.Status) response.Response {
	if status.IsError() {
		return ResponseErrorAuto(status)
	}
	return ResponseSuccessAuto(status)
}

func ResponseSuccessAuto(status *statusPkg.Status) response.Response {
	res := response.Response{
		Status: response.Status{
			Code:    fmt.Sprintf("%d.1", status.Code),
			Message: status.MessageOrDefault(),
		},
		Data: status.Data,
		Meta: ResponseMetaFromStatusMeta(status.Meta),
	}
	return res
}

func ResponseErrorAuto(status *statusPkg.Status) response.Response {
	return response.Response{
		Status: response.Status{
			Code:    fmt.Sprintf("%d.1", status.Code),
			Message: status.MessageOrDefault(),
			Detail:  status.Detail,
		},
	}
}

func ResponseMetaFromStatusMeta(statusMeta *statusPkg.Meta) *response.Meta {
	if statusMeta == nil {
		return nil
	}
	return &response.Meta{
		Total: statusMeta.Total,
		Page:  Page(statusMeta.Total, statusMeta.Limit, statusMeta.PageCurrent),
		Limit: statusMeta.Limit,
	}
}

func Page(total uint64, limit, pageCurrent int) *response.Page {
	current := pageCurrent
	if pageCurrent == 0 {
		current = 1
	}
	res := response.Page{
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

	last := int(int(total) / limit)
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

	return &response.Page{
		First:    1,
		Previous: previous,
		Current:  current,
		Next:     next,
		Last:     last,
	}
}

func LoggerFromStatus(status *statusPkg.Status) response.Log {
	msg := status.MessageOrDefault()
	if err := status.CauseError; err != nil {
		msg = err.Error()
	}
	return response.Log{
		Response: ResponseFromStatus(status),
		Message:  msg,
	}
}

func LoggerErrorAuto(status *statusPkg.Status) response.Log {
	msg := status.MessageOrDefault()
	if err := status.CauseError; err != nil {
		msg = err.Error()
	}
	return response.Log{
		Response: ResponseErrorAuto(status),
		Message:  msg,
	}
}
