// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/db"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	// SystemNameClickhouse is the [ClickHouse].
	//
	// [ClickHouse]: https://clickhouse.com/
	SystemNameClickhouse SystemNameAttr = "clickhouse"
	// SystemNameCockroachdb is the [CockroachDB].
	//
	// [CockroachDB]: https://www.cockroachlabs.com/
	SystemNameCockroachdb SystemNameAttr = "cockroachdb"
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
	// SystemNameFirebirdsql is the [Firebird].
	//
	// [Firebird]: https://www.firebirdsql.org/
	SystemNameFirebirdsql SystemNameAttr = "firebirdsql"
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
	// SystemNameIbmDb2 is the [IBM Db2].
	//
	// [IBM Db2]: https://www.ibm.com/db2
	SystemNameIbmDb2 SystemNameAttr = "ibm.db2"
	// SystemNameIbmInformix is the [IBM Informix].
	//
	// [IBM Informix]: https://www.ibm.com/products/informix
	SystemNameIbmInformix SystemNameAttr = "ibm.informix"
	// SystemNameIbmNetezza is the [IBM Netezza].
	//
	// [IBM Netezza]: https://www.ibm.com/products/netezza
	SystemNameIbmNetezza SystemNameAttr = "ibm.netezza"
	// SystemNameInfluxdb is the [InfluxDB].
	//
	// [InfluxDB]: https://www.influxdata.com/
	SystemNameInfluxdb SystemNameAttr = "influxdb"
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
	// SystemNameOpensearch is the [OpenSearch].
	//
	// [OpenSearch]: https://opensearch.org/
	SystemNameOpensearch SystemNameAttr = "opensearch"
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
	// SystemNameSapHana is the [SAP HANA].
	//
	// [SAP HANA]: https://www.sap.com/products/technology-platform/hana/what-is-sap-hana.html
	SystemNameSapHana SystemNameAttr = "sap.hana"
	// SystemNameSapMaxDB is the [SAP MaxDB].
	//
	// [SAP MaxDB]: https://maxdb.sap.com/
	SystemNameSapMaxDB SystemNameAttr = "sap.maxdb"
	// SystemNameSqlite is the [SQLite].
	//
	// [SQLite]: https://www.sqlite.org/
	SystemNameSqlite SystemNameAttr = "sqlite"
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

// DBClientConnectionCount is an instrument used to record metric values
// conforming to the "db.client.connection.count" semantic conventions. It
// represents the number of connections that are currently in state described by
// the `state` attribute.
type ClientConnectionCount struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionCount returns a new ClientConnectionCount instrument.
func NewClientConnectionCount(m metric.Meter) (ClientConnectionCount, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connection.count",
	    metric.WithDescription("The number of connections that are currently in state described by the `state` attribute"),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionCount{}, err
	}
	return ClientConnectionCount{i}, nil
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
// The dbClientConnectionPoolName is the the name of the connection pool; unique
// within the instrumented application. In case the connection pool
// implementation doesn't provide a name, instrumentation SHOULD use a
// combination of parameters that would make the name unique, for example,
// combining attributes `server.address`, `server.port`, and `db.namespace`,
// formatted as `server.address:server.port/db.namespace`. Instrumentations that
// generate connection pool name following different patterns SHOULD document it.
//
// The dbClientConnectionState is the the state of a connection in the pool
func (m ClientConnectionCount) Add(
    ctx context.Context,
    incr int64,
	clientConnectionPoolName string,
	clientConnectionState ClientConnectionStateAttr,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),
			attribute.String("db.client.connection.state", string(clientConnectionState)),

		),
	)
}

// DBClientConnectionCreateTime is an instrument used to record metric values
// conforming to the "db.client.connection.create_time" semantic conventions. It
// represents the time it took to create a new connection.
type ClientConnectionCreateTime struct {
	inst metric.Float64Histogram
}

