package internal

import (
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"

	"github.com/MarioCarrion/todo-api/internal/envvar"
)

// NewOTExporter instantiates the OpenTelemetry exporters using configuration defined in environment variables.
func NewOTExporter(conf *envvar.Configuration) (*prometheus.Exporter, error) {
	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second)); err != nil {
		return nil, fmt.Errorf("runtime.Start %w", err)
	}

	promExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		return nil, fmt.Errorf("prometheus.NewExportPipeline %w", err)
	}

	global.SetMeterProvider(promExporter.MeterProvider())

	//-

	jaegerEndpoint, _ := conf.Get("JAEGER_ENDPOINT")

	jaegerExporter, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
	)
	if err != nil {
		return nil, fmt.Errorf("jaeger.NewRawExporter %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(jaegerExporter),
		sdktrace.WithResource(resource.NewWithAttributes(attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue("rest-server"),
		})),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return promExporter, nil
}
