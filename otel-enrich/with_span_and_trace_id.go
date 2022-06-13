package otelenrich

import (
	"context"

	"github.com/imylam/zap-ctx-logger/enrich"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	KeySpanId  string = "SpanId"
	KeyTraceId string = "TraceId"
)

func WithSpanIDAndTraceID() enrich.Enrich {
	return func(ctx context.Context) []zapcore.Field {
		var f []zapcore.Field

		span := trace.SpanFromContext(ctx)
		if !span.IsRecording() {
			return f
		}

		if span.SpanContext().HasSpanID() {
			spanID := span.SpanContext().SpanID().String()
			f = append(f, zap.String(KeySpanId, spanID))
		}
		if span.SpanContext().HasTraceID() {
			traceID := span.SpanContext().TraceID().String()
			f = append(f, zap.String(KeyTraceId, traceID))
		}

		return f
	}
}