// NewClientConnectionCreateTime returns a new ClientConnectionCreateTime
// instrument.
func NewClientConnectionCreateTime(m metric.Meter) (ClientConnectionCreateTime, error) {
	i, err := m.Float64Histogram(
	    "db.client.connection.create_time",
	    metric.WithDescription("The time it took to create a new connection"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientConnectionCreateTime{}, err
	}
	return ClientConnectionCreateTime{i}, nil
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

// Record records incr to the existing count.
//
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionIdleMax is an instrument used to record metric values
// conforming to the "db.client.connection.idle.max" semantic conventions. It
// represents the maximum number of idle open connections allowed.
type ClientConnectionIdleMax struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionIdleMax returns a new ClientConnectionIdleMax instrument.
func NewClientConnectionIdleMax(m metric.Meter) (ClientConnectionIdleMax, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connection.idle.max",
	    metric.WithDescription("The maximum number of idle open connections allowed"),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionIdleMax{}, err
	}
	return ClientConnectionIdleMax{i}, nil
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
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionIdleMin is an instrument used to record metric values
// conforming to the "db.client.connection.idle.min" semantic conventions. It
// represents the minimum number of idle open connections allowed.
type ClientConnectionIdleMin struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionIdleMin returns a new ClientConnectionIdleMin instrument.
func NewClientConnectionIdleMin(m metric.Meter) (ClientConnectionIdleMin, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connection.idle.min",
	    metric.WithDescription("The minimum number of idle open connections allowed"),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionIdleMin{}, err
	}
	return ClientConnectionIdleMin{i}, nil
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
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionMax is an instrument used to record metric values conforming
// to the "db.client.connection.max" semantic conventions. It represents the
// maximum number of open connections allowed.
type ClientConnectionMax struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionMax returns a new ClientConnectionMax instrument.
func NewClientConnectionMax(m metric.Meter) (ClientConnectionMax, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connection.max",
	    metric.WithDescription("The maximum number of open connections allowed"),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionMax{}, err
	}
	return ClientConnectionMax{i}, nil
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
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionPendingRequests is an instrument used to record metric
// values conforming to the "db.client.connection.pending_requests" semantic
// conventions. It represents the number of current pending requests for an open
// connection.
type ClientConnectionPendingRequests struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionPendingRequests returns a new
// ClientConnectionPendingRequests instrument.
func NewClientConnectionPendingRequests(m metric.Meter) (ClientConnectionPendingRequests, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connection.pending_requests",
	    metric.WithDescription("The number of current pending requests for an open connection"),
	    metric.WithUnit("{request}"),
	)
	if err != nil {
	    return ClientConnectionPendingRequests{}, err
	}
	return ClientConnectionPendingRequests{i}, nil
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
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionTimeouts is an instrument used to record metric values
// conforming to the "db.client.connection.timeouts" semantic conventions. It
// represents the number of connection timeouts that have occurred trying to
// obtain a connection from the pool.
type ClientConnectionTimeouts struct {
	inst metric.Int64Counter
}

// NewClientConnectionTimeouts returns a new ClientConnectionTimeouts instrument.
func NewClientConnectionTimeouts(m metric.Meter) (ClientConnectionTimeouts, error) {
	i, err := m.Int64Counter(
	    "db.client.connection.timeouts",
	    metric.WithDescription("The number of connection timeouts that have occurred trying to obtain a connection from the pool"),
	    metric.WithUnit("{timeout}"),
	)
	if err != nil {
	    return ClientConnectionTimeouts{}, err
	}
	return ClientConnectionTimeouts{i}, nil
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
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionUseTime is an instrument used to record metric values
// conforming to the "db.client.connection.use_time" semantic conventions. It
// represents the time between borrowing a connection and returning it to the
// pool.
type ClientConnectionUseTime struct {
	inst metric.Float64Histogram
}

// NewClientConnectionUseTime returns a new ClientConnectionUseTime instrument.
func NewClientConnectionUseTime(m metric.Meter) (ClientConnectionUseTime, error) {
	i, err := m.Float64Histogram(
	    "db.client.connection.use_time",
	    metric.WithDescription("The time between borrowing a connection and returning it to the pool"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientConnectionUseTime{}, err
	}
	return ClientConnectionUseTime{i}, nil
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

// Record records incr to the existing count.
//
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionWaitTime is an instrument used to record metric values
// conforming to the "db.client.connection.wait_time" semantic conventions. It
// represents the time it took to obtain an open connection from the pool.
type ClientConnectionWaitTime struct {
	inst metric.Float64Histogram
}

// NewClientConnectionWaitTime returns a new ClientConnectionWaitTime instrument.
func NewClientConnectionWaitTime(m metric.Meter) (ClientConnectionWaitTime, error) {
	i, err := m.Float64Histogram(
	    "db.client.connection.wait_time",
	    metric.WithDescription("The time it took to obtain an open connection from the pool"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientConnectionWaitTime{}, err
	}
	return ClientConnectionWaitTime{i}, nil
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

// Record records incr to the existing count.
//
// The dbClientConnectionPoolName is the the name of the connection pool; unique
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

) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(			attribute.String("db.client.connection.pool.name", clientConnectionPoolName),

		),
	)
}

// DBClientConnectionsCreateTime is an instrument used to record metric values
// conforming to the "db.client.connections.create_time" semantic conventions. It
// represents the deprecated, use `db.client.connection.create_time` instead.
// Note: the unit also changed from `ms` to `s`.
type ClientConnectionsCreateTime struct {
	inst metric.Int64Histogram
}

// NewClientConnectionsCreateTime returns a new ClientConnectionsCreateTime
// instrument.
func NewClientConnectionsCreateTime(m metric.Meter) (ClientConnectionsCreateTime, error) {
	i, err := m.Int64Histogram(
	    "db.client.connections.create_time",
	    metric.WithDescription("Deprecated, use `db.client.connection.create_time` instead. Note: the unit also changed from `ms` to `s`."),
	    metric.WithUnit("ms"),
	)
	if err != nil {
	    return ClientConnectionsCreateTime{}, err
	}
	return ClientConnectionsCreateTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsCreateTime) Name() string {
	return "db.client.connections.create_time"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsCreateTime) Unit() string {
	return "ms"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsCreateTime) Description() string {
	return "Deprecated, use `db.client.connection.create_time` instead. Note: the unit also changed from `ms` to `s`."
}

// Record records incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsCreateTime) Record(
    ctx context.Context,
    val int64,
	clientConnectionsPoolName string,

) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientConnectionsIdleMax is an instrument used to record metric values
// conforming to the "db.client.connections.idle.max" semantic conventions. It
// represents the deprecated, use `db.client.connection.idle.max` instead.
type ClientConnectionsIdleMax struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionsIdleMax returns a new ClientConnectionsIdleMax instrument.
func NewClientConnectionsIdleMax(m metric.Meter) (ClientConnectionsIdleMax, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connections.idle.max",
	    metric.WithDescription("Deprecated, use `db.client.connection.idle.max` instead."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionsIdleMax{}, err
	}
	return ClientConnectionsIdleMax{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsIdleMax) Name() string {
	return "db.client.connections.idle.max"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsIdleMax) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsIdleMax) Description() string {
	return "Deprecated, use `db.client.connection.idle.max` instead."
}

// Add adds incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsIdleMax) Add(
    ctx context.Context,
    incr int64,
	clientConnectionsPoolName string,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientConnectionsIdleMin is an instrument used to record metric values
// conforming to the "db.client.connections.idle.min" semantic conventions. It
// represents the deprecated, use `db.client.connection.idle.min` instead.
type ClientConnectionsIdleMin struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionsIdleMin returns a new ClientConnectionsIdleMin instrument.
func NewClientConnectionsIdleMin(m metric.Meter) (ClientConnectionsIdleMin, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connections.idle.min",
	    metric.WithDescription("Deprecated, use `db.client.connection.idle.min` instead."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionsIdleMin{}, err
	}
	return ClientConnectionsIdleMin{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsIdleMin) Name() string {
	return "db.client.connections.idle.min"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsIdleMin) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsIdleMin) Description() string {
	return "Deprecated, use `db.client.connection.idle.min` instead."
}

// Add adds incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsIdleMin) Add(
    ctx context.Context,
    incr int64,
	clientConnectionsPoolName string,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientConnectionsMax is an instrument used to record metric values
// conforming to the "db.client.connections.max" semantic conventions. It
// represents the deprecated, use `db.client.connection.max` instead.
type ClientConnectionsMax struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionsMax returns a new ClientConnectionsMax instrument.
func NewClientConnectionsMax(m metric.Meter) (ClientConnectionsMax, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connections.max",
	    metric.WithDescription("Deprecated, use `db.client.connection.max` instead."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionsMax{}, err
	}
	return ClientConnectionsMax{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsMax) Name() string {
	return "db.client.connections.max"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsMax) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsMax) Description() string {
	return "Deprecated, use `db.client.connection.max` instead."
}

// Add adds incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsMax) Add(
    ctx context.Context,
    incr int64,
	clientConnectionsPoolName string,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientConnectionsPendingRequests is an instrument used to record metric
// values conforming to the "db.client.connections.pending_requests" semantic
// conventions. It represents the deprecated, use
// `db.client.connection.pending_requests` instead.
type ClientConnectionsPendingRequests struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionsPendingRequests returns a new
// ClientConnectionsPendingRequests instrument.
func NewClientConnectionsPendingRequests(m metric.Meter) (ClientConnectionsPendingRequests, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connections.pending_requests",
	    metric.WithDescription("Deprecated, use `db.client.connection.pending_requests` instead."),
	    metric.WithUnit("{request}"),
	)
	if err != nil {
	    return ClientConnectionsPendingRequests{}, err
	}
	return ClientConnectionsPendingRequests{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsPendingRequests) Name() string {
	return "db.client.connections.pending_requests"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsPendingRequests) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsPendingRequests) Description() string {
	return "Deprecated, use `db.client.connection.pending_requests` instead."
}

// Add adds incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsPendingRequests) Add(
    ctx context.Context,
    incr int64,
	clientConnectionsPoolName string,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientConnectionsTimeouts is an instrument used to record metric values
// conforming to the "db.client.connections.timeouts" semantic conventions. It
// represents the deprecated, use `db.client.connection.timeouts` instead.
type ClientConnectionsTimeouts struct {
	inst metric.Int64Counter
}

// NewClientConnectionsTimeouts returns a new ClientConnectionsTimeouts
// instrument.
func NewClientConnectionsTimeouts(m metric.Meter) (ClientConnectionsTimeouts, error) {
	i, err := m.Int64Counter(
	    "db.client.connections.timeouts",
	    metric.WithDescription("Deprecated, use `db.client.connection.timeouts` instead."),
	    metric.WithUnit("{timeout}"),
	)
	if err != nil {
	    return ClientConnectionsTimeouts{}, err
	}
	return ClientConnectionsTimeouts{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsTimeouts) Name() string {
	return "db.client.connections.timeouts"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsTimeouts) Unit() string {
	return "{timeout}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsTimeouts) Description() string {
	return "Deprecated, use `db.client.connection.timeouts` instead."
}

// Add adds incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsTimeouts) Add(
    ctx context.Context,
    incr int64,
	clientConnectionsPoolName string,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientConnectionsUsage is an instrument used to record metric values
// conforming to the "db.client.connections.usage" semantic conventions. It
// represents the deprecated, use `db.client.connection.count` instead.
type ClientConnectionsUsage struct {
	inst metric.Int64UpDownCounter
}

// NewClientConnectionsUsage returns a new ClientConnectionsUsage instrument.
func NewClientConnectionsUsage(m metric.Meter) (ClientConnectionsUsage, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.connections.usage",
	    metric.WithDescription("Deprecated, use `db.client.connection.count` instead."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientConnectionsUsage{}, err
	}
	return ClientConnectionsUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsUsage) Name() string {
	return "db.client.connections.usage"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsUsage) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsUsage) Description() string {
	return "Deprecated, use `db.client.connection.count` instead."
}

// Add adds incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
//
// The dbClientConnectionsState is the deprecated, use
// `db.client.connection.state` instead.
func (m ClientConnectionsUsage) Add(
    ctx context.Context,
    incr int64,
	clientConnectionsPoolName string,
	clientConnectionsState ClientConnectionsStateAttr,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),
			attribute.String("db.client.connections.state", string(clientConnectionsState)),

		),
	)
}

