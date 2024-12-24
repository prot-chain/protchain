package restapi

import (
	"context"
	"net/http"
	"protchain/internal/tracing"
	"protchain/internal/value"

	"github.com/lucsky/cuid"
	"github.com/pkg/errors"
)

// RequestTracing handles the request tracing context
func RequestTracing(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestSource := r.Header.Get(value.HeaderRequestSource)
		if requestSource == "" {
			errM := errors.New("X-Request-Source is empty")
			writeErrorResponse(w, errM, value.Error, errM.Error())
			return
		}

		requestID := r.Header.Get(value.HeaderRequestID)
		if requestID == "" {
			requestID = cuid.New()
		}

		tracingContext := tracing.Context{
			RequestID:     requestID,
			RequestSource: requestSource,
		}

		ctx = context.WithValue(ctx, value.ContextTracingKey, tracingContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// RequestTracingII handles the request tracing context
func RequestTracingII(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestSource := r.Header.Get(value.HeaderRequestSource)
		if requestSource == "" {
			requestSource = "webhook"
		}

		requestID := r.Header.Get(value.HeaderRequestID)
		if requestID == "" {
			requestID = cuid.New()
		}

		tracingContext := tracing.Context{
			RequestID:     requestID,
			RequestSource: requestSource,
		}

		ctx = context.WithValue(ctx, value.ContextTracingKey, tracingContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
