package logging

import (
	"context"
	"os"
	"path/filepath"
	"ride-sharing/config"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	RequestIDKey  = "x-request-id"
	CorrelationID = "x-correlation-id"
)

var (
	instance *Logger
	once     sync.Once

	standardFields = []zap.Field{
		zap.String("service", config.GetEnv("SERVICE_NAME", "Auth-service")),
		zap.String("environment", config.GetEnv("ENVIRONMENT", "Dev")),
		zap.String("version", config.GetEnv("VERSION", "1.0.0")),
	}

	sensitiveFields = map[string]struct{}{
		"password":         {},
		"confirm_password": {},
		"access_token":     {},
		"refresh_token":    {},
		"token":            {},
		"pin":              {},
		"credit_card":      {},
		"cvv":              {},
		"authorization":    {},
		"set-cookie":       {},
	}
)

type Logger struct {
	*zap.Logger
}

type LogConfig struct {
	Environement string
	Version      string
	ServiceName  string
}

var logger *zap.Logger

func InitLogger(cfg LogConfig) {
	logPath := filepath.Join(cfg.Environement, cfg.Version, "log", "log.log")

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writer,
		zapcore.InfoLevel,
	)

	logger = zap.New(core).With(
		zap.String("service", cfg.ServiceName),
		zap.String("environment", cfg.Environement),
		zap.String("version", cfg.Version),
	)
}

// GetLogger returns the singleton logger instance
func GetLogger() *Logger {
	once.Do(func() {
		var err error
		instance, err = NewLogger(true) // Production defaults to true
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
	})
	return instance
}

func NewLogger(production bool) (*Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Multi-sink setup
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel && lvl >= zapcore.InfoLevel
	})
	logPath := filepath.Join(
		"log",
		os.Getenv("ENVIRONMENT"),
		os.Getenv("VERSION"),
		"log",
	)

	// Core setup
	cores := []zapcore.Core{
		// File output with rotation
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    100, // MB
				MaxBackups: 7,
				MaxAge:     30, // days
				Compress:   true,
			}),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return true }),
		),
		// Stderr for errors
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stderr),
			highPriority,
		),
		// Stdout for info
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			lowPriority,
		),
	}

	logger := zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(standardFields...),
	)

	return &Logger{logger}, nil
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := []zap.Field{
		zap.String("request_id", getStringFromContext(ctx, RequestIDKey)),
		zap.String("correlation_id", getStringFromContext(ctx, CorrelationID)),
	}

	// Add goroutine ID
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields = append(fields, zap.String("goroutine", strings.Fields(string(buf[:n]))[1]))

	return &Logger{l.Logger.With(fields...)}
}

func (l *Logger) Shutdown() error {
	return l.Sync()
}

func getStringFromContext(ctx context.Context, key string) string {
	if val, ok := ctx.Value(key).(string); ok {
		return val
	}
	return ""
}

// MaskSensitiveData recursively masks sensitive fields in any data structure
func MaskSensitiveData(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if _, ok := sensitiveFields[strings.ToLower(key)]; ok {
				v[key] = "****"
			} else {
				v[key] = MaskSensitiveData(val)
			}
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = MaskSensitiveData(val)
		}
		return v
	default:
		return data
	}
}
