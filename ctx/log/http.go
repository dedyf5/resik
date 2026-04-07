// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HTTP struct {
	http.ResponseWriter
	moduleType   configEntity.ModuleType
	context      context.Context
	log          *Log
	start        time.Time
	statusCode   int
	method       string
	url          *url.URL
	contentType  string
	userAgent    string
	requestBody  []byte
	responseBody *bytes.Buffer
}

func NewHTTP(w http.ResponseWriter, moduleType configEntity.ModuleType, c context.Context, log *Log, start time.Time, method string, url *url.URL, contentType, userAgent string, requestBody []byte) *HTTP {
	var buf bytes.Buffer
	return &HTTP{w, moduleType, c, log, start, http.StatusOK, method, url, contentType, userAgent, requestBody, &buf}
}

func (h *HTTP) WriteHeader(code int) {
	h.statusCode = code
	h.ResponseWriter.WriteHeader(code)
}

func (h *HTTP) Write(buf []byte) (int, error) {
	loggerRes := getLogResponse(buf)
	bodyByte, err := json.Marshal(loggerRes.Response)
	if err != nil {
		panic("error encode new body response error: " + err.Error())
	}
	h.responseBody.Write(bodyByte)
	h.writeLogger(loggerRes)
	return h.ResponseWriter.Write(bodyByte)
}

func (h *HTTP) writeLogger(loggerRes *resPkg.Log) {
	msg := ""
	if loggerRes != nil {
		msg = loggerRes.Message
	}

	logger := h.log.logger

	fields := []zap.Field{zap.Inline(h)}

	caller := ""
	if holder, ok := h.context.Value(KeyCallerHolderContext).(*CallerHolder); ok {
		if holder != nil {
			caller = holder.Caller
		}
	}

	if loggerRes.Caller != "" {
		caller = loggerRes.Caller
	}

	if caller != "" {
		fields = append(fields, zap.String("line", caller))
		logger = logger.WithOptions(zap.WithCaller(false))
	}

	switch {
	case h.statusCode >= http.StatusOK && h.statusCode <= http.StatusIMUsed:
		logger.Info(msg, fields...)
	case h.statusCode >= http.StatusInternalServerError && h.statusCode <= http.StatusNetworkAuthenticationRequired:
		if loggerRes.StackTrace != nil {
			fields = append(fields, zap.Strings("stack_trace", loggerRes.StackTrace))
		}
		logger.Error(msg, fields...)
	default:
		logger.Warn(msg, fields...)
	}
}

func (h *HTTP) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("module", h.moduleType.String())
	enc.AddString(KeyCorrelationIDContext.String(), h.log.CorrelationID)
	enc.AddString("method", h.method)
	enc.AddString("path", h.url.Path)
	enc.AddString("query_string", h.url.RawQuery)
	enc.AddString("content_type", h.contentType)
	enc.AddString("user_agent", h.userAgent)
	enc.AddInt("status_code", h.statusCode)
	enc.AddInt64("elapsed_micro", time.Since(h.start).Microseconds())

	if strings.HasPrefix(h.contentType, "multipart/form-data") {
		enc.AddString("request_body", "[MULTIPART: binary data omitted]")
	} else {
		var rawData any
		if err := json.Unmarshal(h.requestBody, &rawData); err == nil {
			cleanReq := maskBinaryFields(rawData)
			reqByte, err := json.Marshal(cleanReq)
			enc.AddString("request_body", string(reqByte))
			if err != nil {
				enc.AddString("request_body_error", err.Error())
			}
		} else {
			bodyStr := string(h.requestBody)
			if len(bodyStr) > 1000 {
				enc.AddString("request_body", bodyStr[:1000]+"[truncated]")
			} else {
				enc.AddString("request_body", bodyStr)
			}
		}
	}

	res := maskBinaryFields(h.responseBody)
	switch v := res.(type) {
	case string:
		enc.AddString("response_body", v)
	default:
		if err := enc.AddReflected("response_body", v); err != nil {
			enc.AddString("response_body_error", err.Error())
		}
	}

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
