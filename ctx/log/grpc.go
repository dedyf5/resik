// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	statusProto "github.com/dedyf5/resik/app/grpc/proto/status"
	configEntity "github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GRPC struct {
	appModule    configEntity.Module
	context      context.Context
	log          *Log
	start        time.Time
	status       *resPkg.Status
	path         string
	requestBody  any
	responseBody any
}

// fieldInfo stores pre-computed metadata for a single struct field to optimize
// subsequent masking operations by avoiding repetitive reflection lookups.
type fieldInfo struct {
	index       int
	name        string
	isBinary    bool
	isSensitive bool
}

// fieldCache serves as a global registry for struct metadata.
var fieldCache sync.Map

// sensitiveFields defines the blacklist of field names that should be masked.
// Names are checked in a case-insensitive manner.
var sensitiveFields = map[string]struct{}{
	"password": {},
	"token":    {},
}

func NewGRPC(appModule configEntity.Module, c context.Context, log *Log, start time.Time, path string, requestBody any, responseBody any, err error) *GRPC {
	status := resPkg.NewStatusCode(http.StatusOK)

	if s := statusProto.Extract(responseBody); s != nil {
		status.Message = s.GetMessage()
	}

	if err != nil {
		if statusErr, ok := errors.AsType[*resPkg.Status](err); ok {
			status = statusErr
		}
	}

	return &GRPC{
		appModule:    appModule,
		context:      c,
		log:          log,
		start:        start,
		status:       status,
		path:         path,
		requestBody:  requestBody,
		responseBody: responseBody,
	}
}

func (h *GRPC) Write() {
	logger := h.log.logger

	fields := []zap.Field{zap.Inline(h)}

	caller := ""
	if holder, ok := h.context.Value(KeyCallerHolderContext).(*CallerHolder); ok {
		caller = *holder.Caller
	}

	if h.status.Caller != "" {
		caller = h.status.Caller
	}

	if caller != "" {
		fields = append(fields, zap.String("line", caller))
		logger = logger.WithOptions(zap.WithCaller(false))
	}

	switch {
	case h.status.Code >= http.StatusOK && h.status.Code <= http.StatusIMUsed:
		logger.Info(h.status.MessageOrDefault(), fields...)
	case h.status.Code >= http.StatusInternalServerError && h.status.Code <= http.StatusNetworkAuthenticationRequired:
		if h.status.StackTrace != nil {
			fields = append(fields, zap.Strings("stack_trace", h.status.StackTrace))
		}
		logger.Error(h.status.CauseErrorMessageOrDefault(), fields...)
	default:
		logger.Warn(h.status.CauseErrorMessageOrDefault(), fields...)
	}
}

func (h *GRPC) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	grpc := h.status.GRPCStatus()

	enc.AddString("module", h.appModule.DirectoryName())
	enc.AddString(KeyCorrelationIDContext.String(), h.log.CorrelationID)
	enc.AddString("path", h.path)
	enc.AddUint32("status_code", uint32(grpc.Code()))
	enc.AddInt64("elapsed_micro", time.Since(h.start).Microseconds())

	cleanReq := maskBinaryFields(h.requestBody)
	reqByte, err := json.Marshal(cleanReq)
	enc.AddString("request_body", string(reqByte))
	if err != nil {
		enc.AddString("request_body_error", err.Error())
	}

	if h.status.IsError() {
		enc.AddString("response_body", h.status.MessageOrDefault())
	} else {
		res := maskBinaryFields(h.responseBody)
		switch v := res.(type) {
		case string:
			enc.AddString("response_body", v)
		default:
			if err := enc.AddReflected("response_body", v); err != nil {
				enc.AddString("response_body_error", err.Error())
			}
		}
	}
	return nil
}
