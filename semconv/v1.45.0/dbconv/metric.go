// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package dbconv provides types and functionality for OpenTelemetry semantic
// conventions in the "db" namespace.
package dbconv

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/internal/metricpool"
)

// ClientConnectionStateAttr is an attribute conforming to the
// db.client.connection.state semantic conventions. It represents the state of a
// connection in the pool.
type ClientConnectionStateAttr string

var (
	// ClientConnectionStateIdle is the standardized value "idle" of
	// ClientConnectionStateAttr.
	ClientConnectionStateIdle ClientConnectionStateAttr = "idle"
	// ClientConnectionStateUsed is the standardized value "used" of
	// ClientConnectionStateAttr.
	ClientConnectionStateUsed ClientConnectionStateAttr = "used"
)

// SystemNameAttr is an attribute conforming to the db.system.name semantic
// conventions. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
type SystemNameAttr string

var (
	// SystemNameOtherSQL is the some other SQL database. Fallback only.
	SystemNameOtherSQL SystemNameAttr = "other_sql"
	// SystemNameSoftwareagAdabas is the [Adabas (Adaptable Database System)].
	//
	// [Adabas (Adaptable Database System)]: https://documentation.softwareag.com/?pf=adabas
	SystemNameSoftwareagAdabas SystemNameAttr = "softwareag.adabas"
	// SystemNameActianIngres is the [Actian Ingres].
	//
	// [Actian Ingres]: https://www.actian.com/databases/ingres/
	SystemNameActianIngres SystemNameAttr = "actian.ingres"
	// SystemNameAWSDynamoDB is the [Amazon DynamoDB].
	//
	// [Amazon DynamoDB]: https://aws.amazon.com/pm/dynamodb/
	SystemNameAWSDynamoDB SystemNameAttr = "aws.dynamodb"
	// SystemNameAWSRedshift is the [Amazon Redshift].
	//
	// [Amazon Redshift]: https://aws.amazon.com/redshift/
	SystemNameAWSRedshift SystemNameAttr = "aws.redshift"
	// SystemNameAzureCosmosDB is the [Azure Cosmos DB].
	//
	// [Azure Cosmos DB]: https://learn.microsoft.com/azure/cosmos-db
	SystemNameAzureCosmosDB SystemNameAttr = "azure.cosmosdb"
	// SystemNameIntersystemsCache is the [InterSystems Caché].
	//
	// [InterSystems Caché]: https://www.intersystems.com/products/cache/
	SystemNameIntersystemsCache SystemNameAttr = "intersystems.cache"
	// SystemNameCassandra is the [Apache Cassandra].
	//
	// [Apache Cassandra]: https://cassandra.apache.org/
	SystemNameCassandra SystemNameAttr = "cassandra"
	// SystemNameClickHouse is the [ClickHouse].
	//
	// [ClickHouse]: https://clickhouse.com/
	SystemNameClickHouse SystemNameAttr = "clickhouse"
	// SystemNameCockroachDB is the [CockroachDB].
	//
	// [CockroachDB]: https://www.cockroachlabs.com/
	SystemNameCockroachDB SystemNameAttr = "cockroachdb"
	// SystemNameCouchbase is the [Couchbase].
	//
	// [Couchbase]: https://www.couchbase.com/
	SystemNameCouchbase SystemNameAttr = "couchbase"
	// SystemNameCouchDB is the [Apache CouchDB].
	//
	// [Apache CouchDB]: https://couchdb.apache.org/
	SystemNameCouchDB SystemNameAttr = "couchdb"
	// SystemNameDerby is the [Apache Derby].
	//
	// [Apache Derby]: https://db.apache.org/derby/
	SystemNameDerby SystemNameAttr = "derby"
	// SystemNameElasticsearch is the [Elasticsearch].
	//
	// [Elasticsearch]: https://www.elastic.co/elasticsearch
	SystemNameElasticsearch SystemNameAttr = "elasticsearch"
	// SystemNameFirebirdSQL is the [Firebird].
	//
	// [Firebird]: https://www.firebirdsql.org/
	SystemNameFirebirdSQL SystemNameAttr = "firebirdsql"
	// SystemNameGCPSpanner is the [Google Cloud Spanner].
	//
	// [Google Cloud Spanner]: https://cloud.google.com/spanner
	SystemNameGCPSpanner SystemNameAttr = "gcp.spanner"
	// SystemNameGeode is the [Apache Geode].
	//
	// [Apache Geode]: https://geode.apache.org/
	SystemNameGeode SystemNameAttr = "geode"
	// SystemNameH2database is the [H2 Database].
	//
	// [H2 Database]: https://h2database.com/
	SystemNameH2database SystemNameAttr = "h2database"
	// SystemNameHBase is the [Apache HBase].
	//
	// [Apache HBase]: https://hbase.apache.org/
	SystemNameHBase SystemNameAttr = "hbase"
	// SystemNameHive is the [Apache Hive].
	//
	// [Apache Hive]: https://hive.apache.org/
	SystemNameHive SystemNameAttr = "hive"
	// SystemNameHSQLDB is the [HyperSQL Database].
	//
	// [HyperSQL Database]: https://hsqldb.org/
	SystemNameHSQLDB SystemNameAttr = "hsqldb"
	// SystemNameIBMDB2 is the [IBM Db2].
	//
	// [IBM Db2]: https://www.ibm.com/products/db2
	SystemNameIBMDB2 SystemNameAttr = "ibm.db2"
	// SystemNameIBMInformix is the [IBM Informix].
	//
	// [IBM Informix]: https://www.ibm.com/products/informix
	SystemNameIBMInformix SystemNameAttr = "ibm.informix"
	// SystemNameIBMNetezza is the [IBM Netezza].
	//
	// [IBM Netezza]: https://www.ibm.com/products/netezza
	SystemNameIBMNetezza SystemNameAttr = "ibm.netezza"
	// SystemNameInfluxDB is the [InfluxDB].
	//
	// [InfluxDB]: https://www.influxdata.com/
	SystemNameInfluxDB SystemNameAttr = "influxdb"
	// SystemNameInstantDB is the [Instant].
	//
	// [Instant]: https://www.instantdb.com/
	SystemNameInstantDB SystemNameAttr = "instantdb"
	// SystemNameMariaDB is the [MariaDB].
	//
	// [MariaDB]: https://mariadb.org/
	SystemNameMariaDB SystemNameAttr = "mariadb"
	// SystemNameMemcached is the [Memcached].
	//
	// [Memcached]: https://memcached.org/
	SystemNameMemcached SystemNameAttr = "memcached"
	// SystemNameMongoDB is the [MongoDB].
	//
	// [MongoDB]: https://www.mongodb.com/
	SystemNameMongoDB SystemNameAttr = "mongodb"
	// SystemNameMicrosoftSQLServer is the [Microsoft SQL Server].
	//
	// [Microsoft SQL Server]: https://www.microsoft.com/sql-server
	SystemNameMicrosoftSQLServer SystemNameAttr = "microsoft.sql_server"
	// SystemNameMySQL is the [MySQL].
	//
	// [MySQL]: https://www.mysql.com/
	SystemNameMySQL SystemNameAttr = "mysql"
	// SystemNameNeo4j is the [Neo4j].
	//
	// [Neo4j]: https://neo4j.com/
	SystemNameNeo4j SystemNameAttr = "neo4j"
	// SystemNameOpenSearch is the [OpenSearch].
	//
	// [OpenSearch]: https://opensearch.org/
	SystemNameOpenSearch SystemNameAttr = "opensearch"
	// SystemNameOracleDB is the [Oracle Database].
	//
	// [Oracle Database]: https://www.oracle.com/database/
	SystemNameOracleDB SystemNameAttr = "oracle.db"
	// SystemNamePostgreSQL is the [PostgreSQL].
	//
	// [PostgreSQL]: https://www.postgresql.org/
	SystemNamePostgreSQL SystemNameAttr = "postgresql"
	// SystemNameRedis is the [Redis].
	//
	// [Redis]: https://redis.io/
	SystemNameRedis SystemNameAttr = "redis"
	// SystemNameSAPHANA is the [SAP HANA].
	//
	// [SAP HANA]: https://www.sap.com/products/technology-platform/hana/what-is-sap-hana.html
	SystemNameSAPHANA SystemNameAttr = "sap.hana"
	// SystemNameSAPMaxDB is the [SAP MaxDB].
	//
	// [SAP MaxDB]: https://maxdb.sap.com/
	SystemNameSAPMaxDB SystemNameAttr = "sap.maxdb"
	// SystemNameSQLite is the [SQLite].
	//
	// [SQLite]: https://www.sqlite.org/
	SystemNameSQLite SystemNameAttr = "sqlite"
	// SystemNameTeradata is the [Teradata].
	//
	// [Teradata]: https://www.teradata.com/
	SystemNameTeradata SystemNameAttr = "teradata"
	// SystemNameTrino is the [Trino].
	//
	// [Trino]: https://trino.io/
	SystemNameTrino SystemNameAttr = "trino"
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

// ClientConnectionCount is an instrument used to record metric values conforming
// to the "db.client.connection.count" semantic conventions. It represents the
// number of connections that are currently in state described by the `state`
// attribute.
//
// Stability: Development
type ClientConnectionCount struct {
	metric.Int64UpDownCounter
}

var newClientConnectionCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of connections that are currently in state described by the `state` attribute."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionCount returns a new ClientConnectionCount instrument.
//
// Stability: Development
func NewClientConnectionCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionCountOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.count",
		opt...,
	)
	if err != nil {
		return ClientConnectionCount{noop.Int64UpDownCounter{}}, err
	}
	return ClientConnectionCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionCount) Name() string {
	return "db.client.connection.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionCount) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionCount) Description() string {
	return "The number of connections that are currently in state described by the `state` attribute."
}

