// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"context"
	"os"
	"sync"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Log struct {
	Logger        *zap.Logger
	CorrelationID string
}

type contextKey struct{}

type CorrelationIDKey string

const (
	CorrelationIDKeyContext CorrelationIDKey = "correlation_id"
	CorrelationIDKeyXHeader CorrelationIDKey = "X-Correlation-ID"
)

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

type ILog interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

func Get(logEntity configEntity.Log) *Log {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)
		encoderConfig := zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			EncodeLevel:   CustomEncodeLevel,
			TimeKey:       "time",
			EncodeTime:    zapcore.TimeEncoderOfLayout(time.RFC3339),
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
			Logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)),
		}
	})
	return log
}

func CustomEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[level])
}

func FromContext(ctx context.Context) *Log {
	if l, ok := ctx.Value(contextKey{}).(*Log); ok {
		return l
	} else if l := log; l != nil {
		return l
	}

	return &Log{
		Logger: zap.NewNop(),
	}
}

func WithContext(ctx context.Context, l *Log) context.Context {
	if lp, ok := ctx.Value(contextKey{}).(*Log); ok {
		if lp == l {
			return ctx
		}
	}

	return context.WithValue(ctx, contextKey{}, l)
}

func (l *Log) Error(service *Service, msg string) {
	l.Logger.Error(msg, zap.String(CorrelationIDKeyContext.String(), l.CorrelationID), zap.Inline(service))
}

func (l *Log) Warn(service *Service, msg string) {
	l.Logger.Warn(msg, zap.String(CorrelationIDKeyContext.String(), l.CorrelationID), zap.Inline(service))
}

func (l *Log) Info(service *Service, msg string) {
	l.Logger.Info(msg, zap.String(CorrelationIDKeyContext.String(), l.CorrelationID), zap.Inline(service))
}

func (l *Log) Debug(service *Service, msg string) {
	l.Logger.Debug(msg, zap.String(CorrelationIDKeyContext.String(), l.CorrelationID), zap.Inline(service))
}

func (l *Log) FromContext(ctx context.Context) *Log {
	return FromContext(ctx)
}

func (l *Log) WithContext(ctx context.Context) context.Context {
	return WithContext(ctx, l)
}

func (k CorrelationIDKey) String() string {
	return string(k)
}
