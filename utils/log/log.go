// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package log

import (
	"os"
	"time"

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

func New() *Log {
	consoleDebugging := zapcore.Lock(os.Stdout)
	cfgConsole := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  CustomEncodeLevel,
		TimeKey:      "time",
		EncodeTime:   zapcore.TimeEncoderOfLayout(time.RFC3339),
		CallerKey:    "line",
		EncodeCaller: zapcore.FullCallerEncoder,
	}
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(cfgConsole), consoleDebugging, zap.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
	// sugar := logger.Sugar()
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
