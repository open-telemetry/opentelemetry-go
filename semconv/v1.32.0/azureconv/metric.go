// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "azure" namespace.
package azureconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
)

// CosmosDBConsistencyLevelAttr is an attribute conforming to the
// azure.cosmosdb.consistency.level semantic conventions. It represents the
// account or request [consistency level].
//
// [consistency level]: https://learn.microsoft.com/azure/cosmos-db/consistency-levels
type CosmosDBConsistencyLevelAttr string

var (
	// CosmosDBConsistencyLevelStrong is the none.
	CosmosDBConsistencyLevelStrong CosmosDBConsistencyLevelAttr = "Strong"
	// CosmosDBConsistencyLevelBoundedStaleness is the none.
	CosmosDBConsistencyLevelBoundedStaleness CosmosDBConsistencyLevelAttr = "BoundedStaleness"
	// CosmosDBConsistencyLevelSession is the none.
	CosmosDBConsistencyLevelSession CosmosDBConsistencyLevelAttr = "Session"
	// CosmosDBConsistencyLevelEventual is the none.
	CosmosDBConsistencyLevelEventual CosmosDBConsistencyLevelAttr = "Eventual"
	// CosmosDBConsistencyLevelConsistentPrefix is the none.
	CosmosDBConsistencyLevelConsistentPrefix CosmosDBConsistencyLevelAttr = "ConsistentPrefix"
)

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the describes a class of error the operation ended
// with.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

// CosmosDBClientActiveInstanceCount is an instrument used to record metric
// values conforming to the "azure.cosmosdb.client.active_instance.count"
// semantic conventions. It represents the number of active client instances.
type CosmosDBClientActiveInstanceCount struct {
	metric.Int64UpDownCounter
}

// NewCosmosDBClientActiveInstanceCount returns a new
// CosmosDBClientActiveInstanceCount instrument.
func NewCosmosDBClientActiveInstanceCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (CosmosDBClientActiveInstanceCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return CosmosDBClientActiveInstanceCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"azure.cosmosdb.client.active_instance.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of active client instances"),
			metric.WithUnit("{instance}"),
		}, opt...)...,
	)
	if err != nil {
	    return CosmosDBClientActiveInstanceCount{noop.Int64UpDownCounter{}}, err
	}
	return CosmosDBClientActiveInstanceCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CosmosDBClientActiveInstanceCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (CosmosDBClientActiveInstanceCount) Name() string {
	return "azure.cosmosdb.client.active_instance.count"
}

// Unit returns the semantic convention unit of the instrument
func (CosmosDBClientActiveInstanceCount) Unit() string {
	return "{instance}"
}

// Description returns the semantic convention description of the instrument
func (CosmosDBClientActiveInstanceCount) Description() string {
	return "Number of active client instances"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m CosmosDBClientActiveInstanceCount) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (CosmosDBClientActiveInstanceCount) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the name of the database host.
func (CosmosDBClientActiveInstanceCount) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// CosmosDBClientOperationRequestCharge is an instrument used to record metric
// values conforming to the "azure.cosmosdb.client.operation.request_charge"
// semantic conventions. It represents the [Request units] consumed by the
// operation.
//
// [Request units]: https://learn.microsoft.com/azure/cosmos-db/request-units
type CosmosDBClientOperationRequestCharge struct {
	metric.Int64Histogram
}

// NewCosmosDBClientOperationRequestCharge returns a new
// CosmosDBClientOperationRequestCharge instrument.
func NewCosmosDBClientOperationRequestCharge(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (CosmosDBClientOperationRequestCharge, error) {
	// Check if the meter is nil.
	if m == nil {
		return CosmosDBClientOperationRequestCharge{noop.Int64Histogram{}}, nil
	}

	i, err := m.Int64Histogram(
		"azure.cosmosdb.client.operation.request_charge",
		append([]metric.Int64HistogramOption{
			metric.WithDescription("[Request units](https://learn.microsoft.com/azure/cosmos-db/request-units) consumed by the operation"),
			metric.WithUnit("{request_unit}"),
		}, opt...)...,
	)
	if err != nil {
	    return CosmosDBClientOperationRequestCharge{noop.Int64Histogram{}}, err
	}
	return CosmosDBClientOperationRequestCharge{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CosmosDBClientOperationRequestCharge) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (CosmosDBClientOperationRequestCharge) Name() string {
	return "azure.cosmosdb.client.operation.request_charge"
}

// Unit returns the semantic convention unit of the instrument
func (CosmosDBClientOperationRequestCharge) Unit() string {
	return "{request_unit}"
}

// Description returns the semantic convention description of the instrument
func (CosmosDBClientOperationRequestCharge) Description() string {
	return "[Request units](https://learn.microsoft.com/azure/cosmos-db/request-units) consumed by the operation"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m CosmosDBClientOperationRequestCharge) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrCosmosDBConsistencyLevel returns an optional attribute for the
// "azure.cosmosdb.consistency.level" semantic convention. It represents the
// account or request [consistency level].
//
// [consistency level]: https://learn.microsoft.com/azure/cosmos-db/consistency-levels
func (CosmosDBClientOperationRequestCharge) AttrCosmosDBConsistencyLevel(val CosmosDBConsistencyLevelAttr) attribute.KeyValue {
	return attribute.String("azure.cosmosdb.consistency.level", string(val))
}

// AttrCosmosDBResponseSubStatusCode returns an optional attribute for the
// "azure.cosmosdb.response.sub_status_code" semantic convention. It represents
// the cosmos DB sub status code.
func (CosmosDBClientOperationRequestCharge) AttrCosmosDBResponseSubStatusCode(val int) attribute.KeyValue {
	return attribute.Int("azure.cosmosdb.response.sub_status_code", val)
}

// AttrDBCollectionName returns an optional attribute for the
// "db.collection.name" semantic convention. It represents the cosmos DB
// container name.
func (CosmosDBClientOperationRequestCharge) AttrDBCollectionName(val string) attribute.KeyValue {
	return attribute.String("db.collection.name", val)
}

// AttrDBNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (CosmosDBClientOperationRequestCharge) AttrDBNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrDBOperationName returns an optional attribute for the "db.operation.name"
// semantic convention. It represents the name of the operation or command being
// executed.
func (CosmosDBClientOperationRequestCharge) AttrDBOperationName(val string) attribute.KeyValue {
	return attribute.String("db.operation.name", val)
}

// AttrDBResponseStatusCode returns an optional attribute for the
// "db.response.status_code" semantic convention. It represents the database
// response status code.
func (CosmosDBClientOperationRequestCharge) AttrDBResponseStatusCode(val string) attribute.KeyValue {
	return attribute.String("db.response.status_code", val)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (CosmosDBClientOperationRequestCharge) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (CosmosDBClientOperationRequestCharge) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrCosmosDBOperationContactedRegions returns an optional attribute for the
// "azure.cosmosdb.operation.contacted_regions" semantic convention. It
// represents the list of regions contacted during operation in the order that
// they were contacted. If there is more than one region listed, it indicates
// that the operation was performed on multiple regions i.e. cross-regional call.
func (CosmosDBClientOperationRequestCharge) AttrCosmosDBOperationContactedRegions(val ...string) attribute.KeyValue {
	return attribute.StringSlice("azure.cosmosdb.operation.contacted_regions", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the name of the database host.
func (CosmosDBClientOperationRequestCharge) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}