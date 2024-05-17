// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// use this Status to wrap error across all apps including non http app
type Status struct {
	// HTTP status codes as registered with IANA. See: https://go.dev/src/net/http/status.go
	Code int `json:"code"`

	// message for ui/ux
	Message string `json:"message"`

	// message for engineer
	CauseError error `json:"cause_error"`

	Data   interface{}       `json:"data"`
	Meta   *Meta             `json:"meta"`
	Detail map[string]string `json:"detail"`
}

type Meta struct {
	Total       uint64 `json:"total"`
	PageCurrent int    `json:"page"`
	Limit       int    `json:"limit"`
}

type IStatus interface {
	IsError() bool
	Error() string
	MessageOrDefault() string
}

func (s *Status) IsError() bool {
	if s.Code >= 400 && s.Code <= 599 {
		return true
	}
	return false
}

func (s *Status) Error() string {
	if s.IsError() {
		return s.Message
	}
	return ""
}

func (s *Status) GRPCStatus() *status.Status {
	switch s.Code {
	case http.StatusOK, http.StatusCreated:
		return status.New(codes.OK, s.MessageOrDefault())
	case http.StatusBadRequest:
		return status.New(codes.InvalidArgument, s.MessageOrDefault())
	case http.StatusUnauthorized:
		return status.New(codes.Unauthenticated, s.MessageOrDefault())
	case http.StatusNotFound:
		return status.New(codes.NotFound, s.MessageOrDefault())
	case http.StatusInternalServerError:
		return status.New(codes.Internal, s.MessageOrDefault())
	default:
		return status.New(codes.Unknown, s.MessageOrDefault())
	}
}

func (s *Status) MessageOrDefault() string {
	if s.Message != "" {
		return s.Message
	}
	return http.StatusText(s.Code)
}

func (s *Status) CauseErrorMessageOrDefault() string {
	if s.CauseError != nil {
		return s.CauseError.Error()
	}
	return s.MessageOrDefault()
}
