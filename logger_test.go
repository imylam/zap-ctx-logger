package zapctx_logger

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/imylam/zap-ctx-logger/enrich"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestNewWithEnriches(t *testing.T) {

	t.Run("GIVEN_logger_WHEN_enrich_once_THEN_logger_contain_one_enrich", func(t *testing.T) {
		zapLogger := zap.NewExample()
		logger := NewWithEnriches(zapLogger, dummyCtxEnrich("dummy", "value"))

		assert.Equal(t, 1, len(logger.ctxEnriches))
	})

	t.Run("GIVEN_logger_WHEN_enrich_twice_THEN_logger_contain_two_enriches", func(t *testing.T) {
		zapLogger := zap.NewExample()
		logger := NewWithEnriches(zapLogger, dummyCtxEnrich("dummy", "value"), dummyCtxEnrich("dummy2", "value2"))

		assert.Equal(t, 2, len(logger.ctxEnriches))
	})
}

func TestLogger(t *testing.T) {
	zapLogger := zap.NewExample()
	logger := NewWithEnriches(zapLogger, dummyCtxEnrich("dummy", "value"))

	assert.Equal(t, zapLogger, logger.Logger())
}

func TestEnrich(t *testing.T) {

	t.Run("GIVEN_logger_WHEN_enrich_once_THEN_logger_contain_one_enrich", func(t *testing.T) {
		zapLogger := zap.NewExample()
		logger := New(zapLogger)

		logger.Enrich(dummyCtxEnrich("dummy", "value"))

		assert.Equal(t, 1, len(logger.ctxEnriches))
	})

	t.Run("GIVEN_logger_WHEN_enrich_twice_THEN_logger_contain_two_enriches", func(t *testing.T) {
		zapLogger := zap.NewExample()
		logger := New(zapLogger)

		logger.Enrich(dummyCtxEnrich("dummy", "value"))
		logger.Enrich(dummyCtxEnrich("dummy2", "value2"))

		assert.Equal(t, 2, len(logger.ctxEnriches))
	})
}

func TestCtx(t *testing.T) {

	t.Run("GIVEN_logger_without_enriches_WHEN_log_THEN_log_msg_has_no_enrich_data", func(t *testing.T) {
		ts := newTestLogSpy(t)

		zapLogger := zaptest.NewLogger(ts)
		logger := New(zapLogger)

		ctx := context.TODO()
		logger.Ctx(ctx).Info("testing")
		ts.AssertMessages("INFO	testing")
	})

	t.Run("GIVEN_logger_with_one_enriches_WHEN_log_THEN_log_msg_has_one_enrich_data", func(t *testing.T) {
		ts := newTestLogSpy(t)

		zapLogger := zaptest.NewLogger(ts)
		logger := New(zapLogger)
		logger.Enrich(dummyCtxEnrich("dummy", "value"))

		ctx := context.TODO()
		logger.Ctx(ctx).Info("testing")
		_ = ts.Messages
		ts.AssertMessages(`INFO	testing	{"dummy": "value"}`)
	})

	t.Run("GIVEN_logger_with_two_enriches_WHEN_log_THEN_log_msg_has_two_enrich_data", func(t *testing.T) {
		ts := newTestLogSpy(t)

		zapLogger := zaptest.NewLogger(ts)
		logger := New(zapLogger)
		logger.Enrich(dummyCtxEnrich("dummy", "value"))
		logger.Enrich(dummyCtxEnrich("dummy2", "value2"))

		ctx := context.TODO()
		logger.Ctx(ctx).Info("testing")
		_ = ts.Messages
		ts.AssertMessages(`INFO	testing	{"dummy": "value", "dummy2": "value2"}`)
	})
}

func dummyCtxEnrich(key string, val string) enrich.Enrich {
	return func(ctx context.Context) []zapcore.Field {
		var f []zapcore.Field
		f = append(f, zap.String(key, val))

		return f
	}
}

// testLogSpy is a testing.TB that captures logged messages.
// taking reference from https://github.com/uber-go/zap/blob/master/zaptest/logger_test.go
type testLogSpy struct {
	testing.TB
	Messages []string
}

func newTestLogSpy(t testing.TB) *testLogSpy {
	return &testLogSpy{TB: t}
}

func (t *testLogSpy) Logf(format string, args ...interface{}) {
	// Log messages are in the format,
	//
	//   2017-10-27T13:03:01.000-0700	DEBUG	your message here	{data here}
	m := fmt.Sprintf(format, args...)
	m = m[strings.IndexByte(m, '\t')+1:]
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *testLogSpy) AssertMessages(msgs ...string) {
	assert.Equal(t.TB, msgs, t.Messages, "logged messages did not match")
}
