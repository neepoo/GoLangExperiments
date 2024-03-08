package main

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func newTraceExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func newMetricExporter() (metric.Exporter, error) {
	return stdoutmetric.New()
}

/*
Our first Signal will be the trace. To interact with this Signal, we must create a Provider, as shown below.
As a parameter, we will have the Exporter that will receive the collected information.
*/
func newTraceProvider(traceExporter trace.SpanExporter) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second)),
	)
	return traceProvider
}

func newMetricProvider(metricExport metric.Exporter) *metric.MeterProvider {
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExport,
			metric.WithInterval(10*time.Second))),
	)
	return meterProvider
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{})
}
