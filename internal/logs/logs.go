package logs

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	loggerKey         struct{}
	loggerEncodingKey struct{}
)

func NewLogger(level, format string) (*zap.Logger, string, error) {
	var (
		parsedLevel zapcore.Level
		err         error
	)

	if parsedLevel, err = zapcore.ParseLevel(level); err != nil {
		return nil, "", errors.Wrap(err, "parsing log level")
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	cfg.Level.SetLevel(parsedLevel)

	// Omit stack traces on errors when log level is not "debug"
	if cfg.Level.Level() != zapcore.DebugLevel {
		cfg.DisableStacktrace = true
	} else {
		cfg.Sampling = nil
		cfg.Development = true
	}

	switch format {
	case "json":
		cfg.Encoding = "json"
	case "text":
		cfg.Encoding = "console"
	default:
		return nil, "", fmt.Errorf("unsupported log format '%s', (expected json or text)", format)
	}

	var logger *zap.Logger
	if logger, err = cfg.Build(); err != nil {
		return nil, cfg.Encoding, errors.Wrap(err, "creating logger")
	}

	return logger, cfg.Encoding, nil
}

func StoreLoggerInContext(ctx context.Context, l *zap.Logger, loggerEncoding string) context.Context {
	ctx = context.WithValue(ctx, loggerKey{}, l)
	ctx = context.WithValue(ctx, loggerEncodingKey{}, loggerEncoding)
	return ctx
}

// LoggerFromContext returns logs from context
func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}

// LoggerEncodingFromContext returns logs's encoding from context
func LoggerEncodingFromContext(ctx context.Context) string {
	if encoding, ok := ctx.Value(loggerEncodingKey{}).(string); ok {
		return encoding
	}
	return "console"
}
