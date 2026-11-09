// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package response provides standard HTTP and gRPC response structures,
// status wrappers, metadata handlers, and pagination setup utilities.
package response

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// Response represents a standard API response structure containing status, data, and metadata.
type Response struct {
	Status ResponseStatus `json:"status"`
	Data   any            `json:"data,omitempty"`
	Meta   *ResponseMeta  `json:"meta,omitempty"`
}

// ResponseStatus represents the status portion of an API response, including code, message, and details.
type ResponseStatus struct {
	Code    string `json:"code" example:"200.1"`
	Message string `json:"message" example:"OK"`
	Details any    `json:"details,omitempty"`
}

// ResponseSuccess represents a standard successful API response without extra metadata.
type ResponseSuccess struct {
	Status ResponseStatusWithoutDetails `json:"status"`
	Data   any                          `json:"data,omitempty"`
}

// ResponseSuccessWithMeta represents a successful API response including pagination or other metadata.
type ResponseSuccessWithMeta struct {
	Status ResponseStatusWithoutDetails `json:"status"`
	Data   any                          `json:"data,omitempty"`
	Meta   *ResponseMeta                `json:"meta,omitempty"`
}

// ResponseErrorWithoutDetails represents a standard error API response without details.
type ResponseErrorWithoutDetails struct {
	Status ResponseStatusWithoutDetails `json:"status"`
}

// ResponseStatusWithoutDetails represents a basic status structure with code and message only.
type ResponseStatusWithoutDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResponseStatusBadRequest represents an API error response for bad request errors.
type ResponseBadRequest struct {
	Status ResponseStatusBadRequest `json:"status"`
}

// ResponseStatusBadRequest details the specific field violations or errors for a bad request response.
type ResponseStatusBadRequest struct {
	Code    string                   `json:"code" example:"400.1"`
	Message string                   `json:"message" example:"Bad Request"`
	Details []*errdetails.BadRequest `json:"details,omitempty"`
}

// Log represents a structured logging structure containing the response details, message, caller, and stack trace.
type Log struct {
	Response   Response `json:"response"`
	Message    string   `json:"message"`
	Caller     string   `json:"caller"`
	StackTrace []string `json:"stack_trace"`
}
