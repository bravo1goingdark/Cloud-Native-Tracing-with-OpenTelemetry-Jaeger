package tracing

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
)

import typeTRACER "go.opentelemetry.io/otel/trace"

func InitTracer(serviceName string) (*trace.TracerProvider, error) {
	jaegerEndpoint := "http://jaeger:14268/api/traces"

	fmt.Println("🚀 Connecting to Jaeger at:", jaegerEndpoint) // Debugging

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		fmt.Println("❌ Failed to create Jaeger exporter:", err) // Debugging
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	fmt.Println("✅ Tracer initialized for", serviceName) // Debugging
	return tp, nil
}

func Tracer() typeTRACER.Tracer {
	return otel.Tracer("jaegaurd-tracer")
}
