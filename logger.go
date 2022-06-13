package zapctx_logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/imylam/zap-ctx-logger/enrich"
)

type Logger struct {
	zapLogger   *zap.Logger
	ctxEnriches []enrich.Enrich
}

func New(zapLogger *zap.Logger) *Logger {
	return &Logger{zapLogger: zapLogger}
}

func NewWithEnriches(zapLogger *zap.Logger, enriches ...enrich.Enrich) *Logger {
	return &Logger{
		zapLogger:   zapLogger,
		ctxEnriches: enriches,
	}
}

func (z *Logger) Logger() *zap.Logger {
	return z.zapLogger
}

func (z *Logger) Enrich(enrich enrich.Enrich) {
	z.ctxEnriches = append(z.ctxEnriches, enrich)
}

func (z *Logger) Ctx(ctx context.Context) *zap.Logger {

	if len(z.ctxEnriches) == 0 {
		return z.zapLogger
	}

	var f []zapcore.Field
	for _, enrich := range z.ctxEnriches {
		f = append(f, enrich(ctx)...)
	}

	return z.zapLogger.With(f...)
}
