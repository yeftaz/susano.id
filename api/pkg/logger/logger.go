package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/yeftaz/susano.id/api/internal/config"
)

type Logger struct {
	*zap.SugaredLogger
}

// New creates a new logger instance
func New(cfg *config.Config) *Logger {
	var zapLogger *zap.Logger

	if cfg.AppEnv == "development" {
		// Development logger: console output with colors
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Level = zap.NewAtomicLevelAt(getLogLevel(cfg.LogLevel))

		zapLogger, _ = config.Build()
	} else {
		// Production logger: file output with rotation
		logDir := "storage/logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic(err)
		}

		logFile := filepath.Join(logDir, "app.log")

		// Lumberjack for log rotation
		lumberjackLogger := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    cfg.LogMaxSize,    // MB
			MaxBackups: cfg.LogMaxBackups, // number of backups
			MaxAge:     cfg.LogMaxAge,     // days
			Compress:   cfg.LogCompress,   // compress rotated files
			LocalTime:  true,
		}

		// Encoder configuration
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		// Create core
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(lumberjackLogger),
			getLogLevel(cfg.LogLevel),
		)

		zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

// getLogLevel converts string log level to zapcore.Level
func getLogLevel(level string) zapcore.Level {
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

// Sync flushes any buffered log entries
func (l *Logger) Sync() {
	_ = l.SugaredLogger.Sync()
}
