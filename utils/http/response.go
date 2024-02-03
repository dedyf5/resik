// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"

	statusCtx "github.com/dedyf5/resik/ctx/status"
	"github.com/dedyf5/resik/entities/response"
)

func ResponseFromStatus(status *statusCtx.Status) response.Response {
	if status.IsError() {
		return ResponseErrorAuto(status)
	}
	return ResponseSuccessAuto(status)
}

func ResponseSuccessAuto(status *statusCtx.Status) response.Response {
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

func ResponseErrorAuto(status *statusCtx.Status) response.Response {
	return response.Response{
		Status: response.Status{
			Code:    fmt.Sprintf("%d.1", status.Code),
			Message: status.MessageOrDefault(),
			Detail:  status.Detail,
		},
	}
}

func ResponseMetaFromStatusMeta(statusMeta *statusCtx.Meta) *response.Meta {
	if statusMeta == nil {
		return nil
	}
	return &response.Meta{
		Total: statusMeta.Total,
		Page:  statusMeta.Page,
		Limit: statusMeta.Limit,
	}
}

func LoggerFromStatus(status *statusCtx.Status) response.Log {
	msg := status.MessageOrDefault()
	if err := status.CauseError; err != nil {
		msg = err.Error()
	}
	return response.Log{
		Response: ResponseFromStatus(status),
		Message:  msg,
	}
}

func LoggerErrorAuto(status *statusCtx.Status) response.Log {
	msg := status.MessageOrDefault()
	if err := status.CauseError; err != nil {
		msg = err.Error()
	}
	return response.Log{
		Response: ResponseErrorAuto(status),
		Message:  msg,
	}
}
