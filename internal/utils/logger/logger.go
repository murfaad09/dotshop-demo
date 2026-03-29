package logger

import (
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLog *zap.Logger
var zapLogSugared *zap.SugaredLogger
var atom zap.AtomicLevel

const TimeFormat = "2006-01-02T15:04:05.000000000Z0700"
const logFile = "_logs/log.jsonl"

// init sets up the logger configuration.
func init() {
	atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	eCfg := myEncCfg()
	eCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewCore(
		zapcore.NewConsoleEncoder(eCfg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		atom,
	)

	zapLog = zap.New(
		consoleEncoder,
		//zap.WrapCore(fileLoggingWrapper(logFile)),
		zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.PanicLevel),
	)

	zapLogSugared = zapLog.Sugar()
}

func translateLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// SetLogLevel sets the log level.
func SetLogLevel(level string) {
	atom.SetLevel(translateLogLevel(level))
}

// SetupFileLogger will also set up the logger to log to a file.
// It should be called before getting integrations, so they can also log to the file.
func SetupFileLogger(logFile string) {
	zapLog = zapLog.WithOptions(zap.WrapCore(fileLoggingWrapper(logFile)))
	zapLogSugared = zapLogSugared.WithOptions(zap.WrapCore(fileLoggingWrapper(logFile)))
}

func myPathEncoder() zapcore.CallerEncoder {
	projectDir, _ := os.Getwd()
	projectDir += "/"

	return func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(strings.TrimPrefix(caller.FullPath(), projectDir))
	}
}

func myEncCfg() zapcore.EncoderConfig {
	eCfg := zap.NewProductionEncoderConfig()
	eCfg.EncodeTime = zapcore.TimeEncoderOfLayout(TimeFormat)
	eCfg.EncodeCaller = myPathEncoder()
	eCfg.CallerKey = "file"

	return eCfg
}

func fileLoggingWrapper(logFile string) func(zapcore.Core) zapcore.Core {
	return func(c zapcore.Core) zapcore.Core {
		// lumberjack.Logger is already safe for concurrent use, so we don't need to`
		// lock it.
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    1, // megabytes
			MaxBackups: 30,
			MaxAge:     28, // days
			LocalTime:  false,
			Compress:   false,
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(myEncCfg()),
			w,
			atom,
		)

		return zapcore.NewTee(c, core)
	}
}
