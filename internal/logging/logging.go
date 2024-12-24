package logging

import (
	"protchain/internal/config"
	"protchain/internal/tracing"
	"protchain/internal/value"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// setupSentryOptions
func setupSentryOptions() zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			if entry.Level == zapcore.ErrorLevel {
				defer sentry.Flush(2 * time.Second)
				// sentry.CaptureException()
			}
			return nil
		})
	})
}

// New initiates a new logger
func New(cfg *config.Config) {
	if cfg.Environment == "production" {
		Log = zap.Must(zap.NewProduction())
		Log.WithOptions(setupSentryOptions())
	} else {
		Log = zap.Must(zap.NewDevelopment())
	}
}

// GetContext returns a list of zap field containing the
// tracing context, and error message
func GetContext(ctx *tracing.Context, err error) []zap.Field {
	fields := make([]zap.Field, 0)
	fields = append(fields, zap.Error(err))
	if ctx != nil {
		fields = append(fields, zap.String(value.HeaderRequestID, ctx.RequestID), zap.String(value.HeaderRequestSource, ctx.RequestSource))
	}
	return fields
}
