# Zap Context Logger

[![Go](https://github.com/imylam/zap-ctx-logger/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/imylam/zap-ctx-logger/actions/workflows/unit-tests.yml)
[![codecov](https://codecov.io/gh/imylam/zap-ctx-logger/branch/main/graph/badge.svg?token=M7AU53MYVJ)](https://codecov.io/gh/imylam/zap-ctx-logger)

Zap-ctx-logger provides a convenient way to to enrich zap logger with context information.

## Installation

`go get github.com/imylam/zap-ctx-logger`

## Quick Start
Let say you have an enrich to retrieve the value of "key" in `context`:

```go
var exampleEnrich = func(ctx context.Context) []zapcore.Field {
    var f []zapcore.Field
    
    value := ctx.Value("key")
    f = append(f, zap.String("key", value.(string)))
    
    return f
}
```

You can instruct `zap-ctx-logger` to enrich zap logs with the value of "key" from `context` as follows:
```go
import zapctx_logger "github.com/imylam/zap-ctx-logger"

func main() {
    zapLogger := zap.NewExample()
    logger := zapctx_logger.New(zapLogger)
    logger.Enrich(exampleEnrich)
    ctx := context.WithValue(context.TODO(), "key", "someValue")
    
    logger.Ctx(ctx).Info("hello")         // {"level":"info","msg":"hello","key":"someValue"}
    logger.Ctx(ctx).Sugar().Info("hello") // {"level":"info","msg":"hello","key":"someValue"}
}
```

## Enrich
Your enrich constructor should have the form of:
```go
func(ctx context.Context) []zapcore.Field
```

You can add enrich to logger during creation:
```go
logger := zapctx_logger.NewWithEnriches(zapLogger, enrich1, enrich2)

logger.Ctx(ctx).Info("hello")
```

Or after creation:
```go
logger := zapctx_logger.New(zapLogger)
logger.Enrich(enrich1)
logger.Enrich(enrich2)

logger.Ctx(ctx).Info("hello")
```

## Enrich zap log with opentelemetry traces information
The module comes with an enrich which adds `spanId` and `traceId` from [opentlemetry traces](https://github.com/open-telemetry/opentelemetry-go).

```go
logger.Enrich(WithSpanIDAndTraceID())

logger.Ctx(ctx).Info("some message") // {"level":"info" "msg":"some message","SpanId":"8b50a3247b315dcd","TraceId":"761679b6ef2e6e20c27c0271f2b07c92"}
```
