package provider

import (
	"context"
	"net/http"

	"github.com/brbranch/go-gae-logger/logger/model"
)

const providerValue = "___log_provider"

// Provider is a provider that provide a Tracer.
type Provider interface {
	// StartSpan start root span of Google Cloud Trace.
	StartSpan(request *http.Request, path string) (context.Context, func())
	// GetSpan get current span information.
	GetSpan(context context.Context) *model.Span
	// CustomSpan start child span of Google Cloud Trace.
	CustomSpan(c context.Context, label string) (context.Context, func())
	// ProjectID get GCP Project ID.
	ProjectID() string
}

func Set(ctx context.Context, provider Provider) context.Context {
	return context.WithValue(ctx, providerValue, provider)
}

func Get(context context.Context) Provider {
	provider, ok := context.Value(providerValue).(Provider)
	if !ok {
		return nil
	}
	return provider
}
