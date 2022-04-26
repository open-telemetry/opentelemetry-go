package metric

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/data"
)

// TODO: incomplete
func TestOutputReuse(t *testing.T) {
	ctx := context.Background()

	rdr := NewManualReader("test")
	provider := New(WithReader(rdr))

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	var reuse data.Metrics

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = rdr.Produce(&reuse)

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = rdr.Produce(&reuse)
}
