// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/azure"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	inst metric.Int64UpDownCounter
}

// NewCosmosDBClientActiveInstanceCount returns a new
// CosmosDBClientActiveInstanceCount instrument.
func NewCosmosDBClientActiveInstanceCount(m metric.Meter) (CosmosDBClientActiveInstanceCount, error) {
	i, err := m.Int64UpDownCounter(
	    "azure.cosmosdb.client.active_instance.count",
	    metric.WithDescription("Number of active client instances"),
	    metric.WithUnit("{instance}"),
	)
	if err != nil {
	    return CosmosDBClientActiveInstanceCount{}, err
	}
	return CosmosDBClientActiveInstanceCount{i}, nil
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
	attrs ...CosmosDBClientActiveInstanceCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m CosmosDBClientActiveInstanceCount) conv(in []CosmosDBClientActiveInstanceCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.cosmosdbClientActiveInstanceCountAttr()
	}
	return out
}

// CosmosDBClientActiveInstanceCountAttr is an optional attribute for the
// CosmosDBClientActiveInstanceCount instrument.
type CosmosDBClientActiveInstanceCountAttr interface {
    cosmosdbClientActiveInstanceCountAttr() attribute.KeyValue
}

type cosmosdbClientActiveInstanceCountAttr struct {
	kv attribute.KeyValue
}

func (a cosmosdbClientActiveInstanceCountAttr) cosmosdbClientActiveInstanceCountAttr() attribute.KeyValue {
    return a.kv
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (CosmosDBClientActiveInstanceCount) ServerPortAttr(val int) CosmosDBClientActiveInstanceCountAttr {
	return cosmosdbClientActiveInstanceCountAttr{kv: attribute.Int("server.port", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the database host.
func (CosmosDBClientActiveInstanceCount) ServerAddressAttr(val string) CosmosDBClientActiveInstanceCountAttr {
	return cosmosdbClientActiveInstanceCountAttr{kv: attribute.String("server.address", val)}
}

// CosmosDBClientOperationRequestCharge is an instrument used to record metric
// values conforming to the "azure.cosmosdb.client.operation.request_charge"
// semantic conventions. It represents the [Request units] consumed by the
// operation.
//
// [Request units]: https://learn.microsoft.com/azure/cosmos-db/request-units
type CosmosDBClientOperationRequestCharge struct {
	inst metric.Int64Histogram
}

// NewCosmosDBClientOperationRequestCharge returns a new
// CosmosDBClientOperationRequestCharge instrument.
func NewCosmosDBClientOperationRequestCharge(m metric.Meter) (CosmosDBClientOperationRequestCharge, error) {
	i, err := m.Int64Histogram(
	    "azure.cosmosdb.client.operation.request_charge",
	    metric.WithDescription("[Request units](https://learn.microsoft.com/azure/cosmos-db/request-units) consumed by the operation"),
	    metric.WithUnit("{request_unit}"),
	)
	if err != nil {
	    return CosmosDBClientOperationRequestCharge{}, err
	}
	return CosmosDBClientOperationRequestCharge{i}, nil
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
	attrs ...CosmosDBClientOperationRequestChargeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m CosmosDBClientOperationRequestCharge) conv(in []CosmosDBClientOperationRequestChargeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.cosmosdbClientOperationRequestChargeAttr()
	}
	return out
}

// CosmosDBClientOperationRequestChargeAttr is an optional attribute for the
// CosmosDBClientOperationRequestCharge instrument.
type CosmosDBClientOperationRequestChargeAttr interface {
    cosmosdbClientOperationRequestChargeAttr() attribute.KeyValue
}

type cosmosdbClientOperationRequestChargeAttr struct {
	kv attribute.KeyValue
}

func (a cosmosdbClientOperationRequestChargeAttr) cosmosdbClientOperationRequestChargeAttr() attribute.KeyValue {
    return a.kv
}

// CosmosDBConsistencyLevel returns an optional attribute for the
// "azure.cosmosdb.consistency.level" semantic convention. It represents the
// account or request [consistency level].
//
// [consistency level]: https://learn.microsoft.com/azure/cosmos-db/consistency-levels
func (CosmosDBClientOperationRequestCharge) CosmosDBConsistencyLevelAttr(val CosmosDBConsistencyLevelAttr) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.String("azure.cosmosdb.consistency.level", string(val))}
}

// CosmosDBResponseSubStatusCode returns an optional attribute for the
// "azure.cosmosdb.response.sub_status_code" semantic convention. It represents
// the cosmos DB sub status code.
func (CosmosDBClientOperationRequestCharge) CosmosDBResponseSubStatusCodeAttr(val int) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.Int("azure.cosmosdb.response.sub_status_code", val)}
}

// DBCollectionName returns an optional attribute for the "db.collection.name"
// semantic convention. It represents the cosmos DB container name.
func (CosmosDBClientOperationRequestCharge) DBCollectionNameAttr(val string) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.String("db.collection.name", val)}
}

// DBNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (CosmosDBClientOperationRequestCharge) DBNamespaceAttr(val string) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.String("db.namespace", val)}
}

// DBOperationName returns an optional attribute for the "db.operation.name"
// semantic convention. It represents the name of the operation or command being
// executed.
func (CosmosDBClientOperationRequestCharge) DBOperationNameAttr(val string) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.String("db.operation.name", val)}
}

// DBResponseStatusCode returns an optional attribute for the
// "db.response.status_code" semantic convention. It represents the database
// response status code.
func (CosmosDBClientOperationRequestCharge) DBResponseStatusCodeAttr(val string) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.String("db.response.status_code", val)}
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (CosmosDBClientOperationRequestCharge) ErrorTypeAttr(val ErrorTypeAttr) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.String("error.type", string(val))}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (CosmosDBClientOperationRequestCharge) ServerPortAttr(val int) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.Int("server.port", val)}
}

// CosmosDBOperationContactedRegions returns an optional attribute for the
// "azure.cosmosdb.operation.contacted_regions" semantic convention. It
// represents the list of regions contacted during operation in the order that
// they were contacted. If there is more than one region listed, it indicates
// that the operation was performed on multiple regions i.e. cross-regional call.
func (CosmosDBClientOperationRequestCharge) CosmosDBOperationContactedRegionsAttr(val ...string) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.StringSlice("azure.cosmosdb.operation.contacted_regions", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the database host.
func (CosmosDBClientOperationRequestCharge) ServerAddressAttr(val string) CosmosDBClientOperationRequestChargeAttr {
	return cosmosdbClientOperationRequestChargeAttr{kv: attribute.String("server.address", val)}
}