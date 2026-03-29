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
	CauseError error    `json:"cause_error"`
	StackTrace []string `json:"stack_trace"`

	Data   any               `json:"data"`
	Meta   *Meta             `json:"meta"`
	Detail map[string]string `json:"detail"`

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

func NewStatusDetail(code int, message string, detail map[string]string) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:    code,
		Message: message,
		Detail:  detail,
		Caller:  fmt.Sprintf("%s:%d", file, line),
	}
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
	case http.StatusServiceUnavailable:
		return status.New(codes.Unavailable, s.MessageOrDefault())
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
