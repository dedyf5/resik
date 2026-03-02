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

	cleanReq := maskBinaryFields(h.requestBody)
	reqByte, _ := json.Marshal(cleanReq)

	enc.AddString("module", h.appModule.DirectoryName())
	enc.AddString(CorrelationIDKeyContext.String(), h.log.CorrelationID)
	enc.AddString("path", h.path)
	enc.AddUint32("status_code", uint32(grpc.Code()))
	enc.AddInt64("elapsed_micro", time.Since(h.start).Microseconds())
	enc.AddString("request_body", string(reqByte))

	if h.status.IsError() {
		enc.AddString("response_body", h.status.MessageOrDefault())
	} else {
		cleanRes := maskBinaryFields(h.responseBody)
		resByte, _ := json.Marshal(cleanRes)
		resString := string(resByte)
		if resString == "null" {
			resString = ""
		}
		enc.AddString("response_body", resString)
	}
	return nil
}

// maskBinaryFields serves as a universal utility to sanitize data payloads for logging.
// It recursively processes various data types (Structs, Maps, Slices, and Buffers)
// to identify and mask both binary data and sensitive information.
//
// Key Features and Optimizations:
//   - Multi-Protocol Support: Handles gRPC structs and REST payloads (maps or *bytes.Buffer).
//   - Sensitive Data Masking: Automatically masks fields like 'password' or 'token' using [MASKED].
//   - Binary Masking: Replaces large []byte data with a descriptive size label to prevent log bloat.
//   - Field Name Resolution: Prioritizes "json" tags, followed by "protobuf" name attributes.
//   - Performance via Caching: Uses sync.Map to store metadata, minimizing reflection overhead.
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
// It leverages a metadata cache and performs case-insensitive sensitivity checks.
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

			fieldName := resolveFieldName(structField)

			infos = append(infos, fieldInfo{
				index:       i,
				name:        fieldName,
				isBinary:    f.Kind() == reflect.Slice && f.Type().Elem().Kind() == reflect.Uint8,
				isSensitive: isSensitive(fieldName), // Pre-compute sensitivity
			})
		}
		fieldCache.Store(t, infos)
	}

	m := make(map[string]any, len(infos))
	for _, info := range infos {
		if info.isSensitive {
			m[info.name] = "[MASKED]"
			continue
		}

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

// handleMapMasking iterates through map keys and applies sensitivity checks
// to keys and recursive masking to values.
func handleMapMasking(v reflect.Value) any {
	m := make(map[string]any, v.Len())
	for _, key := range v.MapKeys() {
		var strKey string
		if key.Kind() == reflect.String {
			strKey = key.String()
		} else {
			strKey = fmt.Sprintf("%v", key.Interface())
		}

		if isSensitive(strKey) {
			m[strKey] = "[MASKED]"
		} else {
			m[strKey] = maskBinaryFields(v.MapIndex(key).Interface())
		}
	}
	return m
}

// handleSliceMasking processes each element of a slice recursively.
func handleSliceMasking(v reflect.Value) any {
	s := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		s[i] = maskBinaryFields(v.Index(i).Interface())
	}
	return s
}

// resolveFieldName determines the final key name for a struct field based on
// priority: "json" tag > "protobuf" name > struct field name.
func resolveFieldName(f reflect.StructField) string {
	if jsonTag := f.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "" {
			return parts[0]
		}
	}

	if protoTag := f.Tag.Get("protobuf"); protoTag != "" {
		for part := range strings.SplitSeq(protoTag, ",") {
			if name, found := strings.CutPrefix(part, "name="); found {
				return name
			}
		}
	}

	return f.Name
}

// isSensitive checks if a field name matches the sensitive blacklist.
func isSensitive(name string) bool {
	_, ok := sensitiveFields[strings.ToLower(name)]
	return ok
}
