// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package response

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

// use this Status to wrap error across all apps including non http app
type Status struct {
	// HTTP status codes as registered with IANA. See: https://go.dev/src/net/http/status.go
	Code int `json:"code"`

	// message for ui/ux
	Message string `json:"message"`

	// message for engineer
	CauseError error    `json:"cause_error"`
	StackTrace []string `json:"stack_trace"`

	Data    any                    `json:"data"`
	Meta    *Meta                  `json:"meta"`
	Details []protoadapt.MessageV1 `json:"details"`

	Caller string `json:"caller"`
}

type Meta struct {
	Total       int64 `json:"total"`
	PageCurrent int32 `json:"page"`
	Limit       int32 `json:"limit"`
}

type IStatus interface {
	IsError() bool
	Error() string
	MessageOrDefault() string
}

type StatusHolder struct {
	Status *Status
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func getProjectDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return ""
}

func getErrorFiles(err error) []string {
	dir := getProjectDir()

	var files []string
	if st, ok := err.(stackTracer); ok {
		for _, frame := range st.StackTrace() {
			path := fmt.Sprintf("%+s:%d", frame, frame)
			if strings.Contains(path, dir) {
				files = append(files, path)
			}
		}
	}

	return files
}

func NewStatusCode(code int) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:   code,
		Caller: fmt.Sprintf("%s:%d", file, line),
	}
}

func NewStatusMessage(code int, message string, err error) *Status {
	_, file, line, _ := runtime.Caller(1)

	errStack := errors.WithStack(err)

	return &Status{
		Code:       code,
		Message:    message,
		CauseError: errStack,
		StackTrace: getErrorFiles(errStack),
		Caller:     fmt.Sprintf("%s:%d", file, line),
	}
}

func NewStatusData(code int, data any) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:   code,
		Data:   data,
		Caller: fmt.Sprintf("%s:%d", file, line),
	}
}

func NewStatusMessageData(code int, message string, data any, err error) *Status {
	_, file, line, _ := runtime.Caller(1)
	status := &Status{
		Code:    code,
		Message: message,
		Data:    data,
		Caller:  fmt.Sprintf("%s:%d", file, line),
	}

	if err != nil {
		errStack := errors.WithStack(err)
		status.CauseError = errStack
		status.StackTrace = getErrorFiles(errStack)
	}

	return status
}

func NewStatusDataMeta(code int, data any, meta *Meta) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:   code,
		Data:   data,
		Meta:   meta,
		Caller: fmt.Sprintf("%s:%d", file, line),
	}
}

func NewStatusSuccess(code int, message string, data any) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:    code,
		Message: message,
		Data:    data,
		Caller:  fmt.Sprintf("%s:%d", file, line),
	}
}

func NewStatusError(code int, err error) *Status {
	_, file, line, _ := runtime.Caller(1)

	errStack := errors.WithStack(err)

	return &Status{
		Code:       code,
		CauseError: errStack,
		StackTrace: getErrorFiles(errStack),
		Caller:     fmt.Sprintf("%s:%d", file, line),
	}
}

func NewStatusDetails(code int, message string, details ...protoadapt.MessageV1) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:    code,
		Message: message,
		Details: details,
		Caller:  fmt.Sprintf("%s:%d", file, line),
	}
}

func NewStatusBadRequest(field, message string) *Status {
	return NewStatusDetails(
		http.StatusBadRequest,
		message,
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       field,
					Description: message,
				},
			},
		},
	)
}

func (s *Status) BadRequests() []*errdetails.BadRequest {
	badReqs := make([]*errdetails.BadRequest, 0, len(s.Details))
	for _, detail := range s.Details {
		if badReq, ok := detail.(*errdetails.BadRequest); ok {
			badReqs = append(badReqs, badReq)
		}
	}
	return badReqs
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
	code := HTTPStatusCodeToGRPCCode(s.Code)
	return status.New(code, s.MessageOrDefault())
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

func HTTPStatusCodeToGRPCCode(code int) codes.Code {
	switch code {
	case http.StatusOK, http.StatusCreated:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
