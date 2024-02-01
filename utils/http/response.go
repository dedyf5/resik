// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import (
	"fmt"

	"github.com/dedyf5/resik/ctx/status"
	"github.com/dedyf5/resik/entities/response"
)

func ResponseFromStatusHTTP(statusHTTP *status.Status) response.Response {
	if statusHTTP.IsError() {
		return ResponseErrorAuto(statusHTTP)
	}
	return ResponseSuccessAuto(statusHTTP)
}

func ResponseSuccessAuto(statusHTTP *status.Status) response.Response {
	res := response.Response{
		Status: response.Status{
			Code:    fmt.Sprintf("%d.1", statusHTTP.Code),
			Message: statusHTTP.MessageOrDefault(),
		},
		Data: statusHTTP.Data,
		Meta: ResponseMetaFromHTTPMeta(statusHTTP.Meta),
	}
	return res
}

func ResponseErrorAuto(statusHTTP *status.Status) response.Response {
	return response.Response{
		Status: response.Status{
			Code:    fmt.Sprintf("%d.1", statusHTTP.Code),
			Message: statusHTTP.MessageOrDefault(),
			Detail:  statusHTTP.Detail,
		},
	}
}

func ResponseMetaFromHTTPMeta(httpMeta *status.Meta) *response.Meta {
	if httpMeta == nil {
		return nil
	}
	return &response.Meta{
		Total: httpMeta.Total,
		Page:  httpMeta.Page,
		Limit: httpMeta.Limit,
	}
}
