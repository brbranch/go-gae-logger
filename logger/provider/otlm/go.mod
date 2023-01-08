module github.com/brbranch/go-gae-logger/provider/otlm

go 1.19

replace github.com/brbranch/go-gae-logger => ../../../

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.10.2
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator v0.34.2
	github.com/brbranch/go-gae-logger v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/contrib/detectors/gcp v1.12.0
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/sdk v1.11.2
	go.opentelemetry.io/otel/trace v1.11.2
)
