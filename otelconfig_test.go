package otelconfig

import (
	"context"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "127.0.0.1:4317")
	os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
	os.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")
	os.Setenv("OTEL_EXPORTER_OTLP_SERVICE_NAME", "my-instrumented-service")

	ctx := context.Background()
	shutdown := ConfigureOpentelemetry(ctx)
	defer shutdown()
}
