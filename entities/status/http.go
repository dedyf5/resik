// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package status

import "net/http"

type HTTP struct {
	Code       int               `json:"code"`        // HTTP status codes as registered with IANA. See: https://go.dev/src/net/http/status.go
	Message    string            `json:"message"`     // message for ui/ux
	LogMessage string            `json:"log_message"` // message for engineer
	Data       interface{}       `json:"data"`
	Meta       *HTTPMeta         `json:"meta"`
	Detail     map[string]string `json:"detail"`
}

type HTTPMeta struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

func (h *HTTP) IsError() bool {
	if h.Code >= 400 && h.Code <= 599 {
		return true
	}
	return false
}

func (h *HTTP) Error() string {
	if h.IsError() {
		return h.Message
	}
	return ""
}

func (h *HTTP) MessageOrDefault() string {
	if h.Message != "" {
		return h.Message
	}
	return http.StatusText(h.Code)
}
