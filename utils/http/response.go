// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"

	resPkg "github.com/dedyf5/resik/pkg/response"
)

func ResponseFromStatus(status *resPkg.Status) resPkg.Response {
	if status.IsError() {
		return ResponseErrorAuto(status)
	}
	return ResponseSuccessAuto(status)
}

func ResponseSuccessAuto(status *resPkg.Status) resPkg.Response {
	var meta *resPkg.ResponseMeta = nil

	if status.Meta != nil {
		meta = resPkg.ResponseMetaSetup(status.Meta.Total, status.Meta.Limit, status.Meta.PageCurrent)
	}

	return resPkg.Response{
		Status: resPkg.ResponseStatus{
			Code:    fmt.Sprintf("%d.1", status.Code),
			Message: status.MessageOrDefault(),
		},
		Data: status.Data,
		Meta: meta,
	}
}

func ResponseErrorAuto(status *resPkg.Status) resPkg.Response {
	var details any = nil
	badRequests := status.BadRequests()
	if len(badRequests) > 0 {
		details = badRequests
	}

	res := resPkg.Response{
		Status: resPkg.ResponseStatus{
			Code:    fmt.Sprintf("%d.1", status.Code),
			Message: status.MessageOrDefault(),
			Details: details,
		},
	}

	if status.Data != nil {
		res.Data = status.Data
	}

	return res
}

func LoggerFromStatus(status *resPkg.Status) resPkg.Log {
	return resPkg.Log{
		Response:   ResponseFromStatus(status),
		Message:    status.CauseErrorMessageOrDefault(),
		Caller:     status.Caller,
		StackTrace: status.StackTrace,
	}
}

func LoggerErrorAuto(status *resPkg.Status) resPkg.Log {
	return resPkg.Log{
		Response:   ResponseErrorAuto(status),
		Message:    status.CauseErrorMessageOrDefault(),
		Caller:     status.Caller,
		StackTrace: status.StackTrace,
	}
}
