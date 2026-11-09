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

// Status is used to wrap errors and response states across all applications, including non-HTTP apps.
type Status struct {
	// HTTP status codes as registered with IANA. See: https://go.dev/src/net/http/status.go
	Code int `json:"code"`

	// Message for UI/UX representation.
	Message string `json:"message"`

	// CauseError stores the underlying error causing this status.
	CauseError error `json:"cause_error"`
	// StackTrace stores the execution stack trace related to the error.
	StackTrace []string `json:"stack_trace"`

	Data    any                    `json:"data"`
	Meta    *Meta                  `json:"meta"`
	Details []protoadapt.MessageV1 `json:"details"`

	Caller string `json:"caller"`
}

// Meta defines basic pagination metadata for a Status response.
type Meta struct {
	Total       int64 `json:"total"`
	PageCurrent int32 `json:"page"`
	Limit       int32 `json:"limit"`
}

// IStatus defines the interface for status objects, providing error handling and message retrieval.
type IStatus interface {
	IsError() bool
	Error() string
	MessageOrDefault() string
}

// StatusHolder wraps a Status pointer.
type StatusHolder struct {
	Status *Status
}

// stackTracer defines the interface for errors that support stack trace extraction.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// getProjectDir retrieves the current working directory of the project.
// It returns an empty string if an error occurs.
func getProjectDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return ""
}

// getErrorFiles extracts and filters file path stack trace frames from an error
// that belong to the current project directory.
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

// NewStatusCode creates a new Status instance with the given HTTP status code.
func NewStatusCode(code int) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:   code,
		Caller: fmt.Sprintf("%s:%d", file, line),
	}
}

// NewStatusMessage creates a new Status instance with a status code, message, and underlying error.
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

// NewStatusData creates a new Status instance with a status code and response data.
func NewStatusData(code int, data any) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:   code,
		Data:   data,
		Caller: fmt.Sprintf("%s:%d", file, line),
	}
}

// NewStatusMessageData creates a new Status instance with a status code, message, data, and underlying error.
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

// NewStatusDataMeta creates a new Status instance with status code, data, and metadata.
func NewStatusDataMeta(code int, data any, meta *Meta) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:   code,
		Data:   data,
		Meta:   meta,
		Caller: fmt.Sprintf("%s:%d", file, line),
	}
}

// NewStatusSuccess creates a new success Status instance with code, message, and data.
func NewStatusSuccess(code int, message string, data any) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:    code,
		Message: message,
		Data:    data,
		Caller:  fmt.Sprintf("%s:%d", file, line),
	}
}

// NewStatusError creates a new error Status instance with code and underlying error.
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

// NewStatusDetails creates a new Status instance with status code, message, and protobuf details.
func NewStatusDetails(code int, message string, details ...protoadapt.MessageV1) *Status {
	_, file, line, _ := runtime.Caller(1)
	return &Status{
		Code:    code,
		Message: message,
		Details: details,
		Caller:  fmt.Sprintf("%s:%d", file, line),
	}
}

// NewStatusBadRequest creates a new Status instance representing a bad request error with field validation details.
func NewStatusBadRequest(langCode, field, message, reason string, err error) *Status {
	return NewStatusDetails(
		http.StatusBadRequest,
		message,
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       field,
					Description: err.Error(),
					Reason:      reason,
					LocalizedMessage: &errdetails.LocalizedMessage{
						Locale:  langCode,
						Message: message,
					},
				},
			},
		},
	)
}

// BadRequests extracts and returns a slice of BadRequest error details from the status details.
func (s *Status) BadRequests() []*errdetails.BadRequest {
	badReqs := make([]*errdetails.BadRequest, 0, len(s.Details))
	for _, detail := range s.Details {
		if badReq, ok := detail.(*errdetails.BadRequest); ok {
			badReqs = append(badReqs, badReq)
		}
	}
	return badReqs
}

// IsError checks if the status code falls within the error range (400 - 599).
func (s *Status) IsError() bool {
	if s.Code >= 400 && s.Code <= 599 {
		return true
	}
	return false
}

// Error returns the status message if it is an error, or an empty string otherwise.
func (s *Status) Error() string {
	if s.IsError() {
		return s.Message
	}
	return ""
}

// GRPCStatus converts the Status instance into a gRPC status.Status object.
func (s *Status) GRPCStatus() *status.Status {
	code := HTTPStatusCodeToGRPCCode(s.Code)
	return status.New(code, s.MessageOrDefault())
}

// MessageOrDefault returns the status message if present, or the default HTTP status text for the code.
func (s *Status) MessageOrDefault() string {
	if s.Message != "" {
		return s.Message
	}
	return http.StatusText(s.Code)
}

// CauseErrorMessageOrDefault returns the underlying cause error message if present, or the default message.
func (s *Status) CauseErrorMessageOrDefault() string {
	if s.CauseError != nil {
		return s.CauseError.Error()
	}
	return s.MessageOrDefault()
}

// HTTPStatusCodeToGRPCCode maps an HTTP status code to its corresponding gRPC status code.
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