// Add adds incr to the existing count for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// # The clientConnectionState is the the state of a connection in the pool
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionCount) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
	clientConnectionState ClientConnectionStateAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			attribute.String("db.client.connection.state", string(clientConnectionState)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
				attribute.String("db.client.connection.state", string(clientConnectionState)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ClientConnectionCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionCount) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionCount) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionCount) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionCount) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionCountObservable is an instrument used to record metric values
// conforming to the "db.client.connection.count" semantic conventions. It
// represents the number of connections that are currently in state described by
// the `state` attribute.
//
// Stability: Development
type ClientConnectionCountObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClientConnectionCountObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of connections that are currently in state described by the `state` attribute."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionCountObservable returns a new
// ClientConnectionCountObservable instrument.
//
// Stability: Development
func NewClientConnectionCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClientConnectionCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionCountObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionCountObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionCountObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"db.client.connection.count",
		opt...,
	)
	if err != nil {
		return ClientConnectionCountObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClientConnectionCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionCountObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionCountObservable) Name() string {
	return "db.client.connection.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionCountObservable) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionCountObservable) Description() string {
	return "The number of connections that are currently in state described by the `state` attribute."
}

// AttrClientConnectionPoolName returns a required attribute for the
// "db.client.connection.pool.name" semantic convention. It represents the name
// of the connection pool. If the connection pool does not have a name,
// instrumentation SHOULD use `db.namespace`; otherwise `server.address`, with
// `server.port` appended using `:` if set. When appending a port to an IPv6
// address, instrumentation MUST enclose the address in square brackets. If
// `server.address` is not available, instrumentation SHOULD use `db.system.name`
// ; otherwise the connection pool library name.
func (ClientConnectionCountObservable) AttrClientConnectionPoolName(val string) attribute.KeyValue {
	return attribute.String("db.client.connection.pool.name", val)
}

