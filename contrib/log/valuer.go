package log

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/trace"
)

func TraceID() log.Valuer {
	return func(ctx context.Context) any {
		return trace.SpanContextFromContext(ctx).TraceID().String()
	}
}
