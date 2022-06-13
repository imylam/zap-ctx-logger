package otelenrich

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imylam/zap-ctx-logger/otel-enrich/mock"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestWithSpanIDAndTraceID(t *testing.T) {
	var enrich = WithSpanIDAndTraceID()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("GIVEN_nil_span_in_context_WHEN_enrich_with_trace_infos_THEN_zapCore_fields_should_be_empty", func(t *testing.T) {
		var ctx = context.TODO()
		ctx = trace.ContextWithSpan(ctx, nil)

		zapFields := enrich(ctx)

		assert.Equal(t, 0, len(zapFields))
	})

	t.Run("GIVEN_span_not_recording_WHEN_enrich_with_trace_infos_THEN_zapCore_fields_should_be_empty", func(t *testing.T) {
		var ctx = context.TODO()
		spanContext := createSpanContext(createValidSpanConfig())

		span := createMockSpan(ctrl, spanContext, false)
		ctx = trace.ContextWithSpan(ctx, span)

		zapFields := enrich(ctx)

		assert.Equal(t, 0, len(zapFields))
	})

	t.Run("GIVEN_span_recording_WHEN_enrich_with_trace_infos_THEN_zapCore_fields_should_contains_spanID_and_traceID", func(t *testing.T) {
		var ctx = context.TODO()
		spanContext := createSpanContext(createValidSpanConfig())

		span := createMockSpan(ctrl, spanContext, true)
		ctx = trace.ContextWithSpan(ctx, span)

		zapFields := enrich(ctx)

		assert.Equal(t, 2, len(zapFields))
		assert.Equal(t, KeySpanId, zapFields[0].Key)
		assert.Equal(t, spanContext.SpanID().String(), zapFields[0].String)
		assert.Equal(t, KeyTraceId, zapFields[1].Key)
		assert.Equal(t, spanContext.TraceID().String(), zapFields[1].String)
	})
}

func createMockSpan(ctrl *gomock.Controller, spanCtx trace.SpanContext, isRecording bool) trace.Span {
	mockSpan := mock.NewMockSpan(ctrl)
	mockSpan.EXPECT().IsRecording().AnyTimes().Return(isRecording)
	mockSpan.EXPECT().SpanContext().AnyTimes().Return(spanCtx)

	return mockSpan
}

func createValidSpanConfig() trace.SpanContextConfig {
	return trace.SpanContextConfig{
		SpanID:  createSpanID(),
		TraceID: createTraceID(),
	}
}

func createSpanContext(config trace.SpanContextConfig) (spanContext trace.SpanContext) {
	spanContext = trace.NewSpanContext(config)

	return
}

func createSpanID() (spanID trace.SpanID) {
	copy(spanID[:], genRandomBytes(8)[:8])
	return
}

func createTraceID() (traceID trace.TraceID) {
	copy(traceID[:], genRandomBytes(16)[:16])
	return
}

func genRandomBytes(size int) (blk []byte) {
	blk = make([]byte, size)
	rand.Read(blk)

	return
}
