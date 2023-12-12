package log

import (
	"os"

	"github.com/wintbiit/ninedns/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func init() {
	logger := NewLogger("main")
	zap.ReplaceGlobals(logger.SugaredLogger.Desugar())
}

func NewLogger(module string) *Logger {
	level := zapcore.InfoLevel
	if utils.C.Debug {
		level = zapcore.DebugLevel
	}

	return NewLoggerWithLevel(module, level)
}

func NewLoggerWithLevel(module string, level zapcore.Level) *Logger {
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(module)
		enc.AppendString(caller.String())
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), zapcore.Lock(os.Stdout), level)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	logger.Named(module)
	defer logger.Sync()

	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}
