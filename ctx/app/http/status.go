// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package http

import "net/http"

type Status struct {
	Code       int               `json:"code"`        // HTTP status codes as registered with IANA. See: https://go.dev/src/net/http/status.go
	Message    string            `json:"message"`     // message for ui/ux
	CauseError error             `json:"cause_error"` // message for engineer
	Data       interface{}       `json:"data"`
	Meta       *StatusMeta       `json:"meta"`
	Detail     map[string]string `json:"detail"`
}

type StatusMeta struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
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
