package observability

import (
	"context"
	"fmt"
	"github.com/openreserveio/dwn/go/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

var SERVICENAME string

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
		resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		log.Fatal("%s: %v", "Failed to create resource", err)
		os.Exit(1)
	}
	return res
}

// Initiates OpenTelemetry provider sending data to OpenTelemetry Collector.
func InitProviderWithOTELExporter(ctx context.Context, serviceName string) (func(context.Context) error, error) {

	log.Info("Creating a new Tracing Provider, should only happen once")

	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(getSampler()),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(newResource(ctx, serviceName)),
	)
	otel.SetTracerProvider(tp)
	SERVICENAME = serviceName

	return tp.Shutdown, nil
}

func Tracer() trace.Tracer {
	return otel.Tracer(SERVICENAME)
}