// AttrClientConnectionState returns a required attribute for the
// "db.client.connection.state" semantic convention. It represents the state of a
// connection in the pool.
func (ClientConnectionCountObservable) AttrClientConnectionState(val ClientConnectionStateAttr) attribute.KeyValue {
	return attribute.String("db.client.connection.state", string(val))
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionCountObservable) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionCountObservable) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionCountObservable) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionCountObservable) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionCreateTime is an instrument used to record metric values
// conforming to the "db.client.connection.create_time" semantic conventions. It
// represents the time it took to create a new connection.
//
// Stability: Development
type ClientConnectionCreateTime struct {
	metric.Float64Histogram
}

var newClientConnectionCreateTimeOpts = []metric.Float64HistogramOption{
	metric.WithDescription("The time it took to create a new connection."),
	metric.WithUnit("s"),
}

// NewClientConnectionCreateTime returns a new ClientConnectionCreateTime
// instrument.
//
// Stability: Development
func NewClientConnectionCreateTime(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientConnectionCreateTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionCreateTime{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionCreateTimeOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionCreateTimeOpts...)
	}

	i, err := m.Float64Histogram(
		"db.client.connection.create_time",
		opt...,
	)
	if err != nil {
		return ClientConnectionCreateTime{noop.Float64Histogram{}}, err
	}
	return ClientConnectionCreateTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionCreateTime) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionCreateTime) Name() string {
	return "db.client.connection.create_time"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionCreateTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionCreateTime) Description() string {
	return "The time it took to create a new connection."
}

// Record records val to the current distribution for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionCreateTime) Record(
	ctx context.Context,
	val float64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientConnectionCreateTime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionCreateTime) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionCreateTime) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionCreateTime) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionCreateTime) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionIdleMax is an instrument used to record metric values
// conforming to the "db.client.connection.idle.max" semantic conventions. It
// represents the maximum number of idle open connections allowed.
//
// Stability: Development
type ClientConnectionIdleMax struct {
	metric.Int64UpDownCounter
}

var newClientConnectionIdleMaxOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The maximum number of idle open connections allowed."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionIdleMax returns a new ClientConnectionIdleMax instrument.
//
// Stability: Development
func NewClientConnectionIdleMax(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionIdleMax, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionIdleMax{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionIdleMaxOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionIdleMaxOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.idle.max",
		opt...,
	)
	if err != nil {
		return ClientConnectionIdleMax{noop.Int64UpDownCounter{}}, err
	}
	return ClientConnectionIdleMax{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionIdleMax) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionIdleMax) Name() string {
	return "db.client.connection.idle.max"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionIdleMax) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionIdleMax) Description() string {
	return "The maximum number of idle open connections allowed."
}

// Add adds incr to the existing count for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionIdleMax) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ClientConnectionIdleMax) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionIdleMax) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionIdleMax) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionIdleMax) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionIdleMax) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionIdleMaxObservable is an instrument used to record metric
// values conforming to the "db.client.connection.idle.max" semantic conventions.
// It represents the maximum number of idle open connections allowed.
//
// Stability: Development
type ClientConnectionIdleMaxObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClientConnectionIdleMaxObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The maximum number of idle open connections allowed."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionIdleMaxObservable returns a new
// ClientConnectionIdleMaxObservable instrument.
//
// Stability: Development
func NewClientConnectionIdleMaxObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClientConnectionIdleMaxObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionIdleMaxObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionIdleMaxObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionIdleMaxObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"db.client.connection.idle.max",
		opt...,
	)
	if err != nil {
		return ClientConnectionIdleMaxObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClientConnectionIdleMaxObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionIdleMaxObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionIdleMaxObservable) Name() string {
	return "db.client.connection.idle.max"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionIdleMaxObservable) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionIdleMaxObservable) Description() string {
	return "The maximum number of idle open connections allowed."
}

