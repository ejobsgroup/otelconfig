package otelconfig

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func ConfigureOpentelemetry(ctx context.Context) func() {
	switch {
	case os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "":
		return configureOTLP(ctx)
	default:
		return configureStdout(ctx)
	}
}

func configureOTLP(ctx context.Context) func() {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(os.Getenv("OTEL_EXPORTER_OTLP_SERVICE_NAME")),
		),
	)
	if err != nil {
		panic(err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)

	protocol := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	switch protocol {
	case "grpc":
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		var conn *grpc.ClientConn
		var err error
		if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true" {
			conn, err = grpc.DialContext(ctx, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithBlock(),
			)
		} else {
			conn, err = grpc.DialContext(ctx, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
				grpc.WithBlock(),
			)
		}
		if err != nil {
			panic(err)
		}
		exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			panic(err)
		}
		bsp := sdktrace.NewBatchSpanProcessor(exp)
		provider.RegisterSpanProcessor(bsp)
	case "http":
		exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")))
		if err != nil {
			panic(err)
		}
		bsp := sdktrace.NewBatchSpanProcessor(exp)
		provider.RegisterSpanProcessor(bsp)
	}

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

func configureStdout(ctx context.Context) func() {
	provider := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(provider)

	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	provider.RegisterSpanProcessor(bsp)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}