// DBClientConnectionsUseTime is an instrument used to record metric values
// conforming to the "db.client.connections.use_time" semantic conventions. It
// represents the deprecated, use `db.client.connection.use_time` instead. Note:
// the unit also changed from `ms` to `s`.
type ClientConnectionsUseTime struct {
	inst metric.Int64Histogram
}

// NewClientConnectionsUseTime returns a new ClientConnectionsUseTime instrument.
func NewClientConnectionsUseTime(m metric.Meter) (ClientConnectionsUseTime, error) {
	i, err := m.Int64Histogram(
	    "db.client.connections.use_time",
	    metric.WithDescription("Deprecated, use `db.client.connection.use_time` instead. Note: the unit also changed from `ms` to `s`."),
	    metric.WithUnit("ms"),
	)
	if err != nil {
	    return ClientConnectionsUseTime{}, err
	}
	return ClientConnectionsUseTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsUseTime) Name() string {
	return "db.client.connections.use_time"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsUseTime) Unit() string {
	return "ms"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsUseTime) Description() string {
	return "Deprecated, use `db.client.connection.use_time` instead. Note: the unit also changed from `ms` to `s`."
}

// Record records incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsUseTime) Record(
    ctx context.Context,
    val int64,
	clientConnectionsPoolName string,

) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientConnectionsWaitTime is an instrument used to record metric values
// conforming to the "db.client.connections.wait_time" semantic conventions. It
// represents the deprecated, use `db.client.connection.wait_time` instead. Note:
// the unit also changed from `ms` to `s`.
type ClientConnectionsWaitTime struct {
	inst metric.Int64Histogram
}

