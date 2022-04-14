package metrictest

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/reader"
)

type Exporter struct {
	lock    sync.Mutex
	metrics *reader.Metrics
}

var _ reader.Exporter = &Exporter{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Export(_ context.Context, metrics reader.Metrics) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.metrics = &metrics

	return nil
}

func (*Exporter) Flush(context.Context) error { return nil }

func (*Exporter) Shutdown(context.Context) error { return nil }

func (e *Exporter) GetByName(name string) (ExportRecord, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	for _, scope := range e.metrics.Scopes {
		for _, inst := range scope.Instruments {
			if inst.Descriptor.Name == name && len(inst.Points) > 0 {
				rec := pointToAggregation(inst.Points[len(inst.Points)-1])
				rec.InstrumentName = name
				rec.InstrumentationLibrary = scope.Library

				return rec, nil
			}
		}
	}
	return ExportRecord{}, fmt.Errorf("record not found")
}

func pointToAggregation(point reader.Point) ExportRecord {
	rec := ExportRecord{
		Attributes:          point.Attributes.ToSlice(),
		AggregationCatagory: point.Aggregation.Category(),
	}

	switch agg := point.Aggregation.(type) {
	case aggregation.Histogram:
		rec.Histogram = agg.Histogram()
		rec.Count = agg.Count()
		rec.Sum = agg.Sum()
	case aggregation.Sum:
		rec.Sum = agg.Sum()
	case aggregation.Count:
		rec.Count = agg.Count()
	case aggregation.Gauge:
		rec.Gauge = agg.Gauge()
	}

	return rec
}

// ExportRecord represents one collected datapoint from the Exporter.
type ExportRecord struct {
	InstrumentName         string
	InstrumentationLibrary instrumentation.Library
	Attributes             []attribute.KeyValue
	AggregationCatagory    aggregation.Category
	Sum                    number.Number
	Count                  uint64
	Histogram              aggregation.Buckets
	Gauge                  number.Number
}

type Reader struct {
	*reader.ManualReader

	exporter Exporter
	producer reader.Producer
}

func NewReader() *Reader {
	exp := Exporter{}
	return &Reader{
		ManualReader: reader.NewManualReader(&exp),

		exporter: exp,
	}
}

func (r *Reader) GetByName(name string) (ExportRecord, error) {
	return r.exporter.GetByName(name)
}

func (r *Reader) Register(prod reader.Producer) {
	r.producer = prod
	r.ManualReader.Register(prod)
}

func (r *Reader) Producer() reader.Producer {
	return r.producer
}
