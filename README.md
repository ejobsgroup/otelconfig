# otelconfig

Use it like this in your service entry point:

```go

import (
    "context"

    "github.com/ejobsgroup/otelconfig"
)

func main() {
	ctx := context.Background()

	shutdown := otelconfig.ConfigureOpentelemetry(ctx)
	defer shutdown()

    // Use the global trace provider which is now configured to what the env vars
    // specified. Example env configs:
    //
    // Note: for GRPC the endpoint shouldn't contain a protocol
    //
    // OTEL_EXPORTER_OTLP_ENDPOINT=127.0.0.1:4317
    // OTEL_EXPORTER_OTLP_PROTOCOL=grpc
    // OTEL_EXPORTER_OTLP_INSECURE=true
    // OTEL_EXPORTER_OTLP_SERVICE_NAME=my-instrumented-service
    //
    // or
    //
    // OTEL_EXPORTER_OTLP_ENDPOINT=http://127.0.0.1:4318
    // OTEL_EXPORTER_OTLP_PROTOCOL=http
    // OTEL_EXPORTER_OTLP_INSECURE=true
    // OTEL_EXPORTER_OTLP_SERVICE_NAME=my-instrumented-service
    //
}
```