// AttrClientConnectionPoolName returns a required attribute for the
// "db.client.connection.pool.name" semantic convention. It represents the name
// of the connection pool. If the connection pool does not have a name,
// instrumentation SHOULD use `db.namespace`; otherwise `server.address`, with
// `server.port` appended using `:` if set. When appending a port to an IPv6
// address, instrumentation MUST enclose the address in square brackets. If
// `server.address` is not available, instrumentation SHOULD use `db.system.name`
// ; otherwise the connection pool library name.
func (ClientConnectionIdleMaxObservable) AttrClientConnectionPoolName(val string) attribute.KeyValue {
	return attribute.String("db.client.connection.pool.name", val)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionIdleMaxObservable) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionIdleMaxObservable) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionIdleMaxObservable) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionIdleMaxObservable) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionIdleMin is an instrument used to record metric values
// conforming to the "db.client.connection.idle.min" semantic conventions. It
// represents the minimum number of idle open connections allowed.
//
// Stability: Development
type ClientConnectionIdleMin struct {
	metric.Int64UpDownCounter
}

var newClientConnectionIdleMinOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The minimum number of idle open connections allowed."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionIdleMin returns a new ClientConnectionIdleMin instrument.
//
// Stability: Development
func NewClientConnectionIdleMin(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionIdleMin, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionIdleMin{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionIdleMinOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionIdleMinOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.idle.min",
		opt...,
	)
	if err != nil {
		return ClientConnectionIdleMin{noop.Int64UpDownCounter{}}, err
	}
	return ClientConnectionIdleMin{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionIdleMin) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionIdleMin) Name() string {
	return "db.client.connection.idle.min"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionIdleMin) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionIdleMin) Description() string {
	return "The minimum number of idle open connections allowed."
}

// Add adds incr to the existing count for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionIdleMin) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ClientConnectionIdleMin) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionIdleMin) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionIdleMin) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionIdleMin) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionIdleMin) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionIdleMinObservable is an instrument used to record metric
// values conforming to the "db.client.connection.idle.min" semantic conventions.
// It represents the minimum number of idle open connections allowed.
//
// Stability: Development
type ClientConnectionIdleMinObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClientConnectionIdleMinObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The minimum number of idle open connections allowed."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionIdleMinObservable returns a new
// ClientConnectionIdleMinObservable instrument.
//
// Stability: Development
func NewClientConnectionIdleMinObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClientConnectionIdleMinObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionIdleMinObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionIdleMinObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionIdleMinObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"db.client.connection.idle.min",
		opt...,
	)
	if err != nil {
		return ClientConnectionIdleMinObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClientConnectionIdleMinObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionIdleMinObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionIdleMinObservable) Name() string {
	return "db.client.connection.idle.min"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionIdleMinObservable) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionIdleMinObservable) Description() string {
	return "The minimum number of idle open connections allowed."
}

// AttrClientConnectionPoolName returns a required attribute for the
// "db.client.connection.pool.name" semantic convention. It represents the name
// of the connection pool. If the connection pool does not have a name,
// instrumentation SHOULD use `db.namespace`; otherwise `server.address`, with
// `server.port` appended using `:` if set. When appending a port to an IPv6
// address, instrumentation MUST enclose the address in square brackets. If
// `server.address` is not available, instrumentation SHOULD use `db.system.name`
// ; otherwise the connection pool library name.
func (ClientConnectionIdleMinObservable) AttrClientConnectionPoolName(val string) attribute.KeyValue {
	return attribute.String("db.client.connection.pool.name", val)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionIdleMinObservable) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionIdleMinObservable) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionIdleMinObservable) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionIdleMinObservable) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionLimit is an instrument used to record metric values conforming
// to the "db.client.connection.limit" semantic conventions. It represents the
// maximum number of open connections allowed.
//
// Stability: Development
type ClientConnectionLimit struct {
	metric.Int64UpDownCounter
}

var newClientConnectionLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The maximum number of open connections allowed."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionLimit returns a new ClientConnectionLimit instrument.
//
// Stability: Development
func NewClientConnectionLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionLimit{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionLimitOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.limit",
		opt...,
	)
	if err != nil {
		return ClientConnectionLimit{noop.Int64UpDownCounter{}}, err
	}
	return ClientConnectionLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionLimit) Name() string {
	return "db.client.connection.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionLimit) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionLimit) Description() string {
	return "The maximum number of open connections allowed."
}

