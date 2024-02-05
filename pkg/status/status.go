// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package status

import "net/http"

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
	Total uint64 `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
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

func (s *Status) MessageOrDefault() string {
	if s.Message != "" {
		return s.Message
	}
	return http.StatusText(s.Code)
}