// NewClientConnectionsWaitTime returns a new ClientConnectionsWaitTime
// instrument.
func NewClientConnectionsWaitTime(m metric.Meter) (ClientConnectionsWaitTime, error) {
	i, err := m.Int64Histogram(
	    "db.client.connections.wait_time",
	    metric.WithDescription("Deprecated, use `db.client.connection.wait_time` instead. Note: the unit also changed from `ms` to `s`."),
	    metric.WithUnit("ms"),
	)
	if err != nil {
	    return ClientConnectionsWaitTime{}, err
	}
	return ClientConnectionsWaitTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionsWaitTime) Name() string {
	return "db.client.connections.wait_time"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionsWaitTime) Unit() string {
	return "ms"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionsWaitTime) Description() string {
	return "Deprecated, use `db.client.connection.wait_time` instead. Note: the unit also changed from `ms` to `s`."
}

// Record records incr to the existing count.
//
// The dbClientConnectionsPoolName is the deprecated, use
// `db.client.connection.pool.name` instead.
func (m ClientConnectionsWaitTime) Record(
    ctx context.Context,
    val int64,
	clientConnectionsPoolName string,

) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(			attribute.String("db.client.connections.pool.name", clientConnectionsPoolName),

		),
	)
}

// DBClientCosmosDBActiveInstanceCount is an instrument used to record metric
// values conforming to the "db.client.cosmosdb.active_instance.count" semantic
// conventions. It represents the deprecated, use
// `azure.cosmosdb.client.active_instance.count` instead.
type ClientCosmosDBActiveInstanceCount struct {
	inst metric.Int64UpDownCounter
}

