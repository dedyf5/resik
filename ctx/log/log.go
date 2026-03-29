// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Log struct {
	logger        *zap.Logger
	AppModule     configEntity.Module
	CorrelationID string
	Path          string
	QueryString   *string
}

type Key string

const (
	KeyCorrelationIDContext Key = "correlation_id"
	KeyXCorrelationIDHeader Key = "X-Correlation-ID"
	KeyCallerHolderContext  Key = "caller"
)

type CallerHolder struct {
	Caller *string
}

var once sync.Once

var log *Log

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "debug",
	zapcore.InfoLevel:   "info",
	zapcore.WarnLevel:   "warning",
	zapcore.ErrorLevel:  "error",
	zapcore.DPanicLevel: "critical",
	zapcore.PanicLevel:  "alert",
	zapcore.FatalLevel:  "emergency",
}

func Get(logEntity configEntity.Log, appModule configEntity.Module) *Log {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		encoderConfig := zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: CustomEncodeLevel,
			TimeKey:     "time",
			EncodeTime: func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
				pae.AppendString(t.UTC().Format("2006-01-02T15:04:05Z"))
			},
			CallerKey:     "line",
			EncodeCaller:  zapcore.FullCallerEncoder,
			StacktraceKey: "stack_trace",
		}

		coreConfig := []zapcore.Core{}
		coreConfig = append(coreConfig, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), stdout, zap.DebugLevel))

		if logEntity.File != "" {
			logFile := zapcore.AddSync(&lumberjack.Logger{
				Filename:   logEntity.File,
				MaxSize:    5,
				MaxBackups: 10,
				MaxAge:     30,
				Compress:   true,
			})
			writer := zapcore.AddSync(logFile)
			coreConfig = append(coreConfig, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, zap.DebugLevel))
		}

		core := zapcore.NewTee(coreConfig...)

		log = &Log{
			logger:    zap.New(core, zap.AddCaller()),
			AppModule: appModule,
		}
	})
	return log
}

func CustomEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[level])
}

func (l *Log) Error(msg string) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, l.ZapFields()...)
}

func (l *Log) Warn(msg string) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, l.ZapFields()...)
}

func (l *Log) Info(msg string) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, l.ZapFields()...)
}

func (l *Log) Debug(msg string) {
	l.logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, l.ZapFields()...)
}

func (l *Log) ZapFields() (fields []zap.Field) {
	fields = append(fields, zap.String("module", l.AppModule.DirectoryName()), zap.String(KeyCorrelationIDContext.String(), l.CorrelationID), zap.String("path", l.Path))
	if l.QueryString != nil {
		fields = append(fields, zap.String("query_string", *l.QueryString))
	}
	return
}

func (k Key) String() string {
	return string(k)
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
