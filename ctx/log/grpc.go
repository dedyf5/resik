// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"encoding/json"
	"net/http"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GRPC struct {
	appModule configEntity.Module
	log       *Log
	start     time.Time
	status    *resPkg.Status
	uri       string
	req       any
	res       any
}

func NewGRPC(appModule configEntity.Module, log *Log, start time.Time, uri string, req any, res any, err error) *GRPC {
	status := &resPkg.Status{
		Code: http.StatusOK,
	}
	if err != nil {
		switch ty := err.(type) {
		case *resPkg.Status:
			status = ty
		}
	}
	return &GRPC{
		appModule: appModule,
		log:       log,
		start:     start,
		status:    status,
		uri:       uri,
		req:       req,
		res:       res,
	}
}

func (h *GRPC) Write() {
	if h.status.Code >= http.StatusOK && h.status.Code <= http.StatusIMUsed {
		h.log.Logger.Info(h.status.MessageOrDefault(), zap.Inline(h))
	} else if h.status.Code >= http.StatusInternalServerError && h.status.Code <= http.StatusNetworkAuthenticationRequired {
		h.log.Logger.Error(h.status.CauseErrorMessageOrDefault(), zap.Inline(h))
	} else {
		h.log.Logger.Warn(h.status.CauseErrorMessageOrDefault(), zap.Inline(h))
	}
}

func (h *GRPC) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	grpc := h.status.GRPCStatus()
	reqByte, _ := json.Marshal(h.req)
	enc.AddString("app", h.appModule.DirectoryName())
	enc.AddString(CorrelationIDKeyContext.String(), h.log.CorrelationID)
	enc.AddString("path", h.uri)
	enc.AddUint32("status_code", uint32(grpc.Code()))
	enc.AddInt64("elapsed_micro", time.Since(h.start).Microseconds())
	enc.AddString("req", string(reqByte))
	if h.status.IsError() {
		enc.AddString("res", h.status.MessageOrDefault())
	} else {
		resByte, _ := json.Marshal(h.res)
		resString := string(resByte)
		if resString == "null" {
			resString = ""
		}
		enc.AddString("res", resString)
	}
	return nil
}
