package response

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/trace"
)

func AppendTraceId(header string) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				tr.ReplyHeader().Set(header, trace.SpanContextFromContext(ctx).TraceID().String())
			}
			return handler(ctx, req)
		}
	}
}
