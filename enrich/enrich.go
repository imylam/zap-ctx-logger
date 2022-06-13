package enrich

import (
	"context"

	"go.uber.org/zap/zapcore"
)

type Enrich func(ctx context.Context) []zapcore.Field
