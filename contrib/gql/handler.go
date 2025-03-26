package gql

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	httptransport "github.com/go-kratos/kratos/v2/transport"
	"github.com/vektah/gqlparser/v2/ast"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func NewServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func TraceAroundResponses(middleware graphql.ResponseMiddleware) graphql.ResponseMiddleware {
	return func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		resp := middleware(ctx, next)
		if resp.Errors != nil {
			trace.SpanFromContext(ctx).SetStatus(codes.Error, resp.Errors.Error())
		} else {
			trace.SpanFromContext(ctx).SetStatus(codes.Ok, "success")
		}
		return resp
	}
}

func TraceAroundOperations(middleware graphql.OperationMiddleware) graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		op := graphql.GetOperationContext(ctx)
		span := trace.SpanFromContext(ctx)

		if s, ok := httptransport.FromServerContext(ctx); ok {
			s.ReplyHeader().Set("X-TRACE-ID", span.SpanContext().TraceID().String())
		}
		if op != nil {
			attrs := []attribute.KeyValue{
				attribute.Key("graphql.document").String(op.RawQuery),
			}
			if op.Operation != nil {
				attrs = append(attrs,
					attribute.Key("graphql.operation.name").String(op.Operation.Name),
					attribute.Key("graphql.operation.type").String(string(op.Operation.Operation)),
				)
				span.SetName(fmt.Sprintf("%s %s", op.Operation.Operation, op.Operation.Name))
			}
			span.SetAttributes(attrs...)
		}
		return middleware(ctx, next)
	}
}
