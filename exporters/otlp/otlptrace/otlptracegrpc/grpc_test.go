package otlptracegrpc

import (
	"context"
	"os"
	"testing"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/otlpconfig"
)

func TestInvalidEndpointPath(t *testing.T) {
	otlpconfig.ResetEnvConfigErr()
	// os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://example.com:4317/v1/traces")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://example.com:4317")
	defer os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	_, err := New(context.Background())
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}