// Add adds incr to the existing count for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionLimit) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ClientConnectionLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionLimit) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionLimit) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionLimit) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionLimit) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionLimitObservable is an instrument used to record metric values
// conforming to the "db.client.connection.limit" semantic conventions. It
// represents the maximum number of open connections allowed.
//
// Stability: Development
type ClientConnectionLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClientConnectionLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The maximum number of open connections allowed."),
	metric.WithUnit("{connection}"),
}

// NewClientConnectionLimitObservable returns a new
// ClientConnectionLimitObservable instrument.
//
// Stability: Development
func NewClientConnectionLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClientConnectionLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionLimitObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"db.client.connection.limit",
		opt...,
	)
	if err != nil {
		return ClientConnectionLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClientConnectionLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionLimitObservable) Name() string {
	return "db.client.connection.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionLimitObservable) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionLimitObservable) Description() string {
	return "The maximum number of open connections allowed."
}

// AttrClientConnectionPoolName returns a required attribute for the
// "db.client.connection.pool.name" semantic convention. It represents the name
// of the connection pool. If the connection pool does not have a name,
// instrumentation SHOULD use `db.namespace`; otherwise `server.address`, with
// `server.port` appended using `:` if set. When appending a port to an IPv6
// address, instrumentation MUST enclose the address in square brackets. If
// `server.address` is not available, instrumentation SHOULD use `db.system.name`
// ; otherwise the connection pool library name.
func (ClientConnectionLimitObservable) AttrClientConnectionPoolName(val string) attribute.KeyValue {
	return attribute.String("db.client.connection.pool.name", val)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionLimitObservable) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionLimitObservable) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionLimitObservable) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionLimitObservable) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionPendingRequests is an instrument used to record metric values
// conforming to the "db.client.connection.pending_requests" semantic
// conventions. It represents the number of current pending requests for an open
// connection.
//
// Stability: Development
type ClientConnectionPendingRequests struct {
	metric.Int64UpDownCounter
}

var newClientConnectionPendingRequestsOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of current pending requests for an open connection."),
	metric.WithUnit("{request}"),
}

// NewClientConnectionPendingRequests returns a new
// ClientConnectionPendingRequests instrument.
//
// Stability: Development
func NewClientConnectionPendingRequests(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionPendingRequests, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionPendingRequests{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionPendingRequestsOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionPendingRequestsOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.pending_requests",
		opt...,
	)
	if err != nil {
		return ClientConnectionPendingRequests{noop.Int64UpDownCounter{}}, err
	}
	return ClientConnectionPendingRequests{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionPendingRequests) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionPendingRequests) Name() string {
	return "db.client.connection.pending_requests"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionPendingRequests) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionPendingRequests) Description() string {
	return "The number of current pending requests for an open connection."
}

// Add adds incr to the existing count for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionPendingRequests) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ClientConnectionPendingRequests) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionPendingRequests) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionPendingRequests) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionPendingRequests) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionPendingRequests) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionPendingRequestsObservable is an instrument used to record
// metric values conforming to the "db.client.connection.pending_requests"
// semantic conventions. It represents the number of current pending requests for
// an open connection.
//
// Stability: Development
type ClientConnectionPendingRequestsObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClientConnectionPendingRequestsObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of current pending requests for an open connection."),
	metric.WithUnit("{request}"),
}

// NewClientConnectionPendingRequestsObservable returns a new
// ClientConnectionPendingRequestsObservable instrument.
//
// Stability: Development
func NewClientConnectionPendingRequestsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClientConnectionPendingRequestsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionPendingRequestsObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionPendingRequestsObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionPendingRequestsObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"db.client.connection.pending_requests",
		opt...,
	)
	if err != nil {
		return ClientConnectionPendingRequestsObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClientConnectionPendingRequestsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionPendingRequestsObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionPendingRequestsObservable) Name() string {
	return "db.client.connection.pending_requests"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionPendingRequestsObservable) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionPendingRequestsObservable) Description() string {
	return "The number of current pending requests for an open connection."
}

// AttrClientConnectionPoolName returns a required attribute for the
// "db.client.connection.pool.name" semantic convention. It represents the name
// of the connection pool. If the connection pool does not have a name,
// instrumentation SHOULD use `db.namespace`; otherwise `server.address`, with
// `server.port` appended using `:` if set. When appending a port to an IPv6
// address, instrumentation MUST enclose the address in square brackets. If
// `server.address` is not available, instrumentation SHOULD use `db.system.name`
// ; otherwise the connection pool library name.
func (ClientConnectionPendingRequestsObservable) AttrClientConnectionPoolName(val string) attribute.KeyValue {
	return attribute.String("db.client.connection.pool.name", val)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionPendingRequestsObservable) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionPendingRequestsObservable) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionPendingRequestsObservable) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionPendingRequestsObservable) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionTimeouts is an instrument used to record metric values
