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

	// DBClientConnectionsUsage is the attribute Key conforming to the
	// "db.client.connections.usage" semantic conventions. It represents the number
	// of connections that are currently in state described by the `state`
	// attribute
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsUsage = attribute.Key("db.client.connections.usage")

	// DBClientConnectionsIdleMax is the attribute Key conforming to the
	// "db.client.connections.idle.max" semantic conventions. It represents the
	// maximum number of idle open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsIdleMax = attribute.Key("db.client.connections.idle.max")

	// DBClientConnectionsIdleMin is the attribute Key conforming to the
	// "db.client.connections.idle.min" semantic conventions. It represents the
	// minimum number of idle open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsIdleMin = attribute.Key("db.client.connections.idle.min")

	// DBClientConnectionsMax is the attribute Key conforming to the
	// "db.client.connections.max" semantic conventions. It represents the maximum
	// number of open connections allowed
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	DBClientConnectionsMax = attribute.Key("db.client.connections.max")

	// DBClientConnectionsPendingRequests is the attribute Key conforming to the
	// "db.client.connections.pending_requests" semantic conventions. It represents
	// the number of pending requests for an open connection, cumulative for the
	// entire pool
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	DBClientConnectionsPendingRequests = attribute.Key("db.client.connections.pending_requests")

	// DBClientConnectionsTimeouts is the attribute Key conforming to the
	// "db.client.connections.timeouts" semantic conventions. It represents the
	// number of connection timeouts that have occurred trying to obtain a
	// connection from the pool
	// Instrument: counter
	// Unit: {timeout}
	// Stability: None
	DBClientConnectionsTimeouts = attribute.Key("db.client.connections.timeouts")

	// DBClientConnectionsCreateTime is the attribute Key conforming to the
	// "db.client.connections.create_time" semantic conventions. It represents the
	// time it took to create a new connection
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	DBClientConnectionsCreateTime = attribute.Key("db.client.connections.create_time")

	// DBClientConnectionsWaitTime is the attribute Key conforming to the
	// "db.client.connections.wait_time" semantic conventions. It represents the
	// time it took to obtain an open connection from the pool
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	DBClientConnectionsWaitTime = attribute.Key("db.client.connections.wait_time")

	// DBClientConnectionsUseTime is the attribute Key conforming to the
	// "db.client.connections.use_time" semantic conventions. It represents the
	// time between borrowing a connection and returning it to the pool
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	DBClientConnectionsUseTime = attribute.Key("db.client.connections.use_time")

	// AspnetcoreRoutingMatchAttempts is the attribute Key conforming to the
	// "aspnetcore.routing.match_attempts" semantic conventions. It represents the
	// number of requests that were attempted to be matched to an endpoint.
	// Instrument: counter
	// Unit: {match_attempt}
	// Stability: None
	AspnetcoreRoutingMatchAttempts = attribute.Key("aspnetcore.routing.match_attempts")

	// AspnetcoreDiagnosticsExceptions is the attribute Key conforming to the
	// "aspnetcore.diagnostics.exceptions" semantic conventions. It represents the
	// number of exceptions caught by exception handling middleware.
	// Instrument: counter
	// Unit: {exception}
	// Stability: None
	AspnetcoreDiagnosticsExceptions = attribute.Key("aspnetcore.diagnostics.exceptions")

	// AspnetcoreRateLimitingActiveRequestLeases is the attribute Key conforming to
	// the "aspnetcore.rate_limiting.active_request_leases" semantic conventions.
	// It represents the number of requests that are currently active on the server
	// that hold a rate limiting lease.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	AspnetcoreRateLimitingActiveRequestLeases = attribute.Key("aspnetcore.rate_limiting.active_request_leases")

	// AspnetcoreRateLimitingRequestLeaseDuration is the attribute Key conforming
	// to the "aspnetcore.rate_limiting.request_lease.duration" semantic
	// conventions. It represents the duration of rate limiting lease held by
	// requests on the server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	AspnetcoreRateLimitingRequestLeaseDuration = attribute.Key("aspnetcore.rate_limiting.request_lease.duration")

	// AspnetcoreRateLimitingRequestTimeInQueue is the attribute Key conforming to
	// the "aspnetcore.rate_limiting.request.time_in_queue" semantic conventions.
	// It represents the time the request spent in a queue waiting to acquire a
	// rate limiting lease.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	AspnetcoreRateLimitingRequestTimeInQueue = attribute.Key("aspnetcore.rate_limiting.request.time_in_queue")

	// AspnetcoreRateLimitingQueuedRequests is the attribute Key conforming to the
	// "aspnetcore.rate_limiting.queued_requests" semantic conventions. It
	// represents the number of requests that are currently queued, waiting to
	// acquire a rate limiting lease.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	AspnetcoreRateLimitingQueuedRequests = attribute.Key("aspnetcore.rate_limiting.queued_requests")

	// AspnetcoreRateLimitingRequests is the attribute Key conforming to the
	// "aspnetcore.rate_limiting.requests" semantic conventions. It represents the
	// number of requests that tried to acquire a rate limiting lease.
	// Instrument: counter
	// Unit: {request}
	// Stability: None
	AspnetcoreRateLimitingRequests = attribute.Key("aspnetcore.rate_limiting.requests")

	// DNSLookupDuration is the attribute Key conforming to the
	// "dns.lookup.duration" semantic conventions. It represents the measures the
	// time taken to perform a DNS lookup.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	DNSLookupDuration = attribute.Key("dns.lookup.duration")

	// HTTPClientOpenConnections is the attribute Key conforming to the
	// "http.client.open_connections" semantic conventions. It represents the
	// number of outbound HTTP connections that are currently active or idle on the
	// client.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	HTTPClientOpenConnections = attribute.Key("http.client.open_connections")

	// HTTPClientConnectionDuration is the attribute Key conforming to the
	// "http.client.connection.duration" semantic conventions. It represents the
	// duration of the successfully established outbound HTTP connections.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	HTTPClientConnectionDuration = attribute.Key("http.client.connection.duration")

	// HTTPClientActiveRequests is the attribute Key conforming to the
	// "http.client.active_requests" semantic conventions. It represents the number
	// of active HTTP requests.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	HTTPClientActiveRequests = attribute.Key("http.client.active_requests")

	// HTTPClientRequestTimeInQueue is the attribute Key conforming to the
	// "http.client.request.time_in_queue" semantic conventions. It represents the
	// amount of time requests spent on a queue waiting for an available
	// connection.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	HTTPClientRequestTimeInQueue = attribute.Key("http.client.request.time_in_queue")

	// KestrelActiveConnections is the attribute Key conforming to the
	// "kestrel.active_connections" semantic conventions. It represents the number
	// of connections that are currently active on the server.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	KestrelActiveConnections = attribute.Key("kestrel.active_connections")

	// KestrelConnectionDuration is the attribute Key conforming to the
	// "kestrel.connection.duration" semantic conventions. It represents the
	// duration of connections on the server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	KestrelConnectionDuration = attribute.Key("kestrel.connection.duration")

	// KestrelRejectedConnections is the attribute Key conforming to the
	// "kestrel.rejected_connections" semantic conventions. It represents the
	// number of connections rejected by the server.
	// Instrument: counter
	// Unit: {connection}
	// Stability: None
	KestrelRejectedConnections = attribute.Key("kestrel.rejected_connections")

	// KestrelQueuedConnections is the attribute Key conforming to the
	// "kestrel.queued_connections" semantic conventions. It represents the number
	// of connections that are currently queued and are waiting to start.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	KestrelQueuedConnections = attribute.Key("kestrel.queued_connections")

	// KestrelQueuedRequests is the attribute Key conforming to the
	// "kestrel.queued_requests" semantic conventions. It represents the number of
	// HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are
	// currently queued and are waiting to start.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	KestrelQueuedRequests = attribute.Key("kestrel.queued_requests")

	// KestrelUpgradedConnections is the attribute Key conforming to the
	// "kestrel.upgraded_connections" semantic conventions. It represents the
	// number of connections that are currently upgraded (WebSockets). .
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	KestrelUpgradedConnections = attribute.Key("kestrel.upgraded_connections")

	// KestrelTLSHandshakeDuration is the attribute Key conforming to the
	// "kestrel.tls_handshake.duration" semantic conventions. It represents the
	// duration of TLS handshakes on the server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	KestrelTLSHandshakeDuration = attribute.Key("kestrel.tls_handshake.duration")

	// KestrelActiveTLSHandshakes is the attribute Key conforming to the
	// "kestrel.active_tls_handshakes" semantic conventions. It represents the
	// number of TLS handshakes that are currently in progress on the server.
	// Instrument: updowncounter
	// Unit: {handshake}
	// Stability: None
	KestrelActiveTLSHandshakes = attribute.Key("kestrel.active_tls_handshakes")

	// SignalrServerConnectionDuration is the attribute Key conforming to the
	// "signalr.server.connection.duration" semantic conventions. It represents the
	// duration of connections on the server.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	SignalrServerConnectionDuration = attribute.Key("signalr.server.connection.duration")

	// SignalrServerActiveConnections is the attribute Key conforming to the
	// "signalr.server.active_connections" semantic conventions. It represents the
	// number of connections that are currently active on the server.
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	SignalrServerActiveConnections = attribute.Key("signalr.server.active_connections")

	// FaaSInvokeDuration is the attribute Key conforming to the
	// "faas.invoke_duration" semantic conventions. It represents the measures the
	// duration of the function's logic execution
	// Instrument: histogram
	// Unit: s
	// Stability: None
	FaaSInvokeDuration = attribute.Key("faas.invoke_duration")

	// FaaSInitDuration is the attribute Key conforming to the "faas.init_duration"
	// semantic conventions. It represents the measures the duration of the
	// function's initialization, such as a cold start
	// Instrument: histogram
	// Unit: s
	// Stability: None
	FaaSInitDuration = attribute.Key("faas.init_duration")

	// FaaSColdstarts is the attribute Key conforming to the "faas.coldstarts"
	// semantic conventions. It represents the number of invocation cold starts
	// Instrument: counter
	// Unit: {coldstart}
	// Stability: None
	FaaSColdstarts = attribute.Key("faas.coldstarts")

	// FaaSErrors is the attribute Key conforming to the "faas.errors" semantic
	// conventions. It represents the number of invocation errors
	// Instrument: counter
	// Unit: {error}
	// Stability: None
	FaaSErrors = attribute.Key("faas.errors")

	// FaaSInvocations is the attribute Key conforming to the "faas.invocations"
	// semantic conventions. It represents the number of successful invocations
	// Instrument: counter
	// Unit: {invocation}
	// Stability: None
	FaaSInvocations = attribute.Key("faas.invocations")

	// FaaSTimeouts is the attribute Key conforming to the "faas.timeouts" semantic
	// conventions. It represents the number of invocation timeouts
	// Instrument: counter
	// Unit: {timeout}
	// Stability: None
	FaaSTimeouts = attribute.Key("faas.timeouts")

	// FaaSMemUsage is the attribute Key conforming to the "faas.mem_usage"
	// semantic conventions. It represents the distribution of max memory usage per
	// invocation
	// Instrument: histogram
	// Unit: By
	// Stability: None
	FaaSMemUsage = attribute.Key("faas.mem_usage")

	// FaaSCPUUsage is the attribute Key conforming to the "faas.cpu_usage"
	// semantic conventions. It represents the distribution of CPU usage per
	// invocation
	// Instrument: histogram
	// Unit: s
	// Stability: None
	FaaSCPUUsage = attribute.Key("faas.cpu_usage")

	// FaaSNetIo is the attribute Key conforming to the "faas.net_io" semantic
	// conventions. It represents the distribution of net I/O usage per invocation
	// Instrument: histogram
	// Unit: By
	// Stability: None
	FaaSNetIo = attribute.Key("faas.net_io")

	// HTTPServerRequestDuration is the attribute Key conforming to the
	// "http.server.request.duration" semantic conventions. It represents the
	// duration of HTTP server requests.
	// Instrument: histogram
	// Unit: s
	// Stability: StabilityLevel.STABLE
	HTTPServerRequestDuration = attribute.Key("http.server.request.duration")

	// HTTPServerActiveRequests is the attribute Key conforming to the
	// "http.server.active_requests" semantic conventions. It represents the number
	// of active HTTP server requests.
	// Instrument: updowncounter
	// Unit: {request}
	// Stability: None
	HTTPServerActiveRequests = attribute.Key("http.server.active_requests")

	// HTTPServerRequestBodySize is the attribute Key conforming to the
	// "http.server.request.body.size" semantic conventions. It represents the size
	// of HTTP server request bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPServerRequestBodySize = attribute.Key("http.server.request.body.size")

	// HTTPServerResponseBodySize is the attribute Key conforming to the
	// "http.server.response.body.size" semantic conventions. It represents the
	// size of HTTP server response bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPServerResponseBodySize = attribute.Key("http.server.response.body.size")

	// HTTPClientRequestDuration is the attribute Key conforming to the
	// "http.client.request.duration" semantic conventions. It represents the
	// duration of HTTP client requests.
	// Instrument: histogram
	// Unit: s
	// Stability: StabilityLevel.STABLE
	HTTPClientRequestDuration = attribute.Key("http.client.request.duration")

	// HTTPClientRequestBodySize is the attribute Key conforming to the
	// "http.client.request.body.size" semantic conventions. It represents the size
	// of HTTP client request bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPClientRequestBodySize = attribute.Key("http.client.request.body.size")

	// HTTPClientResponseBodySize is the attribute Key conforming to the
	// "http.client.response.body.size" semantic conventions. It represents the
	// size of HTTP client response bodies.
	// Instrument: histogram
	// Unit: By
	// Stability: None
	HTTPClientResponseBodySize = attribute.Key("http.client.response.body.size")

	// JvmMemoryInit is the attribute Key conforming to the "jvm.memory.init"
	// semantic conventions. It represents the measure of initial memory requested.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	JvmMemoryInit = attribute.Key("jvm.memory.init")

	// JvmSystemCPUUtilization is the attribute Key conforming to the
	// "jvm.system.cpu.utilization" semantic conventions. It represents the recent
	// CPU utilization for the whole system as reported by the JVM.
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	JvmSystemCPUUtilization = attribute.Key("jvm.system.cpu.utilization")

	// JvmSystemCPULoad1m is the attribute Key conforming to the
	// "jvm.system.cpu.load_1m" semantic conventions. It represents the average CPU
	// load of the whole system for the last minute as reported by the JVM.
	// Instrument: gauge
	// Unit: {run_queue_item}
	// Stability: None
	JvmSystemCPULoad1m = attribute.Key("jvm.system.cpu.load_1m")

	// JvmBufferMemoryUsage is the attribute Key conforming to the
	// "jvm.buffer.memory.usage" semantic conventions. It represents the measure of
	// memory used by buffers.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	JvmBufferMemoryUsage = attribute.Key("jvm.buffer.memory.usage")

	// JvmBufferMemoryLimit is the attribute Key conforming to the
	// "jvm.buffer.memory.limit" semantic conventions. It represents the measure of
	// total memory capacity of buffers.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	JvmBufferMemoryLimit = attribute.Key("jvm.buffer.memory.limit")

	// JvmBufferCount is the attribute Key conforming to the "jvm.buffer.count"
	// semantic conventions. It represents the number of buffers in the pool.
	// Instrument: updowncounter
	// Unit: {buffer}
	// Stability: None
	JvmBufferCount = attribute.Key("jvm.buffer.count")

	// JvmMemoryUsed is the attribute Key conforming to the "jvm.memory.used"
	// semantic conventions. It represents the measure of memory used.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryUsed = attribute.Key("jvm.memory.used")

	// JvmMemoryCommitted is the attribute Key conforming to the
	// "jvm.memory.committed" semantic conventions. It represents the measure of
	// memory committed.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryCommitted = attribute.Key("jvm.memory.committed")

	// JvmMemoryLimit is the attribute Key conforming to the "jvm.memory.limit"
	// semantic conventions. It represents the measure of max obtainable memory.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryLimit = attribute.Key("jvm.memory.limit")

	// JvmMemoryUsedAfterLastGc is the attribute Key conforming to the
	// "jvm.memory.used_after_last_gc" semantic conventions. It represents the
	// measure of memory used, as measured after the most recent garbage collection
	// event on this pool.
	// Instrument: updowncounter
	// Unit: By
	// Stability: StabilityLevel.STABLE
	JvmMemoryUsedAfterLastGc = attribute.Key("jvm.memory.used_after_last_gc")

	// JvmGcDuration is the attribute Key conforming to the "jvm.gc.duration"
	// semantic conventions. It represents the duration of JVM garbage collection
	// actions.
	// Instrument: histogram
	// Unit: s
	// Stability: StabilityLevel.STABLE
	JvmGcDuration = attribute.Key("jvm.gc.duration")

	// JvmThreadCount is the attribute Key conforming to the "jvm.thread.count"
	// semantic conventions. It represents the number of executing platform
	// threads.
	// Instrument: updowncounter
	// Unit: {thread}
	// Stability: StabilityLevel.STABLE
	JvmThreadCount = attribute.Key("jvm.thread.count")

	// JvmClassLoaded is the attribute Key conforming to the "jvm.class.loaded"
	// semantic conventions. It represents the number of classes loaded since JVM
	// start.
	// Instrument: counter
	// Unit: {class}
	// Stability: StabilityLevel.STABLE
	JvmClassLoaded = attribute.Key("jvm.class.loaded")

	// JvmClassUnloaded is the attribute Key conforming to the "jvm.class.unloaded"
	// semantic conventions. It represents the number of classes unloaded since JVM
	// start.
	// Instrument: counter
	// Unit: {class}
	// Stability: StabilityLevel.STABLE
	JvmClassUnloaded = attribute.Key("jvm.class.unloaded")

	// JvmClassCount is the attribute Key conforming to the "jvm.class.count"
	// semantic conventions. It represents the number of classes currently loaded.
	// Instrument: updowncounter
	// Unit: {class}
	// Stability: StabilityLevel.STABLE
	JvmClassCount = attribute.Key("jvm.class.count")

	// JvmCPUCount is the attribute Key conforming to the "jvm.cpu.count" semantic
	// conventions. It represents the number of processors available to the Java
	// virtual machine.
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: StabilityLevel.STABLE
	JvmCPUCount = attribute.Key("jvm.cpu.count")

	// JvmCPUTime is the attribute Key conforming to the "jvm.cpu.time" semantic
	// conventions. It represents the cPU time used by the process as reported by
	// the JVM.
	// Instrument: counter
	// Unit: s
	// Stability: StabilityLevel.STABLE
	JvmCPUTime = attribute.Key("jvm.cpu.time")

	// JvmCPURecentUtilization is the attribute Key conforming to the
	// "jvm.cpu.recent_utilization" semantic conventions. It represents the recent
	// CPU utilization for the process as reported by the JVM.
	// Instrument: gauge
	// Unit: 1
	// Stability: StabilityLevel.STABLE
	JvmCPURecentUtilization = attribute.Key("jvm.cpu.recent_utilization")

	// MessagingPublishDuration is the attribute Key conforming to the
	// "messaging.publish.duration" semantic conventions. It represents the
	// measures the duration of publish operation.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	MessagingPublishDuration = attribute.Key("messaging.publish.duration")

	// MessagingReceiveDuration is the attribute Key conforming to the
	// "messaging.receive.duration" semantic conventions. It represents the
	// measures the duration of receive operation.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	MessagingReceiveDuration = attribute.Key("messaging.receive.duration")

	// MessagingDeliverDuration is the attribute Key conforming to the
	// "messaging.deliver.duration" semantic conventions. It represents the
	// measures the duration of deliver operation.
	// Instrument: histogram
	// Unit: s
	// Stability: None
	MessagingDeliverDuration = attribute.Key("messaging.deliver.duration")

	// MessagingPublishMessages is the attribute Key conforming to the
	// "messaging.publish.messages" semantic conventions. It represents the
	// measures the number of published messages.
	// Instrument: counter
	// Unit: {message}
	// Stability: None
	MessagingPublishMessages = attribute.Key("messaging.publish.messages")

	// MessagingReceiveMessages is the attribute Key conforming to the
	// "messaging.receive.messages" semantic conventions. It represents the
	// measures the number of received messages.
	// Instrument: counter
	// Unit: {message}
	// Stability: None
	MessagingReceiveMessages = attribute.Key("messaging.receive.messages")

	// MessagingDeliverMessages is the attribute Key conforming to the
	// "messaging.deliver.messages" semantic conventions. It represents the
	// measures the number of delivered messages.
	// Instrument: counter
	// Unit: {message}
	// Stability: None
	MessagingDeliverMessages = attribute.Key("messaging.deliver.messages")

	// RPCServerDuration is the attribute Key conforming to the
	// "rpc.server.duration" semantic conventions. It represents the measures the
	// duration of inbound RPC.
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	RPCServerDuration = attribute.Key("rpc.server.duration")

	// RPCServerRequestSize is the attribute Key conforming to the
	// "rpc.server.request.size" semantic conventions. It represents the measures
	// the size of RPC request messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCServerRequestSize = attribute.Key("rpc.server.request.size")

	// RPCServerResponseSize is the attribute Key conforming to the
	// "rpc.server.response.size" semantic conventions. It represents the measures
	// the size of RPC response messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCServerResponseSize = attribute.Key("rpc.server.response.size")

	// RPCServerRequestsPerRPC is the attribute Key conforming to the
	// "rpc.server.requests_per_rpc" semantic conventions. It represents the
	// measures the number of messages received per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCServerRequestsPerRPC = attribute.Key("rpc.server.requests_per_rpc")

	// RPCServerResponsesPerRPC is the attribute Key conforming to the
	// "rpc.server.responses_per_rpc" semantic conventions. It represents the
	// measures the number of messages sent per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCServerResponsesPerRPC = attribute.Key("rpc.server.responses_per_rpc")

	// RPCClientDuration is the attribute Key conforming to the
	// "rpc.client.duration" semantic conventions. It represents the measures the
	// duration of outbound RPC.
	// Instrument: histogram
	// Unit: ms
	// Stability: None
	RPCClientDuration = attribute.Key("rpc.client.duration")

	// RPCClientRequestSize is the attribute Key conforming to the
	// "rpc.client.request.size" semantic conventions. It represents the measures
	// the size of RPC request messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCClientRequestSize = attribute.Key("rpc.client.request.size")

	// RPCClientResponseSize is the attribute Key conforming to the
	// "rpc.client.response.size" semantic conventions. It represents the measures
	// the size of RPC response messages (uncompressed).
	// Instrument: histogram
	// Unit: By
	// Stability: None
	RPCClientResponseSize = attribute.Key("rpc.client.response.size")

	// RPCClientRequestsPerRPC is the attribute Key conforming to the
	// "rpc.client.requests_per_rpc" semantic conventions. It represents the
	// measures the number of messages received per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCClientRequestsPerRPC = attribute.Key("rpc.client.requests_per_rpc")

	// RPCClientResponsesPerRPC is the attribute Key conforming to the
	// "rpc.client.responses_per_rpc" semantic conventions. It represents the
	// measures the number of messages sent per RPC.
	// Instrument: histogram
	// Unit: {count}
	// Stability: None
	RPCClientResponsesPerRPC = attribute.Key("rpc.client.responses_per_rpc")

	// SystemCPUTime is the attribute Key conforming to the "system.cpu.time"
	// semantic conventions. It represents the seconds each logical CPU spent on
	// each mode
	// Instrument: counter
	// Unit: s
	// Stability: None
	SystemCPUTime = attribute.Key("system.cpu.time")

	// SystemCPUUtilization is the attribute Key conforming to the
	// "system.cpu.utilization" semantic conventions. It represents the difference
	// in system.cpu.time since the last measurement, divided by the elapsed time
	// and number of logical CPUs
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemCPUUtilization = attribute.Key("system.cpu.utilization")

	// SystemCPUFrequency is the attribute Key conforming to the
	// "system.cpu.frequency" semantic conventions. It represents the reports the
	// current frequency of the CPU in Hz
	// Instrument: gauge
	// Unit: {Hz}
	// Stability: None
	SystemCPUFrequency = attribute.Key("system.cpu.frequency")

	// SystemCPUPhysicalCount is the attribute Key conforming to the
	// "system.cpu.physical.count" semantic conventions. It represents the reports
	// the number of actual physical processor cores on the hardware
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: None
	SystemCPUPhysicalCount = attribute.Key("system.cpu.physical.count")

	// SystemCPULogicalCount is the attribute Key conforming to the
	// "system.cpu.logical.count" semantic conventions. It represents the reports
	// the number of logical (virtual) processor cores created by the operating
	// system to manage multitasking
	// Instrument: updowncounter
	// Unit: {cpu}
	// Stability: None
	SystemCPULogicalCount = attribute.Key("system.cpu.logical.count")

	// SystemMemoryUsage is the attribute Key conforming to the
	// "system.memory.usage" semantic conventions. It represents the reports memory
	// in use by state.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemMemoryUsage = attribute.Key("system.memory.usage")

	// SystemMemoryLimit is the attribute Key conforming to the
	// "system.memory.limit" semantic conventions. It represents the total memory
	// available in the system.
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemMemoryLimit = attribute.Key("system.memory.limit")

	// SystemMemoryUtilization is the attribute Key conforming to the
	// "system.memory.utilization" semantic conventions. It represents the
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemMemoryUtilization = attribute.Key("system.memory.utilization")

	// SystemPagingUsage is the attribute Key conforming to the
	// "system.paging.usage" semantic conventions. It represents the unix swap or
	// windows pagefile usage
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemPagingUsage = attribute.Key("system.paging.usage")

	// SystemPagingUtilization is the attribute Key conforming to the
	// "system.paging.utilization" semantic conventions. It represents the
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemPagingUtilization = attribute.Key("system.paging.utilization")

	// SystemPagingFaults is the attribute Key conforming to the
	// "system.paging.faults" semantic conventions. It represents the
	// Instrument: counter
	// Unit: {fault}
	// Stability: None
	SystemPagingFaults = attribute.Key("system.paging.faults")

	// SystemPagingOperations is the attribute Key conforming to the
	// "system.paging.operations" semantic conventions. It represents the
	// Instrument: counter
	// Unit: {operation}
	// Stability: None
	SystemPagingOperations = attribute.Key("system.paging.operations")

	// SystemDiskIo is the attribute Key conforming to the "system.disk.io"
	// semantic conventions. It represents the
	// Instrument: counter
	// Unit: By
	// Stability: None
	SystemDiskIo = attribute.Key("system.disk.io")

	// SystemDiskOperations is the attribute Key conforming to the
	// "system.disk.operations" semantic conventions. It represents the
	// Instrument: counter
	// Unit: {operation}
	// Stability: None
	SystemDiskOperations = attribute.Key("system.disk.operations")

	// SystemDiskIoTime is the attribute Key conforming to the
	// "system.disk.io_time" semantic conventions. It represents the time disk
	// spent activated
	// Instrument: counter
	// Unit: s
	// Stability: None
	SystemDiskIoTime = attribute.Key("system.disk.io_time")

	// SystemDiskOperationTime is the attribute Key conforming to the
	// "system.disk.operation_time" semantic conventions. It represents the sum of
	// the time each operation took to complete
	// Instrument: counter
	// Unit: s
	// Stability: None
	SystemDiskOperationTime = attribute.Key("system.disk.operation_time")

	// SystemDiskMerged is the attribute Key conforming to the "system.disk.merged"
	// semantic conventions. It represents the
	// Instrument: counter
	// Unit: {operation}
	// Stability: None
	SystemDiskMerged = attribute.Key("system.disk.merged")

	// SystemFilesystemUsage is the attribute Key conforming to the
	// "system.filesystem.usage" semantic conventions. It represents the
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemFilesystemUsage = attribute.Key("system.filesystem.usage")

	// SystemFilesystemUtilization is the attribute Key conforming to the
	// "system.filesystem.utilization" semantic conventions. It represents the
	// Instrument: gauge
	// Unit: 1
	// Stability: None
	SystemFilesystemUtilization = attribute.Key("system.filesystem.utilization")

	// SystemNetworkDropped is the attribute Key conforming to the
	// "system.network.dropped" semantic conventions. It represents the count of
	// packets that are dropped or discarded even though there was no error
	// Instrument: counter
	// Unit: {packet}
	// Stability: None
	SystemNetworkDropped = attribute.Key("system.network.dropped")

	// SystemNetworkPackets is the attribute Key conforming to the
	// "system.network.packets" semantic conventions. It represents the
	// Instrument: counter
	// Unit: {packet}
	// Stability: None
	SystemNetworkPackets = attribute.Key("system.network.packets")

	// SystemNetworkErrors is the attribute Key conforming to the
	// "system.network.errors" semantic conventions. It represents the count of
	// network errors detected
	// Instrument: counter
	// Unit: {error}
	// Stability: None
	SystemNetworkErrors = attribute.Key("system.network.errors")

	// SystemNetworkIo is the attribute Key conforming to the "system.network.io"
	// semantic conventions. It represents the
	// Instrument: counter
	// Unit: By
	// Stability: None
	SystemNetworkIo = attribute.Key("system.network.io")

	// SystemNetworkConnections is the attribute Key conforming to the
	// "system.network.connections" semantic conventions. It represents the
	// Instrument: updowncounter
	// Unit: {connection}
	// Stability: None
	SystemNetworkConnections = attribute.Key("system.network.connections")

	// SystemProcessesCount is the attribute Key conforming to the
	// "system.processes.count" semantic conventions. It represents the total
	// number of processes in each state
	// Instrument: updowncounter
	// Unit: {process}
	// Stability: None
	SystemProcessesCount = attribute.Key("system.processes.count")

	// SystemProcessesCreated is the attribute Key conforming to the
	// "system.processes.created" semantic conventions. It represents the total
	// number of processes created over uptime of the host
	// Instrument: counter
	// Unit: {process}
	// Stability: None
	SystemProcessesCreated = attribute.Key("system.processes.created")

	// SystemLinuxMemoryAvailable is the attribute Key conforming to the
	// "system.linux.memory.available" semantic conventions. It represents an
	// estimate of how much memory is available for starting new applications,
	// without causing swapping
	// Instrument: updowncounter
	// Unit: By
	// Stability: None
	SystemLinuxMemoryAvailable = attribute.Key("system.linux.memory.available")
)
