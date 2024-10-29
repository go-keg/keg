package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func SetTracerProvider(endpoint string, attributes ...attribute.KeyValue) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpointURL(endpoint),
		otlptracehttp.WithTimeout(3*time.Second),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)

	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(fmt.Errorf("creating OTLP trace exporter: %w", err))
	}

	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewSchemaless(attributes...)),
	))

	otel.SetTextMapPropagator(propagation.TraceContext{})
}