// conforming to the "db.client.connection.timeouts" semantic conventions. It
// represents the number of connection timeouts that have occurred trying to
// obtain a connection from the pool.
//
// Stability: Development
type ClientConnectionTimeouts struct {
	metric.Int64Counter
}

var newClientConnectionTimeoutsOpts = []metric.Int64CounterOption{
	metric.WithDescription("The number of connection timeouts that have occurred trying to obtain a connection from the pool."),
	metric.WithUnit("{timeout}"),
}

// NewClientConnectionTimeouts returns a new ClientConnectionTimeouts instrument.
//
// Stability: Development
func NewClientConnectionTimeouts(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientConnectionTimeouts, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionTimeouts{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionTimeoutsOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionTimeoutsOpts...)
	}

	i, err := m.Int64Counter(
		"db.client.connection.timeouts",
		opt...,
	)
	if err != nil {
		return ClientConnectionTimeouts{noop.Int64Counter{}}, err
	}
	return ClientConnectionTimeouts{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionTimeouts) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionTimeouts) Name() string {
	return "db.client.connection.timeouts"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionTimeouts) Unit() string {
	return "{timeout}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionTimeouts) Description() string {
	return "The number of connection timeouts that have occurred trying to obtain a connection from the pool."
}

// Add adds incr to the existing count for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionTimeouts) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ClientConnectionTimeouts) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionTimeouts) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionTimeouts) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionTimeouts) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionTimeouts) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionTimeoutsObservable is an instrument used to record metric
// values conforming to the "db.client.connection.timeouts" semantic conventions.
// It represents the number of connection timeouts that have occurred trying to
// obtain a connection from the pool.
//
// Stability: Development
type ClientConnectionTimeoutsObservable struct {
	metric.Int64ObservableCounter
}

var newClientConnectionTimeoutsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("The number of connection timeouts that have occurred trying to obtain a connection from the pool."),
	metric.WithUnit("{timeout}"),
}

// NewClientConnectionTimeoutsObservable returns a new
// ClientConnectionTimeoutsObservable instrument.
//
// Stability: Development
func NewClientConnectionTimeoutsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (ClientConnectionTimeoutsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionTimeoutsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionTimeoutsObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionTimeoutsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"db.client.connection.timeouts",
		opt...,
	)
	if err != nil {
		return ClientConnectionTimeoutsObservable{noop.Int64ObservableCounter{}}, err
	}
	return ClientConnectionTimeoutsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionTimeoutsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionTimeoutsObservable) Name() string {
	return "db.client.connection.timeouts"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionTimeoutsObservable) Unit() string {
	return "{timeout}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionTimeoutsObservable) Description() string {
	return "The number of connection timeouts that have occurred trying to obtain a connection from the pool."
}

// AttrClientConnectionPoolName returns a required attribute for the
// "db.client.connection.pool.name" semantic convention. It represents the name
// of the connection pool. If the connection pool does not have a name,
// instrumentation SHOULD use `db.namespace`; otherwise `server.address`, with
// `server.port` appended using `:` if set. When appending a port to an IPv6
// address, instrumentation MUST enclose the address in square brackets. If
// `server.address` is not available, instrumentation SHOULD use `db.system.name`
// ; otherwise the connection pool library name.
func (ClientConnectionTimeoutsObservable) AttrClientConnectionPoolName(val string) attribute.KeyValue {
	return attribute.String("db.client.connection.pool.name", val)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionTimeoutsObservable) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionTimeoutsObservable) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionTimeoutsObservable) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionTimeoutsObservable) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionUseTime is an instrument used to record metric values
// conforming to the "db.client.connection.use_time" semantic conventions. It
// represents the time between borrowing a connection and returning it to the
// pool.
//
// Stability: Development
type ClientConnectionUseTime struct {
	metric.Float64Histogram
}

var newClientConnectionUseTimeOpts = []metric.Float64HistogramOption{
	metric.WithDescription("The time between borrowing a connection and returning it to the pool."),
	metric.WithUnit("s"),
}

// NewClientConnectionUseTime returns a new ClientConnectionUseTime instrument.
//
// Stability: Development
func NewClientConnectionUseTime(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientConnectionUseTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionUseTime{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionUseTimeOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionUseTimeOpts...)
	}

	i, err := m.Float64Histogram(
		"db.client.connection.use_time",
		opt...,
	)
	if err != nil {
		return ClientConnectionUseTime{noop.Float64Histogram{}}, err
	}
	return ClientConnectionUseTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionUseTime) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionUseTime) Name() string {
	return "db.client.connection.use_time"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionUseTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionUseTime) Description() string {
	return "The time between borrowing a connection and returning it to the pool."
}

// Record records val to the current distribution for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionUseTime) Record(
	ctx context.Context,
	val float64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientConnectionUseTime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionUseTime) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionUseTime) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionUseTime) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionUseTime) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientConnectionWaitTime is an instrument used to record metric values
// conforming to the "db.client.connection.wait_time" semantic conventions. It
// represents the time it took to obtain an open connection from the pool.
//
// Stability: Development
type ClientConnectionWaitTime struct {
	metric.Float64Histogram
}