// NewClientCosmosDBActiveInstanceCount returns a new
// ClientCosmosDBActiveInstanceCount instrument.
func NewClientCosmosDBActiveInstanceCount(m metric.Meter) (ClientCosmosDBActiveInstanceCount, error) {
	i, err := m.Int64UpDownCounter(
	    "db.client.cosmosdb.active_instance.count",
	    metric.WithDescription("Deprecated, use `azure.cosmosdb.client.active_instance.count` instead."),
	    metric.WithUnit("{instance}"),
	)
	if err != nil {
	    return ClientCosmosDBActiveInstanceCount{}, err
	}
	return ClientCosmosDBActiveInstanceCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientCosmosDBActiveInstanceCount) Name() string {
	return "db.client.cosmosdb.active_instance.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientCosmosDBActiveInstanceCount) Unit() string {
	return "{instance}"
}

// Description returns the semantic convention description of the instrument
func (ClientCosmosDBActiveInstanceCount) Description() string {
	return "Deprecated, use `azure.cosmosdb.client.active_instance.count` instead."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ClientCosmosDBActiveInstanceCount) Add(
    ctx context.Context,
    incr int64,
	attrs ...ClientCosmosDBActiveInstanceCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ClientCosmosDBActiveInstanceCount) conv(in []ClientCosmosDBActiveInstanceCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientCosmosdbActiveInstanceCountAttr()
	}
	return out
}

// ClientCosmosDBActiveInstanceCountAttr is an optional attribute for the
// ClientCosmosDBActiveInstanceCount instrument.
type ClientCosmosDBActiveInstanceCountAttr interface {
    clientCosmosdbActiveInstanceCountAttr() attribute.KeyValue
}

type clientCosmosdbActiveInstanceCountAttr struct {
	kv attribute.KeyValue
}

