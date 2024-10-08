package metric

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"time"
)

var RequestsCounter metric.Int64Counter
var SecondsHistogram metric.Float64Histogram

func SetMeterProvider(endpoint string, name string) {
	exporter, err := otlpmetrichttp.New(context.Background(), otlpmetrichttp.WithEndpointURL(endpoint))
	if err != nil {
		panic(err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Microsecond*400))),
		sdkmetric.WithView(
			metrics.DefaultSecondsHistogramView(metrics.DefaultServerSecondsHistogramName),
			func(instrument sdkmetric.Instrument) (sdkmetric.Stream, bool) {
				return sdkmetric.Stream{
					Name:        instrument.Name,
					Description: instrument.Description,
					Unit:        instrument.Unit,
					AttributeFilter: func(value attribute.KeyValue) bool {
						return true
					},
				}, true
			}),
	)
	otel.SetMeterProvider(mp)
	meter := mp.Meter(name)
	RequestsCounter, err = metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		panic(err)
	}
	SecondsHistogram, err = metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		panic(err)
	}
}