var newClientConnectionWaitTimeOpts = []metric.Float64HistogramOption{
	metric.WithDescription("The time it took to obtain an open connection from the pool."),
	metric.WithUnit("s"),
}

// NewClientConnectionWaitTime returns a new ClientConnectionWaitTime instrument.
//
// Stability: Development
func NewClientConnectionWaitTime(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientConnectionWaitTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionWaitTime{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientConnectionWaitTimeOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientConnectionWaitTimeOpts...)
	}

	i, err := m.Float64Histogram(
		"db.client.connection.wait_time",
		opt...,
	)
	if err != nil {
		return ClientConnectionWaitTime{noop.Float64Histogram{}}, err
	}
	return ClientConnectionWaitTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionWaitTime) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionWaitTime) Name() string {
	return "db.client.connection.wait_time"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionWaitTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionWaitTime) Description() string {
	return "The time it took to obtain an open connection from the pool."
}

// Record records val to the current distribution for attrs.
//
// The clientConnectionPoolName is the the name of the connection pool. If the
// connection pool does not have a name, instrumentation SHOULD use
// `db.namespace`; otherwise `server.address`, with `server.port` appended using
// `:` if set. When appending a port to an IPv6 address, instrumentation MUST
// enclose the address in square brackets. If `server.address` is not available,
// instrumentation SHOULD use `db.system.name`; otherwise the connection pool
// library name.
//
// All additional attrs passed are included in the recorded value.
func (m ClientConnectionWaitTime) Record(
	ctx context.Context,
	val float64,
	clientConnectionPoolName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val, metric.WithAttributes(
			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
		))
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientConnectionWaitTime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientConnectionWaitTime) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrSystemName returns an optional attribute for the "db.system.name" semantic
// convention. It represents the database management system (DBMS) product as
// identified by the client instrumentation.
func (ClientConnectionWaitTime) AttrSystemName(val SystemNameAttr) attribute.KeyValue {
	return attribute.String("db.system.name", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientConnectionWaitTime) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientConnectionWaitTime) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// ClientOperationDuration is an instrument used to record metric values
// conforming to the "db.client.operation.duration" semantic conventions. It
// represents the duration of database client operations.
//
// Stability: Stable
type ClientOperationDuration struct {
	metric.Float64Histogram
}

var newClientOperationDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("Duration of database client operations."),
	metric.WithUnit("s"),
}

// NewClientOperationDuration returns a new ClientOperationDuration instrument.
//
// Stability: Stable
func NewClientOperationDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientOperationDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientOperationDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientOperationDurationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientOperationDurationOpts...)
	}

	i, err := m.Float64Histogram(
		"db.client.operation.duration",
		opt...,
	)
	if err != nil {
		return ClientOperationDuration{noop.Float64Histogram{}}, err
	}
	return ClientOperationDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientOperationDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientOperationDuration) Name() string {
	return "db.client.operation.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOperationDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientOperationDuration) Description() string {
	return "Duration of database client operations."
}

