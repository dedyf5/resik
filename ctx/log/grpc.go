// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"encoding/json"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPC struct {
	appModule     configEntity.Module
	log           *Log
	start         time.Time
	statusCode    codes.Code
	statusMessage string
	uri           string
	req           any
	res           any
}

func NewGRPC(appModule configEntity.Module, log *Log, start time.Time, uri string, req any, res any, err error) *GRPC {
	statusCode := codes.OK
	statusMessage := statusCode.String()
	if err != nil {
		if statusResponse, ok := status.FromError(err); ok {
			statusCode = statusResponse.Code()
			statusMessage = statusResponse.Message()
		}
	}
	return &GRPC{
		appModule:     appModule,
		log:           log,
		start:         start,
		statusCode:    statusCode,
		statusMessage: statusMessage,
		uri:           uri,
		req:           req,
		res:           res,
	}
}

func (h *GRPC) Write() {
	if h.statusCode == codes.OK {
		h.log.Logger.Info(h.statusMessage, zap.Inline(h))
		return
	} else if h.statusCode == codes.PermissionDenied || h.statusCode == codes.Unauthenticated {
		h.log.Logger.Warn(h.statusMessage, zap.Inline(h))
		return
	}
	h.log.Logger.Error(h.statusMessage, zap.Inline(h))
}

func (h *GRPC) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	reqByte, _ := json.Marshal(h.req)
	resByte, _ := json.Marshal(h.res)
	resString := string(resByte)
	if resString == "null" {
		resString = ""
	}
	enc.AddString("app", h.appModule.DirectoryName())
	enc.AddString(CorrelationIDKeyContext.String(), h.log.CorrelationID)
	enc.AddString("path", h.uri)
	enc.AddUint32("status_code", uint32(h.statusCode))
	enc.AddInt64("elapsed_micro", time.Since(h.start).Microseconds())
	enc.AddString("req", string(reqByte))
	enc.AddString("res", resString)
	return nil
}
