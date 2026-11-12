// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv_test

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/v1.45.0/azureconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/cicdconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/containerconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/dbconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/dnsconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/faasconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/goconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/httpconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/hwconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/k8sconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/messagingconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/nfsconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/openshiftconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/otelconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/processconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/rpcconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/signalrconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/systemconv"
	"go.opentelemetry.io/otel/semconv/v1.45.0/vcsconv"
)

func TestMetricConstructorOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"azureconv/CosmosDBClientActiveInstanceCount", constructorTest(azureconv.NewCosmosDBClientActiveInstanceCount, "Int64UpDownCounterConfig", "azure.cosmosdb.client.active_instance.count", "Number of active client instances.", "{instance}", false, nil)},
		{"azureconv/CosmosDBClientActiveInstanceCountObservable", constructorTest(azureconv.NewCosmosDBClientActiveInstanceCountObservable, "Int64ObservableUpDownCounterConfig", "azure.cosmosdb.client.active_instance.count", "Number of active client instances.", "{instance}", false, nil)},
		{"azureconv/CosmosDBClientOperationRequestCharge", constructorTest(azureconv.NewCosmosDBClientOperationRequestCharge, "Float64HistogramConfig", "azure.cosmosdb.client.operation.request_charge", "[Request units](https://learn.microsoft.com/azure/cosmos-db/request-units) consumed by the operation.", "{request_unit}", true, nil)},
		{"cicdconv/PipelineRunActive", constructorTest(cicdconv.NewPipelineRunActive, "Int64UpDownCounterConfig", "cicd.pipeline.run.active", "The number of pipeline runs currently active in the system by state.", "{run}", false, nil)},
		{"cicdconv/PipelineRunActiveObservable", constructorTest(cicdconv.NewPipelineRunActiveObservable, "Int64ObservableUpDownCounterConfig", "cicd.pipeline.run.active", "The number of pipeline runs currently active in the system by state.", "{run}", false, nil)},
		{"cicdconv/PipelineRunDuration", constructorTest(cicdconv.NewPipelineRunDuration, "Float64HistogramConfig", "cicd.pipeline.run.duration", "Duration of a pipeline run grouped by pipeline, state and result.", "s", true, nil)},
		{"cicdconv/PipelineRunErrors", constructorTest(cicdconv.NewPipelineRunErrors, "Int64CounterConfig", "cicd.pipeline.run.errors", "The number of errors encountered in pipeline runs (eg. compile, test failures).", "{error}", false, nil)},
		{"cicdconv/PipelineRunErrorsObservable", constructorTest(cicdconv.NewPipelineRunErrorsObservable, "Int64ObservableCounterConfig", "cicd.pipeline.run.errors", "The number of errors encountered in pipeline runs (eg. compile, test failures).", "{error}", false, nil)},
		{"cicdconv/SystemErrors", constructorTest(cicdconv.NewSystemErrors, "Int64CounterConfig", "cicd.system.errors", "The number of errors in a component of the CI/CD system (eg. controller, scheduler, agent).", "{error}", false, nil)},
		{"cicdconv/SystemErrorsObservable", constructorTest(cicdconv.NewSystemErrorsObservable, "Int64ObservableCounterConfig", "cicd.system.errors", "The number of errors in a component of the CI/CD system (eg. controller, scheduler, agent).", "{error}", false, nil)},
		{"cicdconv/WorkerCount", constructorTest(cicdconv.NewWorkerCount, "Int64UpDownCounterConfig", "cicd.worker.count", "The number of workers on the CI/CD system by state.", "{worker}", false, nil)},
		{"cicdconv/WorkerCountObservable", constructorTest(cicdconv.NewWorkerCountObservable, "Int64ObservableUpDownCounterConfig", "cicd.worker.count", "The number of workers on the CI/CD system by state.", "{worker}", false, nil)},
		{"containerconv/CPUTime", constructorTest(containerconv.NewCPUTime, "Float64CounterConfig", "container.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"containerconv/CPUTimeObservable", constructorTest(containerconv.NewCPUTimeObservable, "Float64ObservableCounterConfig", "container.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"containerconv/CPUUsage", constructorTest(containerconv.NewCPUUsage, "Float64GaugeConfig", "container.cpu.usage", "Container\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"containerconv/CPUUsageObservable", constructorTest(containerconv.NewCPUUsageObservable, "Float64ObservableGaugeConfig", "container.cpu.usage", "Container\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"containerconv/DiskIO", constructorTest(containerconv.NewDiskIO, "Int64CounterConfig", "container.disk.io", "Disk bytes for the container.", "By", false, nil)},
		{"containerconv/DiskIOObservable", constructorTest(containerconv.NewDiskIOObservable, "Int64ObservableCounterConfig", "container.disk.io", "Disk bytes for the container.", "By", false, nil)},
		{"containerconv/FilesystemAvailable", constructorTest(containerconv.NewFilesystemAvailable, "Int64UpDownCounterConfig", "container.filesystem.available", "Container filesystem available bytes.", "By", false, nil)},
		{"containerconv/FilesystemAvailableObservable", constructorTest(containerconv.NewFilesystemAvailableObservable, "Int64ObservableUpDownCounterConfig", "container.filesystem.available", "Container filesystem available bytes.", "By", false, nil)},
		{"containerconv/FilesystemCapacity", constructorTest(containerconv.NewFilesystemCapacity, "Int64UpDownCounterConfig", "container.filesystem.capacity", "Container filesystem capacity.", "By", false, nil)},
		{"containerconv/FilesystemCapacityObservable", constructorTest(containerconv.NewFilesystemCapacityObservable, "Int64ObservableUpDownCounterConfig", "container.filesystem.capacity", "Container filesystem capacity.", "By", false, nil)},
		{"containerconv/FilesystemUsage", constructorTest(containerconv.NewFilesystemUsage, "Int64UpDownCounterConfig", "container.filesystem.usage", "Container filesystem usage.", "By", false, nil)},
		{"containerconv/FilesystemUsageObservable", constructorTest(containerconv.NewFilesystemUsageObservable, "Int64ObservableUpDownCounterConfig", "container.filesystem.usage", "Container filesystem usage.", "By", false, nil)},
		{"containerconv/MemoryAvailable", constructorTest(containerconv.NewMemoryAvailable, "Int64UpDownCounterConfig", "container.memory.available", "Container memory available.", "By", false, nil)},
		{"containerconv/MemoryAvailableObservable", constructorTest(containerconv.NewMemoryAvailableObservable, "Int64ObservableUpDownCounterConfig", "container.memory.available", "Container memory available.", "By", false, nil)},
		{"containerconv/MemoryRss", constructorTest(containerconv.NewMemoryRss, "Int64UpDownCounterConfig", "container.memory.rss", "Container memory RSS.", "By", false, nil)},
		{"containerconv/MemoryRssObservable", constructorTest(containerconv.NewMemoryRssObservable, "Int64ObservableUpDownCounterConfig", "container.memory.rss", "Container memory RSS.", "By", false, nil)},
		{"containerconv/MemoryUsage", constructorTest(containerconv.NewMemoryUsage, "Int64UpDownCounterConfig", "container.memory.usage", "Memory usage of the container.", "By", false, nil)},
		{"containerconv/MemoryUsageObservable", constructorTest(containerconv.NewMemoryUsageObservable, "Int64ObservableUpDownCounterConfig", "container.memory.usage", "Memory usage of the container.", "By", false, nil)},
		{"containerconv/MemoryWorkingSet", constructorTest(containerconv.NewMemoryWorkingSet, "Int64UpDownCounterConfig", "container.memory.working_set", "Container memory working set.", "By", false, nil)},
		{"containerconv/MemoryWorkingSetObservable", constructorTest(containerconv.NewMemoryWorkingSetObservable, "Int64ObservableUpDownCounterConfig", "container.memory.working_set", "Container memory working set.", "By", false, nil)},
		{"containerconv/NetworkIO", constructorTest(containerconv.NewNetworkIO, "Int64CounterConfig", "container.network.io", "Network bytes for the container.", "By", false, nil)},
		{"containerconv/NetworkIOObservable", constructorTest(containerconv.NewNetworkIOObservable, "Int64ObservableCounterConfig", "container.network.io", "Network bytes for the container.", "By", false, nil)},
		{"containerconv/PagingFaults", constructorTest(containerconv.NewPagingFaults, "Int64CounterConfig", "container.paging.faults", "Container memory paging faults.", "{fault}", false, nil)},
		{"containerconv/PagingFaultsObservable", constructorTest(containerconv.NewPagingFaultsObservable, "Int64ObservableCounterConfig", "container.paging.faults", "Container memory paging faults.", "{fault}", false, nil)},
		{"containerconv/Uptime", constructorTest(containerconv.NewUptime, "Float64GaugeConfig", "container.uptime", "The time the container has been running.", "s", false, nil)},
		{"containerconv/UptimeObservable", constructorTest(containerconv.NewUptimeObservable, "Float64ObservableGaugeConfig", "container.uptime", "The time the container has been running.", "s", false, nil)},
		{"dbconv/ClientConnectionCount", constructorTest(dbconv.NewClientConnectionCount, "Int64UpDownCounterConfig", "db.client.connection.count", "The number of connections that are currently in state described by the `state` attribute.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionCountObservable", constructorTest(dbconv.NewClientConnectionCountObservable, "Int64ObservableUpDownCounterConfig", "db.client.connection.count", "The number of connections that are currently in state described by the `state` attribute.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionCreateTime", constructorTest(dbconv.NewClientConnectionCreateTime, "Float64HistogramConfig", "db.client.connection.create_time", "The time it took to create a new connection.", "s", true, nil)},
		{"dbconv/ClientConnectionIdleMax", constructorTest(dbconv.NewClientConnectionIdleMax, "Int64UpDownCounterConfig", "db.client.connection.idle.max", "The maximum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionIdleMaxObservable", constructorTest(dbconv.NewClientConnectionIdleMaxObservable, "Int64ObservableUpDownCounterConfig", "db.client.connection.idle.max", "The maximum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionIdleMin", constructorTest(dbconv.NewClientConnectionIdleMin, "Int64UpDownCounterConfig", "db.client.connection.idle.min", "The minimum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionIdleMinObservable", constructorTest(dbconv.NewClientConnectionIdleMinObservable, "Int64ObservableUpDownCounterConfig", "db.client.connection.idle.min", "The minimum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionLimit", constructorTest(dbconv.NewClientConnectionLimit, "Int64UpDownCounterConfig", "db.client.connection.limit", "The maximum number of open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionLimitObservable", constructorTest(dbconv.NewClientConnectionLimitObservable, "Int64ObservableUpDownCounterConfig", "db.client.connection.limit", "The maximum number of open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionPendingRequests", constructorTest(dbconv.NewClientConnectionPendingRequests, "Int64UpDownCounterConfig", "db.client.connection.pending_requests", "The number of current pending requests for an open connection.", "{request}", false, nil)},
		{"dbconv/ClientConnectionPendingRequestsObservable", constructorTest(dbconv.NewClientConnectionPendingRequestsObservable, "Int64ObservableUpDownCounterConfig", "db.client.connection.pending_requests", "The number of current pending requests for an open connection.", "{request}", false, nil)},
		{"dbconv/ClientConnectionTimeouts", constructorTest(dbconv.NewClientConnectionTimeouts, "Int64CounterConfig", "db.client.connection.timeouts", "The number of connection timeouts that have occurred trying to obtain a connection from the pool.", "{timeout}", false, nil)},
		{"dbconv/ClientConnectionTimeoutsObservable", constructorTest(dbconv.NewClientConnectionTimeoutsObservable, "Int64ObservableCounterConfig", "db.client.connection.timeouts", "The number of connection timeouts that have occurred trying to obtain a connection from the pool.", "{timeout}", false, nil)},
		{"dbconv/ClientConnectionUseTime", constructorTest(dbconv.NewClientConnectionUseTime, "Float64HistogramConfig", "db.client.connection.use_time", "The time between borrowing a connection and returning it to the pool.", "s", true, nil)},
		{"dbconv/ClientConnectionWaitTime", constructorTest(dbconv.NewClientConnectionWaitTime, "Float64HistogramConfig", "db.client.connection.wait_time", "The time it took to obtain an open connection from the pool.", "s", true, nil)},
		{"dbconv/ClientOperationDuration", constructorTest(dbconv.NewClientOperationDuration, "Float64HistogramConfig", "db.client.operation.duration", "Duration of database client operations.", "s", true, nil)},
		{"dbconv/ClientResponseReturnedRows", constructorTest(dbconv.NewClientResponseReturnedRows, "Int64HistogramConfig", "db.client.response.returned_rows", "The actual number of records returned by the database operation.", "{row}", true, nil)},
		{"dnsconv/LookupDuration", constructorTest(dnsconv.NewLookupDuration, "Float64HistogramConfig", "dns.lookup.duration", "Measures the time taken to perform a DNS lookup.", "s", true, nil)},
		{"faasconv/Coldstarts", constructorTest(faasconv.NewColdstarts, "Int64CounterConfig", "faas.coldstarts", "Number of invocation cold starts.", "{coldstart}", false, nil)},
		{"faasconv/ColdstartsObservable", constructorTest(faasconv.NewColdstartsObservable, "Int64ObservableCounterConfig", "faas.coldstarts", "Number of invocation cold starts.", "{coldstart}", false, nil)},
		{"faasconv/CPUUsage", constructorTest(faasconv.NewCPUUsage, "Float64HistogramConfig", "faas.cpu_usage", "Distribution of CPU usage per invocation.", "s", true, nil)},
		{"faasconv/Errors", constructorTest(faasconv.NewErrors, "Int64CounterConfig", "faas.errors", "Number of invocation errors.", "{error}", false, nil)},
		{"faasconv/ErrorsObservable", constructorTest(faasconv.NewErrorsObservable, "Int64ObservableCounterConfig", "faas.errors", "Number of invocation errors.", "{error}", false, nil)},
		{"faasconv/InitDuration", constructorTest(faasconv.NewInitDuration, "Float64HistogramConfig", "faas.init_duration", "Measures the duration of the function\u0027s initialization, such as a cold start.", "s", true, nil)},
		{"faasconv/Invocations", constructorTest(faasconv.NewInvocations, "Int64CounterConfig", "faas.invocations", "Number of successful invocations.", "{invocation}", false, nil)},
		{"faasconv/InvocationsObservable", constructorTest(faasconv.NewInvocationsObservable, "Int64ObservableCounterConfig", "faas.invocations", "Number of successful invocations.", "{invocation}", false, nil)},
		{"faasconv/InvokeDuration", constructorTest(faasconv.NewInvokeDuration, "Float64HistogramConfig", "faas.invoke_duration", "Measures the duration of the function\u0027s logic execution.", "s", true, nil)},
		{"faasconv/MemUsage", constructorTest(faasconv.NewMemUsage, "Int64HistogramConfig", "faas.mem_usage", "Distribution of max memory usage per invocation.", "By", true, nil)},
		{"faasconv/NetIO", constructorTest(faasconv.NewNetIO, "Int64HistogramConfig", "faas.net_io", "Distribution of net I/O usage per invocation.", "By", true, nil)},
		{"faasconv/Timeouts", constructorTest(faasconv.NewTimeouts, "Int64CounterConfig", "faas.timeouts", "Number of invocation timeouts.", "{timeout}", false, nil)},
		{"faasconv/TimeoutsObservable", constructorTest(faasconv.NewTimeoutsObservable, "Int64ObservableCounterConfig", "faas.timeouts", "Number of invocation timeouts.", "{timeout}", false, nil)},
		{"goconv/ConfigGogc", constructorTest(goconv.NewConfigGogc, "Int64ObservableUpDownCounterConfig", "go.config.gogc", "Heap size target percentage configured by the user, otherwise 100.", "%", false, nil)},
		{"goconv/CPUTime", constructorTest(goconv.NewCPUTime, "Float64CounterConfig", "go.cpu.time", "Estimated CPU time spent by the Go runtime.", "s", false, nil)},
		{"goconv/CPUTimeObservable", constructorTest(goconv.NewCPUTimeObservable, "Float64ObservableCounterConfig", "go.cpu.time", "Estimated CPU time spent by the Go runtime.", "s", false, nil)},
		{"goconv/GoroutineCount", constructorTest(goconv.NewGoroutineCount, "Int64ObservableUpDownCounterConfig", "go.goroutine.count", "Count of live goroutines.", "{goroutine}", false, nil)},
		{"goconv/MemoryAllocated", constructorTest(goconv.NewMemoryAllocated, "Int64ObservableCounterConfig", "go.memory.allocated", "Memory allocated to the heap by the application.", "By", false, nil)},
		{"goconv/MemoryAllocations", constructorTest(goconv.NewMemoryAllocations, "Int64ObservableCounterConfig", "go.memory.allocations", "Count of allocations to the heap by the application.", "{allocation}", false, nil)},
		{"goconv/MemoryGCCycles", constructorTest(goconv.NewMemoryGCCycles, "Int64CounterConfig", "go.memory.gc.cycles", "Number of completed GC cycles.", "{gc_cycle}", false, nil)},
		{"goconv/MemoryGCCyclesObservable", constructorTest(goconv.NewMemoryGCCyclesObservable, "Int64ObservableCounterConfig", "go.memory.gc.cycles", "Number of completed GC cycles.", "{gc_cycle}", false, nil)},
		{"goconv/MemoryGCGoal", constructorTest(goconv.NewMemoryGCGoal, "Int64ObservableUpDownCounterConfig", "go.memory.gc.goal", "Heap size target for the end of the GC cycle.", "By", false, nil)},
		{"goconv/MemoryGCPauseDuration", constructorTest(goconv.NewMemoryGCPauseDuration, "Float64HistogramConfig", "go.memory.gc.pause.duration", "Distribution of individual GC-related stop-the-world pause latencies. This is the time from deciding to stop the world until the world is started again.", "s", true, nil)},
		{"goconv/MemoryLimit", constructorTest(goconv.NewMemoryLimit, "Int64ObservableUpDownCounterConfig", "go.memory.limit", "Go runtime memory limit configured by the user, if a limit exists.", "By", false, nil)},
		{"goconv/MemoryUsed", constructorTest(goconv.NewMemoryUsed, "Int64ObservableUpDownCounterConfig", "go.memory.used", "Memory used by the Go runtime.", "By", false, nil)},
		{"goconv/ProcessorLimit", constructorTest(goconv.NewProcessorLimit, "Int64ObservableUpDownCounterConfig", "go.processor.limit", "The number of OS threads that can execute user-level Go code simultaneously.", "{thread}", false, nil)},
		{"goconv/ScheduleDuration", constructorTest(goconv.NewScheduleDuration, "Float64HistogramConfig", "go.schedule.duration", "The time goroutines have spent in the scheduler in a runnable state before actually running.", "s", true, nil)},
		{"httpconv/ClientActiveRequests", constructorTest(httpconv.NewClientActiveRequests, "Int64UpDownCounterConfig", "http.client.active_requests", "Number of active HTTP requests.", "{request}", false, nil)},
		{"httpconv/ClientActiveRequestsObservable", constructorTest(httpconv.NewClientActiveRequestsObservable, "Int64ObservableUpDownCounterConfig", "http.client.active_requests", "Number of active HTTP requests.", "{request}", false, nil)},
		{"httpconv/ClientConnectionDuration", constructorTest(httpconv.NewClientConnectionDuration, "Float64HistogramConfig", "http.client.connection.duration", "The duration of the successfully established outbound HTTP connections.", "s", true, nil)},
		{"httpconv/ClientOpenConnections", constructorTest(httpconv.NewClientOpenConnections, "Int64UpDownCounterConfig", "http.client.open_connections", "Number of outbound HTTP connections that are currently active or idle on the client.", "{connection}", false, nil)},
		{"httpconv/ClientOpenConnectionsObservable", constructorTest(httpconv.NewClientOpenConnectionsObservable, "Int64ObservableUpDownCounterConfig", "http.client.open_connections", "Number of outbound HTTP connections that are currently active or idle on the client.", "{connection}", false, nil)},
		{"httpconv/ClientRequestBodySize", constructorTest(httpconv.NewClientRequestBodySize, "Int64HistogramConfig", "http.client.request.body.size", "Size of HTTP client request bodies.", "By", true, nil)},
		{"httpconv/ClientRequestDuration", constructorTest(httpconv.NewClientRequestDuration, "Float64HistogramConfig", "http.client.request.duration", "Duration of HTTP client requests.", "s", true, []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10})},
		{"httpconv/ClientResponseBodySize", constructorTest(httpconv.NewClientResponseBodySize, "Int64HistogramConfig", "http.client.response.body.size", "Size of HTTP client response bodies.", "By", true, nil)},
		{"httpconv/ServerActiveRequests", constructorTest(httpconv.NewServerActiveRequests, "Int64UpDownCounterConfig", "http.server.active_requests", "Number of active HTTP server requests.", "{request}", false, nil)},
		{"httpconv/ServerActiveRequestsObservable", constructorTest(httpconv.NewServerActiveRequestsObservable, "Int64ObservableUpDownCounterConfig", "http.server.active_requests", "Number of active HTTP server requests.", "{request}", false, nil)},
		{"httpconv/ServerRequestBodySize", constructorTest(httpconv.NewServerRequestBodySize, "Int64HistogramConfig", "http.server.request.body.size", "Size of HTTP server request bodies.", "By", true, nil)},
		{"httpconv/ServerRequestDuration", constructorTest(httpconv.NewServerRequestDuration, "Float64HistogramConfig", "http.server.request.duration", "Duration of HTTP server requests.", "s", true, []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10})},
		{"httpconv/ServerResponseBodySize", constructorTest(httpconv.NewServerResponseBodySize, "Int64HistogramConfig", "http.server.response.body.size", "Size of HTTP server response bodies.", "By", true, nil)},
		{"hwconv/BatteryCharge", constructorTest(hwconv.NewBatteryCharge, "Float64GaugeConfig", "hw.battery.charge", "Remaining fraction of battery charge.", "1", false, nil)},
		{"hwconv/BatteryChargeObservable", constructorTest(hwconv.NewBatteryChargeObservable, "Float64ObservableGaugeConfig", "hw.battery.charge", "Remaining fraction of battery charge.", "1", false, nil)},
		{"hwconv/BatteryChargeLimit", constructorTest(hwconv.NewBatteryChargeLimit, "Float64GaugeConfig", "hw.battery.charge.limit", "Lower limit of battery charge fraction to ensure proper operation.", "1", false, nil)},
		{"hwconv/BatteryChargeLimitObservable", constructorTest(hwconv.NewBatteryChargeLimitObservable, "Float64ObservableGaugeConfig", "hw.battery.charge.limit", "Lower limit of battery charge fraction to ensure proper operation.", "1", false, nil)},
		{"hwconv/BatteryTimeLeft", constructorTest(hwconv.NewBatteryTimeLeft, "Float64GaugeConfig", "hw.battery.time_left", "Time left before battery is completely charged or discharged.", "s", false, nil)},
		{"hwconv/BatteryTimeLeftObservable", constructorTest(hwconv.NewBatteryTimeLeftObservable, "Float64ObservableGaugeConfig", "hw.battery.time_left", "Time left before battery is completely charged or discharged.", "s", false, nil)},
		{"hwconv/CPUSpeed", constructorTest(hwconv.NewCPUSpeed, "Int64GaugeConfig", "hw.cpu.speed", "CPU current frequency.", "Hz", false, nil)},
		{"hwconv/CPUSpeedObservable", constructorTest(hwconv.NewCPUSpeedObservable, "Int64ObservableGaugeConfig", "hw.cpu.speed", "CPU current frequency.", "Hz", false, nil)},
		{"hwconv/CPUSpeedLimit", constructorTest(hwconv.NewCPUSpeedLimit, "Int64GaugeConfig", "hw.cpu.speed.limit", "CPU maximum frequency.", "Hz", false, nil)},
		{"hwconv/CPUSpeedLimitObservable", constructorTest(hwconv.NewCPUSpeedLimitObservable, "Int64ObservableGaugeConfig", "hw.cpu.speed.limit", "CPU maximum frequency.", "Hz", false, nil)},
		{"hwconv/Energy", constructorTest(hwconv.NewEnergy, "Float64CounterConfig", "hw.energy", "Energy consumed by the component.", "J", false, nil)},
		{"hwconv/EnergyObservable", constructorTest(hwconv.NewEnergyObservable, "Float64ObservableCounterConfig", "hw.energy", "Energy consumed by the component.", "J", false, nil)},
		{"hwconv/Errors", constructorTest(hwconv.NewErrors, "Int64CounterConfig", "hw.errors", "Number of errors encountered by the component.", "{error}", false, nil)},
		{"hwconv/ErrorsObservable", constructorTest(hwconv.NewErrorsObservable, "Int64ObservableCounterConfig", "hw.errors", "Number of errors encountered by the component.", "{error}", false, nil)},
		{"hwconv/FanSpeed", constructorTest(hwconv.NewFanSpeed, "Int64GaugeConfig", "hw.fan.speed", "Fan speed in revolutions per minute.", "rpm", false, nil)},
		{"hwconv/FanSpeedObservable", constructorTest(hwconv.NewFanSpeedObservable, "Int64ObservableGaugeConfig", "hw.fan.speed", "Fan speed in revolutions per minute.", "rpm", false, nil)},
		{"hwconv/FanSpeedLimit", constructorTest(hwconv.NewFanSpeedLimit, "Int64GaugeConfig", "hw.fan.speed.limit", "Speed limit in rpm.", "rpm", false, nil)},
		{"hwconv/FanSpeedLimitObservable", constructorTest(hwconv.NewFanSpeedLimitObservable, "Int64ObservableGaugeConfig", "hw.fan.speed.limit", "Speed limit in rpm.", "rpm", false, nil)},
		{"hwconv/FanSpeedRatio", constructorTest(hwconv.NewFanSpeedRatio, "Float64GaugeConfig", "hw.fan.speed_ratio", "Fan speed expressed as a fraction of its maximum speed.", "1", false, nil)},
		{"hwconv/FanSpeedRatioObservable", constructorTest(hwconv.NewFanSpeedRatioObservable, "Float64ObservableGaugeConfig", "hw.fan.speed_ratio", "Fan speed expressed as a fraction of its maximum speed.", "1", false, nil)},
		{"hwconv/GpuIO", constructorTest(hwconv.NewGpuIO, "Int64CounterConfig", "hw.gpu.io", "Received and transmitted bytes by the GPU.", "By", false, nil)},
		{"hwconv/GpuIOObservable", constructorTest(hwconv.NewGpuIOObservable, "Int64ObservableCounterConfig", "hw.gpu.io", "Received and transmitted bytes by the GPU.", "By", false, nil)},
		{"hwconv/GpuMemoryLimit", constructorTest(hwconv.NewGpuMemoryLimit, "Int64UpDownCounterConfig", "hw.gpu.memory.limit", "Size of the GPU memory.", "By", false, nil)},
		{"hwconv/GpuMemoryLimitObservable", constructorTest(hwconv.NewGpuMemoryLimitObservable, "Int64ObservableUpDownCounterConfig", "hw.gpu.memory.limit", "Size of the GPU memory.", "By", false, nil)},
		{"hwconv/GpuMemoryUsage", constructorTest(hwconv.NewGpuMemoryUsage, "Int64UpDownCounterConfig", "hw.gpu.memory.usage", "GPU memory used.", "By", false, nil)},
		{"hwconv/GpuMemoryUsageObservable", constructorTest(hwconv.NewGpuMemoryUsageObservable, "Int64ObservableUpDownCounterConfig", "hw.gpu.memory.usage", "GPU memory used.", "By", false, nil)},
		{"hwconv/GpuMemoryUtilization", constructorTest(hwconv.NewGpuMemoryUtilization, "Float64GaugeConfig", "hw.gpu.memory.utilization", "Fraction of GPU memory used.", "1", false, nil)},
		{"hwconv/GpuMemoryUtilizationObservable", constructorTest(hwconv.NewGpuMemoryUtilizationObservable, "Float64ObservableGaugeConfig", "hw.gpu.memory.utilization", "Fraction of GPU memory used.", "1", false, nil)},
		{"hwconv/GpuUtilization", constructorTest(hwconv.NewGpuUtilization, "Float64GaugeConfig", "hw.gpu.utilization", "Fraction of time spent in a specific task.", "1", false, nil)},
		{"hwconv/GpuUtilizationObservable", constructorTest(hwconv.NewGpuUtilizationObservable, "Float64ObservableGaugeConfig", "hw.gpu.utilization", "Fraction of time spent in a specific task.", "1", false, nil)},
		{"hwconv/HostAmbientTemperature", constructorTest(hwconv.NewHostAmbientTemperature, "Float64GaugeConfig", "hw.host.ambient_temperature", "Ambient (external) temperature of the physical host.", "Cel", false, nil)},
		{"hwconv/HostAmbientTemperatureObservable", constructorTest(hwconv.NewHostAmbientTemperatureObservable, "Float64ObservableGaugeConfig", "hw.host.ambient_temperature", "Ambient (external) temperature of the physical host.", "Cel", false, nil)},
		{"hwconv/HostEnergy", constructorTest(hwconv.NewHostEnergy, "Int64CounterConfig", "hw.host.energy", "Total energy consumed by the entire physical host, in joules.", "J", false, nil)},
		{"hwconv/HostEnergyObservable", constructorTest(hwconv.NewHostEnergyObservable, "Int64ObservableCounterConfig", "hw.host.energy", "Total energy consumed by the entire physical host, in joules.", "J", false, nil)},
		{"hwconv/HostHeatingMargin", constructorTest(hwconv.NewHostHeatingMargin, "Float64GaugeConfig", "hw.host.heating_margin", "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors.", "Cel", false, nil)},
		{"hwconv/HostHeatingMarginObservable", constructorTest(hwconv.NewHostHeatingMarginObservable, "Float64ObservableGaugeConfig", "hw.host.heating_margin", "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors.", "Cel", false, nil)},
		{"hwconv/HostPower", constructorTest(hwconv.NewHostPower, "Float64GaugeConfig", "hw.host.power", "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred).", "W", false, nil)},
		{"hwconv/HostPowerObservable", constructorTest(hwconv.NewHostPowerObservable, "Float64ObservableGaugeConfig", "hw.host.power", "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred).", "W", false, nil)},
		{"hwconv/LogicalDiskLimit", constructorTest(hwconv.NewLogicalDiskLimit, "Int64UpDownCounterConfig", "hw.logical_disk.limit", "Size of the logical disk.", "By", false, nil)},
		{"hwconv/LogicalDiskLimitObservable", constructorTest(hwconv.NewLogicalDiskLimitObservable, "Int64ObservableUpDownCounterConfig", "hw.logical_disk.limit", "Size of the logical disk.", "By", false, nil)},
		{"hwconv/LogicalDiskUsage", constructorTest(hwconv.NewLogicalDiskUsage, "Int64UpDownCounterConfig", "hw.logical_disk.usage", "Logical disk space usage.", "By", false, nil)},
		{"hwconv/LogicalDiskUsageObservable", constructorTest(hwconv.NewLogicalDiskUsageObservable, "Int64ObservableUpDownCounterConfig", "hw.logical_disk.usage", "Logical disk space usage.", "By", false, nil)},
		{"hwconv/LogicalDiskUtilization", constructorTest(hwconv.NewLogicalDiskUtilization, "Float64GaugeConfig", "hw.logical_disk.utilization", "Logical disk space utilization as a fraction.", "1", false, nil)},
		{"hwconv/LogicalDiskUtilizationObservable", constructorTest(hwconv.NewLogicalDiskUtilizationObservable, "Float64ObservableGaugeConfig", "hw.logical_disk.utilization", "Logical disk space utilization as a fraction.", "1", false, nil)},
		{"hwconv/MemorySize", constructorTest(hwconv.NewMemorySize, "Int64UpDownCounterConfig", "hw.memory.size", "Size of the memory module.", "By", false, nil)},
		{"hwconv/MemorySizeObservable", constructorTest(hwconv.NewMemorySizeObservable, "Int64ObservableUpDownCounterConfig", "hw.memory.size", "Size of the memory module.", "By", false, nil)},
		{"hwconv/NetworkBandwidthLimit", constructorTest(hwconv.NewNetworkBandwidthLimit, "Int64UpDownCounterConfig", "hw.network.bandwidth.limit", "Link speed.", "By/s", false, nil)},
		{"hwconv/NetworkBandwidthLimitObservable", constructorTest(hwconv.NewNetworkBandwidthLimitObservable, "Int64ObservableUpDownCounterConfig", "hw.network.bandwidth.limit", "Link speed.", "By/s", false, nil)},
		{"hwconv/NetworkBandwidthUtilization", constructorTest(hwconv.NewNetworkBandwidthUtilization, "Float64GaugeConfig", "hw.network.bandwidth.utilization", "Utilization of the network bandwidth as a fraction.", "1", false, nil)},
		{"hwconv/NetworkBandwidthUtilizationObservable", constructorTest(hwconv.NewNetworkBandwidthUtilizationObservable, "Float64ObservableGaugeConfig", "hw.network.bandwidth.utilization", "Utilization of the network bandwidth as a fraction.", "1", false, nil)},
		{"hwconv/NetworkIO", constructorTest(hwconv.NewNetworkIO, "Int64CounterConfig", "hw.network.io", "Received and transmitted network traffic in bytes.", "By", false, nil)},
		{"hwconv/NetworkIOObservable", constructorTest(hwconv.NewNetworkIOObservable, "Int64ObservableCounterConfig", "hw.network.io", "Received and transmitted network traffic in bytes.", "By", false, nil)},
		{"hwconv/NetworkPackets", constructorTest(hwconv.NewNetworkPackets, "Int64CounterConfig", "hw.network.packets", "Received and transmitted network traffic in packets (or frames).", "{packet}", false, nil)},
		{"hwconv/NetworkPacketsObservable", constructorTest(hwconv.NewNetworkPacketsObservable, "Int64ObservableCounterConfig", "hw.network.packets", "Received and transmitted network traffic in packets (or frames).", "{packet}", false, nil)},
		{"hwconv/NetworkUp", constructorTest(hwconv.NewNetworkUp, "Int64UpDownCounterConfig", "hw.network.up", "Link status: `1` (up) or `0` (down).", "1", false, nil)},
		{"hwconv/NetworkUpObservable", constructorTest(hwconv.NewNetworkUpObservable, "Int64ObservableUpDownCounterConfig", "hw.network.up", "Link status: `1` (up) or `0` (down).", "1", false, nil)},
		{"hwconv/PhysicalDiskEnduranceUtilization", constructorTest(hwconv.NewPhysicalDiskEnduranceUtilization, "Float64GaugeConfig", "hw.physical_disk.endurance_utilization", "Endurance remaining for this SSD disk.", "1", false, nil)},
		{"hwconv/PhysicalDiskEnduranceUtilizationObservable", constructorTest(hwconv.NewPhysicalDiskEnduranceUtilizationObservable, "Float64ObservableGaugeConfig", "hw.physical_disk.endurance_utilization", "Endurance remaining for this SSD disk.", "1", false, nil)},
		{"hwconv/PhysicalDiskSize", constructorTest(hwconv.NewPhysicalDiskSize, "Int64UpDownCounterConfig", "hw.physical_disk.size", "Size of the disk.", "By", false, nil)},
		{"hwconv/PhysicalDiskSizeObservable", constructorTest(hwconv.NewPhysicalDiskSizeObservable, "Int64ObservableUpDownCounterConfig", "hw.physical_disk.size", "Size of the disk.", "By", false, nil)},
		{"hwconv/PhysicalDiskSmart", constructorTest(hwconv.NewPhysicalDiskSmart, "Int64GaugeConfig", "hw.physical_disk.smart", "Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute.", "1", false, nil)},
		{"hwconv/PhysicalDiskSmartObservable", constructorTest(hwconv.NewPhysicalDiskSmartObservable, "Int64ObservableGaugeConfig", "hw.physical_disk.smart", "Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute.", "1", false, nil)},
		{"hwconv/Power", constructorTest(hwconv.NewPower, "Float64GaugeConfig", "hw.power", "Instantaneous power consumed by the component.", "W", false, nil)},
		{"hwconv/PowerObservable", constructorTest(hwconv.NewPowerObservable, "Float64ObservableGaugeConfig", "hw.power", "Instantaneous power consumed by the component.", "W", false, nil)},
		{"hwconv/PowerSupplyLimit", constructorTest(hwconv.NewPowerSupplyLimit, "Int64UpDownCounterConfig", "hw.power_supply.limit", "Maximum power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyLimitObservable", constructorTest(hwconv.NewPowerSupplyLimitObservable, "Int64ObservableUpDownCounterConfig", "hw.power_supply.limit", "Maximum power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyUsage", constructorTest(hwconv.NewPowerSupplyUsage, "Float64UpDownCounterConfig", "hw.power_supply.usage", "Current power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyUsageObservable", constructorTest(hwconv.NewPowerSupplyUsageObservable, "Float64ObservableUpDownCounterConfig", "hw.power_supply.usage", "Current power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyUtilization", constructorTest(hwconv.NewPowerSupplyUtilization, "Float64GaugeConfig", "hw.power_supply.utilization", "Utilization of the power supply as a fraction of its maximum output.", "1", false, nil)},
		{"hwconv/PowerSupplyUtilizationObservable", constructorTest(hwconv.NewPowerSupplyUtilizationObservable, "Float64ObservableGaugeConfig", "hw.power_supply.utilization", "Utilization of the power supply as a fraction of its maximum output.", "1", false, nil)},
		{"hwconv/Status", constructorTest(hwconv.NewStatus, "Int64UpDownCounterConfig", "hw.status", "Operational status: `1` (true) or `0` (false) for each of the possible states.", "1", false, nil)},
		{"hwconv/StatusObservable", constructorTest(hwconv.NewStatusObservable, "Int64ObservableUpDownCounterConfig", "hw.status", "Operational status: `1` (true) or `0` (false) for each of the possible states.", "1", false, nil)},
		{"hwconv/TapeDriveOperations", constructorTest(hwconv.NewTapeDriveOperations, "Int64CounterConfig", "hw.tape_drive.operations", "Operations performed by the tape drive.", "{operation}", false, nil)},
		{"hwconv/TapeDriveOperationsObservable", constructorTest(hwconv.NewTapeDriveOperationsObservable, "Int64ObservableCounterConfig", "hw.tape_drive.operations", "Operations performed by the tape drive.", "{operation}", false, nil)},
		{"hwconv/Temperature", constructorTest(hwconv.NewTemperature, "Float64GaugeConfig", "hw.temperature", "Temperature in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/TemperatureObservable", constructorTest(hwconv.NewTemperatureObservable, "Float64ObservableGaugeConfig", "hw.temperature", "Temperature in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/TemperatureLimit", constructorTest(hwconv.NewTemperatureLimit, "Float64GaugeConfig", "hw.temperature.limit", "Temperature limit in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/TemperatureLimitObservable", constructorTest(hwconv.NewTemperatureLimitObservable, "Float64ObservableGaugeConfig", "hw.temperature.limit", "Temperature limit in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/Voltage", constructorTest(hwconv.NewVoltage, "Float64GaugeConfig", "hw.voltage", "Voltage measured by the sensor.", "V", false, nil)},
		{"hwconv/VoltageObservable", constructorTest(hwconv.NewVoltageObservable, "Float64ObservableGaugeConfig", "hw.voltage", "Voltage measured by the sensor.", "V", false, nil)},
		{"hwconv/VoltageLimit", constructorTest(hwconv.NewVoltageLimit, "Float64GaugeConfig", "hw.voltage.limit", "Voltage limit in Volts.", "V", false, nil)},
		{"hwconv/VoltageLimitObservable", constructorTest(hwconv.NewVoltageLimitObservable, "Float64ObservableGaugeConfig", "hw.voltage.limit", "Voltage limit in Volts.", "V", false, nil)},
		{"hwconv/VoltageNominal", constructorTest(hwconv.NewVoltageNominal, "Float64GaugeConfig", "hw.voltage.nominal", "Nominal (expected) voltage.", "V", false, nil)},
		{"hwconv/VoltageNominalObservable", constructorTest(hwconv.NewVoltageNominalObservable, "Float64ObservableGaugeConfig", "hw.voltage.nominal", "Nominal (expected) voltage.", "V", false, nil)},
		{"k8sconv/ContainerCPULimitCurrent", constructorTest(k8sconv.NewContainerCPULimitCurrent, "Float64UpDownCounterConfig", "k8s.container.cpu.limit.current", "Maximum CPU resource limit currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitCurrentObservable", constructorTest(k8sconv.NewContainerCPULimitCurrentObservable, "Float64ObservableUpDownCounterConfig", "k8s.container.cpu.limit.current", "Maximum CPU resource limit currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitDesired", constructorTest(k8sconv.NewContainerCPULimitDesired, "Float64UpDownCounterConfig", "k8s.container.cpu.limit.desired", "Maximum CPU resource limit as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitDesiredObservable", constructorTest(k8sconv.NewContainerCPULimitDesiredObservable, "Float64ObservableUpDownCounterConfig", "k8s.container.cpu.limit.desired", "Maximum CPU resource limit as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitUtilization", constructorTest(k8sconv.NewContainerCPULimitUtilization, "Float64GaugeConfig", "k8s.container.cpu.limit.utilization", "The ratio of container CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/ContainerCPULimitUtilizationObservable", constructorTest(k8sconv.NewContainerCPULimitUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.container.cpu.limit.utilization", "The ratio of container CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/ContainerCPURequestCurrent", constructorTest(k8sconv.NewContainerCPURequestCurrent, "Float64UpDownCounterConfig", "k8s.container.cpu.request.current", "CPU resource requested currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestCurrentObservable", constructorTest(k8sconv.NewContainerCPURequestCurrentObservable, "Float64ObservableUpDownCounterConfig", "k8s.container.cpu.request.current", "CPU resource requested currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestDesired", constructorTest(k8sconv.NewContainerCPURequestDesired, "Float64UpDownCounterConfig", "k8s.container.cpu.request.desired", "CPU resource requested as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestDesiredObservable", constructorTest(k8sconv.NewContainerCPURequestDesiredObservable, "Float64ObservableUpDownCounterConfig", "k8s.container.cpu.request.desired", "CPU resource requested as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestUtilization", constructorTest(k8sconv.NewContainerCPURequestUtilization, "Float64GaugeConfig", "k8s.container.cpu.request.utilization", "The ratio of container CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/ContainerCPURequestUtilizationObservable", constructorTest(k8sconv.NewContainerCPURequestUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.container.cpu.request.utilization", "The ratio of container CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/ContainerEphemeralStorageLimit", constructorTest(k8sconv.NewContainerEphemeralStorageLimit, "Int64UpDownCounterConfig", "k8s.container.ephemeral_storage.limit", "Maximum ephemeral storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageLimitObservable", constructorTest(k8sconv.NewContainerEphemeralStorageLimitObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.ephemeral_storage.limit", "Maximum ephemeral storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageRequest", constructorTest(k8sconv.NewContainerEphemeralStorageRequest, "Int64UpDownCounterConfig", "k8s.container.ephemeral_storage.request", "Ephemeral storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageRequestObservable", constructorTest(k8sconv.NewContainerEphemeralStorageRequestObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.ephemeral_storage.request", "Ephemeral storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageUsage", constructorTest(k8sconv.NewContainerEphemeralStorageUsage, "Int64UpDownCounterConfig", "k8s.container.ephemeral_storage.usage", "The ephemeral storage used by a container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageUsageObservable", constructorTest(k8sconv.NewContainerEphemeralStorageUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.ephemeral_storage.usage", "The ephemeral storage used by a container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitCurrent", constructorTest(k8sconv.NewContainerMemoryLimitCurrent, "Int64UpDownCounterConfig", "k8s.container.memory.limit.current", "Maximum memory resource limit currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitCurrentObservable", constructorTest(k8sconv.NewContainerMemoryLimitCurrentObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.memory.limit.current", "Maximum memory resource limit currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitDesired", constructorTest(k8sconv.NewContainerMemoryLimitDesired, "Int64UpDownCounterConfig", "k8s.container.memory.limit.desired", "Maximum memory resource limit as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitDesiredObservable", constructorTest(k8sconv.NewContainerMemoryLimitDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.memory.limit.desired", "Maximum memory resource limit as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitUtilization", constructorTest(k8sconv.NewContainerMemoryLimitUtilization, "Float64GaugeConfig", "k8s.container.memory.limit.utilization", "The ratio of container memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/ContainerMemoryLimitUtilizationObservable", constructorTest(k8sconv.NewContainerMemoryLimitUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.container.memory.limit.utilization", "The ratio of container memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/ContainerMemoryRequestCurrent", constructorTest(k8sconv.NewContainerMemoryRequestCurrent, "Int64UpDownCounterConfig", "k8s.container.memory.request.current", "Memory resource request currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestCurrentObservable", constructorTest(k8sconv.NewContainerMemoryRequestCurrentObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.memory.request.current", "Memory resource request currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestDesired", constructorTest(k8sconv.NewContainerMemoryRequestDesired, "Int64UpDownCounterConfig", "k8s.container.memory.request.desired", "Memory resource requested as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestDesiredObservable", constructorTest(k8sconv.NewContainerMemoryRequestDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.memory.request.desired", "Memory resource requested as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestUtilization", constructorTest(k8sconv.NewContainerMemoryRequestUtilization, "Float64GaugeConfig", "k8s.container.memory.request.utilization", "The ratio of container memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/ContainerMemoryRequestUtilizationObservable", constructorTest(k8sconv.NewContainerMemoryRequestUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.container.memory.request.utilization", "The ratio of container memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/ContainerReady", constructorTest(k8sconv.NewContainerReady, "Int64UpDownCounterConfig", "k8s.container.ready", "Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready).", "{container}", false, nil)},
		{"k8sconv/ContainerReadyObservable", constructorTest(k8sconv.NewContainerReadyObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.ready", "Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready).", "{container}", false, nil)},
		{"k8sconv/ContainerRestartCount", constructorTest(k8sconv.NewContainerRestartCount, "Int64UpDownCounterConfig", "k8s.container.restart.count", "Describes how many times the container has restarted (since the last counter reset).", "{restart}", false, nil)},
		{"k8sconv/ContainerRestartCountObservable", constructorTest(k8sconv.NewContainerRestartCountObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.restart.count", "Describes how many times the container has restarted (since the last counter reset).", "{restart}", false, nil)},
		{"k8sconv/ContainerStatusReason", constructorTest(k8sconv.NewContainerStatusReason, "Int64UpDownCounterConfig", "k8s.container.status.reason", "Describes the number of K8s containers that are currently in a state for a given reason.", "{container}", false, nil)},
		{"k8sconv/ContainerStatusReasonObservable", constructorTest(k8sconv.NewContainerStatusReasonObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.status.reason", "Describes the number of K8s containers that are currently in a state for a given reason.", "{container}", false, nil)},
		{"k8sconv/ContainerStatusState", constructorTest(k8sconv.NewContainerStatusState, "Int64UpDownCounterConfig", "k8s.container.status.state", "Describes the number of K8s containers that are currently in a given state.", "{container}", false, nil)},
		{"k8sconv/ContainerStatusStateObservable", constructorTest(k8sconv.NewContainerStatusStateObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.status.state", "Describes the number of K8s containers that are currently in a given state.", "{container}", false, nil)},
		{"k8sconv/ContainerStorageLimit", constructorTest(k8sconv.NewContainerStorageLimit, "Int64UpDownCounterConfig", "k8s.container.storage.limit", "Maximum storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerStorageLimitObservable", constructorTest(k8sconv.NewContainerStorageLimitObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.storage.limit", "Maximum storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerStorageRequest", constructorTest(k8sconv.NewContainerStorageRequest, "Int64UpDownCounterConfig", "k8s.container.storage.request", "Storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/ContainerStorageRequestObservable", constructorTest(k8sconv.NewContainerStorageRequestObservable, "Int64ObservableUpDownCounterConfig", "k8s.container.storage.request", "Storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/CronJobJobActive", constructorTest(k8sconv.NewCronJobJobActive, "Int64UpDownCounterConfig", "k8s.cronjob.job.active", "The number of actively running jobs for a cronjob.", "{job}", false, nil)},
		{"k8sconv/CronJobJobActiveObservable", constructorTest(k8sconv.NewCronJobJobActiveObservable, "Int64ObservableUpDownCounterConfig", "k8s.cronjob.job.active", "The number of actively running jobs for a cronjob.", "{job}", false, nil)},
		{"k8sconv/DaemonSetNodeCurrentScheduled", constructorTest(k8sconv.NewDaemonSetNodeCurrentScheduled, "Int64UpDownCounterConfig", "k8s.daemonset.node.current_scheduled", "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeCurrentScheduledObservable", constructorTest(k8sconv.NewDaemonSetNodeCurrentScheduledObservable, "Int64ObservableUpDownCounterConfig", "k8s.daemonset.node.current_scheduled", "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeDesiredScheduled", constructorTest(k8sconv.NewDaemonSetNodeDesiredScheduled, "Int64UpDownCounterConfig", "k8s.daemonset.node.desired_scheduled", "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod).", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeDesiredScheduledObservable", constructorTest(k8sconv.NewDaemonSetNodeDesiredScheduledObservable, "Int64ObservableUpDownCounterConfig", "k8s.daemonset.node.desired_scheduled", "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod).", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeMisscheduled", constructorTest(k8sconv.NewDaemonSetNodeMisscheduled, "Int64UpDownCounterConfig", "k8s.daemonset.node.misscheduled", "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeMisscheduledObservable", constructorTest(k8sconv.NewDaemonSetNodeMisscheduledObservable, "Int64ObservableUpDownCounterConfig", "k8s.daemonset.node.misscheduled", "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeReady", constructorTest(k8sconv.NewDaemonSetNodeReady, "Int64UpDownCounterConfig", "k8s.daemonset.node.ready", "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeReadyObservable", constructorTest(k8sconv.NewDaemonSetNodeReadyObservable, "Int64ObservableUpDownCounterConfig", "k8s.daemonset.node.ready", "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.", "{node}", false, nil)},
		{"k8sconv/DeploymentPodAvailable", constructorTest(k8sconv.NewDeploymentPodAvailable, "Int64UpDownCounterConfig", "k8s.deployment.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment.", "{pod}", false, nil)},
		{"k8sconv/DeploymentPodAvailableObservable", constructorTest(k8sconv.NewDeploymentPodAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.deployment.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment.", "{pod}", false, nil)},
		{"k8sconv/DeploymentPodDesired", constructorTest(k8sconv.NewDeploymentPodDesired, "Int64UpDownCounterConfig", "k8s.deployment.pod.desired", "Number of desired replica pods in this deployment.", "{pod}", false, nil)},
		{"k8sconv/DeploymentPodDesiredObservable", constructorTest(k8sconv.NewDeploymentPodDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.deployment.pod.desired", "Number of desired replica pods in this deployment.", "{pod}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageUtilization", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageUtilization, "Int64GaugeConfig", "k8s.hpa.metric.target.cpu.average_utilization", "Target average utilization, in percentage, for CPU resource in HPA config.", "1", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageUtilizationObservable", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageUtilizationObservable, "Int64ObservableGaugeConfig", "k8s.hpa.metric.target.cpu.average_utilization", "Target average utilization, in percentage, for CPU resource in HPA config.", "1", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageValue", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageValue, "Int64GaugeConfig", "k8s.hpa.metric.target.cpu.average_value", "Target average value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageValueObservable", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageValueObservable, "Int64ObservableGaugeConfig", "k8s.hpa.metric.target.cpu.average_value", "Target average value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUValue", constructorTest(k8sconv.NewHPAMetricTargetCPUValue, "Int64GaugeConfig", "k8s.hpa.metric.target.cpu.value", "Target value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUValueObservable", constructorTest(k8sconv.NewHPAMetricTargetCPUValueObservable, "Int64ObservableGaugeConfig", "k8s.hpa.metric.target.cpu.value", "Target value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAPodCurrent", constructorTest(k8sconv.NewHPAPodCurrent, "Int64UpDownCounterConfig", "k8s.hpa.pod.current", "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodCurrentObservable", constructorTest(k8sconv.NewHPAPodCurrentObservable, "Int64ObservableUpDownCounterConfig", "k8s.hpa.pod.current", "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodDesired", constructorTest(k8sconv.NewHPAPodDesired, "Int64UpDownCounterConfig", "k8s.hpa.pod.desired", "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodDesiredObservable", constructorTest(k8sconv.NewHPAPodDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.hpa.pod.desired", "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMax", constructorTest(k8sconv.NewHPAPodMax, "Int64UpDownCounterConfig", "k8s.hpa.pod.max", "The upper limit for the number of replica pods to which the autoscaler can scale up.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMaxObservable", constructorTest(k8sconv.NewHPAPodMaxObservable, "Int64ObservableUpDownCounterConfig", "k8s.hpa.pod.max", "The upper limit for the number of replica pods to which the autoscaler can scale up.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMin", constructorTest(k8sconv.NewHPAPodMin, "Int64UpDownCounterConfig", "k8s.hpa.pod.min", "The lower limit for the number of replica pods to which the autoscaler can scale down.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMinObservable", constructorTest(k8sconv.NewHPAPodMinObservable, "Int64ObservableUpDownCounterConfig", "k8s.hpa.pod.min", "The lower limit for the number of replica pods to which the autoscaler can scale down.", "{pod}", false, nil)},
		{"k8sconv/JobPodActive", constructorTest(k8sconv.NewJobPodActive, "Int64UpDownCounterConfig", "k8s.job.pod.active", "The number of pending and actively running pods for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodActiveObservable", constructorTest(k8sconv.NewJobPodActiveObservable, "Int64ObservableUpDownCounterConfig", "k8s.job.pod.active", "The number of pending and actively running pods for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodDesiredSuccessful", constructorTest(k8sconv.NewJobPodDesiredSuccessful, "Int64UpDownCounterConfig", "k8s.job.pod.desired_successful", "The desired number of successfully finished pods the job should be run with.", "{pod}", false, nil)},
		{"k8sconv/JobPodDesiredSuccessfulObservable", constructorTest(k8sconv.NewJobPodDesiredSuccessfulObservable, "Int64ObservableUpDownCounterConfig", "k8s.job.pod.desired_successful", "The desired number of successfully finished pods the job should be run with.", "{pod}", false, nil)},
		{"k8sconv/JobPodFailed", constructorTest(k8sconv.NewJobPodFailed, "Int64UpDownCounterConfig", "k8s.job.pod.failed", "The number of pods which reached phase Failed for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodFailedObservable", constructorTest(k8sconv.NewJobPodFailedObservable, "Int64ObservableUpDownCounterConfig", "k8s.job.pod.failed", "The number of pods which reached phase Failed for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodMaxParallel", constructorTest(k8sconv.NewJobPodMaxParallel, "Int64UpDownCounterConfig", "k8s.job.pod.max_parallel", "The max desired number of pods the job should run at any given time.", "{pod}", false, nil)},
		{"k8sconv/JobPodMaxParallelObservable", constructorTest(k8sconv.NewJobPodMaxParallelObservable, "Int64ObservableUpDownCounterConfig", "k8s.job.pod.max_parallel", "The max desired number of pods the job should run at any given time.", "{pod}", false, nil)},
		{"k8sconv/JobPodSuccessful", constructorTest(k8sconv.NewJobPodSuccessful, "Int64UpDownCounterConfig", "k8s.job.pod.successful", "The number of pods which reached phase Succeeded for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodSuccessfulObservable", constructorTest(k8sconv.NewJobPodSuccessfulObservable, "Int64ObservableUpDownCounterConfig", "k8s.job.pod.successful", "The number of pods which reached phase Succeeded for a job.", "{pod}", false, nil)},
		{"k8sconv/NamespacePhase", constructorTest(k8sconv.NewNamespacePhase, "Int64UpDownCounterConfig", "k8s.namespace.phase", "Describes number of K8s namespaces that are currently in a given phase.", "{namespace}", false, nil)},
		{"k8sconv/NamespacePhaseObservable", constructorTest(k8sconv.NewNamespacePhaseObservable, "Int64ObservableUpDownCounterConfig", "k8s.namespace.phase", "Describes number of K8s namespaces that are currently in a given phase.", "{namespace}", false, nil)},
		{"k8sconv/NodeConditionStatus", constructorTest(k8sconv.NewNodeConditionStatus, "Int64UpDownCounterConfig", "k8s.node.condition.status", "Describes the condition of a particular Node.", "{node}", false, nil)},
		{"k8sconv/NodeConditionStatusObservable", constructorTest(k8sconv.NewNodeConditionStatusObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.condition.status", "Describes the condition of a particular Node.", "{node}", false, nil)},
		{"k8sconv/NodeCPUAllocatable", constructorTest(k8sconv.NewNodeCPUAllocatable, "Float64UpDownCounterConfig", "k8s.node.cpu.allocatable", "Amount of CPU allocatable on the node.", "{cpu}", false, nil)},
		{"k8sconv/NodeCPUAllocatableObservable", constructorTest(k8sconv.NewNodeCPUAllocatableObservable, "Float64ObservableUpDownCounterConfig", "k8s.node.cpu.allocatable", "Amount of CPU allocatable on the node.", "{cpu}", false, nil)},
		{"k8sconv/NodeCPUTime", constructorTest(k8sconv.NewNodeCPUTime, "Float64CounterConfig", "k8s.node.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/NodeCPUTimeObservable", constructorTest(k8sconv.NewNodeCPUTimeObservable, "Float64ObservableCounterConfig", "k8s.node.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/NodeCPUUsage", constructorTest(k8sconv.NewNodeCPUUsage, "Float64GaugeConfig", "k8s.node.cpu.usage", "Node\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeCPUUsageObservable", constructorTest(k8sconv.NewNodeCPUUsageObservable, "Float64ObservableGaugeConfig", "k8s.node.cpu.usage", "Node\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeEphemeralStorageAllocatable", constructorTest(k8sconv.NewNodeEphemeralStorageAllocatable, "Int64UpDownCounterConfig", "k8s.node.ephemeral_storage.allocatable", "Amount of ephemeral-storage allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeEphemeralStorageAllocatableObservable", constructorTest(k8sconv.NewNodeEphemeralStorageAllocatableObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.ephemeral_storage.allocatable", "Amount of ephemeral-storage allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeFilesystemAvailable", constructorTest(k8sconv.NewNodeFilesystemAvailable, "Int64UpDownCounterConfig", "k8s.node.filesystem.available", "Node filesystem available bytes.", "By", false, nil)},
		{"k8sconv/NodeFilesystemAvailableObservable", constructorTest(k8sconv.NewNodeFilesystemAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.filesystem.available", "Node filesystem available bytes.", "By", false, nil)},
		{"k8sconv/NodeFilesystemCapacity", constructorTest(k8sconv.NewNodeFilesystemCapacity, "Int64UpDownCounterConfig", "k8s.node.filesystem.capacity", "Node filesystem capacity.", "By", false, nil)},
		{"k8sconv/NodeFilesystemCapacityObservable", constructorTest(k8sconv.NewNodeFilesystemCapacityObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.filesystem.capacity", "Node filesystem capacity.", "By", false, nil)},
		{"k8sconv/NodeFilesystemInodeCount", constructorTest(k8sconv.NewNodeFilesystemInodeCount, "Int64UpDownCounterConfig", "k8s.node.filesystem.inode.count", "The total inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemInodeCountObservable", constructorTest(k8sconv.NewNodeFilesystemInodeCountObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.filesystem.inode.count", "The total inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemInodeFree", constructorTest(k8sconv.NewNodeFilesystemInodeFree, "Int64UpDownCounterConfig", "k8s.node.filesystem.inode.free", "The free inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemInodeFreeObservable", constructorTest(k8sconv.NewNodeFilesystemInodeFreeObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.filesystem.inode.free", "The free inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemUsage", constructorTest(k8sconv.NewNodeFilesystemUsage, "Int64UpDownCounterConfig", "k8s.node.filesystem.usage", "Node filesystem usage.", "By", false, nil)},
		{"k8sconv/NodeFilesystemUsageObservable", constructorTest(k8sconv.NewNodeFilesystemUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.filesystem.usage", "Node filesystem usage.", "By", false, nil)},
		{"k8sconv/NodeMemoryAllocatable", constructorTest(k8sconv.NewNodeMemoryAllocatable, "Int64UpDownCounterConfig", "k8s.node.memory.allocatable", "Amount of memory allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeMemoryAllocatableObservable", constructorTest(k8sconv.NewNodeMemoryAllocatableObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.memory.allocatable", "Amount of memory allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeMemoryAvailable", constructorTest(k8sconv.NewNodeMemoryAvailable, "Int64UpDownCounterConfig", "k8s.node.memory.available", "Node memory available.", "By", false, nil)},
		{"k8sconv/NodeMemoryAvailableObservable", constructorTest(k8sconv.NewNodeMemoryAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.memory.available", "Node memory available.", "By", false, nil)},
		{"k8sconv/NodeMemoryRss", constructorTest(k8sconv.NewNodeMemoryRss, "Int64UpDownCounterConfig", "k8s.node.memory.rss", "Node memory RSS.", "By", false, nil)},
		{"k8sconv/NodeMemoryRssObservable", constructorTest(k8sconv.NewNodeMemoryRssObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.memory.rss", "Node memory RSS.", "By", false, nil)},
		{"k8sconv/NodeMemoryUsage", constructorTest(k8sconv.NewNodeMemoryUsage, "Int64UpDownCounterConfig", "k8s.node.memory.usage", "Memory usage of the Node.", "By", false, nil)},
		{"k8sconv/NodeMemoryUsageObservable", constructorTest(k8sconv.NewNodeMemoryUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.memory.usage", "Memory usage of the Node.", "By", false, nil)},
		{"k8sconv/NodeMemoryWorkingSet", constructorTest(k8sconv.NewNodeMemoryWorkingSet, "Int64UpDownCounterConfig", "k8s.node.memory.working_set", "Node memory working set.", "By", false, nil)},
		{"k8sconv/NodeMemoryWorkingSetObservable", constructorTest(k8sconv.NewNodeMemoryWorkingSetObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.memory.working_set", "Node memory working set.", "By", false, nil)},
		{"k8sconv/NodeNetworkErrors", constructorTest(k8sconv.NewNodeNetworkErrors, "Int64CounterConfig", "k8s.node.network.errors", "Node network errors.", "{error}", false, nil)},
		{"k8sconv/NodeNetworkErrorsObservable", constructorTest(k8sconv.NewNodeNetworkErrorsObservable, "Int64ObservableCounterConfig", "k8s.node.network.errors", "Node network errors.", "{error}", false, nil)},
		{"k8sconv/NodeNetworkIO", constructorTest(k8sconv.NewNodeNetworkIO, "Int64CounterConfig", "k8s.node.network.io", "Network bytes for the Node.", "By", false, nil)},
		{"k8sconv/NodeNetworkIOObservable", constructorTest(k8sconv.NewNodeNetworkIOObservable, "Int64ObservableCounterConfig", "k8s.node.network.io", "Network bytes for the Node.", "By", false, nil)},
		{"k8sconv/NodePagingFaults", constructorTest(k8sconv.NewNodePagingFaults, "Int64CounterConfig", "k8s.node.paging.faults", "Node memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/NodePagingFaultsObservable", constructorTest(k8sconv.NewNodePagingFaultsObservable, "Int64ObservableCounterConfig", "k8s.node.paging.faults", "Node memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/NodePIDLimit", constructorTest(k8sconv.NewNodePIDLimit, "Int64UpDownCounterConfig", "k8s.node.pid.limit", "Total number of processes/threads allowed by the operating system on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePIDLimitObservable", constructorTest(k8sconv.NewNodePIDLimitObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.pid.limit", "Total number of processes/threads allowed by the operating system on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePIDUsage", constructorTest(k8sconv.NewNodePIDUsage, "Int64UpDownCounterConfig", "k8s.node.pid.usage", "Total number of existing processes and threads on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePIDUsageObservable", constructorTest(k8sconv.NewNodePIDUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.pid.usage", "Total number of existing processes and threads on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePodAllocatable", constructorTest(k8sconv.NewNodePodAllocatable, "Int64UpDownCounterConfig", "k8s.node.pod.allocatable", "Amount of pods allocatable on the node.", "{pod}", false, nil)},
		{"k8sconv/NodePodAllocatableObservable", constructorTest(k8sconv.NewNodePodAllocatableObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.pod.allocatable", "Amount of pods allocatable on the node.", "{pod}", false, nil)},
		{"k8sconv/NodeSystemContainerCPUTime", constructorTest(k8sconv.NewNodeSystemContainerCPUTime, "Float64CounterConfig", "k8s.node.system_container.cpu.time", "Node\u0027s system container CPU time.", "s", false, nil)},
		{"k8sconv/NodeSystemContainerCPUTimeObservable", constructorTest(k8sconv.NewNodeSystemContainerCPUTimeObservable, "Float64ObservableCounterConfig", "k8s.node.system_container.cpu.time", "Node\u0027s system container CPU time.", "s", false, nil)},
		{"k8sconv/NodeSystemContainerCPUUsage", constructorTest(k8sconv.NewNodeSystemContainerCPUUsage, "Float64GaugeConfig", "k8s.node.system_container.cpu.usage", "Node\u0027s system container CPU usage, measured in CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeSystemContainerCPUUsageObservable", constructorTest(k8sconv.NewNodeSystemContainerCPUUsageObservable, "Float64ObservableGaugeConfig", "k8s.node.system_container.cpu.usage", "Node\u0027s system container CPU usage, measured in CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryUsage", constructorTest(k8sconv.NewNodeSystemContainerMemoryUsage, "Int64UpDownCounterConfig", "k8s.node.system_container.memory.usage", "Node\u0027s system container memory usage.", "By", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryUsageObservable", constructorTest(k8sconv.NewNodeSystemContainerMemoryUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.system_container.memory.usage", "Node\u0027s system container memory usage.", "By", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryWorkingSet", constructorTest(k8sconv.NewNodeSystemContainerMemoryWorkingSet, "Int64UpDownCounterConfig", "k8s.node.system_container.memory.working_set", "The amount of working set memory.", "By", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryWorkingSetObservable", constructorTest(k8sconv.NewNodeSystemContainerMemoryWorkingSetObservable, "Int64ObservableUpDownCounterConfig", "k8s.node.system_container.memory.working_set", "The amount of working set memory.", "By", false, nil)},
		{"k8sconv/NodeUptime", constructorTest(k8sconv.NewNodeUptime, "Float64GaugeConfig", "k8s.node.uptime", "The time the Node has been running.", "s", false, nil)},
		{"k8sconv/NodeUptimeObservable", constructorTest(k8sconv.NewNodeUptimeObservable, "Float64ObservableGaugeConfig", "k8s.node.uptime", "The time the Node has been running.", "s", false, nil)},
		{"k8sconv/PersistentvolumeStatusPhase", constructorTest(k8sconv.NewPersistentvolumeStatusPhase, "Int64UpDownCounterConfig", "k8s.persistentvolume.status.phase", "Number of PersistentVolumes in a given phase.", "{persistentvolume}", false, nil)},
		{"k8sconv/PersistentvolumeStatusPhaseObservable", constructorTest(k8sconv.NewPersistentvolumeStatusPhaseObservable, "Int64ObservableUpDownCounterConfig", "k8s.persistentvolume.status.phase", "Number of PersistentVolumes in a given phase.", "{persistentvolume}", false, nil)},
		{"k8sconv/PersistentvolumeStorageCapacity", constructorTest(k8sconv.NewPersistentvolumeStorageCapacity, "Int64UpDownCounterConfig", "k8s.persistentvolume.storage.capacity", "The storage capacity of the PersistentVolume.", "By", false, nil)},
		{"k8sconv/PersistentvolumeStorageCapacityObservable", constructorTest(k8sconv.NewPersistentvolumeStorageCapacityObservable, "Int64ObservableUpDownCounterConfig", "k8s.persistentvolume.storage.capacity", "The storage capacity of the PersistentVolume.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStatusPhase", constructorTest(k8sconv.NewPersistentvolumeclaimStatusPhase, "Int64UpDownCounterConfig", "k8s.persistentvolumeclaim.status.phase", "Number of PersistentVolumeClaims in a given phase.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/PersistentvolumeclaimStatusPhaseObservable", constructorTest(k8sconv.NewPersistentvolumeclaimStatusPhaseObservable, "Int64ObservableUpDownCounterConfig", "k8s.persistentvolumeclaim.status.phase", "Number of PersistentVolumeClaims in a given phase.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageCapacity", constructorTest(k8sconv.NewPersistentvolumeclaimStorageCapacity, "Int64UpDownCounterConfig", "k8s.persistentvolumeclaim.storage.capacity", "The actual storage capacity provisioned for the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageCapacityObservable", constructorTest(k8sconv.NewPersistentvolumeclaimStorageCapacityObservable, "Int64ObservableUpDownCounterConfig", "k8s.persistentvolumeclaim.storage.capacity", "The actual storage capacity provisioned for the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageRequest", constructorTest(k8sconv.NewPersistentvolumeclaimStorageRequest, "Int64UpDownCounterConfig", "k8s.persistentvolumeclaim.storage.request", "The storage requested by the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageRequestObservable", constructorTest(k8sconv.NewPersistentvolumeclaimStorageRequestObservable, "Int64ObservableUpDownCounterConfig", "k8s.persistentvolumeclaim.storage.request", "The storage requested by the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PodCPULimitCurrent", constructorTest(k8sconv.NewPodCPULimitCurrent, "Float64UpDownCounterConfig", "k8s.pod.cpu.limit.current", "Maximum CPU resource limit currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitCurrentObservable", constructorTest(k8sconv.NewPodCPULimitCurrentObservable, "Float64ObservableUpDownCounterConfig", "k8s.pod.cpu.limit.current", "Maximum CPU resource limit currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitDesired", constructorTest(k8sconv.NewPodCPULimitDesired, "Float64UpDownCounterConfig", "k8s.pod.cpu.limit.desired", "Maximum CPU resource limit as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitDesiredObservable", constructorTest(k8sconv.NewPodCPULimitDesiredObservable, "Float64ObservableUpDownCounterConfig", "k8s.pod.cpu.limit.desired", "Maximum CPU resource limit as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitUtilization", constructorTest(k8sconv.NewPodCPULimitUtilization, "Float64GaugeConfig", "k8s.pod.cpu.limit.utilization", "The ratio of Pod CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/PodCPULimitUtilizationObservable", constructorTest(k8sconv.NewPodCPULimitUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.pod.cpu.limit.utilization", "The ratio of Pod CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/PodCPURequestCurrent", constructorTest(k8sconv.NewPodCPURequestCurrent, "Float64UpDownCounterConfig", "k8s.pod.cpu.request.current", "CPU resource requested currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestCurrentObservable", constructorTest(k8sconv.NewPodCPURequestCurrentObservable, "Float64ObservableUpDownCounterConfig", "k8s.pod.cpu.request.current", "CPU resource requested currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestDesired", constructorTest(k8sconv.NewPodCPURequestDesired, "Float64UpDownCounterConfig", "k8s.pod.cpu.request.desired", "CPU resource requested as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestDesiredObservable", constructorTest(k8sconv.NewPodCPURequestDesiredObservable, "Float64ObservableUpDownCounterConfig", "k8s.pod.cpu.request.desired", "CPU resource requested as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestUtilization", constructorTest(k8sconv.NewPodCPURequestUtilization, "Float64GaugeConfig", "k8s.pod.cpu.request.utilization", "The ratio of Pod CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/PodCPURequestUtilizationObservable", constructorTest(k8sconv.NewPodCPURequestUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.pod.cpu.request.utilization", "The ratio of Pod CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/PodCPUTime", constructorTest(k8sconv.NewPodCPUTime, "Float64CounterConfig", "k8s.pod.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/PodCPUTimeObservable", constructorTest(k8sconv.NewPodCPUTimeObservable, "Float64ObservableCounterConfig", "k8s.pod.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/PodCPUUsage", constructorTest(k8sconv.NewPodCPUUsage, "Float64GaugeConfig", "k8s.pod.cpu.usage", "Pod\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/PodCPUUsageObservable", constructorTest(k8sconv.NewPodCPUUsageObservable, "Float64ObservableGaugeConfig", "k8s.pod.cpu.usage", "Pod\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/PodFilesystemAvailable", constructorTest(k8sconv.NewPodFilesystemAvailable, "Int64UpDownCounterConfig", "k8s.pod.filesystem.available", "Pod filesystem available bytes.", "By", false, nil)},
		{"k8sconv/PodFilesystemAvailableObservable", constructorTest(k8sconv.NewPodFilesystemAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.filesystem.available", "Pod filesystem available bytes.", "By", false, nil)},
		{"k8sconv/PodFilesystemCapacity", constructorTest(k8sconv.NewPodFilesystemCapacity, "Int64UpDownCounterConfig", "k8s.pod.filesystem.capacity", "Pod filesystem capacity.", "By", false, nil)},
		{"k8sconv/PodFilesystemCapacityObservable", constructorTest(k8sconv.NewPodFilesystemCapacityObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.filesystem.capacity", "Pod filesystem capacity.", "By", false, nil)},
		{"k8sconv/PodFilesystemUsage", constructorTest(k8sconv.NewPodFilesystemUsage, "Int64UpDownCounterConfig", "k8s.pod.filesystem.usage", "Pod filesystem usage.", "By", false, nil)},
		{"k8sconv/PodFilesystemUsageObservable", constructorTest(k8sconv.NewPodFilesystemUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.filesystem.usage", "Pod filesystem usage.", "By", false, nil)},
		{"k8sconv/PodMemoryAvailable", constructorTest(k8sconv.NewPodMemoryAvailable, "Int64UpDownCounterConfig", "k8s.pod.memory.available", "Pod memory available.", "By", false, nil)},
		{"k8sconv/PodMemoryAvailableObservable", constructorTest(k8sconv.NewPodMemoryAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.available", "Pod memory available.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitCurrent", constructorTest(k8sconv.NewPodMemoryLimitCurrent, "Int64UpDownCounterConfig", "k8s.pod.memory.limit.current", "Maximum memory resource limit currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitCurrentObservable", constructorTest(k8sconv.NewPodMemoryLimitCurrentObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.limit.current", "Maximum memory resource limit currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitDesired", constructorTest(k8sconv.NewPodMemoryLimitDesired, "Int64UpDownCounterConfig", "k8s.pod.memory.limit.desired", "Maximum memory resource limit as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitDesiredObservable", constructorTest(k8sconv.NewPodMemoryLimitDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.limit.desired", "Maximum memory resource limit as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitUtilization", constructorTest(k8sconv.NewPodMemoryLimitUtilization, "Float64GaugeConfig", "k8s.pod.memory.limit.utilization", "The ratio of Pod memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/PodMemoryLimitUtilizationObservable", constructorTest(k8sconv.NewPodMemoryLimitUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.pod.memory.limit.utilization", "The ratio of Pod memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/PodMemoryRequestCurrent", constructorTest(k8sconv.NewPodMemoryRequestCurrent, "Int64UpDownCounterConfig", "k8s.pod.memory.request.current", "Memory resource request currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestCurrentObservable", constructorTest(k8sconv.NewPodMemoryRequestCurrentObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.request.current", "Memory resource request currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestDesired", constructorTest(k8sconv.NewPodMemoryRequestDesired, "Int64UpDownCounterConfig", "k8s.pod.memory.request.desired", "Memory resource requested as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestDesiredObservable", constructorTest(k8sconv.NewPodMemoryRequestDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.request.desired", "Memory resource requested as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestUtilization", constructorTest(k8sconv.NewPodMemoryRequestUtilization, "Float64GaugeConfig", "k8s.pod.memory.request.utilization", "The ratio of Pod memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/PodMemoryRequestUtilizationObservable", constructorTest(k8sconv.NewPodMemoryRequestUtilizationObservable, "Float64ObservableGaugeConfig", "k8s.pod.memory.request.utilization", "The ratio of Pod memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/PodMemoryRss", constructorTest(k8sconv.NewPodMemoryRss, "Int64UpDownCounterConfig", "k8s.pod.memory.rss", "Pod memory RSS.", "By", false, nil)},
		{"k8sconv/PodMemoryRssObservable", constructorTest(k8sconv.NewPodMemoryRssObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.rss", "Pod memory RSS.", "By", false, nil)},
		{"k8sconv/PodMemoryUsage", constructorTest(k8sconv.NewPodMemoryUsage, "Int64UpDownCounterConfig", "k8s.pod.memory.usage", "Memory usage of the Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryUsageObservable", constructorTest(k8sconv.NewPodMemoryUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.usage", "Memory usage of the Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryWorkingSet", constructorTest(k8sconv.NewPodMemoryWorkingSet, "Int64UpDownCounterConfig", "k8s.pod.memory.working_set", "Pod memory working set.", "By", false, nil)},
		{"k8sconv/PodMemoryWorkingSetObservable", constructorTest(k8sconv.NewPodMemoryWorkingSetObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.memory.working_set", "Pod memory working set.", "By", false, nil)},
		{"k8sconv/PodNetworkErrors", constructorTest(k8sconv.NewPodNetworkErrors, "Int64CounterConfig", "k8s.pod.network.errors", "Pod network errors.", "{error}", false, nil)},
		{"k8sconv/PodNetworkErrorsObservable", constructorTest(k8sconv.NewPodNetworkErrorsObservable, "Int64ObservableCounterConfig", "k8s.pod.network.errors", "Pod network errors.", "{error}", false, nil)},
		{"k8sconv/PodNetworkIO", constructorTest(k8sconv.NewPodNetworkIO, "Int64CounterConfig", "k8s.pod.network.io", "Network bytes for the Pod.", "By", false, nil)},
		{"k8sconv/PodNetworkIOObservable", constructorTest(k8sconv.NewPodNetworkIOObservable, "Int64ObservableCounterConfig", "k8s.pod.network.io", "Network bytes for the Pod.", "By", false, nil)},
		{"k8sconv/PodPagingFaults", constructorTest(k8sconv.NewPodPagingFaults, "Int64CounterConfig", "k8s.pod.paging.faults", "Pod memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/PodPagingFaultsObservable", constructorTest(k8sconv.NewPodPagingFaultsObservable, "Int64ObservableCounterConfig", "k8s.pod.paging.faults", "Pod memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/PodStatusPhase", constructorTest(k8sconv.NewPodStatusPhase, "Int64UpDownCounterConfig", "k8s.pod.status.phase", "Describes number of K8s Pods that are currently in a given phase.", "{pod}", false, nil)},
		{"k8sconv/PodStatusPhaseObservable", constructorTest(k8sconv.NewPodStatusPhaseObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.status.phase", "Describes number of K8s Pods that are currently in a given phase.", "{pod}", false, nil)},
		{"k8sconv/PodStatusReason", constructorTest(k8sconv.NewPodStatusReason, "Int64UpDownCounterConfig", "k8s.pod.status.reason", "Describes the number of K8s Pods that are currently in a state for a given reason.", "{pod}", false, nil)},
		{"k8sconv/PodStatusReasonObservable", constructorTest(k8sconv.NewPodStatusReasonObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.status.reason", "Describes the number of K8s Pods that are currently in a state for a given reason.", "{pod}", false, nil)},
		{"k8sconv/PodUptime", constructorTest(k8sconv.NewPodUptime, "Float64GaugeConfig", "k8s.pod.uptime", "The time the Pod has been running.", "s", false, nil)},
		{"k8sconv/PodUptimeObservable", constructorTest(k8sconv.NewPodUptimeObservable, "Float64ObservableGaugeConfig", "k8s.pod.uptime", "The time the Pod has been running.", "s", false, nil)},
		{"k8sconv/PodVolumeAvailable", constructorTest(k8sconv.NewPodVolumeAvailable, "Int64UpDownCounterConfig", "k8s.pod.volume.available", "Pod volume storage space available.", "By", false, nil)},
		{"k8sconv/PodVolumeAvailableObservable", constructorTest(k8sconv.NewPodVolumeAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.volume.available", "Pod volume storage space available.", "By", false, nil)},
		{"k8sconv/PodVolumeCapacity", constructorTest(k8sconv.NewPodVolumeCapacity, "Int64UpDownCounterConfig", "k8s.pod.volume.capacity", "Pod volume total capacity.", "By", false, nil)},
		{"k8sconv/PodVolumeCapacityObservable", constructorTest(k8sconv.NewPodVolumeCapacityObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.volume.capacity", "Pod volume total capacity.", "By", false, nil)},
		{"k8sconv/PodVolumeInodeCount", constructorTest(k8sconv.NewPodVolumeInodeCount, "Int64UpDownCounterConfig", "k8s.pod.volume.inode.count", "The total inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeCountObservable", constructorTest(k8sconv.NewPodVolumeInodeCountObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.volume.inode.count", "The total inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeFree", constructorTest(k8sconv.NewPodVolumeInodeFree, "Int64UpDownCounterConfig", "k8s.pod.volume.inode.free", "The free inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeFreeObservable", constructorTest(k8sconv.NewPodVolumeInodeFreeObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.volume.inode.free", "The free inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeUsed", constructorTest(k8sconv.NewPodVolumeInodeUsed, "Int64UpDownCounterConfig", "k8s.pod.volume.inode.used", "The inodes used by the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeUsedObservable", constructorTest(k8sconv.NewPodVolumeInodeUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.volume.inode.used", "The inodes used by the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeUsage", constructorTest(k8sconv.NewPodVolumeUsage, "Int64UpDownCounterConfig", "k8s.pod.volume.usage", "Pod volume usage.", "By", false, nil)},
		{"k8sconv/PodVolumeUsageObservable", constructorTest(k8sconv.NewPodVolumeUsageObservable, "Int64ObservableUpDownCounterConfig", "k8s.pod.volume.usage", "Pod volume usage.", "By", false, nil)},
		{"k8sconv/ReplicaSetPodAvailable", constructorTest(k8sconv.NewReplicaSetPodAvailable, "Int64UpDownCounterConfig", "k8s.replicaset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicaSetPodAvailableObservable", constructorTest(k8sconv.NewReplicaSetPodAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.replicaset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicaSetPodDesired", constructorTest(k8sconv.NewReplicaSetPodDesired, "Int64UpDownCounterConfig", "k8s.replicaset.pod.desired", "Number of desired replica pods in this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicaSetPodDesiredObservable", constructorTest(k8sconv.NewReplicaSetPodDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.replicaset.pod.desired", "Number of desired replica pods in this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodAvailable", constructorTest(k8sconv.NewReplicationControllerPodAvailable, "Int64UpDownCounterConfig", "k8s.replicationcontroller.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodAvailableObservable", constructorTest(k8sconv.NewReplicationControllerPodAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.replicationcontroller.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodDesired", constructorTest(k8sconv.NewReplicationControllerPodDesired, "Int64UpDownCounterConfig", "k8s.replicationcontroller.pod.desired", "Number of desired replica pods in this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodDesiredObservable", constructorTest(k8sconv.NewReplicationControllerPodDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.replicationcontroller.pod.desired", "Number of desired replica pods in this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitHard", constructorTest(k8sconv.NewResourceQuotaCPULimitHard, "Float64UpDownCounterConfig", "k8s.resourcequota.cpu.limit.hard", "The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitHardObservable", constructorTest(k8sconv.NewResourceQuotaCPULimitHardObservable, "Float64ObservableUpDownCounterConfig", "k8s.resourcequota.cpu.limit.hard", "The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitUsed", constructorTest(k8sconv.NewResourceQuotaCPULimitUsed, "Float64UpDownCounterConfig", "k8s.resourcequota.cpu.limit.used", "The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitUsedObservable", constructorTest(k8sconv.NewResourceQuotaCPULimitUsedObservable, "Float64ObservableUpDownCounterConfig", "k8s.resourcequota.cpu.limit.used", "The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestHard", constructorTest(k8sconv.NewResourceQuotaCPURequestHard, "Float64UpDownCounterConfig", "k8s.resourcequota.cpu.request.hard", "The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestHardObservable", constructorTest(k8sconv.NewResourceQuotaCPURequestHardObservable, "Float64ObservableUpDownCounterConfig", "k8s.resourcequota.cpu.request.hard", "The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestUsed", constructorTest(k8sconv.NewResourceQuotaCPURequestUsed, "Float64UpDownCounterConfig", "k8s.resourcequota.cpu.request.used", "The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaCPURequestUsedObservable, "Float64ObservableUpDownCounterConfig", "k8s.resourcequota.cpu.request.used", "The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitHard", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitHard, "Int64UpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.limit.hard", "The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitHardObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.limit.hard", "The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitUsed", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.limit.used", "The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitUsedObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.limit.used", "The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestHard", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestHard, "Int64UpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.request.hard", "The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.request.hard", "The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestUsed", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.request.used", "The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.ephemeral_storage.request.used", "The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestHard", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestHard, "Int64UpDownCounterConfig", "k8s.resourcequota.hugepage_count.request.hard", "The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.hugepage_count.request.hard", "The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestUsed", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.hugepage_count.request.used", "The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.hugepage_count.request.used", "The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitHard", constructorTest(k8sconv.NewResourceQuotaMemoryLimitHard, "Int64UpDownCounterConfig", "k8s.resourcequota.memory.limit.hard", "The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitHardObservable", constructorTest(k8sconv.NewResourceQuotaMemoryLimitHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.memory.limit.hard", "The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitUsed", constructorTest(k8sconv.NewResourceQuotaMemoryLimitUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.memory.limit.used", "The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitUsedObservable", constructorTest(k8sconv.NewResourceQuotaMemoryLimitUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.memory.limit.used", "The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestHard", constructorTest(k8sconv.NewResourceQuotaMemoryRequestHard, "Int64UpDownCounterConfig", "k8s.resourcequota.memory.request.hard", "The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaMemoryRequestHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.memory.request.hard", "The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestUsed", constructorTest(k8sconv.NewResourceQuotaMemoryRequestUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.memory.request.used", "The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaMemoryRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.memory.request.used", "The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountHard", constructorTest(k8sconv.NewResourceQuotaObjectCountHard, "Int64UpDownCounterConfig", "k8s.resourcequota.object_count.hard", "The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountHardObservable", constructorTest(k8sconv.NewResourceQuotaObjectCountHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.object_count.hard", "The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountUsed", constructorTest(k8sconv.NewResourceQuotaObjectCountUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.object_count.used", "The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountUsedObservable", constructorTest(k8sconv.NewResourceQuotaObjectCountUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.object_count.used", "The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountHard", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountHard, "Int64UpDownCounterConfig", "k8s.resourcequota.persistentvolumeclaim_count.hard", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountHardObservable", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.persistentvolumeclaim_count.hard", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountUsed", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.persistentvolumeclaim_count.used", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountUsedObservable", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.persistentvolumeclaim_count.used", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestHard", constructorTest(k8sconv.NewResourceQuotaStorageRequestHard, "Int64UpDownCounterConfig", "k8s.resourcequota.storage.request.hard", "The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaStorageRequestHardObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.storage.request.hard", "The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestUsed", constructorTest(k8sconv.NewResourceQuotaStorageRequestUsed, "Int64UpDownCounterConfig", "k8s.resourcequota.storage.request.used", "The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaStorageRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "k8s.resourcequota.storage.request.used", "The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ServiceEndpointCount", constructorTest(k8sconv.NewServiceEndpointCount, "Int64GaugeConfig", "k8s.service.endpoint.count", "Number of endpoints for a service by condition and address type.", "{endpoint}", false, nil)},
		{"k8sconv/ServiceEndpointCountObservable", constructorTest(k8sconv.NewServiceEndpointCountObservable, "Int64ObservableGaugeConfig", "k8s.service.endpoint.count", "Number of endpoints for a service by condition and address type.", "{endpoint}", false, nil)},
		{"k8sconv/ServiceLoadBalancerIngressCount", constructorTest(k8sconv.NewServiceLoadBalancerIngressCount, "Int64GaugeConfig", "k8s.service.load_balancer.ingress.count", "Number of load balancer ingress points (external IPs/hostnames) assigned to the service.", "{ingress}", false, nil)},
		{"k8sconv/ServiceLoadBalancerIngressCountObservable", constructorTest(k8sconv.NewServiceLoadBalancerIngressCountObservable, "Int64ObservableGaugeConfig", "k8s.service.load_balancer.ingress.count", "Number of load balancer ingress points (external IPs/hostnames) assigned to the service.", "{ingress}", false, nil)},
		{"k8sconv/StatefulSetPodAvailable", constructorTest(k8sconv.NewStatefulSetPodAvailable, "Int64UpDownCounterConfig", "k8s.statefulset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodAvailableObservable", constructorTest(k8sconv.NewStatefulSetPodAvailableObservable, "Int64ObservableUpDownCounterConfig", "k8s.statefulset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodCurrent", constructorTest(k8sconv.NewStatefulSetPodCurrent, "Int64UpDownCounterConfig", "k8s.statefulset.pod.current", "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodCurrentObservable", constructorTest(k8sconv.NewStatefulSetPodCurrentObservable, "Int64ObservableUpDownCounterConfig", "k8s.statefulset.pod.current", "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodDesired", constructorTest(k8sconv.NewStatefulSetPodDesired, "Int64UpDownCounterConfig", "k8s.statefulset.pod.desired", "Number of desired replica pods in this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodDesiredObservable", constructorTest(k8sconv.NewStatefulSetPodDesiredObservable, "Int64ObservableUpDownCounterConfig", "k8s.statefulset.pod.desired", "Number of desired replica pods in this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodReady", constructorTest(k8sconv.NewStatefulSetPodReady, "Int64UpDownCounterConfig", "k8s.statefulset.pod.ready", "The number of replica pods created for this statefulset with a Ready Condition.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodReadyObservable", constructorTest(k8sconv.NewStatefulSetPodReadyObservable, "Int64ObservableUpDownCounterConfig", "k8s.statefulset.pod.ready", "The number of replica pods created for this statefulset with a Ready Condition.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodUpdated", constructorTest(k8sconv.NewStatefulSetPodUpdated, "Int64UpDownCounterConfig", "k8s.statefulset.pod.updated", "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodUpdatedObservable", constructorTest(k8sconv.NewStatefulSetPodUpdatedObservable, "Int64ObservableUpDownCounterConfig", "k8s.statefulset.pod.updated", "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision.", "{pod}", false, nil)},
		{"messagingconv/ClientConsumedMessages", constructorTest(messagingconv.NewClientConsumedMessages, "Int64CounterConfig", "messaging.client.consumed.messages", "Number of messages that were delivered to the application.", "{message}", false, nil)},
		{"messagingconv/ClientConsumedMessagesObservable", constructorTest(messagingconv.NewClientConsumedMessagesObservable, "Int64ObservableCounterConfig", "messaging.client.consumed.messages", "Number of messages that were delivered to the application.", "{message}", false, nil)},
		{"messagingconv/ClientOperationDuration", constructorTest(messagingconv.NewClientOperationDuration, "Float64HistogramConfig", "messaging.client.operation.duration", "Duration of messaging operation initiated by a producer or consumer client.", "s", true, nil)},
		{"messagingconv/ClientSentMessages", constructorTest(messagingconv.NewClientSentMessages, "Int64CounterConfig", "messaging.client.sent.messages", "Number of messages producer attempted to send to the broker.", "{message}", false, nil)},
		{"messagingconv/ClientSentMessagesObservable", constructorTest(messagingconv.NewClientSentMessagesObservable, "Int64ObservableCounterConfig", "messaging.client.sent.messages", "Number of messages producer attempted to send to the broker.", "{message}", false, nil)},
		{"messagingconv/ProcessDuration", constructorTest(messagingconv.NewProcessDuration, "Float64HistogramConfig", "messaging.process.duration", "Duration of processing operation.", "s", true, nil)},
		{"nfsconv/ClientNetCount", constructorTest(nfsconv.NewClientNetCount, "Int64CounterConfig", "nfs.client.net.count", "Reports the count of kernel NFS client TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ClientNetCountObservable", constructorTest(nfsconv.NewClientNetCountObservable, "Int64ObservableCounterConfig", "nfs.client.net.count", "Reports the count of kernel NFS client TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ClientNetTCPConnectionAccepted", constructorTest(nfsconv.NewClientNetTCPConnectionAccepted, "Int64CounterConfig", "nfs.client.net.tcp.connection.accepted", "Reports the count of kernel NFS client TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ClientNetTCPConnectionAcceptedObservable", constructorTest(nfsconv.NewClientNetTCPConnectionAcceptedObservable, "Int64ObservableCounterConfig", "nfs.client.net.tcp.connection.accepted", "Reports the count of kernel NFS client TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ClientOperationCount", constructorTest(nfsconv.NewClientOperationCount, "Int64CounterConfig", "nfs.client.operation.count", "Reports the count of kernel NFSv4+ client operations.", "{operation}", false, nil)},
		{"nfsconv/ClientOperationCountObservable", constructorTest(nfsconv.NewClientOperationCountObservable, "Int64ObservableCounterConfig", "nfs.client.operation.count", "Reports the count of kernel NFSv4+ client operations.", "{operation}", false, nil)},
		{"nfsconv/ClientProcedureCount", constructorTest(nfsconv.NewClientProcedureCount, "Int64CounterConfig", "nfs.client.procedure.count", "Reports the count of kernel NFS client procedures.", "{procedure}", false, nil)},
		{"nfsconv/ClientProcedureCountObservable", constructorTest(nfsconv.NewClientProcedureCountObservable, "Int64ObservableCounterConfig", "nfs.client.procedure.count", "Reports the count of kernel NFS client procedures.", "{procedure}", false, nil)},
		{"nfsconv/ClientRPCAuthrefreshCount", constructorTest(nfsconv.NewClientRPCAuthrefreshCount, "Int64CounterConfig", "nfs.client.rpc.authrefresh.count", "Reports the count of kernel NFS client RPC authentication refreshes.", "{authrefresh}", false, nil)},
		{"nfsconv/ClientRPCAuthrefreshCountObservable", constructorTest(nfsconv.NewClientRPCAuthrefreshCountObservable, "Int64ObservableCounterConfig", "nfs.client.rpc.authrefresh.count", "Reports the count of kernel NFS client RPC authentication refreshes.", "{authrefresh}", false, nil)},
		{"nfsconv/ClientRPCCount", constructorTest(nfsconv.NewClientRPCCount, "Int64CounterConfig", "nfs.client.rpc.count", "Reports the count of kernel NFS client RPCs sent, regardless of whether they\u0027re accepted/rejected by the server.", "{request}", false, nil)},
		{"nfsconv/ClientRPCCountObservable", constructorTest(nfsconv.NewClientRPCCountObservable, "Int64ObservableCounterConfig", "nfs.client.rpc.count", "Reports the count of kernel NFS client RPCs sent, regardless of whether they\u0027re accepted/rejected by the server.", "{request}", false, nil)},
		{"nfsconv/ClientRPCRetransmitCount", constructorTest(nfsconv.NewClientRPCRetransmitCount, "Int64CounterConfig", "nfs.client.rpc.retransmit.count", "Reports the count of kernel NFS client RPC retransmits.", "{retransmit}", false, nil)},
		{"nfsconv/ClientRPCRetransmitCountObservable", constructorTest(nfsconv.NewClientRPCRetransmitCountObservable, "Int64ObservableCounterConfig", "nfs.client.rpc.retransmit.count", "Reports the count of kernel NFS client RPC retransmits.", "{retransmit}", false, nil)},
		{"nfsconv/ServerFhStaleCount", constructorTest(nfsconv.NewServerFhStaleCount, "Int64CounterConfig", "nfs.server.fh.stale.count", "Reports the count of kernel NFS server stale file handles.", "{fh}", false, nil)},
		{"nfsconv/ServerFhStaleCountObservable", constructorTest(nfsconv.NewServerFhStaleCountObservable, "Int64ObservableCounterConfig", "nfs.server.fh.stale.count", "Reports the count of kernel NFS server stale file handles.", "{fh}", false, nil)},
		{"nfsconv/ServerIO", constructorTest(nfsconv.NewServerIO, "Int64CounterConfig", "nfs.server.io", "Reports the count of kernel NFS server bytes returned to receive and transmit (read and write) requests.", "By", false, nil)},
		{"nfsconv/ServerIOObservable", constructorTest(nfsconv.NewServerIOObservable, "Int64ObservableCounterConfig", "nfs.server.io", "Reports the count of kernel NFS server bytes returned to receive and transmit (read and write) requests.", "By", false, nil)},
		{"nfsconv/ServerNetCount", constructorTest(nfsconv.NewServerNetCount, "Int64CounterConfig", "nfs.server.net.count", "Reports the count of kernel NFS server TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ServerNetCountObservable", constructorTest(nfsconv.NewServerNetCountObservable, "Int64ObservableCounterConfig", "nfs.server.net.count", "Reports the count of kernel NFS server TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ServerNetTCPConnectionAccepted", constructorTest(nfsconv.NewServerNetTCPConnectionAccepted, "Int64CounterConfig", "nfs.server.net.tcp.connection.accepted", "Reports the count of kernel NFS server TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ServerNetTCPConnectionAcceptedObservable", constructorTest(nfsconv.NewServerNetTCPConnectionAcceptedObservable, "Int64ObservableCounterConfig", "nfs.server.net.tcp.connection.accepted", "Reports the count of kernel NFS server TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ServerOperationCount", constructorTest(nfsconv.NewServerOperationCount, "Int64CounterConfig", "nfs.server.operation.count", "Reports the count of kernel NFSv4+ server operations.", "{operation}", false, nil)},
		{"nfsconv/ServerOperationCountObservable", constructorTest(nfsconv.NewServerOperationCountObservable, "Int64ObservableCounterConfig", "nfs.server.operation.count", "Reports the count of kernel NFSv4+ server operations.", "{operation}", false, nil)},
		{"nfsconv/ServerProcedureCount", constructorTest(nfsconv.NewServerProcedureCount, "Int64CounterConfig", "nfs.server.procedure.count", "Reports the count of kernel NFS server procedures.", "{procedure}", false, nil)},
		{"nfsconv/ServerProcedureCountObservable", constructorTest(nfsconv.NewServerProcedureCountObservable, "Int64ObservableCounterConfig", "nfs.server.procedure.count", "Reports the count of kernel NFS server procedures.", "{procedure}", false, nil)},
		{"nfsconv/ServerRepcacheRequests", constructorTest(nfsconv.NewServerRepcacheRequests, "Int64CounterConfig", "nfs.server.repcache.requests", "Reports the kernel NFS server reply cache request count by cache hit status.", "{request}", false, nil)},
		{"nfsconv/ServerRepcacheRequestsObservable", constructorTest(nfsconv.NewServerRepcacheRequestsObservable, "Int64ObservableCounterConfig", "nfs.server.repcache.requests", "Reports the kernel NFS server reply cache request count by cache hit status.", "{request}", false, nil)},
		{"nfsconv/ServerRPCCount", constructorTest(nfsconv.NewServerRPCCount, "Int64CounterConfig", "nfs.server.rpc.count", "Reports the count of kernel NFS server RPCs handled.", "{request}", false, nil)},
		{"nfsconv/ServerRPCCountObservable", constructorTest(nfsconv.NewServerRPCCountObservable, "Int64ObservableCounterConfig", "nfs.server.rpc.count", "Reports the count of kernel NFS server RPCs handled.", "{request}", false, nil)},
		{"nfsconv/ServerThreadCount", constructorTest(nfsconv.NewServerThreadCount, "Int64UpDownCounterConfig", "nfs.server.thread.count", "Reports the count of kernel NFS server available threads.", "{thread}", false, nil)},
		{"nfsconv/ServerThreadCountObservable", constructorTest(nfsconv.NewServerThreadCountObservable, "Int64ObservableUpDownCounterConfig", "nfs.server.thread.count", "Reports the count of kernel NFS server available threads.", "{thread}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitHard", constructorTest(openshiftconv.NewClusterquotaCPULimitHard, "Float64UpDownCounterConfig", "openshift.clusterquota.cpu.limit.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitHardObservable", constructorTest(openshiftconv.NewClusterquotaCPULimitHardObservable, "Float64ObservableUpDownCounterConfig", "openshift.clusterquota.cpu.limit.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitUsed", constructorTest(openshiftconv.NewClusterquotaCPULimitUsed, "Float64UpDownCounterConfig", "openshift.clusterquota.cpu.limit.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitUsedObservable", constructorTest(openshiftconv.NewClusterquotaCPULimitUsedObservable, "Float64ObservableUpDownCounterConfig", "openshift.clusterquota.cpu.limit.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestHard", constructorTest(openshiftconv.NewClusterquotaCPURequestHard, "Float64UpDownCounterConfig", "openshift.clusterquota.cpu.request.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestHardObservable", constructorTest(openshiftconv.NewClusterquotaCPURequestHardObservable, "Float64ObservableUpDownCounterConfig", "openshift.clusterquota.cpu.request.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestUsed", constructorTest(openshiftconv.NewClusterquotaCPURequestUsed, "Float64UpDownCounterConfig", "openshift.clusterquota.cpu.request.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaCPURequestUsedObservable, "Float64ObservableUpDownCounterConfig", "openshift.clusterquota.cpu.request.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitHard", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitHard, "Int64UpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitHardObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitUsed", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitUsedObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestHard", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestHard, "Int64UpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestUsed", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.ephemeral_storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestHard", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestHard, "Int64UpDownCounterConfig", "openshift.clusterquota.hugepage_count.request.hard", "The enforced hard limit of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.hugepage_count.request.hard", "The enforced hard limit of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestUsed", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.hugepage_count.request.used", "The current observed total usage of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.hugepage_count.request.used", "The current observed total usage of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitHard", constructorTest(openshiftconv.NewClusterquotaMemoryLimitHard, "Int64UpDownCounterConfig", "openshift.clusterquota.memory.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitHardObservable", constructorTest(openshiftconv.NewClusterquotaMemoryLimitHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.memory.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitUsed", constructorTest(openshiftconv.NewClusterquotaMemoryLimitUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.memory.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitUsedObservable", constructorTest(openshiftconv.NewClusterquotaMemoryLimitUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.memory.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestHard", constructorTest(openshiftconv.NewClusterquotaMemoryRequestHard, "Int64UpDownCounterConfig", "openshift.clusterquota.memory.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaMemoryRequestHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.memory.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestUsed", constructorTest(openshiftconv.NewClusterquotaMemoryRequestUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.memory.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaMemoryRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.memory.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountHard", constructorTest(openshiftconv.NewClusterquotaObjectCountHard, "Int64UpDownCounterConfig", "openshift.clusterquota.object_count.hard", "The enforced hard limit of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountHardObservable", constructorTest(openshiftconv.NewClusterquotaObjectCountHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.object_count.hard", "The enforced hard limit of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountUsed", constructorTest(openshiftconv.NewClusterquotaObjectCountUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.object_count.used", "The current observed total usage of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountUsedObservable", constructorTest(openshiftconv.NewClusterquotaObjectCountUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.object_count.used", "The current observed total usage of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountHard", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountHard, "Int64UpDownCounterConfig", "openshift.clusterquota.persistentvolumeclaim_count.hard", "The enforced hard limit of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountHardObservable", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.persistentvolumeclaim_count.hard", "The enforced hard limit of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountUsed", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.persistentvolumeclaim_count.used", "The current observed total usage of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountUsedObservable", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.persistentvolumeclaim_count.used", "The current observed total usage of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestHard", constructorTest(openshiftconv.NewClusterquotaStorageRequestHard, "Int64UpDownCounterConfig", "openshift.clusterquota.storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaStorageRequestHardObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestUsed", constructorTest(openshiftconv.NewClusterquotaStorageRequestUsed, "Int64UpDownCounterConfig", "openshift.clusterquota.storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaStorageRequestUsedObservable, "Int64ObservableUpDownCounterConfig", "openshift.clusterquota.storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"otelconv/SDKExporterLogExported", constructorTest(otelconv.NewSDKExporterLogExported, "Int64CounterConfig", "otel.sdk.exporter.log.exported", "The number of log records for which the export has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKExporterLogExportedObservable", constructorTest(otelconv.NewSDKExporterLogExportedObservable, "Int64ObservableCounterConfig", "otel.sdk.exporter.log.exported", "The number of log records for which the export has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKExporterLogInflight", constructorTest(otelconv.NewSDKExporterLogInflight, "Int64UpDownCounterConfig", "otel.sdk.exporter.log.inflight", "The number of log records which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{log_record}", false, nil)},
		{"otelconv/SDKExporterLogInflightObservable", constructorTest(otelconv.NewSDKExporterLogInflightObservable, "Int64ObservableUpDownCounterConfig", "otel.sdk.exporter.log.inflight", "The number of log records which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{log_record}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointExported", constructorTest(otelconv.NewSDKExporterMetricDataPointExported, "Int64CounterConfig", "otel.sdk.exporter.metric_data_point.exported", "The number of metric data points for which the export has finished, either successful or failed.", "{data_point}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointExportedObservable", constructorTest(otelconv.NewSDKExporterMetricDataPointExportedObservable, "Int64ObservableCounterConfig", "otel.sdk.exporter.metric_data_point.exported", "The number of metric data points for which the export has finished, either successful or failed.", "{data_point}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointInflight", constructorTest(otelconv.NewSDKExporterMetricDataPointInflight, "Int64UpDownCounterConfig", "otel.sdk.exporter.metric_data_point.inflight", "The number of metric data points which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{data_point}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointInflightObservable", constructorTest(otelconv.NewSDKExporterMetricDataPointInflightObservable, "Int64ObservableUpDownCounterConfig", "otel.sdk.exporter.metric_data_point.inflight", "The number of metric data points which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{data_point}", false, nil)},
		{"otelconv/SDKExporterOperationDuration", constructorTest(otelconv.NewSDKExporterOperationDuration, "Float64HistogramConfig", "otel.sdk.exporter.operation.duration", "The duration of exporting a batch of telemetry records.", "s", true, nil)},
		{"otelconv/SDKExporterSpanExported", constructorTest(otelconv.NewSDKExporterSpanExported, "Int64CounterConfig", "otel.sdk.exporter.span.exported", "The number of spans for which the export has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKExporterSpanExportedObservable", constructorTest(otelconv.NewSDKExporterSpanExportedObservable, "Int64ObservableCounterConfig", "otel.sdk.exporter.span.exported", "The number of spans for which the export has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKExporterSpanInflight", constructorTest(otelconv.NewSDKExporterSpanInflight, "Int64UpDownCounterConfig", "otel.sdk.exporter.span.inflight", "The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{span}", false, nil)},
		{"otelconv/SDKExporterSpanInflightObservable", constructorTest(otelconv.NewSDKExporterSpanInflightObservable, "Int64ObservableUpDownCounterConfig", "otel.sdk.exporter.span.inflight", "The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{span}", false, nil)},
		{"otelconv/SDKLogCreated", constructorTest(otelconv.NewSDKLogCreated, "Int64CounterConfig", "otel.sdk.log.created", "The number of log records submitted to an enabled `Logger`.", "{log_record}", false, nil)},
		{"otelconv/SDKLogCreatedObservable", constructorTest(otelconv.NewSDKLogCreatedObservable, "Int64ObservableCounterConfig", "otel.sdk.log.created", "The number of log records submitted to an enabled `Logger`.", "{log_record}", false, nil)},
		{"otelconv/SDKMetricReaderCollectionDuration", constructorTest(otelconv.NewSDKMetricReaderCollectionDuration, "Float64HistogramConfig", "otel.sdk.metric_reader.collection.duration", "The duration of the collect operation of the metric reader.", "s", true, nil)},
		{"otelconv/SDKProcessorLogProcessed", constructorTest(otelconv.NewSDKProcessorLogProcessed, "Int64CounterConfig", "otel.sdk.processor.log.processed", "The number of log records for which the processing has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorLogProcessedObservable", constructorTest(otelconv.NewSDKProcessorLogProcessedObservable, "Int64ObservableCounterConfig", "otel.sdk.processor.log.processed", "The number of log records for which the processing has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorLogQueueCapacity", constructorTest(otelconv.NewSDKProcessorLogQueueCapacity, "Int64ObservableUpDownCounterConfig", "otel.sdk.processor.log.queue.capacity", "The maximum number of log records the queue of a given instance of an SDK Log Record processor can hold.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorLogQueueSize", constructorTest(otelconv.NewSDKProcessorLogQueueSize, "Int64ObservableUpDownCounterConfig", "otel.sdk.processor.log.queue.size", "The number of log records in the queue of a given instance of an SDK log processor.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorSpanProcessed", constructorTest(otelconv.NewSDKProcessorSpanProcessed, "Int64CounterConfig", "otel.sdk.processor.span.processed", "The number of spans for which the processing has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKProcessorSpanProcessedObservable", constructorTest(otelconv.NewSDKProcessorSpanProcessedObservable, "Int64ObservableCounterConfig", "otel.sdk.processor.span.processed", "The number of spans for which the processing has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKProcessorSpanQueueCapacity", constructorTest(otelconv.NewSDKProcessorSpanQueueCapacity, "Int64ObservableUpDownCounterConfig", "otel.sdk.processor.span.queue.capacity", "The maximum number of spans the queue of a given instance of an SDK span processor can hold.", "{span}", false, nil)},
		{"otelconv/SDKProcessorSpanQueueSize", constructorTest(otelconv.NewSDKProcessorSpanQueueSize, "Int64ObservableUpDownCounterConfig", "otel.sdk.processor.span.queue.size", "The number of spans in the queue of a given instance of an SDK span processor.", "{span}", false, nil)},
		{"otelconv/SDKSpanLive", constructorTest(otelconv.NewSDKSpanLive, "Int64UpDownCounterConfig", "otel.sdk.span.live", "The number of created spans with `recording=true` for which the end operation has not been called yet.", "{span}", false, nil)},
		{"otelconv/SDKSpanLiveObservable", constructorTest(otelconv.NewSDKSpanLiveObservable, "Int64ObservableUpDownCounterConfig", "otel.sdk.span.live", "The number of created spans with `recording=true` for which the end operation has not been called yet.", "{span}", false, nil)},
		{"otelconv/SDKSpanStarted", constructorTest(otelconv.NewSDKSpanStarted, "Int64CounterConfig", "otel.sdk.span.started", "The number of created spans.", "{span}", false, nil)},
		{"otelconv/SDKSpanStartedObservable", constructorTest(otelconv.NewSDKSpanStartedObservable, "Int64ObservableCounterConfig", "otel.sdk.span.started", "The number of created spans.", "{span}", false, nil)},
		{"processconv/ContextSwitches", constructorTest(processconv.NewContextSwitches, "Int64CounterConfig", "process.context_switches", "Number of times the process has been context switched.", "{context_switch}", false, nil)},
		{"processconv/ContextSwitchesObservable", constructorTest(processconv.NewContextSwitchesObservable, "Int64ObservableCounterConfig", "process.context_switches", "Number of times the process has been context switched.", "{context_switch}", false, nil)},
		{"processconv/CPUTime", constructorTest(processconv.NewCPUTime, "Float64ObservableCounterConfig", "process.cpu.time", "Total CPU seconds broken down by different CPU modes.", "s", false, nil)},
		{"processconv/CPUUtilization", constructorTest(processconv.NewCPUUtilization, "Float64GaugeConfig", "process.cpu.utilization", "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process.", "1", false, nil)},
		{"processconv/CPUUtilizationObservable", constructorTest(processconv.NewCPUUtilizationObservable, "Float64ObservableGaugeConfig", "process.cpu.utilization", "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process.", "1", false, nil)},
		{"processconv/DiskIO", constructorTest(processconv.NewDiskIO, "Int64CounterConfig", "process.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"processconv/DiskIOObservable", constructorTest(processconv.NewDiskIOObservable, "Int64ObservableCounterConfig", "process.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"processconv/DiskOperations", constructorTest(processconv.NewDiskOperations, "Int64CounterConfig", "process.disk.operations", "Number of disk operations performed by the process.", "{operation}", false, nil)},
		{"processconv/DiskOperationsObservable", constructorTest(processconv.NewDiskOperationsObservable, "Int64ObservableCounterConfig", "process.disk.operations", "Number of disk operations performed by the process.", "{operation}", false, nil)},
		{"processconv/MemoryUsage", constructorTest(processconv.NewMemoryUsage, "Int64UpDownCounterConfig", "process.memory.usage", "The amount of physical memory in use.", "By", false, nil)},
		{"processconv/MemoryUsageObservable", constructorTest(processconv.NewMemoryUsageObservable, "Int64ObservableUpDownCounterConfig", "process.memory.usage", "The amount of physical memory in use.", "By", false, nil)},
		{"processconv/MemoryUtilization", constructorTest(processconv.NewMemoryUtilization, "Float64GaugeConfig", "process.memory.utilization", "Percentage of total physical memory that is used by the process.", "1", false, nil)},
		{"processconv/MemoryUtilizationObservable", constructorTest(processconv.NewMemoryUtilizationObservable, "Float64ObservableGaugeConfig", "process.memory.utilization", "Percentage of total physical memory that is used by the process.", "1", false, nil)},
		{"processconv/MemoryVirtual", constructorTest(processconv.NewMemoryVirtual, "Int64UpDownCounterConfig", "process.memory.virtual", "The amount of committed virtual memory.", "By", false, nil)},
		{"processconv/MemoryVirtualObservable", constructorTest(processconv.NewMemoryVirtualObservable, "Int64ObservableUpDownCounterConfig", "process.memory.virtual", "The amount of committed virtual memory.", "By", false, nil)},
		{"processconv/NetworkIO", constructorTest(processconv.NewNetworkIO, "Int64CounterConfig", "process.network.io", "Network bytes transferred.", "By", false, nil)},
		{"processconv/NetworkIOObservable", constructorTest(processconv.NewNetworkIOObservable, "Int64ObservableCounterConfig", "process.network.io", "Network bytes transferred.", "By", false, nil)},
		{"processconv/PagingFaults", constructorTest(processconv.NewPagingFaults, "Int64CounterConfig", "process.paging.faults", "Number of page faults the process has made.", "{fault}", false, nil)},
		{"processconv/PagingFaultsObservable", constructorTest(processconv.NewPagingFaultsObservable, "Int64ObservableCounterConfig", "process.paging.faults", "Number of page faults the process has made.", "{fault}", false, nil)},
		{"processconv/SignalsPending", constructorTest(processconv.NewSignalsPending, "Int64UpDownCounterConfig", "process.signals_pending", "Number of pending signals for the process.", "{signal}", false, nil)},
		{"processconv/SignalsPendingObservable", constructorTest(processconv.NewSignalsPendingObservable, "Int64ObservableUpDownCounterConfig", "process.signals_pending", "Number of pending signals for the process.", "{signal}", false, nil)},
		{"processconv/ThreadCount", constructorTest(processconv.NewThreadCount, "Int64UpDownCounterConfig", "process.thread.count", "Process threads count.", "{thread}", false, nil)},
		{"processconv/ThreadCountObservable", constructorTest(processconv.NewThreadCountObservable, "Int64ObservableUpDownCounterConfig", "process.thread.count", "Process threads count.", "{thread}", false, nil)},
		{"processconv/UnixFileDescriptorCount", constructorTest(processconv.NewUnixFileDescriptorCount, "Int64UpDownCounterConfig", "process.unix.file_descriptor.count", "Number of UNIX file descriptors in use by the process.", "{file_descriptor}", false, nil)},
		{"processconv/UnixFileDescriptorCountObservable", constructorTest(processconv.NewUnixFileDescriptorCountObservable, "Int64ObservableUpDownCounterConfig", "process.unix.file_descriptor.count", "Number of UNIX file descriptors in use by the process.", "{file_descriptor}", false, nil)},
		{"processconv/Uptime", constructorTest(processconv.NewUptime, "Float64GaugeConfig", "process.uptime", "The time the process has been running.", "s", false, nil)},
		{"processconv/UptimeObservable", constructorTest(processconv.NewUptimeObservable, "Float64ObservableGaugeConfig", "process.uptime", "The time the process has been running.", "s", false, nil)},
		{"processconv/WindowsHandleCount", constructorTest(processconv.NewWindowsHandleCount, "Int64UpDownCounterConfig", "process.windows.handle.count", "Number of handles held by the process.", "{handle}", false, nil)},
		{"processconv/WindowsHandleCountObservable", constructorTest(processconv.NewWindowsHandleCountObservable, "Int64ObservableUpDownCounterConfig", "process.windows.handle.count", "Number of handles held by the process.", "{handle}", false, nil)},
		{"rpcconv/ClientCallDuration", constructorTest(rpcconv.NewClientCallDuration, "Float64HistogramConfig", "rpc.client.call.duration", "Measures the duration of an outgoing Remote Procedure Call (RPC).", "s", true, []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 60})},
		{"rpcconv/ServerCallDuration", constructorTest(rpcconv.NewServerCallDuration, "Float64HistogramConfig", "rpc.server.call.duration", "Measures the duration of an incoming Remote Procedure Call (RPC).", "s", true, []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 60})},
		{"signalrconv/ServerActiveConnections", constructorTest(signalrconv.NewServerActiveConnections, "Int64UpDownCounterConfig", "signalr.server.active_connections", "Number of connections that are currently active on the server.", "{connection}", false, nil)},
		{"signalrconv/ServerActiveConnectionsObservable", constructorTest(signalrconv.NewServerActiveConnectionsObservable, "Int64ObservableUpDownCounterConfig", "signalr.server.active_connections", "Number of connections that are currently active on the server.", "{connection}", false, nil)},
		{"signalrconv/ServerConnectionDuration", constructorTest(signalrconv.NewServerConnectionDuration, "Float64HistogramConfig", "signalr.server.connection.duration", "The duration of connections on the server.", "s", true, nil)},
		{"systemconv/CPUFrequency", constructorTest(systemconv.NewCPUFrequency, "Int64GaugeConfig", "system.cpu.frequency", "Operating frequency of the logical CPU in Hertz.", "Hz", false, nil)},
		{"systemconv/CPUFrequencyObservable", constructorTest(systemconv.NewCPUFrequencyObservable, "Int64ObservableGaugeConfig", "system.cpu.frequency", "Operating frequency of the logical CPU in Hertz.", "Hz", false, nil)},
		{"systemconv/CPULogicalCount", constructorTest(systemconv.NewCPULogicalCount, "Int64UpDownCounterConfig", "system.cpu.logical.count", "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking.", "{cpu}", false, nil)},
		{"systemconv/CPULogicalCountObservable", constructorTest(systemconv.NewCPULogicalCountObservable, "Int64ObservableUpDownCounterConfig", "system.cpu.logical.count", "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking.", "{cpu}", false, nil)},
		{"systemconv/CPUPhysicalCount", constructorTest(systemconv.NewCPUPhysicalCount, "Int64UpDownCounterConfig", "system.cpu.physical.count", "Reports the number of actual physical processor cores on the hardware.", "{cpu}", false, nil)},
		{"systemconv/CPUPhysicalCountObservable", constructorTest(systemconv.NewCPUPhysicalCountObservable, "Int64ObservableUpDownCounterConfig", "system.cpu.physical.count", "Reports the number of actual physical processor cores on the hardware.", "{cpu}", false, nil)},
		{"systemconv/CPUTime", constructorTest(systemconv.NewCPUTime, "Float64ObservableCounterConfig", "system.cpu.time", "Seconds each logical CPU spent on each mode.", "s", false, nil)},
		{"systemconv/CPUUtilization", constructorTest(systemconv.NewCPUUtilization, "Float64GaugeConfig", "system.cpu.utilization", "For each logical CPU, the utilization is calculated as the change in cumulative CPU time (cpu.time) over a measurement interval, divided by the elapsed time.", "1", false, nil)},
		{"systemconv/CPUUtilizationObservable", constructorTest(systemconv.NewCPUUtilizationObservable, "Float64ObservableGaugeConfig", "system.cpu.utilization", "For each logical CPU, the utilization is calculated as the change in cumulative CPU time (cpu.time) over a measurement interval, divided by the elapsed time.", "1", false, nil)},
		{"systemconv/DiskIO", constructorTest(systemconv.NewDiskIO, "Int64CounterConfig", "system.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"systemconv/DiskIOObservable", constructorTest(systemconv.NewDiskIOObservable, "Int64ObservableCounterConfig", "system.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"systemconv/DiskIOTime", constructorTest(systemconv.NewDiskIOTime, "Float64CounterConfig", "system.disk.io_time", "Time disk spent activated.", "s", false, nil)},
		{"systemconv/DiskIOTimeObservable", constructorTest(systemconv.NewDiskIOTimeObservable, "Float64ObservableCounterConfig", "system.disk.io_time", "Time disk spent activated.", "s", false, nil)},
		{"systemconv/DiskLimit", constructorTest(systemconv.NewDiskLimit, "Int64UpDownCounterConfig", "system.disk.limit", "The total storage capacity of the disk.", "By", false, nil)},
		{"systemconv/DiskLimitObservable", constructorTest(systemconv.NewDiskLimitObservable, "Int64ObservableUpDownCounterConfig", "system.disk.limit", "The total storage capacity of the disk.", "By", false, nil)},
		{"systemconv/DiskMerged", constructorTest(systemconv.NewDiskMerged, "Int64CounterConfig", "system.disk.merged", "The number of disk reads/writes merged into single physical disk access operations.", "{operation}", false, nil)},
		{"systemconv/DiskMergedObservable", constructorTest(systemconv.NewDiskMergedObservable, "Int64ObservableCounterConfig", "system.disk.merged", "The number of disk reads/writes merged into single physical disk access operations.", "{operation}", false, nil)},
		{"systemconv/DiskOperationTime", constructorTest(systemconv.NewDiskOperationTime, "Float64CounterConfig", "system.disk.operation_time", "Sum of the time each operation took to complete.", "s", false, nil)},
		{"systemconv/DiskOperationTimeObservable", constructorTest(systemconv.NewDiskOperationTimeObservable, "Float64ObservableCounterConfig", "system.disk.operation_time", "Sum of the time each operation took to complete.", "s", false, nil)},
		{"systemconv/DiskOperations", constructorTest(systemconv.NewDiskOperations, "Int64CounterConfig", "system.disk.operations", "Disk operations count.", "{operation}", false, nil)},
		{"systemconv/DiskOperationsObservable", constructorTest(systemconv.NewDiskOperationsObservable, "Int64ObservableCounterConfig", "system.disk.operations", "Disk operations count.", "{operation}", false, nil)},
		{"systemconv/FilesystemLimit", constructorTest(systemconv.NewFilesystemLimit, "Int64UpDownCounterConfig", "system.filesystem.limit", "The total storage capacity of the filesystem.", "By", false, nil)},
		{"systemconv/FilesystemLimitObservable", constructorTest(systemconv.NewFilesystemLimitObservable, "Int64ObservableUpDownCounterConfig", "system.filesystem.limit", "The total storage capacity of the filesystem.", "By", false, nil)},
		{"systemconv/FilesystemLockCount", constructorTest(systemconv.NewFilesystemLockCount, "Int64UpDownCounterConfig", "system.filesystem.lock.count", "Filesystem lock counts.", "{lock}", false, nil)},
		{"systemconv/FilesystemLockCountObservable", constructorTest(systemconv.NewFilesystemLockCountObservable, "Int64ObservableUpDownCounterConfig", "system.filesystem.lock.count", "Filesystem lock counts.", "{lock}", false, nil)},
		{"systemconv/FilesystemUsage", constructorTest(systemconv.NewFilesystemUsage, "Int64UpDownCounterConfig", "system.filesystem.usage", "Reports a filesystem\u0027s space usage across different states.", "By", false, nil)},
		{"systemconv/FilesystemUsageObservable", constructorTest(systemconv.NewFilesystemUsageObservable, "Int64ObservableUpDownCounterConfig", "system.filesystem.usage", "Reports a filesystem\u0027s space usage across different states.", "By", false, nil)},
		{"systemconv/FilesystemUtilization", constructorTest(systemconv.NewFilesystemUtilization, "Float64GaugeConfig", "system.filesystem.utilization", "Fraction of filesystem bytes used.", "1", false, nil)},
		{"systemconv/FilesystemUtilizationObservable", constructorTest(systemconv.NewFilesystemUtilizationObservable, "Float64ObservableGaugeConfig", "system.filesystem.utilization", "Fraction of filesystem bytes used.", "1", false, nil)},
		{"systemconv/MemoryLimit", constructorTest(systemconv.NewMemoryLimit, "Int64UpDownCounterConfig", "system.memory.limit", "Total virtual memory available in the system.", "By", false, nil)},
		{"systemconv/MemoryLimitObservable", constructorTest(systemconv.NewMemoryLimitObservable, "Int64ObservableUpDownCounterConfig", "system.memory.limit", "Total virtual memory available in the system.", "By", false, nil)},
		{"systemconv/MemoryLinuxAvailable", constructorTest(systemconv.NewMemoryLinuxAvailable, "Int64UpDownCounterConfig", "system.memory.linux.available", "An estimate of how much memory is available for starting new applications, without causing swapping.", "By", false, nil)},
		{"systemconv/MemoryLinuxAvailableObservable", constructorTest(systemconv.NewMemoryLinuxAvailableObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.available", "An estimate of how much memory is available for starting new applications, without causing swapping.", "By", false, nil)},
		{"systemconv/MemoryLinuxHugepagesLimit", constructorTest(systemconv.NewMemoryLinuxHugepagesLimit, "Int64UpDownCounterConfig", "system.memory.linux.hugepages.limit", "Total number of hugepages available.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesLimitObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesLimitObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.hugepages.limit", "Total number of hugepages available.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesPageSize", constructorTest(systemconv.NewMemoryLinuxHugepagesPageSize, "Int64UpDownCounterConfig", "system.memory.linux.hugepages.page_size", "System hugepage size in bytes.", "By", false, nil)},
		{"systemconv/MemoryLinuxHugepagesPageSizeObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesPageSizeObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.hugepages.page_size", "System hugepage size in bytes.", "By", false, nil)},
		{"systemconv/MemoryLinuxHugepagesReserved", constructorTest(systemconv.NewMemoryLinuxHugepagesReserved, "Int64UpDownCounterConfig", "system.memory.linux.hugepages.reserved", "Number of reserved hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesReservedObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesReservedObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.hugepages.reserved", "Number of reserved hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesSurplus", constructorTest(systemconv.NewMemoryLinuxHugepagesSurplus, "Int64UpDownCounterConfig", "system.memory.linux.hugepages.surplus", "Number of surplus hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesSurplusObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesSurplusObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.hugepages.surplus", "Number of surplus hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUsage", constructorTest(systemconv.NewMemoryLinuxHugepagesUsage, "Int64UpDownCounterConfig", "system.memory.linux.hugepages.usage", "Number of hugepages in use by state.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUsageObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesUsageObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.hugepages.usage", "Number of hugepages in use by state.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUtilization", constructorTest(systemconv.NewMemoryLinuxHugepagesUtilization, "Float64GaugeConfig", "system.memory.linux.hugepages.utilization", "Percentage of hugepages in use by state.", "1", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUtilizationObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesUtilizationObservable, "Float64ObservableGaugeConfig", "system.memory.linux.hugepages.utilization", "Percentage of hugepages in use by state.", "1", false, nil)},
		{"systemconv/MemoryLinuxShared", constructorTest(systemconv.NewMemoryLinuxShared, "Int64UpDownCounterConfig", "system.memory.linux.shared", "Shared memory used (mostly by tmpfs).", "By", false, nil)},
		{"systemconv/MemoryLinuxSharedObservable", constructorTest(systemconv.NewMemoryLinuxSharedObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.shared", "Shared memory used (mostly by tmpfs).", "By", false, nil)},
		{"systemconv/MemoryLinuxSlabUsage", constructorTest(systemconv.NewMemoryLinuxSlabUsage, "Int64UpDownCounterConfig", "system.memory.linux.slab.usage", "Reports the memory used by the Linux kernel for managing caches of frequently used objects.", "By", false, nil)},
		{"systemconv/MemoryLinuxSlabUsageObservable", constructorTest(systemconv.NewMemoryLinuxSlabUsageObservable, "Int64ObservableUpDownCounterConfig", "system.memory.linux.slab.usage", "Reports the memory used by the Linux kernel for managing caches of frequently used objects.", "By", false, nil)},
		{"systemconv/MemoryUsage", constructorTest(systemconv.NewMemoryUsage, "Int64ObservableUpDownCounterConfig", "system.memory.usage", "Reports memory by state.", "By", false, nil)},
		{"systemconv/MemoryUtilization", constructorTest(systemconv.NewMemoryUtilization, "Float64ObservableGaugeConfig", "system.memory.utilization", "Fraction of memory bytes by state.", "1", false, nil)},
		{"systemconv/NetworkConnectionCount", constructorTest(systemconv.NewNetworkConnectionCount, "Int64UpDownCounterConfig", "system.network.connection.count", "The number of connections.", "{connection}", false, nil)},
		{"systemconv/NetworkConnectionCountObservable", constructorTest(systemconv.NewNetworkConnectionCountObservable, "Int64ObservableUpDownCounterConfig", "system.network.connection.count", "The number of connections.", "{connection}", false, nil)},
		{"systemconv/NetworkErrors", constructorTest(systemconv.NewNetworkErrors, "Int64CounterConfig", "system.network.errors", "Count of network errors detected.", "{error}", false, nil)},
		{"systemconv/NetworkErrorsObservable", constructorTest(systemconv.NewNetworkErrorsObservable, "Int64ObservableCounterConfig", "system.network.errors", "Count of network errors detected.", "{error}", false, nil)},
		{"systemconv/NetworkIO", constructorTest(systemconv.NewNetworkIO, "Int64ObservableCounterConfig", "system.network.io", "The number of bytes transmitted and received.", "By", false, nil)},
		{"systemconv/NetworkPacketCount", constructorTest(systemconv.NewNetworkPacketCount, "Int64CounterConfig", "system.network.packet.count", "The number of packets transferred.", "{packet}", false, nil)},
		{"systemconv/NetworkPacketCountObservable", constructorTest(systemconv.NewNetworkPacketCountObservable, "Int64ObservableCounterConfig", "system.network.packet.count", "The number of packets transferred.", "{packet}", false, nil)},
		{"systemconv/NetworkPacketDropped", constructorTest(systemconv.NewNetworkPacketDropped, "Int64CounterConfig", "system.network.packet.dropped", "Count of packets that are dropped or discarded even though there was no error.", "{packet}", false, nil)},
		{"systemconv/NetworkPacketDroppedObservable", constructorTest(systemconv.NewNetworkPacketDroppedObservable, "Int64ObservableCounterConfig", "system.network.packet.dropped", "Count of packets that are dropped or discarded even though there was no error.", "{packet}", false, nil)},
		{"systemconv/PagingFaults", constructorTest(systemconv.NewPagingFaults, "Int64CounterConfig", "system.paging.faults", "The number of page faults.", "{fault}", false, nil)},
		{"systemconv/PagingFaultsObservable", constructorTest(systemconv.NewPagingFaultsObservable, "Int64ObservableCounterConfig", "system.paging.faults", "The number of page faults.", "{fault}", false, nil)},
		{"systemconv/PagingOperations", constructorTest(systemconv.NewPagingOperations, "Int64CounterConfig", "system.paging.operations", "The number of paging operations.", "{operation}", false, nil)},
		{"systemconv/PagingOperationsObservable", constructorTest(systemconv.NewPagingOperationsObservable, "Int64ObservableCounterConfig", "system.paging.operations", "The number of paging operations.", "{operation}", false, nil)},
		{"systemconv/PagingUsage", constructorTest(systemconv.NewPagingUsage, "Int64UpDownCounterConfig", "system.paging.usage", "UNIX swap or windows pagefile usage.", "By", false, nil)},
		{"systemconv/PagingUsageObservable", constructorTest(systemconv.NewPagingUsageObservable, "Int64ObservableUpDownCounterConfig", "system.paging.usage", "UNIX swap or windows pagefile usage.", "By", false, nil)},
		{"systemconv/PagingUtilization", constructorTest(systemconv.NewPagingUtilization, "Float64GaugeConfig", "system.paging.utilization", "Swap (UNIX) or pagefile (windows) utilization.", "1", false, nil)},
		{"systemconv/PagingUtilizationObservable", constructorTest(systemconv.NewPagingUtilizationObservable, "Float64ObservableGaugeConfig", "system.paging.utilization", "Swap (UNIX) or pagefile (windows) utilization.", "1", false, nil)},
		{"systemconv/ProcessCount", constructorTest(systemconv.NewProcessCount, "Int64UpDownCounterConfig", "system.process.count", "Total number of processes in each state.", "{process}", false, nil)},
		{"systemconv/ProcessCountObservable", constructorTest(systemconv.NewProcessCountObservable, "Int64ObservableUpDownCounterConfig", "system.process.count", "Total number of processes in each state.", "{process}", false, nil)},
		{"systemconv/ProcessCreated", constructorTest(systemconv.NewProcessCreated, "Int64CounterConfig", "system.process.created", "Total number of processes created over uptime of the host.", "{process}", false, nil)},
		{"systemconv/ProcessCreatedObservable", constructorTest(systemconv.NewProcessCreatedObservable, "Int64ObservableCounterConfig", "system.process.created", "Total number of processes created over uptime of the host.", "{process}", false, nil)},
		{"systemconv/ProcessLimit", constructorTest(systemconv.NewProcessLimit, "Int64UpDownCounterConfig", "system.process.limit", "The maximum number of concurrent processes/tasks allowed by the operating system.", "{thread}", false, nil)},
		{"systemconv/ProcessLimitObservable", constructorTest(systemconv.NewProcessLimitObservable, "Int64ObservableUpDownCounterConfig", "system.process.limit", "The maximum number of concurrent processes/tasks allowed by the operating system.", "{thread}", false, nil)},
		{"systemconv/Uptime", constructorTest(systemconv.NewUptime, "Float64GaugeConfig", "system.uptime", "The time the system has been running.", "s", false, nil)},
		{"systemconv/UptimeObservable", constructorTest(systemconv.NewUptimeObservable, "Float64ObservableGaugeConfig", "system.uptime", "The time the system has been running.", "s", false, nil)},
		{"vcsconv/ChangeCount", constructorTest(vcsconv.NewChangeCount, "Int64UpDownCounterConfig", "vcs.change.count", "The number of changes (pull requests/merge requests/changelists) in a repository, categorized by their state (e.g. open or merged).", "{change}", false, nil)},
		{"vcsconv/ChangeCountObservable", constructorTest(vcsconv.NewChangeCountObservable, "Int64ObservableUpDownCounterConfig", "vcs.change.count", "The number of changes (pull requests/merge requests/changelists) in a repository, categorized by their state (e.g. open or merged).", "{change}", false, nil)},
		{"vcsconv/ChangeDuration", constructorTest(vcsconv.NewChangeDuration, "Float64GaugeConfig", "vcs.change.duration", "The time duration a change (pull request/merge request/changelist) has been in a given state.", "s", false, nil)},
		{"vcsconv/ChangeDurationObservable", constructorTest(vcsconv.NewChangeDurationObservable, "Float64ObservableGaugeConfig", "vcs.change.duration", "The time duration a change (pull request/merge request/changelist) has been in a given state.", "s", false, nil)},
		{"vcsconv/ChangeTimeToApproval", constructorTest(vcsconv.NewChangeTimeToApproval, "Float64GaugeConfig", "vcs.change.time_to_approval", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get the first approval.", "s", false, nil)},
		{"vcsconv/ChangeTimeToApprovalObservable", constructorTest(vcsconv.NewChangeTimeToApprovalObservable, "Float64ObservableGaugeConfig", "vcs.change.time_to_approval", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get the first approval.", "s", false, nil)},
		{"vcsconv/ChangeTimeToMerge", constructorTest(vcsconv.NewChangeTimeToMerge, "Float64GaugeConfig", "vcs.change.time_to_merge", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get merged into the target(base) ref.", "s", false, nil)},
		{"vcsconv/ChangeTimeToMergeObservable", constructorTest(vcsconv.NewChangeTimeToMergeObservable, "Float64ObservableGaugeConfig", "vcs.change.time_to_merge", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get merged into the target(base) ref.", "s", false, nil)},
		{"vcsconv/ContributorCount", constructorTest(vcsconv.NewContributorCount, "Int64GaugeConfig", "vcs.contributor.count", "The number of unique contributors to a repository.", "{contributor}", false, nil)},
		{"vcsconv/ContributorCountObservable", constructorTest(vcsconv.NewContributorCountObservable, "Int64ObservableGaugeConfig", "vcs.contributor.count", "The number of unique contributors to a repository.", "{contributor}", false, nil)},
		{"vcsconv/RefCount", constructorTest(vcsconv.NewRefCount, "Int64UpDownCounterConfig", "vcs.ref.count", "The number of refs of type branch or tag in a repository.", "{ref}", false, nil)},
		{"vcsconv/RefCountObservable", constructorTest(vcsconv.NewRefCountObservable, "Int64ObservableUpDownCounterConfig", "vcs.ref.count", "The number of refs of type branch or tag in a repository.", "{ref}", false, nil)},
		{"vcsconv/RefLinesDelta", constructorTest(vcsconv.NewRefLinesDelta, "Int64GaugeConfig", "vcs.ref.lines_delta", "The number of lines added/removed in a ref (branch) relative to the ref from the `vcs.ref.base.name` attribute.", "{line}", false, nil)},
		{"vcsconv/RefLinesDeltaObservable", constructorTest(vcsconv.NewRefLinesDeltaObservable, "Int64ObservableGaugeConfig", "vcs.ref.lines_delta", "The number of lines added/removed in a ref (branch) relative to the ref from the `vcs.ref.base.name` attribute.", "{line}", false, nil)},
		{"vcsconv/RefRevisionsDelta", constructorTest(vcsconv.NewRefRevisionsDelta, "Int64GaugeConfig", "vcs.ref.revisions_delta", "The number of revisions (commits) a ref (branch) is ahead/behind the branch from the `vcs.ref.base.name` attribute.", "{revision}", false, nil)},
		{"vcsconv/RefRevisionsDeltaObservable", constructorTest(vcsconv.NewRefRevisionsDeltaObservable, "Int64ObservableGaugeConfig", "vcs.ref.revisions_delta", "The number of revisions (commits) a ref (branch) is ahead/behind the branch from the `vcs.ref.base.name` attribute.", "{revision}", false, nil)},
		{"vcsconv/RefTime", constructorTest(vcsconv.NewRefTime, "Float64GaugeConfig", "vcs.ref.time", "Time a ref (branch) created from the default branch (trunk) has existed. The `ref.type` attribute will always be `branch`.", "s", false, nil)},
		{"vcsconv/RefTimeObservable", constructorTest(vcsconv.NewRefTimeObservable, "Float64ObservableGaugeConfig", "vcs.ref.time", "Time a ref (branch) created from the default branch (trunk) has existed. The `ref.type` attribute will always be `branch`.", "s", false, nil)},
		{"vcsconv/RepositoryCount", constructorTest(vcsconv.NewRepositoryCount, "Int64UpDownCounterConfig", "vcs.repository.count", "The number of repositories in an organization.", "{repository}", false, nil)},
		{"vcsconv/RepositoryCountObservable", constructorTest(vcsconv.NewRepositoryCountObservable, "Int64ObservableUpDownCounterConfig", "vcs.repository.count", "The number of repositories in an organization.", "{repository}", false, nil)},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func constructorTest[O any, I interface {
	Name() string
	Unit() string
}](
	newInstrument func(metric.Meter, ...O) (I, error),
	kind, name, description, unit string,
	histogram bool,
	advice []float64,
) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		for _, callerOptions := range []bool{false, true} {
			label := "defaults"
			if callerOptions {
				label = "caller_options"
			}
			t.Run(label, func(t *testing.T) {
				m := &constructorMeter{Meter: noop.NewMeterProvider().Meter("test")}
				var options, storage, before []O
				wantBoundaries := advice
				if callerOptions {
					// Keep sentinels beyond len(options) to detect writes into spare capacity.
					storage = make([]O, 16)
					for j := range storage {
						storage[j] = any(metric.WithUnit("sentinel")).(O)
					}
					options = storage[:2]
					options[0] = any(metric.WithDescription("caller description")).(O)
					options[1] = any(metric.WithUnit("caller unit")).(O)
					if histogram {
						wantBoundaries = []float64{1, 2, 3}
						options = append(options, any(metric.WithExplicitBucketBoundaries(wantBoundaries...)).(O))
					}
					before = slices.Clone(storage)
				}
				instrument, err := newInstrument(m, options...)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(storage, before) {
					t.Error("constructor changed caller option storage")
				}
				if m.config == nil {
					t.Fatal("constructor did not create an instrument")
				}
				if got := reflect.TypeOf(m.config).Name(); got != kind {
					t.Errorf("instrument config = %s, want %s", got, kind)
				}
				if got := instrument.Name(); got != name {
					t.Errorf("instrument name = %q, want %q", got, name)
				}
				if got := instrument.Unit(); got != unit {
					t.Errorf("instrument unit = %q, want %q", got, unit)
				}
				if m.name != name {
					t.Errorf("name = %q, want %q", m.name, name)
				}
				if got := m.config.Description(); got != description {
					t.Errorf("description = %q, want %q", got, description)
				}
				if got := m.config.Unit(); got != unit {
					t.Errorf("unit = %q, want %q", got, unit)
				}
				if histogram {
					config, ok := m.config.(interface{ ExplicitBucketBoundaries() []float64 })
					if !ok {
						t.Fatal("constructor did not create a histogram")
					}
					if got := config.ExplicitBucketBoundaries(); !slices.Equal(got, wantBoundaries) {
						t.Errorf("bucket boundaries = %v, want %v", got, wantBoundaries)
					}
				}
			})
		}
	}
}

type constructorMeter struct {
	metric.Meter
	name   string
	config interface {
		Description() string
		Unit() string
	}
}

func (m *constructorMeter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	m.name = name
	m.config = metric.NewInt64CounterConfig(options...)
	return m.Meter.Int64Counter(name, options...)
}

func (m *constructorMeter) Int64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	m.name = name
	m.config = metric.NewInt64UpDownCounterConfig(options...)
	return m.Meter.Int64UpDownCounter(name, options...)
}

func (m *constructorMeter) Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	m.name = name
	m.config = metric.NewInt64HistogramConfig(options...)
	return m.Meter.Int64Histogram(name, options...)
}

func (m *constructorMeter) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	m.name = name
	m.config = metric.NewInt64GaugeConfig(options...)
	return m.Meter.Int64Gauge(name, options...)
}

func (m *constructorMeter) Int64ObservableCounter(name string, options ...metric.Int64ObservableCounterOption) (metric.Int64ObservableCounter, error) {
	m.name = name
	m.config = metric.NewInt64ObservableCounterConfig(options...)
	return m.Meter.Int64ObservableCounter(name, options...)
}

func (m *constructorMeter) Int64ObservableUpDownCounter(name string, options ...metric.Int64ObservableUpDownCounterOption) (metric.Int64ObservableUpDownCounter, error) {
	m.name = name
	m.config = metric.NewInt64ObservableUpDownCounterConfig(options...)
	return m.Meter.Int64ObservableUpDownCounter(name, options...)
}

func (m *constructorMeter) Int64ObservableGauge(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error) {
	m.name = name
	m.config = metric.NewInt64ObservableGaugeConfig(options...)
	return m.Meter.Int64ObservableGauge(name, options...)
}

func (m *constructorMeter) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	m.name = name
	m.config = metric.NewFloat64CounterConfig(options...)
	return m.Meter.Float64Counter(name, options...)
}

func (m *constructorMeter) Float64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	m.name = name
	m.config = metric.NewFloat64UpDownCounterConfig(options...)
	return m.Meter.Float64UpDownCounter(name, options...)
}

func (m *constructorMeter) Float64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	m.name = name
	m.config = metric.NewFloat64HistogramConfig(options...)
	return m.Meter.Float64Histogram(name, options...)
}

func (m *constructorMeter) Float64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	m.name = name
	m.config = metric.NewFloat64GaugeConfig(options...)
	return m.Meter.Float64Gauge(name, options...)
}

func (m *constructorMeter) Float64ObservableCounter(name string, options ...metric.Float64ObservableCounterOption) (metric.Float64ObservableCounter, error) {
	m.name = name
	m.config = metric.NewFloat64ObservableCounterConfig(options...)
	return m.Meter.Float64ObservableCounter(name, options...)
}

func (m *constructorMeter) Float64ObservableUpDownCounter(name string, options ...metric.Float64ObservableUpDownCounterOption) (metric.Float64ObservableUpDownCounter, error) {
	m.name = name
	m.config = metric.NewFloat64ObservableUpDownCounterConfig(options...)
	return m.Meter.Float64ObservableUpDownCounter(name, options...)
}

func (m *constructorMeter) Float64ObservableGauge(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error) {
	m.name = name
	m.config = metric.NewFloat64ObservableGaugeConfig(options...)
	return m.Meter.Float64ObservableGauge(name, options...)
}

func TestMetricMeasurements(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{"azureconv/CosmosDBClientActiveInstanceCount", func(t *testing.T) {
			capture := &Int64UpDownCounterCapture{}
			instrument := azureconv.CosmosDBClientActiveInstanceCount{Int64UpDownCounter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{
				instrument.AttrServerPort(7),
				instrument.AttrServerAddress("value"),
			}, []attribute.KeyValue{
				attribute.Int("server.port", 7),
				attribute.String("server.address", "value"),
			}, 7)
		}},
		{"azureconv/CosmosDBClientOperationRequestCharge", func(t *testing.T) {
			capture := &Float64HistogramCapture{}
			instrument := azureconv.CosmosDBClientOperationRequestCharge{Float64Histogram: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, "value", attrs...)
			}, []attribute.KeyValue{
				attribute.String("db.operation.name", "value"),
			}, []attribute.KeyValue{
				instrument.AttrCosmosDBConsistencyLevel(azureconv.CosmosDBConsistencyLevelAttr("Strong")),
				instrument.AttrCosmosDBResponseSubStatusCode(7),
				instrument.AttrDBCollectionName("value"),
			}, []attribute.KeyValue{
				attribute.String("azure.cosmosdb.consistency.level", "Strong"),
				attribute.Int("azure.cosmosdb.response.sub_status_code", 7),
				attribute.String("db.collection.name", "value"),
			}, 0.125, 1.5)
		}},
		{"cicdconv/PipelineRunActive", func(t *testing.T) {
			capture := &Int64UpDownCounterCapture{}
			instrument := cicdconv.PipelineRunActive{Int64UpDownCounter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, "value", cicdconv.PipelineRunStateAttr("pending"), attrs...)
			}, []attribute.KeyValue{
				attribute.String("cicd.pipeline.name", "value"),
				attribute.String("cicd.pipeline.run.state", "pending"),
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, 7)
		}},
		{"cicdconv/PipelineRunDuration", func(t *testing.T) {
			capture := &Float64HistogramCapture{}
			instrument := cicdconv.PipelineRunDuration{Float64Histogram: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, "value", cicdconv.PipelineRunStateAttr("pending"), attrs...)
			}, []attribute.KeyValue{
				attribute.String("cicd.pipeline.name", "value"),
				attribute.String("cicd.pipeline.run.state", "pending"),
			}, []attribute.KeyValue{
				instrument.AttrPipelineResult(cicdconv.PipelineResultAttr("success")),
			}, []attribute.KeyValue{
				attribute.String("cicd.pipeline.result", "success"),
			}, 0.125, 1.5)
		}},
		{"cicdconv/PipelineRunErrors", func(t *testing.T) {
			capture := &Int64CounterCapture{}
			instrument := cicdconv.PipelineRunErrors{Int64Counter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, "value", cicdconv.ErrorTypeAttr("_OTHER"), attrs...)
			}, []attribute.KeyValue{
				attribute.String("cicd.pipeline.name", "value"),
				attribute.String("error.type", "_OTHER"),
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, 7)
		}},
		{"containerconv/CPUTime", func(t *testing.T) {
			capture := &Float64CounterCapture{}
			instrument := containerconv.CPUTime{Float64Counter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{
				instrument.AttrCPUMode(containerconv.CPUModeAttr("user")),
			}, []attribute.KeyValue{
				attribute.String("cpu.mode", "user"),
			}, 0.125, 1.5)
		}},
		{"containerconv/CPUUsage", func(t *testing.T) {
			capture := &Float64GaugeCapture{}
			instrument := containerconv.CPUUsage{Float64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{
				instrument.AttrCPUMode(containerconv.CPUModeAttr("user")),
			}, []attribute.KeyValue{
				attribute.String("cpu.mode", "user"),
			}, 0.125, 1.5)
		}},
		{"containerconv/DiskIO", func(t *testing.T) {
			capture := &Int64CounterCapture{}
			instrument := containerconv.DiskIO{Int64Counter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{
				instrument.AttrDiskIODirection(containerconv.DiskIODirectionAttr("read")),
				instrument.AttrSystemDeviceName("value"),
			}, []attribute.KeyValue{
				attribute.String("disk.io.direction", "read"),
				attribute.String("system.device.name", "value"),
			}, 7)
		}},
		{"containerconv/FilesystemAvailable", func(t *testing.T) {
			capture := &Int64UpDownCounterCapture{}
			instrument := containerconv.FilesystemAvailable{Int64UpDownCounter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, []attribute.KeyValue{}, 7)
		}},
		{"containerconv/Uptime", func(t *testing.T) {
			capture := &Float64GaugeCapture{}
			instrument := containerconv.Uptime{Float64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, []attribute.KeyValue{}, 0.125, 1.5)
		}},
		{"dbconv/ClientConnectionIdleMax", func(t *testing.T) {
			capture := &Int64UpDownCounterCapture{}
			instrument := dbconv.ClientConnectionIdleMax{Int64UpDownCounter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, "value", attrs...)
			}, []attribute.KeyValue{
				attribute.String("db.client.connection.pool.name", "value"),
			}, []attribute.KeyValue{
				instrument.AttrNamespace("value"),
				instrument.AttrSystemName(dbconv.SystemNameAttr("other_sql")),
				instrument.AttrServerPort(7),
			}, []attribute.KeyValue{
				attribute.String("db.namespace", "value"),
				attribute.String("db.system.name", "other_sql"),
				attribute.Int("server.port", 7),
			}, 7)
		}},
		{"dbconv/ClientConnectionTimeouts", func(t *testing.T) {
			capture := &Int64CounterCapture{}
			instrument := dbconv.ClientConnectionTimeouts{Int64Counter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, "value", attrs...)
			}, []attribute.KeyValue{
				attribute.String("db.client.connection.pool.name", "value"),
			}, []attribute.KeyValue{
				instrument.AttrNamespace("value"),
				instrument.AttrSystemName(dbconv.SystemNameAttr("other_sql")),
				instrument.AttrServerPort(7),
			}, []attribute.KeyValue{
				attribute.String("db.namespace", "value"),
				attribute.String("db.system.name", "other_sql"),
				attribute.Int("server.port", 7),
			}, 7)
		}},
		{"dbconv/ClientResponseReturnedRows", func(t *testing.T) {
			capture := &Int64HistogramCapture{}
			instrument := dbconv.ClientResponseReturnedRows{Int64Histogram: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, dbconv.SystemNameAttr("other_sql"), attrs...)
			}, []attribute.KeyValue{
				attribute.String("db.system.name", "other_sql"),
			}, []attribute.KeyValue{
				instrument.AttrCollectionName("value"),
				instrument.AttrErrorType(dbconv.ErrorTypeAttr("_OTHER")),
				instrument.AttrServerPort(7),
			}, []attribute.KeyValue{
				attribute.String("db.collection.name", "value"),
				attribute.String("error.type", "_OTHER"),
				attribute.Int("server.port", 7),
			}, 7)
		}},
		{"faasconv/CPUUsage", func(t *testing.T) {
			capture := &Float64HistogramCapture{}
			instrument := faasconv.CPUUsage{Float64Histogram: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{
				instrument.AttrTrigger(faasconv.TriggerAttr("datasource")),
			}, []attribute.KeyValue{
				attribute.String("faas.trigger", "datasource"),
			}, 0.125, 1.5)
		}},
		{"faasconv/MemUsage", func(t *testing.T) {
			capture := &Int64HistogramCapture{}
			instrument := faasconv.MemUsage{Int64Histogram: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{
				instrument.AttrTrigger(faasconv.TriggerAttr("datasource")),
			}, []attribute.KeyValue{
				attribute.String("faas.trigger", "datasource"),
			}, 7)
		}},
		{"goconv/CPUTime", func(t *testing.T) {
			capture := &Float64CounterCapture{}
			instrument := goconv.CPUTime{Float64Counter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, goconv.CPUStateAttr("user"), attrs...)
			}, []attribute.KeyValue{
				attribute.String("go.cpu.state", "user"),
			}, []attribute.KeyValue{
				instrument.AttrCPUDetailedState("value"),
			}, []attribute.KeyValue{
				attribute.String("go.cpu.detailed_state", "value"),
			}, 0.125, 1.5)
		}},
		{"goconv/MemoryGCCycles", func(t *testing.T) {
			capture := &Int64CounterCapture{}
			instrument := goconv.MemoryGCCycles{Int64Counter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, []attribute.KeyValue{}, 7)
		}},
		{"goconv/MemoryGCPauseDuration", func(t *testing.T) {
			capture := &Float64HistogramCapture{}
			instrument := goconv.MemoryGCPauseDuration{Float64Histogram: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, []attribute.KeyValue{}, 0.125, 1.5)
		}},
		{"hwconv/BatteryCharge", func(t *testing.T) {
			capture := &Float64GaugeCapture{}
			instrument := hwconv.BatteryCharge{Float64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, "value", attrs...)
			}, []attribute.KeyValue{
				attribute.String("hw.id", "value"),
			}, []attribute.KeyValue{
				instrument.AttrBatteryCapacity("value"),
			}, []attribute.KeyValue{
				attribute.String("hw.battery.capacity", "value"),
			}, 0.125, 1.5)
		}},
		{"hwconv/BatteryTimeLeft", func(t *testing.T) {
			capture := &Float64GaugeCapture{}
			instrument := hwconv.BatteryTimeLeft{Float64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, "value", hwconv.StateAttr("degraded"), attrs...)
			}, []attribute.KeyValue{
				attribute.String("hw.id", "value"),
				attribute.String("hw.state", "degraded"),
			}, []attribute.KeyValue{
				instrument.AttrBatteryState(hwconv.BatteryStateAttr("charging")),
				instrument.AttrBatteryCapacity("value"),
			}, []attribute.KeyValue{
				attribute.String("hw.battery.state", "charging"),
				attribute.String("hw.battery.capacity", "value"),
			}, 0.125, 1.5)
		}},
		{"hwconv/CPUSpeed", func(t *testing.T) {
			capture := &Int64GaugeCapture{}
			instrument := hwconv.CPUSpeed{Int64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, "value", attrs...)
			}, []attribute.KeyValue{
				attribute.String("hw.id", "value"),
			}, []attribute.KeyValue{
				instrument.AttrModel("value"),
			}, []attribute.KeyValue{
				attribute.String("hw.model", "value"),
			}, 7)
		}},
		{"hwconv/PowerSupplyUsage", func(t *testing.T) {
			capture := &Float64UpDownCounterCapture{}
			instrument := hwconv.PowerSupplyUsage{Float64UpDownCounter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, "value", attrs...)
			}, []attribute.KeyValue{
				attribute.String("hw.id", "value"),
			}, []attribute.KeyValue{
				instrument.AttrModel("value"),
			}, []attribute.KeyValue{
				attribute.String("hw.model", "value"),
			}, 0.125, 1.5)
		}},
		{"k8sconv/ContainerCPULimitCurrent", func(t *testing.T) {
			capture := &Float64UpDownCounterCapture{}
			instrument := k8sconv.ContainerCPULimitCurrent{Float64UpDownCounter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, []attribute.KeyValue{}, 0.125, 1.5)
		}},
		{"k8sconv/HPAMetricTargetCPUAverageUtilization", func(t *testing.T) {
			capture := &Int64GaugeCapture{}
			instrument := k8sconv.HPAMetricTargetCPUAverageUtilization{Int64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{
				instrument.AttrContainerName("value"),
			}, []attribute.KeyValue{
				attribute.String("k8s.container.name", "value"),
			}, 7)
		}},
		{"k8sconv/NodeCPUTime", func(t *testing.T) {
			capture := &Float64CounterCapture{}
			instrument := k8sconv.NodeCPUTime{Float64Counter: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
				instrument.Add(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, []attribute.KeyValue{}, 0.125, 1.5)
		}},
		{"k8sconv/ServiceEndpointCount", func(t *testing.T) {
			capture := &Int64GaugeCapture{}
			instrument := k8sconv.ServiceEndpointCount{Int64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, k8sconv.ServiceEndpointAddressTypeAttr("IPv4"), k8sconv.ServiceEndpointConditionAttr("ready"), attrs...)
			}, []attribute.KeyValue{
				attribute.String("k8s.service.endpoint.address_type", "IPv4"),
				attribute.String("k8s.service.endpoint.condition", "ready"),
			}, []attribute.KeyValue{
				instrument.AttrServiceEndpointZone("value"),
			}, []attribute.KeyValue{
				attribute.String("k8s.service.endpoint.zone", "value"),
			}, 7)
		}},
		{"k8sconv/ServiceLoadBalancerIngressCount", func(t *testing.T) {
			capture := &Int64GaugeCapture{}
			instrument := k8sconv.ServiceLoadBalancerIngressCount{Int64Gauge: capture}
			measurementTest(t, &capture.measurementCapture, func(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
				instrument.Record(ctx, value, attrs...)
			}, []attribute.KeyValue{}, []attribute.KeyValue{}, []attribute.KeyValue{}, 7)
		}},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func measurementTest[V int64 | float64](
	t *testing.T,
	capture *measurementCapture[V],
	measure func(context.Context, V, ...attribute.KeyValue),
	required, optional, wantOptional []attribute.KeyValue,
	values ...V,
) {
	t.Helper()
	t.Parallel()
	optional = append(optional, attribute.String("extra", "value"))
	wantOptional = append(wantOptional, attribute.String("extra", "value"))
	for _, value := range values {
		for _, withOptional := range []bool{false, true} {
			var attrs []attribute.KeyValue
			wantAttrs := slices.Clone(required)
			if withOptional {
				attrs = optional
				wantAttrs = append(wantAttrs, wantOptional...)
			}
			*capture = measurementCapture[V]{}
			measure(t.Context(), value, attrs...)
			want := attribute.NewSet(wantAttrs...)
			if capture.calls != 1 || capture.value != value || !capture.attrs.Equals(&want) {
				t.Errorf("measurement with optional attributes %t = %d calls, value %v, attributes %v; want 1 call, value %v, attributes %v", withOptional, capture.calls, capture.value, capture.attrs, value, want)
			}
		}
	}
}

type measurementCapture[V int64 | float64] struct {
	calls int
	value V
	attrs attribute.Set
}

type Int64CounterCapture struct {
	noop.Int64Counter
	measurementCapture[int64]
}

func (*Int64CounterCapture) Enabled(context.Context) bool { return true }

func (c *Int64CounterCapture) Add(_ context.Context, value int64, options ...metric.AddOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewAddConfig(options).Attributes()
}

type Int64UpDownCounterCapture struct {
	noop.Int64UpDownCounter
	measurementCapture[int64]
}

func (*Int64UpDownCounterCapture) Enabled(context.Context) bool { return true }

func (c *Int64UpDownCounterCapture) Add(_ context.Context, value int64, options ...metric.AddOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewAddConfig(options).Attributes()
}

type Int64HistogramCapture struct {
	noop.Int64Histogram
	measurementCapture[int64]
}

func (*Int64HistogramCapture) Enabled(context.Context) bool { return true }

func (c *Int64HistogramCapture) Record(_ context.Context, value int64, options ...metric.RecordOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewRecordConfig(options).Attributes()
}

type Int64GaugeCapture struct {
	noop.Int64Gauge
	measurementCapture[int64]
}

func (*Int64GaugeCapture) Enabled(context.Context) bool { return true }

func (c *Int64GaugeCapture) Record(_ context.Context, value int64, options ...metric.RecordOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewRecordConfig(options).Attributes()
}

type Float64CounterCapture struct {
	noop.Float64Counter
	measurementCapture[float64]
}

func (*Float64CounterCapture) Enabled(context.Context) bool { return true }

func (c *Float64CounterCapture) Add(_ context.Context, value float64, options ...metric.AddOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewAddConfig(options).Attributes()
}

type Float64UpDownCounterCapture struct {
	noop.Float64UpDownCounter
	measurementCapture[float64]
}

func (*Float64UpDownCounterCapture) Enabled(context.Context) bool { return true }

func (c *Float64UpDownCounterCapture) Add(_ context.Context, value float64, options ...metric.AddOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewAddConfig(options).Attributes()
}

type Float64HistogramCapture struct {
	noop.Float64Histogram
	measurementCapture[float64]
}

func (*Float64HistogramCapture) Enabled(context.Context) bool { return true }

func (c *Float64HistogramCapture) Record(_ context.Context, value float64, options ...metric.RecordOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewRecordConfig(options).Attributes()
}

type Float64GaugeCapture struct {
	noop.Float64Gauge
	measurementCapture[float64]
}

func (*Float64GaugeCapture) Enabled(context.Context) bool { return true }

func (c *Float64GaugeCapture) Record(_ context.Context, value float64, options ...metric.RecordOption) {
	c.calls++
	c.value = value
	c.attrs = metric.NewRecordConfig(options).Attributes()
}