func (a clientCosmosdbActiveInstanceCountAttr) clientCosmosdbActiveInstanceCountAttr() attribute.KeyValue {
    return a.kv
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientCosmosDBActiveInstanceCount) ServerPort(val int) ClientCosmosDBActiveInstanceCountAttr {
	return clientCosmosdbActiveInstanceCountAttr{kv: attribute.Int("server.port", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the database host.
func (ClientCosmosDBActiveInstanceCount) ServerAddress(val string) ClientCosmosDBActiveInstanceCountAttr {
	return clientCosmosdbActiveInstanceCountAttr{kv: attribute.String("server.address", val)}
}

// DBClientCosmosDBOperationRequestCharge is an instrument used to record metric
// values conforming to the "db.client.cosmosdb.operation.request_charge"
// semantic conventions. It represents the deprecated, use
// `azure.cosmosdb.client.operation.request_charge` instead.
type ClientCosmosDBOperationRequestCharge struct {
	inst metric.Int64Histogram
}

// NewClientCosmosDBOperationRequestCharge returns a new
// ClientCosmosDBOperationRequestCharge instrument.
func NewClientCosmosDBOperationRequestCharge(m metric.Meter) (ClientCosmosDBOperationRequestCharge, error) {
	i, err := m.Int64Histogram(
	    "db.client.cosmosdb.operation.request_charge",
	    metric.WithDescription("Deprecated, use `azure.cosmosdb.client.operation.request_charge` instead."),
	    metric.WithUnit("{request_unit}"),
	)
	if err != nil {
	    return ClientCosmosDBOperationRequestCharge{}, err
	}
	return ClientCosmosDBOperationRequestCharge{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientCosmosDBOperationRequestCharge) Name() string {
	return "db.client.cosmosdb.operation.request_charge"
}

// Unit returns the semantic convention unit of the instrument
func (ClientCosmosDBOperationRequestCharge) Unit() string {
	return "{request_unit}"
}

// Description returns the semantic convention description of the instrument
func (ClientCosmosDBOperationRequestCharge) Description() string {
	return "Deprecated, use `azure.cosmosdb.client.operation.request_charge` instead."
}

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ClientCosmosDBOperationRequestCharge) Record(
    ctx context.Context,
    val int64,
	attrs ...ClientCosmosDBOperationRequestChargeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ClientCosmosDBOperationRequestCharge) conv(in []ClientCosmosDBOperationRequestChargeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientCosmosdbOperationRequestChargeAttr()
	}
	return out
}

// ClientCosmosDBOperationRequestChargeAttr is an optional attribute for the
// ClientCosmosDBOperationRequestCharge instrument.
type ClientCosmosDBOperationRequestChargeAttr interface {
    clientCosmosdbOperationRequestChargeAttr() attribute.KeyValue
}

type clientCosmosdbOperationRequestChargeAttr struct {
	kv attribute.KeyValue
}

func (a clientCosmosdbOperationRequestChargeAttr) clientCosmosdbOperationRequestChargeAttr() attribute.KeyValue {
    return a.kv
}

// CollectionName returns an optional attribute for the "db.collection.name"
// semantic convention. It represents the cosmos DB container name.
func (ClientCosmosDBOperationRequestCharge) CollectionName(val string) ClientCosmosDBOperationRequestChargeAttr {
	return clientCosmosdbOperationRequestChargeAttr{kv: attribute.String("db.collection.name", val)}
}

// CosmosDBConsistencyLevel returns an optional attribute for the
// "db.cosmosdb.consistency_level" semantic convention. It represents the
// deprecated, use `cosmosdb.consistency.level` instead.
func (ClientCosmosDBOperationRequestCharge) CosmosDBConsistencyLevel(val CosmosDBConsistencyLevelAttr) ClientCosmosDBOperationRequestChargeAttr {
	return clientCosmosdbOperationRequestChargeAttr{kv: attribute.String("db.cosmosdb.consistency_level", string(val))}
}

// CosmosDBSubStatusCode returns an optional attribute for the
// "db.cosmosdb.sub_status_code" semantic convention. It represents the
// deprecated, use `azure.cosmosdb.response.sub_status_code` instead.
func (ClientCosmosDBOperationRequestCharge) CosmosDBSubStatusCode(val int) ClientCosmosDBOperationRequestChargeAttr {
	return clientCosmosdbOperationRequestChargeAttr{kv: attribute.Int("db.cosmosdb.sub_status_code", val)}
}

// Namespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientCosmosDBOperationRequestCharge) Namespace(val string) ClientCosmosDBOperationRequestChargeAttr {
	return clientCosmosdbOperationRequestChargeAttr{kv: attribute.String("db.namespace", val)}
}

// OperationName returns an optional attribute for the "db.operation.name"
// semantic convention. It represents the name of the operation or command being
// executed.
func (ClientCosmosDBOperationRequestCharge) OperationName(val string) ClientCosmosDBOperationRequestChargeAttr {
	return clientCosmosdbOperationRequestChargeAttr{kv: attribute.String("db.operation.name", val)}
}

// CosmosDBRegionsContacted returns an optional attribute for the
// "db.cosmosdb.regions_contacted" semantic convention. It represents the
// deprecated, use `azure.cosmosdb.operation.contacted_regions` instead.
func (ClientCosmosDBOperationRequestCharge) CosmosDBRegionsContacted(val ...string) ClientCosmosDBOperationRequestChargeAttr {
	return clientCosmosdbOperationRequestChargeAttr{kv: attribute.StringSlice("db.cosmosdb.regions_contacted", val)}
}

// DBClientOperationDuration is an instrument used to record metric values
// conforming to the "db.client.operation.duration" semantic conventions. It
// represents the duration of database client operations.
type ClientOperationDuration struct {
	inst metric.Float64Histogram
}

// NewClientOperationDuration returns a new ClientOperationDuration instrument.
func NewClientOperationDuration(m metric.Meter) (ClientOperationDuration, error) {
	i, err := m.Float64Histogram(
	    "db.client.operation.duration",
	    metric.WithDescription("Duration of database client operations."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientOperationDuration{}, err
	}
	return ClientOperationDuration{i}, nil
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

// Record records incr to the existing count.
//
// The dbSystemName is the the database management system (DBMS) product as
// identified by the client instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientOperationDuration) Record(
    ctx context.Context,
    val float64,
	systemName SystemNameAttr,
	attrs ...ClientOperationDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("db.system.name", string(systemName)),
			)...,
		),
	)
}

