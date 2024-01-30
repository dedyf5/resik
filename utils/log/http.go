// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HTTP struct {
	http.ResponseWriter
	log        *Log
	start      time.Time
	statusCode int
	method     string
	uri        string
	userAgent  string
}

func NewHTTP(w http.ResponseWriter, log *Log, start time.Time, method, uri, userAgent string) *HTTP {
	return &HTTP{w, log, start, http.StatusOK, method, uri, userAgent}
}

func (h *HTTP) WriteHeader(code int) {
	h.statusCode = code
	h.ResponseWriter.WriteHeader(code)
	h.writeLogger()
}

func (h *HTTP) writeLogger() {
	msg := fmt.Sprintf(
		"%s request to %s completed",
		h.method,
		h.uri,
	)
	if h.statusCode >= http.StatusOK && h.statusCode <= http.StatusIMUsed {
		h.log.Logger.Info(msg, zap.Inline(h))
	} else if h.statusCode >= http.StatusInternalServerError && h.statusCode <= http.StatusNetworkAuthenticationRequired {
		h.log.Logger.Error(msg, zap.Inline(h))
	} else {
		h.log.Logger.Warn(msg, zap.Inline(h))
	}
}

func (h *HTTP) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("app", "http")
	enc.AddString(CorrelationIDKeyContext.String(), h.log.CorrelationID)
	enc.AddString("method", h.method)
	enc.AddString("path", h.uri)
	enc.AddString("user_agent", h.userAgent)
	enc.AddInt("status_code", h.statusCode)
	enc.AddDuration("elapsed_micro", time.Since(h.start))
	return nil
}
