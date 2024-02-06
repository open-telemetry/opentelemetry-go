// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.24.0"

import "go.opentelemetry.io/otel/attribute"

const (

	// DBClientConnectionsUsage is the metric conforming to the
	// "db.client.connections.usage" semantic conventions. It represents the number
	// of connections that are currently in state described by the `state`
	// attribute
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsUsageName        = attribute.Key("db.client.connections.usage")
	DBClientConnectionsUsageUnit        = attribute.Key("{connection}")
	DBClientConnectionsUsageDescription = attribute.Key("The number of connections that are currently in state described by the `state` attribute")

	// DBClientConnectionsIdleMax is the metric conforming to the
	// "db.client.connections.idle.max" semantic conventions. It represents the
	// maximum number of idle open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsIdleMaxName        = attribute.Key("db.client.connections.idle.max")
	DBClientConnectionsIdleMaxUnit        = attribute.Key("{connection}")
	DBClientConnectionsIdleMaxDescription = attribute.Key("The maximum number of idle open connections allowed")

	// DBClientConnectionsIdleMin is the metric conforming to the
	// "db.client.connections.idle.min" semantic conventions. It represents the
	// minimum number of idle open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsIdleMinName        = attribute.Key("db.client.connections.idle.min")
	DBClientConnectionsIdleMinUnit        = attribute.Key("{connection}")
	DBClientConnectionsIdleMinDescription = attribute.Key("The minimum number of idle open connections allowed")

	// DBClientConnectionsMax is the metric conforming to the
	// "db.client.connections.max" semantic conventions. It represents the maximum
	// number of open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsMaxName        = attribute.Key("db.client.connections.max")
	DBClientConnectionsMaxUnit        = attribute.Key("{connection}")
	DBClientConnectionsMaxDescription = attribute.Key("The maximum number of open connections allowed")

	// DBClientConnectionsPendingRequests is the metric conforming to the
	// "db.client.connections.pending_requests" semantic conventions. It represents
	// the number of pending requests for an open connection, cumulative for the
	// entire pool
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	DBClientConnectionsPendingRequestsName        = attribute.Key("db.client.connections.pending_requests")
	DBClientConnectionsPendingRequestsUnit        = attribute.Key("{request}")
	DBClientConnectionsPendingRequestsDescription = attribute.Key("The number of pending requests for an open connection, cumulative for the entire pool")

	// DBClientConnectionsTimeouts is the metric conforming to the
	// "db.client.connections.timeouts" semantic conventions. It represents the
	// number of connection timeouts that have occurred trying to obtain a
	// connection from the pool
	// Instrument: counter
	// Unit: {timeout}
	// Stability: None
	DBClientConnectionsTimeoutsName        = attribute.Key("db.client.connections.timeouts")
	DBClientConnectionsTimeoutsUnit        = attribute.Key("{timeout}")
	DBClientConnectionsTimeoutsDescription = attribute.Key("The number of connection timeouts that have occurred trying to obtain a connection from the pool")

	// DBClientConnectionsCreateTime is the metric conforming to the
	// "db.client.connections.create_time" semantic conventions. It represents the
	// time it took to create a new connection
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	DBClientConnectionsCreateTimeName        = attribute.Key("db.client.connections.create_time")
	DBClientConnectionsCreateTimeUnit        = attribute.Key("ms")
	DBClientConnectionsCreateTimeDescription = attribute.Key("The time it took to create a new connection")

	// DBClientConnectionsWaitTime is the metric conforming to the
	// "db.client.connections.wait_time" semantic conventions. It represents the
	// time it took to obtain an open connection from the pool
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	DBClientConnectionsWaitTimeName        = attribute.Key("db.client.connections.wait_time")
	DBClientConnectionsWaitTimeUnit        = attribute.Key("ms")
	DBClientConnectionsWaitTimeDescription = attribute.Key("The time it took to obtain an open connection from the pool")

	// DBClientConnectionsUseTime is the metric conforming to the
	// "db.client.connections.use_time" semantic conventions. It represents the
	// time between borrowing a connection and returning it to the pool
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	DBClientConnectionsUseTimeName        = attribute.Key("db.client.connections.use_time")
	DBClientConnectionsUseTimeUnit        = attribute.Key("ms")
	DBClientConnectionsUseTimeDescription = attribute.Key("The time between borrowing a connection and returning it to the pool")

	// AspnetcoreRoutingMatchAttempts is the metric conforming to the
	// "aspnetcore.routing.match_attempts" semantic conventions. It represents the
	// number of requests that were attempted to be matched to an endpoint.
	// Instrument: counter
	// Unit: {match_attempt}
	// Stability: None
	AspnetcoreRoutingMatchAttemptsName        = attribute.Key("aspnetcore.routing.match_attempts")
	AspnetcoreRoutingMatchAttemptsUnit        = attribute.Key("{match_attempt}")
	AspnetcoreRoutingMatchAttemptsDescription = attribute.Key("Number of requests that were attempted to be matched to an endpoint.")

	// AspnetcoreDiagnosticsExceptions is the metric conforming to the
	// "aspnetcore.diagnostics.exceptions" semantic conventions. It represents the
	// number of exceptions caught by exception handling middleware.
	// Instrument: counter
	// Unit: {exception}
	// Stability: None
	AspnetcoreDiagnosticsExceptionsName        = attribute.Key("aspnetcore.diagnostics.exceptions")
	AspnetcoreDiagnosticsExceptionsUnit        = attribute.Key("{exception}")
	AspnetcoreDiagnosticsExceptionsDescription = attribute.Key("Number of exceptions caught by exception handling middleware.")

	// AspnetcoreRateLimitingActiveRequestLeases is the metric conforming to the
	// "aspnetcore.rate_limiting.active_request_leases" semantic conventions. It
	// represents the number of requests that are currently active on the server
	// that hold a rate limiting lease.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	AspnetcoreRateLimitingActiveRequestLeasesName        = attribute.Key("aspnetcore.rate_limiting.active_request_leases")
	AspnetcoreRateLimitingActiveRequestLeasesUnit        = attribute.Key("{request}")
	AspnetcoreRateLimitingActiveRequestLeasesDescription = attribute.Key("Number of requests that are currently active on the server that hold a rate limiting lease.")

	// AspnetcoreRateLimitingRequestLeaseDuration is the metric conforming to the
	// "aspnetcore.rate_limiting.request_lease.duration" semantic conventions. It
	// represents the duration of rate limiting lease held by requests on the
	// server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	AspnetcoreRateLimitingRequestLeaseDurationName        = attribute.Key("aspnetcore.rate_limiting.request_lease.duration")
	AspnetcoreRateLimitingRequestLeaseDurationUnit        = attribute.Key("s")
	AspnetcoreRateLimitingRequestLeaseDurationDescription = attribute.Key("The duration of rate limiting lease held by requests on the server.")

	// AspnetcoreRateLimitingRequestTimeInQueue is the metric conforming to the
	// "aspnetcore.rate_limiting.request.time_in_queue" semantic conventions. It
	// represents the time the request spent in a queue waiting to acquire a rate
	// limiting lease.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	AspnetcoreRateLimitingRequestTimeInQueueName        = attribute.Key("aspnetcore.rate_limiting.request.time_in_queue")
	AspnetcoreRateLimitingRequestTimeInQueueUnit        = attribute.Key("s")
	AspnetcoreRateLimitingRequestTimeInQueueDescription = attribute.Key("The time the request spent in a queue waiting to acquire a rate limiting lease.")

	// AspnetcoreRateLimitingQueuedRequests is the metric conforming to the
	// "aspnetcore.rate_limiting.queued_requests" semantic conventions. It
	// represents the number of requests that are currently queued, waiting to
	// acquire a rate limiting lease.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	AspnetcoreRateLimitingQueuedRequestsName        = attribute.Key("aspnetcore.rate_limiting.queued_requests")
	AspnetcoreRateLimitingQueuedRequestsUnit        = attribute.Key("{request}")
	AspnetcoreRateLimitingQueuedRequestsDescription = attribute.Key("Number of requests that are currently queued, waiting to acquire a rate limiting lease.")

	// AspnetcoreRateLimitingRequests is the metric conforming to the
	// "aspnetcore.rate_limiting.requests" semantic conventions. It represents the
	// number of requests that tried to acquire a rate limiting lease.
	// Instrument: counter
	// Unit: {request}
	// Stability: None
	AspnetcoreRateLimitingRequestsName        = attribute.Key("aspnetcore.rate_limiting.requests")
	AspnetcoreRateLimitingRequestsUnit        = attribute.Key("{request}")
	AspnetcoreRateLimitingRequestsDescription = attribute.Key("Number of requests that tried to acquire a rate limiting lease.")

	// DNSLookupDuration is the metric conforming to the "dns.lookup.duration"
	// semantic conventions. It represents the measures the time taken to perform a
	// DNS lookup.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	DNSLookupDurationName        = attribute.Key("dns.lookup.duration")
	DNSLookupDurationUnit        = attribute.Key("s")
	DNSLookupDurationDescription = attribute.Key("Measures the time taken to perform a DNS lookup.")

	// HTTPClientOpenConnections is the metric conforming to the
	// "http.client.open_connections" semantic conventions. It represents the
	// number of outbound HTTP connections that are currently active or idle on the
	// client.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	HTTPClientOpenConnectionsName        = attribute.Key("http.client.open_connections")
	HTTPClientOpenConnectionsUnit        = attribute.Key("{connection}")
	HTTPClientOpenConnectionsDescription = attribute.Key("Number of outbound HTTP connections that are currently active or idle on the client.")

	// HTTPClientConnectionDuration is the metric conforming to the
	// "http.client.connection.duration" semantic conventions. It represents the
	// duration of the successfully established outbound HTTP connections.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	HTTPClientConnectionDurationName        = attribute.Key("http.client.connection.duration")
	HTTPClientConnectionDurationUnit        = attribute.Key("s")
	HTTPClientConnectionDurationDescription = attribute.Key("The duration of the successfully established outbound HTTP connections.")

	// HTTPClientActiveRequests is the metric conforming to the
	// "http.client.active_requests" semantic conventions. It represents the number
	// of active HTTP requests.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	HTTPClientActiveRequestsName        = attribute.Key("http.client.active_requests")
	HTTPClientActiveRequestsUnit        = attribute.Key("{request}")
	HTTPClientActiveRequestsDescription = attribute.Key("Number of active HTTP requests.")

	// HTTPClientRequestTimeInQueue is the metric conforming to the
	// "http.client.request.time_in_queue" semantic conventions. It represents the
	// amount of time requests spent on a queue waiting for an available
	// connection.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	HTTPClientRequestTimeInQueueName        = attribute.Key("http.client.request.time_in_queue")
	HTTPClientRequestTimeInQueueUnit        = attribute.Key("s")
	HTTPClientRequestTimeInQueueDescription = attribute.Key("The amount of time requests spent on a queue waiting for an available connection.")

	// KestrelActiveConnections is the metric conforming to the
	// "kestrel.active_connections" semantic conventions. It represents the number
	// of connections that are currently active on the server.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	KestrelActiveConnectionsName        = attribute.Key("kestrel.active_connections")
	KestrelActiveConnectionsUnit        = attribute.Key("{connection}")
	KestrelActiveConnectionsDescription = attribute.Key("Number of connections that are currently active on the server.")

	// KestrelConnectionDuration is the metric conforming to the
	// "kestrel.connection.duration" semantic conventions. It represents the
	// duration of connections on the server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	KestrelConnectionDurationName        = attribute.Key("kestrel.connection.duration")
	KestrelConnectionDurationUnit        = attribute.Key("s")
	KestrelConnectionDurationDescription = attribute.Key("The duration of connections on the server.")

	// KestrelRejectedConnections is the metric conforming to the
	// "kestrel.rejected_connections" semantic conventions. It represents the
	// number of connections rejected by the server.
	// Instrument: counter
	// Unit: {connection}
	// Stability: None
	KestrelRejectedConnectionsName        = attribute.Key("kestrel.rejected_connections")
	KestrelRejectedConnectionsUnit        = attribute.Key("{connection}")
	KestrelRejectedConnectionsDescription = attribute.Key("Number of connections rejected by the server.")

	// KestrelQueuedConnections is the metric conforming to the
	// "kestrel.queued_connections" semantic conventions. It represents the number
	// of connections that are currently queued and are waiting to start.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	KestrelQueuedConnectionsName        = attribute.Key("kestrel.queued_connections")
	KestrelQueuedConnectionsUnit        = attribute.Key("{connection}")
	KestrelQueuedConnectionsDescription = attribute.Key("Number of connections that are currently queued and are waiting to start.")

	// KestrelQueuedRequests is the metric conforming to the
	// "kestrel.queued_requests" semantic conventions. It represents the number of
	// HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are
	// currently queued and are waiting to start.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	KestrelQueuedRequestsName        = attribute.Key("kestrel.queued_requests")
	KestrelQueuedRequestsUnit        = attribute.Key("{request}")
	KestrelQueuedRequestsDescription = attribute.Key("Number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are currently queued and are waiting to start.")

	// KestrelUpgradedConnections is the metric conforming to the
	// "kestrel.upgraded_connections" semantic conventions. It represents the
	// number of connections that are currently upgraded (WebSockets). .
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	KestrelUpgradedConnectionsName        = attribute.Key("kestrel.upgraded_connections")
	KestrelUpgradedConnectionsUnit        = attribute.Key("{connection}")
	KestrelUpgradedConnectionsDescription = attribute.Key("Number of connections that are currently upgraded (WebSockets). .")

	// KestrelTLSHandshakeDuration is the metric conforming to the
	// "kestrel.tls_handshake.duration" semantic conventions. It represents the
	// duration of TLS handshakes on the server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	KestrelTLSHandshakeDurationName        = attribute.Key("kestrel.tls_handshake.duration")
	KestrelTLSHandshakeDurationUnit        = attribute.Key("s")
	KestrelTLSHandshakeDurationDescription = attribute.Key("The duration of TLS handshakes on the server.")

	// KestrelActiveTLSHandshakes is the metric conforming to the
	// "kestrel.active_tls_handshakes" semantic conventions. It represents the
	// number of TLS handshakes that are currently in progress on the server.
	// Instrument: updowncounter
	// Unit: {handshake}
	// Stability: None
	KestrelActiveTLSHandshakesName        = attribute.Key("kestrel.active_tls_handshakes")
	KestrelActiveTLSHandshakesUnit        = attribute.Key("{handshake}")
	KestrelActiveTLSHandshakesDescription = attribute.Key("Number of TLS handshakes that are currently in progress on the server.")

	// SignalrServerConnectionDuration is the metric conforming to the
	// "signalr.server.connection.duration" semantic conventions. It represents the
	// duration of connections on the server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	SignalrServerConnectionDurationName        = attribute.Key("signalr.server.connection.duration")
	SignalrServerConnectionDurationUnit        = attribute.Key("s")
	SignalrServerConnectionDurationDescription = attribute.Key("The duration of connections on the server.")

	// SignalrServerActiveConnections is the metric conforming to the
	// "signalr.server.active_connections" semantic conventions. It represents the
	// number of connections that are currently active on the server.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	SignalrServerActiveConnectionsName        = attribute.Key("signalr.server.active_connections")
	SignalrServerActiveConnectionsUnit        = attribute.Key("{connection}")
	SignalrServerActiveConnectionsDescription = attribute.Key("Number of connections that are currently active on the server.")

	// FaaSInvokeDuration is the metric conforming to the "faas.invoke_duration"
	// semantic conventions. It represents the measures the duration of the
	// function's logic execution
	// Instrument: histogram
	// Unit: s
	// Stability: None
	FaaSInvokeDurationName        = attribute.Key("faas.invoke_duration")
	FaaSInvokeDurationUnit        = attribute.Key("s")
	FaaSInvokeDurationDescription = attribute.Key("Measures the duration of the function's logic execution")

	// FaaSInitDuration is the metric conforming to the "faas.init_duration"
	// semantic conventions. It represents the measures the duration of the
	// function's initialization, such as a cold start
	// Instrument: histogram
	// Unit: s
	// Stability: None
	FaaSInitDurationName        = attribute.Key("faas.init_duration")
	FaaSInitDurationUnit        = attribute.Key("s")
	FaaSInitDurationDescription = attribute.Key("Measures the duration of the function's initialization, such as a cold start")

	// FaaSColdstarts is the metric conforming to the "faas.coldstarts" semantic
	// conventions. It represents the number of invocation cold starts
	// Instrument: counter
	// Unit: {coldstart}
	// Stability: None
	FaaSColdstartsName        = attribute.Key("faas.coldstarts")
	FaaSColdstartsUnit        = attribute.Key("{coldstart}")
	FaaSColdstartsDescription = attribute.Key("Number of invocation cold starts")

	// FaaSErrors is the metric conforming to the "faas.errors" semantic
	// conventions. It represents the number of invocation errors
	// Instrument: counter
	// Unit: {error}
	// Stability: None
	FaaSErrorsName        = attribute.Key("faas.errors")
	FaaSErrorsUnit        = attribute.Key("{error}")
	FaaSErrorsDescription = attribute.Key("Number of invocation errors")

	// FaaSInvocations is the metric conforming to the "faas.invocations" semantic
	// conventions. It represents the number of successful invocations
	// Instrument: counter
	// Unit: {invocation}
	// Stability: None
	FaaSInvocationsName        = attribute.Key("faas.invocations")
	FaaSInvocationsUnit        = attribute.Key("{invocation}")
	FaaSInvocationsDescription = attribute.Key("Number of successful invocations")

	// FaaSTimeouts is the metric conforming to the "faas.timeouts" semantic
	// conventions. It represents the number of invocation timeouts
	// Instrument: counter
	// Unit: {timeout}
	// Stability: None
	FaaSTimeoutsName        = attribute.Key("faas.timeouts")
	FaaSTimeoutsUnit        = attribute.Key("{timeout}")
	FaaSTimeoutsDescription = attribute.Key("Number of invocation timeouts")

	// FaaSMemUsage is the metric conforming to the "faas.mem_usage" semantic
	// conventions. It represents the distribution of max memory usage per
	// invocation
	// Instrument: histogram
	// Unit: By
	// Stability: None
	FaaSMemUsageName        = attribute.Key("faas.mem_usage")
	FaaSMemUsageUnit        = attribute.Key("By")
	FaaSMemUsageDescription = attribute.Key("Distribution of max memory usage per invocation")

	// FaaSCPUUsage is the metric conforming to the "faas.cpu_usage" semantic
	// conventions. It represents the distribution of CPU usage per invocation
	// Instrument: histogram
	// Unit: s
	// Stability: None
	FaaSCPUUsageName        = attribute.Key("faas.cpu_usage")
	FaaSCPUUsageUnit        = attribute.Key("s")
	FaaSCPUUsageDescription = attribute.Key("Distribution of CPU usage per invocation")

	// FaaSNetIo is the metric conforming to the "faas.net_io" semantic
	// conventions. It represents the distribution of net I/O usage per invocation
	// Instrument: histogram
	// Unit: By
	// Stability: None
	FaaSNetIoName        = attribute.Key("faas.net_io")
	FaaSNetIoUnit        = attribute.Key("By")
	FaaSNetIoDescription = attribute.Key("Distribution of net I/O usage per invocation")

	// HTTPServerRequestDuration is the metric conforming to the
	// "http.server.request.duration" semantic conventions. It represents the
	// duration of HTTP server requests.
	// Instrument: histogram
	// Unit: s
	// Stability: StabilityLevel.STABLE
	HTTPServerRequestDurationName        = attribute.Key("http.server.request.duration")
	HTTPServerRequestDurationUnit        = attribute.Key("s")
	HTTPServerRequestDurationDescription = attribute.Key("Duration of HTTP server requests.")

	// HTTPServerActiveRequests is the metric conforming to the
	// "http.server.active_requests" semantic conventions. It represents the number
	// of active HTTP server requests.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	HTTPServerActiveRequestsName        = attribute.Key("http.server.active_requests")
	HTTPServerActiveRequestsUnit        = attribute.Key("{request}")
	HTTPServerActiveRequestsDescription = attribute.Key("Number of active HTTP server requests.")

	// HTTPServerRequestBodySize is the metric conforming to the
	// "http.server.request.body.size" semantic conventions. It represents the size
	// of HTTP server request bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPServerRequestBodySizeName        = attribute.Key("http.server.request.body.size")
	HTTPServerRequestBodySizeUnit        = attribute.Key("By")
	HTTPServerRequestBodySizeDescription = attribute.Key("Size of HTTP server request bodies.")

	// HTTPServerResponseBodySize is the metric conforming to the
	// "http.server.response.body.size" semantic conventions. It represents the
	// size of HTTP server response bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPServerResponseBodySizeName        = attribute.Key("http.server.response.body.size")
	HTTPServerResponseBodySizeUnit        = attribute.Key("By")
	HTTPServerResponseBodySizeDescription = attribute.Key("Size of HTTP server response bodies.")

	// HTTPClientRequestDuration is the metric conforming to the
	// "http.client.request.duration" semantic conventions. It represents the
	// duration of HTTP client requests.
	// Instrument: histogram
	// Unit: s
	// Stability: StabilityLevel.STABLE
	HTTPClientRequestDurationName        = attribute.Key("http.client.request.duration")
	HTTPClientRequestDurationUnit        = attribute.Key("s")
	HTTPClientRequestDurationDescription = attribute.Key("Duration of HTTP client requests.")

	// HTTPClientRequestBodySize is the metric conforming to the
	// "http.client.request.body.size" semantic conventions. It represents the size
	// of HTTP client request bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPClientRequestBodySizeName        = attribute.Key("http.client.request.body.size")
	HTTPClientRequestBodySizeUnit        = attribute.Key("By")
	HTTPClientRequestBodySizeDescription = attribute.Key("Size of HTTP client request bodies.")

	// HTTPClientResponseBodySize is the metric conforming to the
	// "http.client.response.body.size" semantic conventions. It represents the
	// size of HTTP client response bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPClientResponseBodySizeName        = attribute.Key("http.client.response.body.size")
	HTTPClientResponseBodySizeUnit        = attribute.Key("By")
	HTTPClientResponseBodySizeDescription = attribute.Key("Size of HTTP client response bodies.")

	// JvmMemoryInit is the metric conforming to the "jvm.memory.init" semantic
	// conventions. It represents the measure of initial memory requested.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	JvmMemoryInitName        = attribute.Key("jvm.memory.init")
	JvmMemoryInitUnit        = attribute.Key("By")
	JvmMemoryInitDescription = attribute.Key("Measure of initial memory requested.")

	// JvmSystemCPUUtilization is the metric conforming to the
	// "jvm.system.cpu.utilization" semantic conventions. It represents the recent
	// CPU utilization for the whole system as reported by the JVM.
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	JvmSystemCPUUtilizationName        = attribute.Key("jvm.system.cpu.utilization")
	JvmSystemCPUUtilizationUnit        = attribute.Key("1")
	JvmSystemCPUUtilizationDescription = attribute.Key("Recent CPU utilization for the whole system as reported by the JVM.")

	// JvmSystemCPULoad1m is the metric conforming to the "jvm.system.cpu.load_1m"
	// semantic conventions. It represents the average CPU load of the whole system
	// for the last minute as reported by the JVM.
	// Instrument: gauge
	// Unit: {run_queue_item}
	// Stability: None
	JvmSystemCPULoad1mName        = attribute.Key("jvm.system.cpu.load_1m")
	JvmSystemCPULoad1mUnit        = attribute.Key("{run_queue_item}")
	JvmSystemCPULoad1mDescription = attribute.Key("Average CPU load of the whole system for the last minute as reported by the JVM.")

	// JvmBufferMemoryUsage is the metric conforming to the
	// "jvm.buffer.memory.usage" semantic conventions. It represents the measure of
	// memory used by buffers.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	JvmBufferMemoryUsageName        = attribute.Key("jvm.buffer.memory.usage")
	JvmBufferMemoryUsageUnit        = attribute.Key("By")
	JvmBufferMemoryUsageDescription = attribute.Key("Measure of memory used by buffers.")

	// JvmBufferMemoryLimit is the metric conforming to the
	// "jvm.buffer.memory.limit" semantic conventions. It represents the measure of
	// total memory capacity of buffers.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	JvmBufferMemoryLimitName        = attribute.Key("jvm.buffer.memory.limit")
	JvmBufferMemoryLimitUnit        = attribute.Key("By")
	JvmBufferMemoryLimitDescription = attribute.Key("Measure of total memory capacity of buffers.")

	// JvmBufferCount is the metric conforming to the "jvm.buffer.count" semantic
	// conventions. It represents the number of buffers in the pool.
	// Instrument: updowncounter
	// Unit: {buffer}
	// Stability: None
	JvmBufferCountName        = attribute.Key("jvm.buffer.count")
	JvmBufferCountUnit        = attribute.Key("{buffer}")
	JvmBufferCountDescription = attribute.Key("Number of buffers in the pool.")

	// JvmMemoryUsed is the metric conforming to the "jvm.memory.used" semantic
	// conventions. It represents the measure of memory used.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryUsedName        = attribute.Key("jvm.memory.used")
	JvmMemoryUsedUnit        = attribute.Key("By")
	JvmMemoryUsedDescription = attribute.Key("Measure of memory used.")

	// JvmMemoryCommitted is the metric conforming to the "jvm.memory.committed"
	// semantic conventions. It represents the measure of memory committed.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryCommittedName        = attribute.Key("jvm.memory.committed")
	JvmMemoryCommittedUnit        = attribute.Key("By")
	JvmMemoryCommittedDescription = attribute.Key("Measure of memory committed.")

	// JvmMemoryLimit is the metric conforming to the "jvm.memory.limit" semantic
	// conventions. It represents the measure of max obtainable memory.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryLimitName        = attribute.Key("jvm.memory.limit")
	JvmMemoryLimitUnit        = attribute.Key("By")
	JvmMemoryLimitDescription = attribute.Key("Measure of max obtainable memory.")

	// JvmMemoryUsedAfterLastGc is the metric conforming to the
	// "jvm.memory.used_after_last_gc" semantic conventions. It represents the
	// measure of memory used, as measured after the most recent garbage collection
	// event on this pool.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryUsedAfterLastGcName        = attribute.Key("jvm.memory.used_after_last_gc")
	JvmMemoryUsedAfterLastGcUnit        = attribute.Key("By")
	JvmMemoryUsedAfterLastGcDescription = attribute.Key("Measure of memory used, as measured after the most recent garbage collection event on this pool.")

	// JvmGcDuration is the metric conforming to the "jvm.gc.duration" semantic
	// conventions. It represents the duration of JVM garbage collection actions.
	// Instrument: histogram
	// Unit: s
	// Stability: StabilityLevel.STABLE
	JvmGcDurationName        = attribute.Key("jvm.gc.duration")
	JvmGcDurationUnit        = attribute.Key("s")
	JvmGcDurationDescription = attribute.Key("Duration of JVM garbage collection actions.")

	// JvmThreadCount is the metric conforming to the "jvm.thread.count" semantic
	// conventions. It represents the number of executing platform threads.
	// Instrument: updowncounter
	// Unit: {thread}
	// Stability: StabilityLevel.STABLE
	JvmThreadCountName        = attribute.Key("jvm.thread.count")
	JvmThreadCountUnit        = attribute.Key("{thread}")
	JvmThreadCountDescription = attribute.Key("Number of executing platform threads.")

	// JvmClassLoaded is the metric conforming to the "jvm.class.loaded" semantic
	// conventions. It represents the number of classes loaded since JVM start.
	// Instrument: counter
	// Unit: {class}
	// Stability: StabilityLevel.STABLE
	JvmClassLoadedName        = attribute.Key("jvm.class.loaded")
	JvmClassLoadedUnit        = attribute.Key("{class}")
	JvmClassLoadedDescription = attribute.Key("Number of classes loaded since JVM start.")

	// JvmClassUnloaded is the metric conforming to the "jvm.class.unloaded"
	// semantic conventions. It represents the number of classes unloaded since JVM
	// start.
	// Instrument: counter
	// Unit: {class}
	// Stability: StabilityLevel.STABLE
	JvmClassUnloadedName        = attribute.Key("jvm.class.unloaded")
	JvmClassUnloadedUnit        = attribute.Key("{class}")
	JvmClassUnloadedDescription = attribute.Key("Number of classes unloaded since JVM start.")

	// JvmClassCount is the metric conforming to the "jvm.class.count" semantic
	// conventions. It represents the number of classes currently loaded.
	// Instrument: updowncounter
	// Unit: {class}
	// Stability: StabilityLevel.STABLE
	JvmClassCountName        = attribute.Key("jvm.class.count")
	JvmClassCountUnit        = attribute.Key("{class}")
	JvmClassCountDescription = attribute.Key("Number of classes currently loaded.")

	// JvmCPUCount is the metric conforming to the "jvm.cpu.count" semantic
	// conventions. It represents the number of processors available to the Java
	// virtual machine.
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: StabilityLevel.STABLE
	JvmCPUCountName        = attribute.Key("jvm.cpu.count")
	JvmCPUCountUnit        = attribute.Key("{cpu}")
	JvmCPUCountDescription = attribute.Key("Number of processors available to the Java virtual machine.")

	// JvmCPUTime is the metric conforming to the "jvm.cpu.time" semantic
	// conventions. It represents the cPU time used by the process as reported by
	// the JVM.
	// Instrument: counter
	// Unit: s
	// Stability: StabilityLevel.STABLE
	JvmCPUTimeName        = attribute.Key("jvm.cpu.time")
	JvmCPUTimeUnit        = attribute.Key("s")
	JvmCPUTimeDescription = attribute.Key("CPU time used by the process as reported by the JVM.")

	// JvmCPURecentUtilization is the metric conforming to the
	// "jvm.cpu.recent_utilization" semantic conventions. It represents the recent
	// CPU utilization for the process as reported by the JVM.
	// Instrument: gauge
	// Unit: 1
	// Stability: StabilityLevel.STABLE
	JvmCPURecentUtilizationName        = attribute.Key("jvm.cpu.recent_utilization")
	JvmCPURecentUtilizationUnit        = attribute.Key("1")
	JvmCPURecentUtilizationDescription = attribute.Key("Recent CPU utilization for the process as reported by the JVM.")

	// MessagingPublishDuration is the metric conforming to the
	// "messaging.publish.duration" semantic conventions. It represents the
	// measures the duration of publish operation.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	MessagingPublishDurationName        = attribute.Key("messaging.publish.duration")
	MessagingPublishDurationUnit        = attribute.Key("s")
	MessagingPublishDurationDescription = attribute.Key("Measures the duration of publish operation.")

	// MessagingReceiveDuration is the metric conforming to the
	// "messaging.receive.duration" semantic conventions. It represents the
	// measures the duration of receive operation.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	MessagingReceiveDurationName        = attribute.Key("messaging.receive.duration")
	MessagingReceiveDurationUnit        = attribute.Key("s")
	MessagingReceiveDurationDescription = attribute.Key("Measures the duration of receive operation.")

	// MessagingDeliverDuration is the metric conforming to the
	// "messaging.deliver.duration" semantic conventions. It represents the
	// measures the duration of deliver operation.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	MessagingDeliverDurationName        = attribute.Key("messaging.deliver.duration")
	MessagingDeliverDurationUnit        = attribute.Key("s")
	MessagingDeliverDurationDescription = attribute.Key("Measures the duration of deliver operation.")

	// MessagingPublishMessages is the metric conforming to the
	// "messaging.publish.messages" semantic conventions. It represents the
	// measures the number of published messages.
	// Instrument: counter
	// Unit: {message}
	// Stability: None
	MessagingPublishMessagesName        = attribute.Key("messaging.publish.messages")
	MessagingPublishMessagesUnit        = attribute.Key("{message}")
	MessagingPublishMessagesDescription = attribute.Key("Measures the number of published messages.")

	// MessagingReceiveMessages is the metric conforming to the
	// "messaging.receive.messages" semantic conventions. It represents the
	// measures the number of received messages.
	// Instrument: counter
	// Unit: {message}
	// Stability: None
	MessagingReceiveMessagesName        = attribute.Key("messaging.receive.messages")
	MessagingReceiveMessagesUnit        = attribute.Key("{message}")
	MessagingReceiveMessagesDescription = attribute.Key("Measures the number of received messages.")

	// MessagingDeliverMessages is the metric conforming to the
	// "messaging.deliver.messages" semantic conventions. It represents the
	// measures the number of delivered messages.
	// Instrument: counter
	// Unit: {message}
	// Stability: None
	MessagingDeliverMessagesName        = attribute.Key("messaging.deliver.messages")
	MessagingDeliverMessagesUnit        = attribute.Key("{message}")
	MessagingDeliverMessagesDescription = attribute.Key("Measures the number of delivered messages.")

	// RPCServerDuration is the metric conforming to the "rpc.server.duration"
	// semantic conventions. It represents the measures the duration of inbound
	// RPC.
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	RPCServerDurationName        = attribute.Key("rpc.server.duration")
	RPCServerDurationUnit        = attribute.Key("ms")
	RPCServerDurationDescription = attribute.Key("Measures the duration of inbound RPC.")

	// RPCServerRequestSize is the metric conforming to the
	// "rpc.server.request.size" semantic conventions. It represents the measures
	// the size of RPC request messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCServerRequestSizeName        = attribute.Key("rpc.server.request.size")
	RPCServerRequestSizeUnit        = attribute.Key("By")
	RPCServerRequestSizeDescription = attribute.Key("Measures the size of RPC request messages (uncompressed).")

	// RPCServerResponseSize is the metric conforming to the
	// "rpc.server.response.size" semantic conventions. It represents the measures
	// the size of RPC response messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCServerResponseSizeName        = attribute.Key("rpc.server.response.size")
	RPCServerResponseSizeUnit        = attribute.Key("By")
	RPCServerResponseSizeDescription = attribute.Key("Measures the size of RPC response messages (uncompressed).")

	// RPCServerRequestsPerRPC is the metric conforming to the
	// "rpc.server.requests_per_rpc" semantic conventions. It represents the
	// measures the number of messages received per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCServerRequestsPerRPCName        = attribute.Key("rpc.server.requests_per_rpc")
	RPCServerRequestsPerRPCUnit        = attribute.Key("{count}")
	RPCServerRequestsPerRPCDescription = attribute.Key("Measures the number of messages received per RPC.")

	// RPCServerResponsesPerRPC is the metric conforming to the
	// "rpc.server.responses_per_rpc" semantic conventions. It represents the
	// measures the number of messages sent per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCServerResponsesPerRPCName        = attribute.Key("rpc.server.responses_per_rpc")
	RPCServerResponsesPerRPCUnit        = attribute.Key("{count}")
	RPCServerResponsesPerRPCDescription = attribute.Key("Measures the number of messages sent per RPC.")

	// RPCClientDuration is the metric conforming to the "rpc.client.duration"
	// semantic conventions. It represents the measures the duration of outbound
	// RPC.
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	RPCClientDurationName        = attribute.Key("rpc.client.duration")
	RPCClientDurationUnit        = attribute.Key("ms")
	RPCClientDurationDescription = attribute.Key("Measures the duration of outbound RPC.")

	// RPCClientRequestSize is the metric conforming to the
	// "rpc.client.request.size" semantic conventions. It represents the measures
	// the size of RPC request messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCClientRequestSizeName        = attribute.Key("rpc.client.request.size")
	RPCClientRequestSizeUnit        = attribute.Key("By")
	RPCClientRequestSizeDescription = attribute.Key("Measures the size of RPC request messages (uncompressed).")

	// RPCClientResponseSize is the metric conforming to the
	// "rpc.client.response.size" semantic conventions. It represents the measures
	// the size of RPC response messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCClientResponseSizeName        = attribute.Key("rpc.client.response.size")
	RPCClientResponseSizeUnit        = attribute.Key("By")
	RPCClientResponseSizeDescription = attribute.Key("Measures the size of RPC response messages (uncompressed).")

	// RPCClientRequestsPerRPC is the metric conforming to the
	// "rpc.client.requests_per_rpc" semantic conventions. It represents the
	// measures the number of messages received per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCClientRequestsPerRPCName        = attribute.Key("rpc.client.requests_per_rpc")
	RPCClientRequestsPerRPCUnit        = attribute.Key("{count}")
	RPCClientRequestsPerRPCDescription = attribute.Key("Measures the number of messages received per RPC.")

	// RPCClientResponsesPerRPC is the metric conforming to the
	// "rpc.client.responses_per_rpc" semantic conventions. It represents the
	// measures the number of messages sent per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCClientResponsesPerRPCName        = attribute.Key("rpc.client.responses_per_rpc")
	RPCClientResponsesPerRPCUnit        = attribute.Key("{count}")
	RPCClientResponsesPerRPCDescription = attribute.Key("Measures the number of messages sent per RPC.")

	// SystemCPUTime is the metric conforming to the "system.cpu.time" semantic
	// conventions. It represents the seconds each logical CPU spent on each mode
	// Instrument: counter
	// Unit: s
	// Stability: None
	SystemCPUTimeName        = attribute.Key("system.cpu.time")
	SystemCPUTimeUnit        = attribute.Key("s")
	SystemCPUTimeDescription = attribute.Key("Seconds each logical CPU spent on each mode")

	// SystemCPUUtilization is the metric conforming to the
	// "system.cpu.utilization" semantic conventions. It represents the difference
	// in system.cpu.time since the last measurement, divided by the elapsed time
	// and number of logical CPUs
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemCPUUtilizationName        = attribute.Key("system.cpu.utilization")
	SystemCPUUtilizationUnit        = attribute.Key("1")
	SystemCPUUtilizationDescription = attribute.Key("Difference in system.cpu.time since the last measurement, divided by the elapsed time and number of logical CPUs")

	// SystemCPUFrequency is the metric conforming to the "system.cpu.frequency"
	// semantic conventions. It represents the reports the current frequency of the
	// CPU in Hz
	// Instrument: gauge
	// Unit: {Hz}
	// Stability: None
	SystemCPUFrequencyName        = attribute.Key("system.cpu.frequency")
	SystemCPUFrequencyUnit        = attribute.Key("{Hz}")
	SystemCPUFrequencyDescription = attribute.Key("Reports the current frequency of the CPU in Hz")

	// SystemCPUPhysicalCount is the metric conforming to the
	// "system.cpu.physical.count" semantic conventions. It represents the reports
	// the number of actual physical processor cores on the hardware
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: None
	SystemCPUPhysicalCountName        = attribute.Key("system.cpu.physical.count")
	SystemCPUPhysicalCountUnit        = attribute.Key("{cpu}")
	SystemCPUPhysicalCountDescription = attribute.Key("Reports the number of actual physical processor cores on the hardware")

	// SystemCPULogicalCount is the metric conforming to the
	// "system.cpu.logical.count" semantic conventions. It represents the reports
	// the number of logical (virtual) processor cores created by the operating
	// system to manage multitasking
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: None
	SystemCPULogicalCountName        = attribute.Key("system.cpu.logical.count")
	SystemCPULogicalCountUnit        = attribute.Key("{cpu}")
	SystemCPULogicalCountDescription = attribute.Key("Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking")

	// SystemMemoryUsage is the metric conforming to the "system.memory.usage"
	// semantic conventions. It represents the reports memory in use by state.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemMemoryUsageName        = attribute.Key("system.memory.usage")
	SystemMemoryUsageUnit        = attribute.Key("By")
	SystemMemoryUsageDescription = attribute.Key("Reports memory in use by state.")

	// SystemMemoryLimit is the metric conforming to the "system.memory.limit"
	// semantic conventions. It represents the total memory available in the
	// system.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemMemoryLimitName        = attribute.Key("system.memory.limit")
	SystemMemoryLimitUnit        = attribute.Key("By")
	SystemMemoryLimitDescription = attribute.Key("Total memory available in the system.")

	// SystemMemoryUtilization is the metric conforming to the
	// "system.memory.utilization" semantic conventions. It represents the
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemMemoryUtilizationName        = attribute.Key("system.memory.utilization")
	SystemMemoryUtilizationUnit        = attribute.Key("1")
	SystemMemoryUtilizationDescription = attribute.Key("")

	// SystemPagingUsage is the metric conforming to the "system.paging.usage"
	// semantic conventions. It represents the unix swap or windows pagefile usage
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemPagingUsageName        = attribute.Key("system.paging.usage")
	SystemPagingUsageUnit        = attribute.Key("By")
	SystemPagingUsageDescription = attribute.Key("Unix swap or windows pagefile usage")

	// SystemPagingUtilization is the metric conforming to the
	// "system.paging.utilization" semantic conventions. It represents the
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemPagingUtilizationName        = attribute.Key("system.paging.utilization")
	SystemPagingUtilizationUnit        = attribute.Key("1")
	SystemPagingUtilizationDescription = attribute.Key("")

	// SystemPagingFaults is the metric conforming to the "system.paging.faults"
	// semantic conventions. It represents the
	// Instrument: counter
	// Unit: {fault}
	// Stability: None
	SystemPagingFaultsName        = attribute.Key("system.paging.faults")
	SystemPagingFaultsUnit        = attribute.Key("{fault}")
	SystemPagingFaultsDescription = attribute.Key("")

	// SystemPagingOperations is the metric conforming to the
	// "system.paging.operations" semantic conventions. It represents the
	// Instrument: counter
	// Unit: {operation}
	// Stability: None
	SystemPagingOperationsName        = attribute.Key("system.paging.operations")
	SystemPagingOperationsUnit        = attribute.Key("{operation}")
	SystemPagingOperationsDescription = attribute.Key("")

	// SystemDiskIo is the metric conforming to the "system.disk.io" semantic
	// conventions. It represents the
	// Instrument: counter
	// Unit: By
	// Stability: None
	SystemDiskIoName        = attribute.Key("system.disk.io")
	SystemDiskIoUnit        = attribute.Key("By")
	SystemDiskIoDescription = attribute.Key("")

	// SystemDiskOperations is the metric conforming to the
	// "system.disk.operations" semantic conventions. It represents the
	// Instrument: counter
	// Unit: {operation}
	// Stability: None
	SystemDiskOperationsName        = attribute.Key("system.disk.operations")
	SystemDiskOperationsUnit        = attribute.Key("{operation}")
	SystemDiskOperationsDescription = attribute.Key("")

	// SystemDiskIoTime is the metric conforming to the "system.disk.io_time"
	// semantic conventions. It represents the time disk spent activated
	// Instrument: counter
	// Unit: s
	// Stability: None
	SystemDiskIoTimeName        = attribute.Key("system.disk.io_time")
	SystemDiskIoTimeUnit        = attribute.Key("s")
	SystemDiskIoTimeDescription = attribute.Key("Time disk spent activated")

	// SystemDiskOperationTime is the metric conforming to the
	// "system.disk.operation_time" semantic conventions. It represents the sum of
	// the time each operation took to complete
	// Instrument: counter
	// Unit: s
	// Stability: None
	SystemDiskOperationTimeName        = attribute.Key("system.disk.operation_time")
	SystemDiskOperationTimeUnit        = attribute.Key("s")
	SystemDiskOperationTimeDescription = attribute.Key("Sum of the time each operation took to complete")

	// SystemDiskMerged is the metric conforming to the "system.disk.merged"
	// semantic conventions. It represents the
	// Instrument: counter
	// Unit: {operation}
	// Stability: None
	SystemDiskMergedName        = attribute.Key("system.disk.merged")
	SystemDiskMergedUnit        = attribute.Key("{operation}")
	SystemDiskMergedDescription = attribute.Key("")

	// SystemFilesystemUsage is the metric conforming to the
	// "system.filesystem.usage" semantic conventions. It represents the
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemFilesystemUsageName        = attribute.Key("system.filesystem.usage")
	SystemFilesystemUsageUnit        = attribute.Key("By")
	SystemFilesystemUsageDescription = attribute.Key("")

	// SystemFilesystemUtilization is the metric conforming to the
	// "system.filesystem.utilization" semantic conventions. It represents the
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemFilesystemUtilizationName        = attribute.Key("system.filesystem.utilization")
	SystemFilesystemUtilizationUnit        = attribute.Key("1")
	SystemFilesystemUtilizationDescription = attribute.Key("")

	// SystemNetworkDropped is the metric conforming to the
	// "system.network.dropped" semantic conventions. It represents the count of
	// packets that are dropped or discarded even though there was no error
	// Instrument: counter
	// Unit: {packet}
	// Stability: None
	SystemNetworkDroppedName        = attribute.Key("system.network.dropped")
	SystemNetworkDroppedUnit        = attribute.Key("{packet}")
	SystemNetworkDroppedDescription = attribute.Key("Count of packets that are dropped or discarded even though there was no error")

	// SystemNetworkPackets is the metric conforming to the
	// "system.network.packets" semantic conventions. It represents the
	// Instrument: counter
	// Unit: {packet}
	// Stability: None
	SystemNetworkPacketsName        = attribute.Key("system.network.packets")
	SystemNetworkPacketsUnit        = attribute.Key("{packet}")
	SystemNetworkPacketsDescription = attribute.Key("")

	// SystemNetworkErrors is the metric conforming to the "system.network.errors"
	// semantic conventions. It represents the count of network errors detected
	// Instrument: counter
	// Unit: {error}
	// Stability: None
	SystemNetworkErrorsName        = attribute.Key("system.network.errors")
	SystemNetworkErrorsUnit        = attribute.Key("{error}")
	SystemNetworkErrorsDescription = attribute.Key("Count of network errors detected")

	// SystemNetworkIo is the metric conforming to the "system.network.io" semantic
	// conventions. It represents the
	// Instrument: counter
	// Unit: By
	// Stability: None
	SystemNetworkIoName        = attribute.Key("system.network.io")
	SystemNetworkIoUnit        = attribute.Key("By")
	SystemNetworkIoDescription = attribute.Key("")

	// SystemNetworkConnections is the metric conforming to the
	// "system.network.connections" semantic conventions. It represents the
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	SystemNetworkConnectionsName        = attribute.Key("system.network.connections")
	SystemNetworkConnectionsUnit        = attribute.Key("{connection}")
	SystemNetworkConnectionsDescription = attribute.Key("")

	// SystemProcessesCount is the metric conforming to the
	// "system.processes.count" semantic conventions. It represents the total
	// number of processes in each state
	// Instrument: updowncounter
	// Unit: {process}
	// Stability: None
	SystemProcessesCountName        = attribute.Key("system.processes.count")
	SystemProcessesCountUnit        = attribute.Key("{process}")
	SystemProcessesCountDescription = attribute.Key("Total number of processes in each state")

	// SystemProcessesCreated is the metric conforming to the
	// "system.processes.created" semantic conventions. It represents the total
	// number of processes created over uptime of the host
	// Instrument: counter
	// Unit: {process}
	// Stability: None
	SystemProcessesCreatedName        = attribute.Key("system.processes.created")
	SystemProcessesCreatedUnit        = attribute.Key("{process}")
	SystemProcessesCreatedDescription = attribute.Key("Total number of processes created over uptime of the host")

	// SystemLinuxMemoryAvailable is the metric conforming to the
	// "system.linux.memory.available" semantic conventions. It represents an
	// estimate of how much memory is available for starting new applications,
	// without causing swapping
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemLinuxMemoryAvailableName        = attribute.Key("system.linux.memory.available")
	SystemLinuxMemoryAvailableUnit        = attribute.Key("By")
	SystemLinuxMemoryAvailableDescription = attribute.Key("An estimate of how much memory is available for starting new applications, without causing swapping")
)
