// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"os"
	"time"

	configEntity "github.com/dedyf5/resik/entities/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	Logger *zap.Logger
}

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "debug",
	zapcore.InfoLevel:   "info",
	zapcore.WarnLevel:   "warning",
	zapcore.ErrorLevel:  "error",
	zapcore.DPanicLevel: "critical",
	zapcore.PanicLevel:  "alert",
	zapcore.FatalLevel:  "emergency",
}

func New(logEntity configEntity.Log) *Log {
	consoleDebugging := zapcore.Lock(os.Stdout)
	cfgConsole := zapcore.EncoderConfig{
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
	coreConfig = append(coreConfig, zapcore.NewCore(zapcore.NewJSONEncoder(cfgConsole), consoleDebugging, zap.DebugLevel))
	if logEntity.File != "" {
		logFile, err := os.OpenFile(logEntity.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err.Error())
		}
		writer := zapcore.AddSync(logFile)
		coreConfig = append(coreConfig, zapcore.NewCore(zapcore.NewJSONEncoder(cfgConsole), writer, zap.DebugLevel))
	}
	core := zapcore.NewTee(coreConfig...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()
	return &Log{
		Logger: logger,
	}
}

func CustomEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[level])
}

func CustomLevelFileEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + logLevelSeverity[level] + "]")
}
