// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
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

// fieldInfo stores pre-computed metadata for a single struct field.
type fieldInfo struct {
	index    int
	name     string
	isBinary bool
}

// fieldCache serves as a global registry for struct metadata to avoid
// expensive reflection calls on every log entry.
var fieldCache sync.Map

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

	cleanReq := maskBinaryFields(h.req)
	reqByte, _ := json.Marshal(cleanReq)

	enc.AddString("app", h.appModule.DirectoryName())
	enc.AddString(CorrelationIDKeyContext.String(), h.log.CorrelationID)
	enc.AddString("path", h.uri)
	enc.AddUint32("status_code", uint32(grpc.Code()))
	enc.AddInt64("elapsed_micro", time.Since(h.start).Microseconds())
	enc.AddString("req", string(reqByte))

	if h.status.IsError() {
		enc.AddString("res", h.status.MessageOrDefault())
	} else {
		cleanRes := maskBinaryFields(h.res)
		resByte, _ := json.Marshal(cleanRes)
		resString := string(resByte)
		if resString == "null" {
			resString = ""
		}
		enc.AddString("res", resString)
	}
	return nil
}

// maskBinaryFields processes a struct (typically a gRPC request or response) and returns
// a map representation where all binary fields ([]byte) are replaced with a descriptive
// string indicating their size.
//
// It performs the following optimizations and features:
//   - Field Name Prioritization: Resolves keys using "json" tag first, then "protobuf"
//     name tag, and finally the struct field name.
//   - Caching: Uses a global sync.Map to cache struct metadata (field indices, names,
//     and types) to minimize reflection overhead on subsequent calls.
//   - Memory Efficiency: Leverages Go 1.24+ iterators (strings.SplitSeq) for zero-allocation
//     tag parsing during initial metadata discovery.
//   - Safety: Gracefully handles pointers, nil values, and unexported fields to prevent
//     runtime panics.
//
// This is primarily used in logging interceptors to prevent large binary payloads
// from bloating logs or causing performance degradation during JSON marshaling.
func maskBinaryFields(data any) any {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return data
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return data
	}

	t := v.Type()

	var infos []fieldInfo
	if val, ok := fieldCache.Load(t); ok {
		infos = val.([]fieldInfo)
	} else {
		infos = []fieldInfo{}
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			structField := t.Field(i)

			if !f.CanInterface() {
				continue
			}

			fieldName := structField.Name
			if jsonTag := structField.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				parts := strings.Split(jsonTag, ",")
				if parts[0] != "" {
					fieldName = parts[0]
				}
			} else if protoTag := structField.Tag.Get("protobuf"); protoTag != "" {
				for part := range strings.SplitSeq(protoTag, ",") {
					if name, found := strings.CutPrefix(part, "name="); found {
						fieldName = name
						break
					}
				}
			}

			isBinary := f.Kind() == reflect.Slice && f.Type().Elem().Kind() == reflect.Uint8

			infos = append(infos, fieldInfo{
				index:    i,
				name:     fieldName,
				isBinary: isBinary,
			})
		}
		fieldCache.Store(t, infos)
	}

	m := make(map[string]any, len(infos))
	for _, info := range infos {
		valField := v.Field(info.index)
		if info.isBinary {
			if !valField.IsNil() {
				m[info.name] = fmt.Sprintf("[BINARY: %d bytes]", valField.Len())
			} else {
				m[info.name] = nil
			}
		} else {
			m[info.name] = valField.Interface()
		}
	}

	return m
}