func (m ClientOperationDuration) conv(in []ClientOperationDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientOperationDurationAttr()
	}
	return out
}

// ClientOperationDurationAttr is an optional attribute for the
// ClientOperationDuration instrument.
type ClientOperationDurationAttr interface {
    clientOperationDurationAttr() attribute.KeyValue
}

type clientOperationDurationAttr struct {
	kv attribute.KeyValue
}

func (a clientOperationDurationAttr) clientOperationDurationAttr() attribute.KeyValue {
    return a.kv
}

// CollectionName returns an optional attribute for the "db.collection.name"
// semantic convention. It represents the name of a collection (table, container)
// within the database.
func (ClientOperationDuration) CollectionName(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("db.collection.name", val)}
}

// Namespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientOperationDuration) Namespace(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("db.namespace", val)}
}

// OperationName returns an optional attribute for the "db.operation.name"
// semantic convention. It represents the name of the operation or command being
// executed.
func (ClientOperationDuration) OperationName(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("db.operation.name", val)}
}

// ResponseStatusCode returns an optional attribute for the
// "db.response.status_code" semantic convention. It represents the database
// response status code.
func (ClientOperationDuration) ResponseStatusCode(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("db.response.status_code", val)}
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientOperationDuration) ErrorType(val ErrorTypeAttr) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("error.type", string(val))}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientOperationDuration) ServerPort(val int) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.Int("server.port", val)}
}

// QuerySummary returns an optional attribute for the "db.query.summary" semantic
// convention. It represents the low cardinality representation of a database
// query text.
func (ClientOperationDuration) QuerySummary(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("db.query.summary", val)}
}

// NetworkPeerAddress returns an optional attribute for the
// "network.peer.address" semantic convention. It represents the peer address of
// the database node where the operation was performed.
func (ClientOperationDuration) NetworkPeerAddress(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("network.peer.address", val)}
}

// NetworkPeerPort returns an optional attribute for the "network.peer.port"
// semantic convention. It represents the peer port number of the network
// connection.
func (ClientOperationDuration) NetworkPeerPort(val int) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.Int("network.peer.port", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the database host.
func (ClientOperationDuration) ServerAddress(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("server.address", val)}
}

