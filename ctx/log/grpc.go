// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"bytes"
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

// maskBinaryFields serves as a universal utility to sanitize data payloads for logging.
// It recursively processes various data types (Structs, Maps, Slices, and Buffers)
// to identify and mask binary data, ensuring logs remain concise and human-readable.
//
// Key Features and Optimizations:
//   - Multi-Protocol Support: Seamlessly handles gRPC structs (using reflection/tags)
//     and REST payloads (unmarshaled maps or raw *bytes.Buffer).
//   - Field Name Resolution: Prioritizes "json" tags, followed by "protobuf" name
//     attributes, falling back to original struct field names for consistency.
//   - Performance via Caching: Utilizes a global sync.Map to store struct metadata
//     (field indices and tags), significantly reducing reflection overhead in high-throughput environments.
//   - Smart Buffer Handling: Automatically attempts to unmarshal *bytes.Buffer
//     if it contains JSON; otherwise, it masks the entire buffer as a binary label.
//   - Recursive Sanitization: Deeply traverses nested maps, slices, and pointers
//     to ensure all binary content is identified and masked.
//
// This utility is essential for logging interceptors/middleware to prevent
// large binary payloads (e.g., file uploads, protobuf bytes) from causing
// memory pressure or bloating log storage.
func maskBinaryFields(data any) any {
	if data == nil {
		return nil
	}

	if buf, ok := data.(*bytes.Buffer); ok {
		if buf == nil {
			return nil
		}
		var rawData any
		if err := json.Unmarshal(buf.Bytes(), &rawData); err == nil {
			return maskBinaryFields(rawData)
		}
		return fmt.Sprintf("[BINARY: %d bytes]", buf.Len())
	}

	v := reflect.ValueOf(data)

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return handleStructMasking(v)

	case reflect.Map:
		return handleMapMasking(v)

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return fmt.Sprintf("[BINARY: %d bytes]", v.Len())
		}
		return handleSliceMasking(v)

	default:
		return v.Interface()
	}
}

// handleStructMasking extracts and masks fields from a reflect.Value of kind Struct.
// It leverages a metadata cache to maintain high performance and resolves field
// names according to the framework's tag priority (JSON > Protobuf > FieldName).
func handleStructMasking(v reflect.Value) any {
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
		fieldVal := v.Field(info.index)
		if info.isBinary {
			if !fieldVal.IsNil() {
				m[info.name] = fmt.Sprintf("[BINARY: %d bytes]", fieldVal.Len())
			} else {
				m[info.name] = nil
			}
		} else {
			m[info.name] = maskBinaryFields(fieldVal.Interface())
		}
	}
	return m
}

// handleMapMasking iterates through map keys and recursively applies masking
// to their values. It is primarily used for processing REST request/response
// payloads that have been unmarshaled into map[string]any.
func handleMapMasking(v reflect.Value) any {
	m := make(map[string]any, v.Len())
	for _, key := range v.MapKeys() {
		strKey := fmt.Sprintf("%v", key.Interface())
		m[strKey] = maskBinaryFields(v.MapIndex(key).Interface())
	}
	return m
}

// handleSliceMasking processes slice or array elements. If the slice is a
// byte slice ([]uint8), it returns a summary string. For other types,
// it recursively processes each element to ensure nested binary data is masked.
func handleSliceMasking(v reflect.Value) any {
	s := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		s[i] = maskBinaryFields(v.Index(i).Interface())
	}
	return s
}
