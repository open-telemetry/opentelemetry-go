// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "db" namespace.
package dbconv

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

// ClientConnectionStateAttr is an attribute conforming to the
// db.client.connection.state semantic conventions. It represents the state of a
// connection in the pool.
type ClientConnectionStateAttr string

var (
	// ClientConnectionStateIdle is the none.
	ClientConnectionStateIdle ClientConnectionStateAttr = "idle"
	// ClientConnectionStateUsed is the none.
	ClientConnectionStateUsed ClientConnectionStateAttr = "used"
)

// ClientConnectionsStateAttr is an attribute conforming to the
// db.client.connections.state semantic conventions. It represents the
// deprecated, use `db.client.connection.state` instead.
type ClientConnectionsStateAttr string

var (
	// ClientConnectionsStateIdle is the none.
	ClientConnectionsStateIdle ClientConnectionsStateAttr = "idle"
	// ClientConnectionsStateUsed is the none.
	ClientConnectionsStateUsed ClientConnectionsStateAttr = "used"
)

// CosmosDBConsistencyLevelAttr is an attribute conforming to the
// db.cosmosdb.consistency_level semantic conventions. It represents the
// deprecated, use `cosmosdb.consistency.level` instead.
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
	// [IBM Db2]: https://www.ibm.com/db2
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
type ClientConnectionCount struct {
	metric.Int64UpDownCounter
}

// NewClientConnectionCount returns a new ClientConnectionCount instrument.
func NewClientConnectionCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of connections that are currently in state described by the `state` attribute"),
			metric.WithUnit("{connection}"),
		}, opt...)...,
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
	return "The number of connections that are currently in state described by the `state` attribute"
}

// Add adds incr to the existing count.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
//
// The clientConnectionState is the the state of a connection in the pool
func (m ClientConnectionCount) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
	clientConnectionState ClientConnectionStateAttr,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
				attribute.String("db.client.connection.state", string(clientConnectionState)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ClientConnectionCreateTime is an instrument used to record metric values
// conforming to the "db.client.connection.create_time" semantic conventions. It
// represents the time it took to create a new connection.
type ClientConnectionCreateTime struct {
	metric.Float64Histogram
}

// NewClientConnectionCreateTime returns a new ClientConnectionCreateTime
// instrument.
func NewClientConnectionCreateTime(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientConnectionCreateTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionCreateTime{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"db.client.connection.create_time",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("The time it took to create a new connection"),
			metric.WithUnit("s"),
		}, opt...)...,
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
	return "The time it took to create a new connection"
}

// Record records val to the current distribution.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionCreateTime) Record(
	ctx context.Context,
	val float64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// ClientConnectionIdleMax is an instrument used to record metric values
// conforming to the "db.client.connection.idle.max" semantic conventions. It
// represents the maximum number of idle open connections allowed.
type ClientConnectionIdleMax struct {
	metric.Int64UpDownCounter
}

// NewClientConnectionIdleMax returns a new ClientConnectionIdleMax instrument.
func NewClientConnectionIdleMax(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionIdleMax, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionIdleMax{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.idle.max",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The maximum number of idle open connections allowed"),
			metric.WithUnit("{connection}"),
		}, opt...)...,
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
	return "The maximum number of idle open connections allowed"
}

// Add adds incr to the existing count.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionIdleMax) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ClientConnectionIdleMin is an instrument used to record metric values
// conforming to the "db.client.connection.idle.min" semantic conventions. It
// represents the minimum number of idle open connections allowed.
type ClientConnectionIdleMin struct {
	metric.Int64UpDownCounter
}

// NewClientConnectionIdleMin returns a new ClientConnectionIdleMin instrument.
func NewClientConnectionIdleMin(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionIdleMin, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionIdleMin{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.idle.min",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The minimum number of idle open connections allowed"),
			metric.WithUnit("{connection}"),
		}, opt...)...,
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
	return "The minimum number of idle open connections allowed"
}

// Add adds incr to the existing count.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionIdleMin) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ClientConnectionMax is an instrument used to record metric values conforming
// to the "db.client.connection.max" semantic conventions. It represents the
// maximum number of open connections allowed.
type ClientConnectionMax struct {
	metric.Int64UpDownCounter
}

// NewClientConnectionMax returns a new ClientConnectionMax instrument.
func NewClientConnectionMax(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionMax, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionMax{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.max",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The maximum number of open connections allowed"),
			metric.WithUnit("{connection}"),
		}, opt...)...,
	)
	if err != nil {
	    return ClientConnectionMax{noop.Int64UpDownCounter{}}, err
	}
	return ClientConnectionMax{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientConnectionMax) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionMax) Name() string {
	return "db.client.connection.max"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionMax) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionMax) Description() string {
	return "The maximum number of open connections allowed"
}

