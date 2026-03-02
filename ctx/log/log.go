// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"os"
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

func Get(logEntity configEntity.Log, appModule configEntity.Module) *Log {
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
			logger:    zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)),
			AppModule: appModule,
		}
	})
	return log
}

func CustomEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[level])
}

func (l *Log) Error(msg string) {
	l.logger.Error(msg, l.ZapFields()...)
}

func (l *Log) Warn(msg string) {
	l.logger.Warn(msg, l.ZapFields()...)
}

func (l *Log) Info(msg string) {
	l.logger.Info(msg, l.ZapFields()...)
}

func (l *Log) Debug(msg string) {
	l.logger.Debug(msg, l.ZapFields()...)
}

func (l *Log) ZapFields() (fields []zap.Field) {
	fields = append(fields, zap.String("module", l.AppModule.DirectoryName()), zap.String(CorrelationIDKeyContext.String(), l.CorrelationID), zap.String("path", l.Path))
	if l.QueryString != nil {
		fields = append(fields, zap.String("query_string", *l.QueryString))
	}
	return
}

func (k CorrelationIDKey) String() string {
	return string(k)
}
