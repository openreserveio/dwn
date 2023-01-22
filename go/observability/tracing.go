package observability

import (
	"context"
	"github.com/openreserveio/dwn/go/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

var TP *tracesdk.TracerProvider
var Tracer trace.Tracer

// Helper function to define sampling.
// When in development mode, AlwaysSample is defined,
// otherwise, sample based on Parent and IDRatio will be used.
func getSampler() tracesdk.Sampler {
	ENV := os.Getenv("GO_ENV")
	switch ENV {
	case "development":
		return tracesdk.AlwaysSample()
	case "production":
		return tracesdk.ParentBased(tracesdk.TraceIDRatioBased(0.5))
	default:
		return tracesdk.AlwaysSample()
	}
}

// Returns a new OpenTelemetry resource describing this application.
func newResource(ctx context.Context, serviceName string) *resource.Resource {
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", os.Getenv("GO_ENV")),
		),
	)
	if err != nil {
		log.Fatal("%s: %v", "Failed to create resource", err)
		os.Exit(1)
	}
	return res
}

// Creates Jaeger exporter
func exporterToJaeger() (*jaeger.Exporter, error) {
	collUrl := os.Getenv("OPEN_TELEMETRY_COLLECTOR_URL")
	log.Info("Open Telemetry Collector URL:  %s", collUrl)
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(collUrl)))
}

// Maybe this will work?
func exporterToCollector(ctx context.Context) (tracesdk.SpanExporter, error) {
	collUrl := os.Getenv("OPEN_TELEMETRY_COLLECTOR_URL")
	log.Info("Open Telemetry Collector URL:  %s", collUrl)

	client := otlptracehttp.NewClient()
	return otlptrace.New(ctx, client)
}

// Initiates OpenTelemetry provider sending data to OpenTelemetry Collector.
func InitProviderWithJaegerExporter(ctx context.Context, serviceName string) (func(context.Context) error, error) {

	if TP != nil {
		return TP.Shutdown, nil
	}

	log.Info("Creating a new Tracing Provider, should only happen once")

	exp, err := exporterToCollector(ctx)
	if err != nil {
		log.Fatal("error: %s", err.Error())
		os.Exit(1)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(getSampler()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(newResource(ctx, serviceName)),
	)
	otel.SetTracerProvider(tp)
	TP = tp
	Tracer = tp.Tracer(serviceName)
	return tp.Shutdown, nil
}