// QueryText returns an optional attribute for the "db.query.text" semantic
// convention. It represents the database query being executed.
func (ClientOperationDuration) QueryText(val string) ClientOperationDurationAttr {
	return clientOperationDurationAttr{kv: attribute.String("db.query.text", val)}
}

// DBClientResponseReturnedRows is an instrument used to record metric values
// conforming to the "db.client.response.returned_rows" semantic conventions. It
// represents the actual number of records returned by the database operation.
type ClientResponseReturnedRows struct {
	inst metric.Int64Histogram
}

// NewClientResponseReturnedRows returns a new ClientResponseReturnedRows
// instrument.
func NewClientResponseReturnedRows(m metric.Meter) (ClientResponseReturnedRows, error) {
	i, err := m.Int64Histogram(
	    "db.client.response.returned_rows",
	    metric.WithDescription("The actual number of records returned by the database operation."),
	    metric.WithUnit("{row}"),
	)
	if err != nil {
	    return ClientResponseReturnedRows{}, err
	}
	return ClientResponseReturnedRows{i}, nil
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

// Record records incr to the existing count.
//
// The dbSystemName is the the database management system (DBMS) product as
// identified by the client instrumentation.
//
// All additional attrs passed are included in the recorded value.
func (m ClientResponseReturnedRows) Record(
    ctx context.Context,
    val int64,
	systemName SystemNameAttr,
	attrs ...ClientResponseReturnedRowsAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("db.system.name", string(systemName)),
			)...,
		),
	)
}

func (m ClientResponseReturnedRows) conv(in []ClientResponseReturnedRowsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientResponseReturnedRowsAttr()
	}
	return out
}

// ClientResponseReturnedRowsAttr is an optional attribute for the
// ClientResponseReturnedRows instrument.
type ClientResponseReturnedRowsAttr interface {
    clientResponseReturnedRowsAttr() attribute.KeyValue
}

type clientResponseReturnedRowsAttr struct {
	kv attribute.KeyValue
}

func (a clientResponseReturnedRowsAttr) clientResponseReturnedRowsAttr() attribute.KeyValue {
    return a.kv
}

// CollectionName returns an optional attribute for the "db.collection.name"
// semantic convention. It represents the name of a collection (table, container)
// within the database.
func (ClientResponseReturnedRows) CollectionName(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("db.collection.name", val)}
}

// Namespace returns an optional attribute for the "db.namespace" semantic
// convention. It represents the name of the database, fully qualified within the
// server address and port.
func (ClientResponseReturnedRows) Namespace(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("db.namespace", val)}
}

// OperationName returns an optional attribute for the "db.operation.name"
// semantic convention. It represents the name of the operation or command being
// executed.
func (ClientResponseReturnedRows) OperationName(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("db.operation.name", val)}
}

// ResponseStatusCode returns an optional attribute for the
// "db.response.status_code" semantic convention. It represents the database
// response status code.
func (ClientResponseReturnedRows) ResponseStatusCode(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("db.response.status_code", val)}
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientResponseReturnedRows) ErrorType(val ErrorTypeAttr) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("error.type", string(val))}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientResponseReturnedRows) ServerPort(val int) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.Int("server.port", val)}
}

// QuerySummary returns an optional attribute for the "db.query.summary" semantic
// convention. It represents the low cardinality representation of a database
// query text.
func (ClientResponseReturnedRows) QuerySummary(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("db.query.summary", val)}
}

// NetworkPeerAddress returns an optional attribute for the
// "network.peer.address" semantic convention. It represents the peer address of
// the database node where the operation was performed.
func (ClientResponseReturnedRows) NetworkPeerAddress(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("network.peer.address", val)}
}

// NetworkPeerPort returns an optional attribute for the "network.peer.port"
// semantic convention. It represents the peer port number of the network
// connection.
func (ClientResponseReturnedRows) NetworkPeerPort(val int) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.Int("network.peer.port", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the database host.
func (ClientResponseReturnedRows) ServerAddress(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("server.address", val)}
}

// QueryText returns an optional attribute for the "db.query.text" semantic
// convention. It represents the database query being executed.
func (ClientResponseReturnedRows) QueryText(val string) ClientResponseReturnedRowsAttr {
	return clientResponseReturnedRowsAttr{kv: attribute.String("db.query.text", val)}
}