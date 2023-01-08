// Package otlm supports to provide Tracer using OpenTelemetry.
package otlm

import (
	"context"
	"fmt"
	"log"
	"net/http"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/brbranch/go-gae-logger/logger/model"
)

// OpenTelemetryProvider implements Provider interface that provide OpenTelemetry Tracer.
type OpenTelemetryProvider struct {
	name      string
	projectId string
	isLocal   bool
}

// Create a OpenTelemetryProvider.
func Create(traceName string, projectId string) *OpenTelemetryProvider {
	return &OpenTelemetryProvider{name: traceName, projectId: projectId, isLocal: false}
}

func (o *OpenTelemetryProvider) GetSpan(ctx context.Context) *model.Span {
	span := trace.SpanFromContext(ctx).SpanContext()
	return &model.Span{
		SpanID:  span.SpanID().String(),
		TraceID: span.TraceID().String(),
		Valid:   span.IsValid(),
	}
}

func (o *OpenTelemetryProvider) LocalMode() {
	log.Print("OpenTelemetryProvider - Set localMode")
	o.isLocal = true
}

func (o *OpenTelemetryProvider) createExporter(ctx context.Context, path string) (context.Context, error, func()) {
	projectID := o.projectId
	if o.isLocal {
		return ctx, nil, func() {}
	}
	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(projectID))

	if err != nil {
		return nil, err, nil
	}

	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(o.name),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	tracer := tp.Tracer(o.name)
	spanCtx, span := tracer.Start(ctx, path)
	otel.SetTracerProvider(tp)

	return spanCtx, nil, func() {
		_ = tp.Shutdown(ctx)
		span.End()
	}
}

func (o *OpenTelemetryProvider) ProjectID() string {
	return o.projectId
}

func (o *OpenTelemetryProvider) CustomSpan(c context.Context, label string) (context.Context, func()) {
	tracer := otel.GetTracerProvider().Tracer(o.name)
	ctx, span := tracer.Start(c, label)
	return ctx, func() {
		span.End()
	}
}

func (o *OpenTelemetryProvider) StartSpan(request *http.Request, path string) (context.Context, func()) {
	ctx := request.Context()
	if sc, _ := propagator.SpanContextFromRequest(request); sc.IsValid() {
		ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
	} else {
		prop := propagation.TraceContext{}
		ctx = prop.Extract(ctx, propagation.HeaderCarrier(request.Header))
	}

	newCtx, err, cf := o.createExporter(ctx, path)
	if err != nil {
		panic(fmt.Sprintf("failed to startspan: %v", err))
	}
	return newCtx, cf
}
