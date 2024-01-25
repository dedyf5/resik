// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"

	statusEntity "github.com/dedyf5/resik/entities/status"
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

func ResponseFromStatusHTTP(statusHTTP *statusEntity.HTTP) Response {
	if statusHTTP.IsError() {
		return ResponseErrorAuto(statusHTTP)
	}
	return ResponseSuccessAuto(statusHTTP)
}

func ResponseSuccessAuto(statusHTTP *statusEntity.HTTP) Response {
	res := Response{
		Status: ResponseStatus{
			Code:    fmt.Sprintf("%d001", statusHTTP.Code),
			Message: statusHTTP.MessageOrDefault(),
		},
		Data: statusHTTP.Data,
		Meta: ResponseMetaFromHTTPMeta(statusHTTP.Meta),
	}
	return res
}

func ResponseErrorAuto(statusHTTP *statusEntity.HTTP) Response {
	return Response{
		Status: ResponseStatus{
			Code:    fmt.Sprintf("%d001", statusHTTP.Code),
			Message: statusHTTP.MessageOrDefault(),
			Detail:  statusHTTP.Detail,
		},
	}
}

func ResponseMetaFromHTTPMeta(httpMeta *statusEntity.HTTPMeta) *ResponseMeta {
	if httpMeta == nil {
		return nil
	}
	return &ResponseMeta{
		Total: httpMeta.Total,
		Page:  httpMeta.Page,
		Limit: httpMeta.Limit,
	}
}