// Add adds incr to the existing count.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionMax) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ClientConnectionPendingRequests is an instrument used to record metric values
// conforming to the "db.client.connection.pending_requests" semantic
// conventions. It represents the number of current pending requests for an open
// connection.
type ClientConnectionPendingRequests struct {
	metric.Int64UpDownCounter
}

// NewClientConnectionPendingRequests returns a new
// ClientConnectionPendingRequests instrument.
func NewClientConnectionPendingRequests(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClientConnectionPendingRequests, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionPendingRequests{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"db.client.connection.pending_requests",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of current pending requests for an open connection"),
			metric.WithUnit("{request}"),
		}, opt...)...,
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
	return "The number of current pending requests for an open connection"
}

// Add adds incr to the existing count.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionPendingRequests) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ClientConnectionTimeouts is an instrument used to record metric values
// conforming to the "db.client.connection.timeouts" semantic conventions. It
// represents the number of connection timeouts that have occurred trying to
// obtain a connection from the pool.
type ClientConnectionTimeouts struct {
	metric.Int64Counter
}

// NewClientConnectionTimeouts returns a new ClientConnectionTimeouts instrument.
func NewClientConnectionTimeouts(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientConnectionTimeouts, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionTimeouts{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"db.client.connection.timeouts",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of connection timeouts that have occurred trying to obtain a connection from the pool"),
			metric.WithUnit("{timeout}"),
		}, opt...)...,
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
	return "The number of connection timeouts that have occurred trying to obtain a connection from the pool"
}

// Add adds incr to the existing count.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionTimeouts) Add(
	ctx context.Context,
	incr int64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// ClientConnectionUseTime is an instrument used to record metric values
// conforming to the "db.client.connection.use_time" semantic conventions. It
// represents the time between borrowing a connection and returning it to the
// pool.
type ClientConnectionUseTime struct {
	metric.Float64Histogram
}

// NewClientConnectionUseTime returns a new ClientConnectionUseTime instrument.
func NewClientConnectionUseTime(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientConnectionUseTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionUseTime{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"db.client.connection.use_time",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("The time between borrowing a connection and returning it to the pool"),
			metric.WithUnit("s"),
		}, opt...)...,
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
	return "The time between borrowing a connection and returning it to the pool"
}

// Record records val to the current distribution.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionUseTime) Record(
	ctx context.Context,
	val float64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// ClientConnectionWaitTime is an instrument used to record metric values
// conforming to the "db.client.connection.wait_time" semantic conventions. It
// represents the time it took to obtain an open connection from the pool.
type ClientConnectionWaitTime struct {
	metric.Float64Histogram
}

// NewClientConnectionWaitTime returns a new ClientConnectionWaitTime instrument.
func NewClientConnectionWaitTime(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientConnectionWaitTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientConnectionWaitTime{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"db.client.connection.wait_time",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("The time it took to obtain an open connection from the pool"),
			metric.WithUnit("s"),
		}, opt...)...,
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
	return "The time it took to obtain an open connection from the pool"
}

// Record records val to the current distribution.
//
// The clientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
func (m ClientConnectionWaitTime) Record(
	ctx context.Context,
	val float64,
	clientConnectionPoolName string,
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
			append(
				attrs,
				attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// ClientOperationDuration is an instrument used to record metric values
// conforming to the "db.client.operation.duration" semantic conventions. It
// represents the duration of database client operations.
type ClientOperationDuration struct {
	metric.Float64Histogram
}

// NewClientOperationDuration returns a new ClientOperationDuration instrument.
func NewClientOperationDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientOperationDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientOperationDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"db.client.operation.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Duration of database client operations."),
			metric.WithUnit("s"),
		}, opt...)...,
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

// Record records val to the current distribution.
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
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("db.system.name", string(systemName)),
			)...,
		),
	)

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
// semantic convention. It represents the low cardinality representation of a
// database query text.
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
// semantic convention. It represents the name of the database host.
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
type ClientResponseReturnedRows struct {
	metric.Int64Histogram
}

// NewClientResponseReturnedRows returns a new ClientResponseReturnedRows
// instrument.
func NewClientResponseReturnedRows(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (ClientResponseReturnedRows, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientResponseReturnedRows{noop.Int64Histogram{}}, nil
	}

	i, err := m.Int64Histogram(
		"db.client.response.returned_rows",
		append([]metric.Int64HistogramOption{
			metric.WithDescription("The actual number of records returned by the database operation."),
			metric.WithUnit("{row}"),
		}, opt...)...,
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

// Record records val to the current distribution.
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
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("db.system.name", string(systemName)),
			)...,
		),
	)

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
// semantic convention. It represents the low cardinality representation of a
// database query text.
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
// semantic convention. It represents the name of the database host.
func (ClientResponseReturnedRows) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrQueryText returns an optional attribute for the "db.query.text" semantic
// convention. It represents the database query being executed.
func (ClientResponseReturnedRows) AttrQueryText(val string) attribute.KeyValue {
	return attribute.String("db.query.text", val)
}