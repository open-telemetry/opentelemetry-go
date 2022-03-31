package metric

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

type testExporter struct {
	producer reader.Producer
}

func (t *testExporter) Register(producer reader.Producer) {
	t.producer = producer
}

func (*testExporter) Flush(context.Context) error { return nil }

func (*testExporter) Shutdown(context.Context) error { return nil }

// TODO: incomplete
func TestOutputReuse(t *testing.T) {
	ctx := context.Background()
	exp := &testExporter{}
	rdr := reader.New(exp)
	provider := New(WithReader(rdr))

	cntr, _ := provider.Meter("test").SyncInt64().Counter("hello")

	var reuse reader.Metrics

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = exp.producer.Produce(&reuse)

	cntr.Add(ctx, 1, attribute.Int("K", 1))

	reuse = exp.producer.Produce(&reuse)
}
