// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"encoding/json"
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

func NewGRPC(appModule configEntity.Module, log *Log, start time.Time, path string, requestBody any, responseBody any, err error) *GRPC {
	status := &resPkg.Status{
		Code: http.StatusOK,
	}

	if s := statusProto.Extract(responseBody); s != nil {
		status.Message = s.GetMessage()
	}

	if err != nil {
		switch ty := err.(type) {
		case *resPkg.Status:
			status = ty
		}
	}

	return &GRPC{
		appModule:    appModule,
		log:          log,
		start:        start,
		status:       status,
		path:         path,
		requestBody:  requestBody,
		responseBody: responseBody,
	}
}

func (h *GRPC) Write() {
	if h.status.Code >= http.StatusOK && h.status.Code <= http.StatusIMUsed {
		h.log.logger.Info(h.status.MessageOrDefault(), zap.Inline(h))
	} else if h.status.Code >= http.StatusInternalServerError && h.status.Code <= http.StatusNetworkAuthenticationRequired {
		h.log.logger.Error(h.status.CauseErrorMessageOrDefault(), zap.Inline(h))
	} else {
		h.log.logger.Warn(h.status.CauseErrorMessageOrDefault(), zap.Inline(h))
	}
}

func (h *GRPC) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	grpc := h.status.GRPCStatus()

	enc.AddString("module", h.appModule.DirectoryName())
	enc.AddString(CorrelationIDKeyContext.String(), h.log.CorrelationID)
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