// Record records val to the current distribution for attrs.
//
// The systemName is the the database management system (DBMS) product as
// identified by the client instrumentation.
//
// All additional attrs passed are included in the recorded value.
//
// Batch operations SHOULD be recorded as a single operation.
func (m ClientOperationDuration) Record(
	ctx context.Context,
	val float64,
	systemName SystemNameAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val, metric.WithAttributes(
			attribute.String("db.system.name", string(systemName)),
		))
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.system.name", string(systemName)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Batch operations SHOULD be recorded as a single operation.
func (m ClientOperationDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Histogram.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrCollectionName returns an optional attribute for the "db.collection.name"
// semantic convention. It represents the name of a collection (table, container)
// within the database.
func (ClientOperationDuration) AttrCollectionName(val string) attribute.KeyValue {
	return attribute.String("db.collection.name", val)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientOperationDuration) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrOperationName returns an optional attribute for the "db.operation.name"
// semantic convention. It represents the name of the operation or command being
// executed.
func (ClientOperationDuration) AttrOperationName(val string) attribute.KeyValue {
	return attribute.String("db.operation.name", val)
}

// AttrResponseStatusCode returns an optional attribute for the
// "db.response.status_code" semantic convention. It represents the database
// response status code.
func (ClientOperationDuration) AttrResponseStatusCode(val string) attribute.KeyValue {
	return attribute.String("db.response.status_code", val)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientOperationDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientOperationDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrQuerySummary returns an optional attribute for the "db.query.summary"
// semantic convention. It represents the low cardinality summary of a database
// query.
func (ClientOperationDuration) AttrQuerySummary(val string) attribute.KeyValue {
	return attribute.String("db.query.summary", val)
}

// AttrStoredProcedureName returns an optional attribute for the
// "db.stored_procedure.name" semantic convention. It represents the name of a
// stored procedure within the database.
func (ClientOperationDuration) AttrStoredProcedureName(val string) attribute.KeyValue {
	return attribute.String("db.stored_procedure.name", val)
}

// AttrNetworkPeerAddress returns an optional attribute for the
// "network.peer.address" semantic convention. It represents the peer address of
// the database node where the operation was performed.
func (ClientOperationDuration) AttrNetworkPeerAddress(val string) attribute.KeyValue {
	return attribute.String("network.peer.address", val)
}

// AttrNetworkPeerPort returns an optional attribute for the "network.peer.port"
// semantic convention. It represents the peer port number of the network
// connection.
func (ClientOperationDuration) AttrNetworkPeerPort(val int) attribute.KeyValue {
	return attribute.Int("network.peer.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientOperationDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrQueryText returns an optional attribute for the "db.query.text" semantic
// convention. It represents the database query being executed.
func (ClientOperationDuration) AttrQueryText(val string) attribute.KeyValue {
	return attribute.String("db.query.text", val)
}

// ClientResponseReturnedRows is an instrument used to record metric values
// conforming to the "db.client.response.returned_rows" semantic conventions. It
// represents the actual number of records returned by the database operation.
//
// Stability: Development
type ClientResponseReturnedRows struct {
	metric.Int64Histogram
}

var newClientResponseReturnedRowsOpts = []metric.Int64HistogramOption{
	metric.WithDescription("The actual number of records returned by the database operation."),
	metric.WithUnit("{row}"),
}

// NewClientResponseReturnedRows returns a new ClientResponseReturnedRows
// instrument.
//
// Stability: Development
func NewClientResponseReturnedRows(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (ClientResponseReturnedRows, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientResponseReturnedRows{noop.Int64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientResponseReturnedRowsOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newClientResponseReturnedRowsOpts...)
	}

	i, err := m.Int64Histogram(
		"db.client.response.returned_rows",
		opt...,
	)
	if err != nil {
		return ClientResponseReturnedRows{noop.Int64Histogram{}}, err
	}
	return ClientResponseReturnedRows{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientResponseReturnedRows) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientResponseReturnedRows) Name() string {
	return "db.client.response.returned_rows"
}

// Unit returns the semantic convention unit of the instrument
func (ClientResponseReturnedRows) Unit() string {
	return "{row}"
}

// Description returns the semantic convention description of the instrument
func (ClientResponseReturnedRows) Description() string {
	return "The actual number of records returned by the database operation."
}

// Record records val to the current distribution for attrs.
//
// The systemName is the the database management system (DBMS) product as
// identified by the client instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientResponseReturnedRows) Record(
	ctx context.Context,
	val int64,
	systemName SystemNameAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Histogram.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Histogram.Record(ctx, val, metric.WithAttributes(
			attribute.String("db.system.name", string(systemName)),
		))
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("db.system.name", string(systemName)),
			)...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientResponseReturnedRows) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if !m.Int64Histogram.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Histogram.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrCollectionName returns an optional attribute for the "db.collection.name"
// semantic convention. It represents the name of a collection (table, container)
// within the database.
func (ClientResponseReturnedRows) AttrCollectionName(val string) attribute.KeyValue {
	return attribute.String("db.collection.name", val)
}

// AttrNamespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientResponseReturnedRows) AttrNamespace(val string) attribute.KeyValue {
	return attribute.String("db.namespace", val)
}

// AttrOperationName returns an optional attribute for the "db.operation.name"
// semantic convention. It represents the name of the operation or command being
// executed.
func (ClientResponseReturnedRows) AttrOperationName(val string) attribute.KeyValue {
	return attribute.String("db.operation.name", val)
}

// AttrResponseStatusCode returns an optional attribute for the
// "db.response.status_code" semantic convention. It represents the database
// response status code.
func (ClientResponseReturnedRows) AttrResponseStatusCode(val string) attribute.KeyValue {
	return attribute.String("db.response.status_code", val)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientResponseReturnedRows) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientResponseReturnedRows) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrQuerySummary returns an optional attribute for the "db.query.summary"
// semantic convention. It represents the low cardinality summary of a database
// query.
func (ClientResponseReturnedRows) AttrQuerySummary(val string) attribute.KeyValue {
	return attribute.String("db.query.summary", val)
}

// AttrNetworkPeerAddress returns an optional attribute for the
// "network.peer.address" semantic convention. It represents the peer address of
// the database node where the operation was performed.
func (ClientResponseReturnedRows) AttrNetworkPeerAddress(val string) attribute.KeyValue {
	return attribute.String("network.peer.address", val)
}

// AttrNetworkPeerPort returns an optional attribute for the "network.peer.port"
// semantic convention. It represents the peer port number of the network
// connection.
func (ClientResponseReturnedRows) AttrNetworkPeerPort(val int) attribute.KeyValue {
	return attribute.Int("network.peer.port", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the database address specified in the
// client configuration.
func (ClientResponseReturnedRows) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrQueryText returns an optional attribute for the "db.query.text" semantic
// convention. It represents the database query being executed.
func (ClientResponseReturnedRows) AttrQueryText(val string) attribute.KeyValue {
	return attribute.String("db.query.text", val)
}
