// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"

	"github.com/dedyf5/resik/ctx/status"
)

type Response struct {
	Status ResponseStatus `json:"status"`
	Data   interface{}    `json:"data,omitempty"`
	Meta   *ResponseMeta  `json:"meta,omitempty"`
}

type ResponseStatus struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail,omitempty"`
}

type ResponseMeta struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

func ResponseFromStatusHTTP(statusHTTP *status.Status) Response {
	if statusHTTP.IsError() {
		return ResponseErrorAuto(statusHTTP)
	}
	return ResponseSuccessAuto(statusHTTP)
}

func ResponseSuccessAuto(statusHTTP *status.Status) Response {
	res := Response{
		Status: ResponseStatus{
			Code:    fmt.Sprintf("%d.1", statusHTTP.Code),
			Message: statusHTTP.MessageOrDefault(),
		},
		Data: statusHTTP.Data,
		Meta: ResponseMetaFromHTTPMeta(statusHTTP.Meta),
	}
	return res
}

func ResponseErrorAuto(statusHTTP *status.Status) Response {
	return Response{
		Status: ResponseStatus{
			Code:    fmt.Sprintf("%d.1", statusHTTP.Code),
			Message: statusHTTP.MessageOrDefault(),
			Detail:  statusHTTP.Detail,
		},
	}
}

func ResponseMetaFromHTTPMeta(httpMeta *status.Meta) *ResponseMeta {
	if httpMeta == nil {
		return nil
	}
	return &ResponseMeta{
		Total: httpMeta.Total,
		Page:  httpMeta.Page,
		Limit: httpMeta.Limit,
	}
}
