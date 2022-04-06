package metric

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

// TODO: incomplete
func TestOutputReuse(t *testing.T) {
	ctx := context.Background()
	exp := metrictest.NewExporter()
	rdr := reader.New(exp)
	provider := New(WithReader(rdr))

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	var reuse reader.Metrics

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = exp.Produce(&reuse)

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = exp.Produce(&reuse)
}
