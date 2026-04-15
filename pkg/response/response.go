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
