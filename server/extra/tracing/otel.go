package tracing

import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

func SetTraceProvider(name string) (*trace.TracerProvider, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(),
		otlptracehttp.WithInsecure(),
	)
}
