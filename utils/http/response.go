// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"

	resPkg "github.com/dedyf5/resik/pkg/response"
)

// ResponseFromStatus convert a *resPkg.Status into a standardized resPkg.Response,
// automatically routing to error or success formatting based on the status condition.
func ResponseFromStatus(status *resPkg.Status) resPkg.Response {
	if status.IsError() {
		return ResponseErrorAuto(status)
	}
	return ResponseSuccessAuto(status)
}

// ResponseSuccessAuto constructs a standardized success resPkg.Response from a given *resPkg.Status,
// including metadata configuration and response payload data if available.
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

// ResponseErrorAuto constructs a standardized error resPkg.Response from a given *resPkg.Status,
// formatting error codes, default messages, validation details (bad requests), and additional data if present.
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

// LoggerFromStatus generates a structured resPkg.Log configuration from a *resPkg.Status,
// incorporating the mapped HTTP response, cause error messages, caller info, and stact trace.
func LoggerFromStatus(status *resPkg.Status) resPkg.Log {
	return resPkg.Log{
		Response:   ResponseFromStatus(status),
		Message:    status.CauseErrorMessageOrDefault(),
		Caller:     status.Caller,
		StackTrace: status.StackTrace,
	}
}

// LoggerErrorAuto generates a structured error-specific resPkg.Log configuration from a *resPkg.Status,
// utilizing the automated error response mapping alongside debugging metadata.
func LoggerErrorAuto(status *resPkg.Status) resPkg.Log {
	return resPkg.Log{
		Response:   ResponseErrorAuto(status),
		Message:    status.CauseErrorMessageOrDefault(),
		Caller:     status.Caller,
		StackTrace: status.StackTrace,
	}
}
