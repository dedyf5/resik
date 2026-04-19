// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type Response struct {
	Status ResponseStatus `json:"status"`
	Data   any            `json:"data,omitempty"`
	Meta   *ResponseMeta  `json:"meta,omitempty"`
}

type ResponseStatus struct {
	Code    string `json:"code" example:"200.1"`
	Message string `json:"message" example:"OK"`
	Details any    `json:"details,omitempty"`
}

type ResponseSuccess struct {
	Status ResponseStatusWithoutDetails `json:"status"`
	Data   any                          `json:"data,omitempty"`
}

type ResponseSuccessWithMeta struct {
	Status ResponseStatusWithoutDetails `json:"status"`
	Data   any                          `json:"data,omitempty"`
	Meta   *ResponseMeta                `json:"meta,omitempty"`
}

type ResponseErrorWithoutDetails struct {
	Status ResponseStatusWithoutDetails `json:"status"`
}

type ResponseStatusWithoutDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResponseBadRequest struct {
	Status ResponseStatusBadRequest `json:"status"`
}

type ResponseStatusBadRequest struct {
	Code    string                   `json:"code" example:"400.1"`
	Message string                   `json:"message" example:"Bad Request"`
	Details []*errdetails.BadRequest `json:"details,omitempty"`
}

type Log struct {
	Response   Response `json:"response"`
	Message    string   `json:"message"`
	Caller     string   `json:"caller"`
	StackTrace []string `json:"stack_trace"`
}
