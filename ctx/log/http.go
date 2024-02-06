// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	resPkg "github.com/dedyf5/resik/pkg/response"
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
	body       *bytes.Buffer
}

func NewHTTP(w http.ResponseWriter, log *Log, start time.Time, method, uri, userAgent string) *HTTP {
	var buf bytes.Buffer
	return &HTTP{w, log, start, http.StatusOK, method, uri, userAgent, &buf}
}

func (h *HTTP) WriteHeader(code int) {
	h.statusCode = code
	h.ResponseWriter.WriteHeader(code)
}

func (h *HTTP) Write(buf []byte) (int, error) {
	loggerRes := getLogResponse(buf)
	bodyByte, err := json.Marshal(loggerRes.Response)
	if err != nil {
		panic(fmt.Sprintf("error encode new body response error: %s", err.Error()))
	}
	h.body.Write(bodyByte)
	h.writeLogger(loggerRes)
	return h.ResponseWriter.Write(bodyByte)
}

func (h *HTTP) writeLogger(loggerRes *resPkg.Log) {
	msg := ""
	if loggerRes != nil {
		msg = loggerRes.Message
	}
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
	enc.AddString("response", h.body.String())
	return nil
}

func getLogResponse(buf []byte) *resPkg.Log {
	var response resPkg.Log
	err := json.Unmarshal(buf, &response)
	if err != nil {
		return nil
	}
	return &response
}
