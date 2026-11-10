// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package semconv_test

import (
	"reflect"
	"slices"
	"testing"

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
		{"azureconv/CosmosDBClientActiveInstanceCount", constructorTest(azureconv.NewCosmosDBClientActiveInstanceCount, "azure.cosmosdb.client.active_instance.count", "Number of active client instances.", "{instance}", false, nil)},
		{"azureconv/CosmosDBClientActiveInstanceCountObservable", constructorTest(azureconv.NewCosmosDBClientActiveInstanceCountObservable, "azure.cosmosdb.client.active_instance.count", "Number of active client instances.", "{instance}", false, nil)},
		{"azureconv/CosmosDBClientOperationRequestCharge", constructorTest(azureconv.NewCosmosDBClientOperationRequestCharge, "azure.cosmosdb.client.operation.request_charge", "[Request units](https://learn.microsoft.com/azure/cosmos-db/request-units) consumed by the operation.", "{request_unit}", true, nil)},
		{"cicdconv/PipelineRunActive", constructorTest(cicdconv.NewPipelineRunActive, "cicd.pipeline.run.active", "The number of pipeline runs currently active in the system by state.", "{run}", false, nil)},
		{"cicdconv/PipelineRunActiveObservable", constructorTest(cicdconv.NewPipelineRunActiveObservable, "cicd.pipeline.run.active", "The number of pipeline runs currently active in the system by state.", "{run}", false, nil)},
		{"cicdconv/PipelineRunDuration", constructorTest(cicdconv.NewPipelineRunDuration, "cicd.pipeline.run.duration", "Duration of a pipeline run grouped by pipeline, state and result.", "s", true, nil)},
		{"cicdconv/PipelineRunErrors", constructorTest(cicdconv.NewPipelineRunErrors, "cicd.pipeline.run.errors", "The number of errors encountered in pipeline runs (eg. compile, test failures).", "{error}", false, nil)},
		{"cicdconv/PipelineRunErrorsObservable", constructorTest(cicdconv.NewPipelineRunErrorsObservable, "cicd.pipeline.run.errors", "The number of errors encountered in pipeline runs (eg. compile, test failures).", "{error}", false, nil)},
		{"cicdconv/SystemErrors", constructorTest(cicdconv.NewSystemErrors, "cicd.system.errors", "The number of errors in a component of the CI/CD system (eg. controller, scheduler, agent).", "{error}", false, nil)},
		{"cicdconv/SystemErrorsObservable", constructorTest(cicdconv.NewSystemErrorsObservable, "cicd.system.errors", "The number of errors in a component of the CI/CD system (eg. controller, scheduler, agent).", "{error}", false, nil)},
		{"cicdconv/WorkerCount", constructorTest(cicdconv.NewWorkerCount, "cicd.worker.count", "The number of workers on the CI/CD system by state.", "{worker}", false, nil)},
		{"cicdconv/WorkerCountObservable", constructorTest(cicdconv.NewWorkerCountObservable, "cicd.worker.count", "The number of workers on the CI/CD system by state.", "{worker}", false, nil)},
		{"containerconv/CPUTime", constructorTest(containerconv.NewCPUTime, "container.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"containerconv/CPUTimeObservable", constructorTest(containerconv.NewCPUTimeObservable, "container.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"containerconv/CPUUsage", constructorTest(containerconv.NewCPUUsage, "container.cpu.usage", "Container\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"containerconv/CPUUsageObservable", constructorTest(containerconv.NewCPUUsageObservable, "container.cpu.usage", "Container\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"containerconv/DiskIO", constructorTest(containerconv.NewDiskIO, "container.disk.io", "Disk bytes for the container.", "By", false, nil)},
		{"containerconv/DiskIOObservable", constructorTest(containerconv.NewDiskIOObservable, "container.disk.io", "Disk bytes for the container.", "By", false, nil)},
		{"containerconv/FilesystemAvailable", constructorTest(containerconv.NewFilesystemAvailable, "container.filesystem.available", "Container filesystem available bytes.", "By", false, nil)},
		{"containerconv/FilesystemAvailableObservable", constructorTest(containerconv.NewFilesystemAvailableObservable, "container.filesystem.available", "Container filesystem available bytes.", "By", false, nil)},
		{"containerconv/FilesystemCapacity", constructorTest(containerconv.NewFilesystemCapacity, "container.filesystem.capacity", "Container filesystem capacity.", "By", false, nil)},
		{"containerconv/FilesystemCapacityObservable", constructorTest(containerconv.NewFilesystemCapacityObservable, "container.filesystem.capacity", "Container filesystem capacity.", "By", false, nil)},
		{"containerconv/FilesystemUsage", constructorTest(containerconv.NewFilesystemUsage, "container.filesystem.usage", "Container filesystem usage.", "By", false, nil)},
		{"containerconv/FilesystemUsageObservable", constructorTest(containerconv.NewFilesystemUsageObservable, "container.filesystem.usage", "Container filesystem usage.", "By", false, nil)},
		{"containerconv/MemoryAvailable", constructorTest(containerconv.NewMemoryAvailable, "container.memory.available", "Container memory available.", "By", false, nil)},
		{"containerconv/MemoryAvailableObservable", constructorTest(containerconv.NewMemoryAvailableObservable, "container.memory.available", "Container memory available.", "By", false, nil)},
		{"containerconv/MemoryRss", constructorTest(containerconv.NewMemoryRss, "container.memory.rss", "Container memory RSS.", "By", false, nil)},
		{"containerconv/MemoryRssObservable", constructorTest(containerconv.NewMemoryRssObservable, "container.memory.rss", "Container memory RSS.", "By", false, nil)},
		{"containerconv/MemoryUsage", constructorTest(containerconv.NewMemoryUsage, "container.memory.usage", "Memory usage of the container.", "By", false, nil)},
		{"containerconv/MemoryUsageObservable", constructorTest(containerconv.NewMemoryUsageObservable, "container.memory.usage", "Memory usage of the container.", "By", false, nil)},
		{"containerconv/MemoryWorkingSet", constructorTest(containerconv.NewMemoryWorkingSet, "container.memory.working_set", "Container memory working set.", "By", false, nil)},
		{"containerconv/MemoryWorkingSetObservable", constructorTest(containerconv.NewMemoryWorkingSetObservable, "container.memory.working_set", "Container memory working set.", "By", false, nil)},
		{"containerconv/NetworkIO", constructorTest(containerconv.NewNetworkIO, "container.network.io", "Network bytes for the container.", "By", false, nil)},
		{"containerconv/NetworkIOObservable", constructorTest(containerconv.NewNetworkIOObservable, "container.network.io", "Network bytes for the container.", "By", false, nil)},
		{"containerconv/PagingFaults", constructorTest(containerconv.NewPagingFaults, "container.paging.faults", "Container memory paging faults.", "{fault}", false, nil)},
		{"containerconv/PagingFaultsObservable", constructorTest(containerconv.NewPagingFaultsObservable, "container.paging.faults", "Container memory paging faults.", "{fault}", false, nil)},
		{"containerconv/Uptime", constructorTest(containerconv.NewUptime, "container.uptime", "The time the container has been running.", "s", false, nil)},
		{"containerconv/UptimeObservable", constructorTest(containerconv.NewUptimeObservable, "container.uptime", "The time the container has been running.", "s", false, nil)},
		{"dbconv/ClientConnectionCount", constructorTest(dbconv.NewClientConnectionCount, "db.client.connection.count", "The number of connections that are currently in state described by the `state` attribute.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionCountObservable", constructorTest(dbconv.NewClientConnectionCountObservable, "db.client.connection.count", "The number of connections that are currently in state described by the `state` attribute.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionCreateTime", constructorTest(dbconv.NewClientConnectionCreateTime, "db.client.connection.create_time", "The time it took to create a new connection.", "s", true, nil)},
		{"dbconv/ClientConnectionIdleMax", constructorTest(dbconv.NewClientConnectionIdleMax, "db.client.connection.idle.max", "The maximum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionIdleMaxObservable", constructorTest(dbconv.NewClientConnectionIdleMaxObservable, "db.client.connection.idle.max", "The maximum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionIdleMin", constructorTest(dbconv.NewClientConnectionIdleMin, "db.client.connection.idle.min", "The minimum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionIdleMinObservable", constructorTest(dbconv.NewClientConnectionIdleMinObservable, "db.client.connection.idle.min", "The minimum number of idle open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionLimit", constructorTest(dbconv.NewClientConnectionLimit, "db.client.connection.limit", "The maximum number of open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionLimitObservable", constructorTest(dbconv.NewClientConnectionLimitObservable, "db.client.connection.limit", "The maximum number of open connections allowed.", "{connection}", false, nil)},
		{"dbconv/ClientConnectionPendingRequests", constructorTest(dbconv.NewClientConnectionPendingRequests, "db.client.connection.pending_requests", "The number of current pending requests for an open connection.", "{request}", false, nil)},
		{"dbconv/ClientConnectionPendingRequestsObservable", constructorTest(dbconv.NewClientConnectionPendingRequestsObservable, "db.client.connection.pending_requests", "The number of current pending requests for an open connection.", "{request}", false, nil)},
		{"dbconv/ClientConnectionTimeouts", constructorTest(dbconv.NewClientConnectionTimeouts, "db.client.connection.timeouts", "The number of connection timeouts that have occurred trying to obtain a connection from the pool.", "{timeout}", false, nil)},
		{"dbconv/ClientConnectionTimeoutsObservable", constructorTest(dbconv.NewClientConnectionTimeoutsObservable, "db.client.connection.timeouts", "The number of connection timeouts that have occurred trying to obtain a connection from the pool.", "{timeout}", false, nil)},
		{"dbconv/ClientConnectionUseTime", constructorTest(dbconv.NewClientConnectionUseTime, "db.client.connection.use_time", "The time between borrowing a connection and returning it to the pool.", "s", true, nil)},
		{"dbconv/ClientConnectionWaitTime", constructorTest(dbconv.NewClientConnectionWaitTime, "db.client.connection.wait_time", "The time it took to obtain an open connection from the pool.", "s", true, nil)},
		{"dbconv/ClientOperationDuration", constructorTest(dbconv.NewClientOperationDuration, "db.client.operation.duration", "Duration of database client operations.", "s", true, nil)},
		{"dbconv/ClientResponseReturnedRows", constructorTest(dbconv.NewClientResponseReturnedRows, "db.client.response.returned_rows", "The actual number of records returned by the database operation.", "{row}", true, nil)},
		{"dnsconv/LookupDuration", constructorTest(dnsconv.NewLookupDuration, "dns.lookup.duration", "Measures the time taken to perform a DNS lookup.", "s", true, nil)},
		{"faasconv/Coldstarts", constructorTest(faasconv.NewColdstarts, "faas.coldstarts", "Number of invocation cold starts.", "{coldstart}", false, nil)},
		{"faasconv/ColdstartsObservable", constructorTest(faasconv.NewColdstartsObservable, "faas.coldstarts", "Number of invocation cold starts.", "{coldstart}", false, nil)},
		{"faasconv/CPUUsage", constructorTest(faasconv.NewCPUUsage, "faas.cpu_usage", "Distribution of CPU usage per invocation.", "s", true, nil)},
		{"faasconv/Errors", constructorTest(faasconv.NewErrors, "faas.errors", "Number of invocation errors.", "{error}", false, nil)},
		{"faasconv/ErrorsObservable", constructorTest(faasconv.NewErrorsObservable, "faas.errors", "Number of invocation errors.", "{error}", false, nil)},
		{"faasconv/InitDuration", constructorTest(faasconv.NewInitDuration, "faas.init_duration", "Measures the duration of the function\u0027s initialization, such as a cold start.", "s", true, nil)},
		{"faasconv/Invocations", constructorTest(faasconv.NewInvocations, "faas.invocations", "Number of successful invocations.", "{invocation}", false, nil)},
		{"faasconv/InvocationsObservable", constructorTest(faasconv.NewInvocationsObservable, "faas.invocations", "Number of successful invocations.", "{invocation}", false, nil)},
		{"faasconv/InvokeDuration", constructorTest(faasconv.NewInvokeDuration, "faas.invoke_duration", "Measures the duration of the function\u0027s logic execution.", "s", true, nil)},
		{"faasconv/MemUsage", constructorTest(faasconv.NewMemUsage, "faas.mem_usage", "Distribution of max memory usage per invocation.", "By", true, nil)},
		{"faasconv/NetIO", constructorTest(faasconv.NewNetIO, "faas.net_io", "Distribution of net I/O usage per invocation.", "By", true, nil)},
		{"faasconv/Timeouts", constructorTest(faasconv.NewTimeouts, "faas.timeouts", "Number of invocation timeouts.", "{timeout}", false, nil)},
		{"faasconv/TimeoutsObservable", constructorTest(faasconv.NewTimeoutsObservable, "faas.timeouts", "Number of invocation timeouts.", "{timeout}", false, nil)},
		{"goconv/ConfigGogc", constructorTest(goconv.NewConfigGogc, "go.config.gogc", "Heap size target percentage configured by the user, otherwise 100.", "%", false, nil)},
		{"goconv/CPUTime", constructorTest(goconv.NewCPUTime, "go.cpu.time", "Estimated CPU time spent by the Go runtime.", "s", false, nil)},
		{"goconv/CPUTimeObservable", constructorTest(goconv.NewCPUTimeObservable, "go.cpu.time", "Estimated CPU time spent by the Go runtime.", "s", false, nil)},
		{"goconv/GoroutineCount", constructorTest(goconv.NewGoroutineCount, "go.goroutine.count", "Count of live goroutines.", "{goroutine}", false, nil)},
		{"goconv/MemoryAllocated", constructorTest(goconv.NewMemoryAllocated, "go.memory.allocated", "Memory allocated to the heap by the application.", "By", false, nil)},
		{"goconv/MemoryAllocations", constructorTest(goconv.NewMemoryAllocations, "go.memory.allocations", "Count of allocations to the heap by the application.", "{allocation}", false, nil)},
		{"goconv/MemoryGCCycles", constructorTest(goconv.NewMemoryGCCycles, "go.memory.gc.cycles", "Number of completed GC cycles.", "{gc_cycle}", false, nil)},
		{"goconv/MemoryGCCyclesObservable", constructorTest(goconv.NewMemoryGCCyclesObservable, "go.memory.gc.cycles", "Number of completed GC cycles.", "{gc_cycle}", false, nil)},
		{"goconv/MemoryGCGoal", constructorTest(goconv.NewMemoryGCGoal, "go.memory.gc.goal", "Heap size target for the end of the GC cycle.", "By", false, nil)},
		{"goconv/MemoryGCPauseDuration", constructorTest(goconv.NewMemoryGCPauseDuration, "go.memory.gc.pause.duration", "Distribution of individual GC-related stop-the-world pause latencies. This is the time from deciding to stop the world until the world is started again.", "s", true, nil)},
		{"goconv/MemoryLimit", constructorTest(goconv.NewMemoryLimit, "go.memory.limit", "Go runtime memory limit configured by the user, if a limit exists.", "By", false, nil)},
		{"goconv/MemoryUsed", constructorTest(goconv.NewMemoryUsed, "go.memory.used", "Memory used by the Go runtime.", "By", false, nil)},
		{"goconv/ProcessorLimit", constructorTest(goconv.NewProcessorLimit, "go.processor.limit", "The number of OS threads that can execute user-level Go code simultaneously.", "{thread}", false, nil)},
		{"goconv/ScheduleDuration", constructorTest(goconv.NewScheduleDuration, "go.schedule.duration", "The time goroutines have spent in the scheduler in a runnable state before actually running.", "s", true, nil)},
		{"httpconv/ClientActiveRequests", constructorTest(httpconv.NewClientActiveRequests, "http.client.active_requests", "Number of active HTTP requests.", "{request}", false, nil)},
		{"httpconv/ClientActiveRequestsObservable", constructorTest(httpconv.NewClientActiveRequestsObservable, "http.client.active_requests", "Number of active HTTP requests.", "{request}", false, nil)},
		{"httpconv/ClientConnectionDuration", constructorTest(httpconv.NewClientConnectionDuration, "http.client.connection.duration", "The duration of the successfully established outbound HTTP connections.", "s", true, nil)},
		{"httpconv/ClientOpenConnections", constructorTest(httpconv.NewClientOpenConnections, "http.client.open_connections", "Number of outbound HTTP connections that are currently active or idle on the client.", "{connection}", false, nil)},
		{"httpconv/ClientOpenConnectionsObservable", constructorTest(httpconv.NewClientOpenConnectionsObservable, "http.client.open_connections", "Number of outbound HTTP connections that are currently active or idle on the client.", "{connection}", false, nil)},
		{"httpconv/ClientRequestBodySize", constructorTest(httpconv.NewClientRequestBodySize, "http.client.request.body.size", "Size of HTTP client request bodies.", "By", true, nil)},
		{"httpconv/ClientRequestDuration", constructorTest(httpconv.NewClientRequestDuration, "http.client.request.duration", "Duration of HTTP client requests.", "s", true, []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10})},
		{"httpconv/ClientResponseBodySize", constructorTest(httpconv.NewClientResponseBodySize, "http.client.response.body.size", "Size of HTTP client response bodies.", "By", true, nil)},
		{"httpconv/ServerActiveRequests", constructorTest(httpconv.NewServerActiveRequests, "http.server.active_requests", "Number of active HTTP server requests.", "{request}", false, nil)},
		{"httpconv/ServerActiveRequestsObservable", constructorTest(httpconv.NewServerActiveRequestsObservable, "http.server.active_requests", "Number of active HTTP server requests.", "{request}", false, nil)},
		{"httpconv/ServerRequestBodySize", constructorTest(httpconv.NewServerRequestBodySize, "http.server.request.body.size", "Size of HTTP server request bodies.", "By", true, nil)},
		{"httpconv/ServerRequestDuration", constructorTest(httpconv.NewServerRequestDuration, "http.server.request.duration", "Duration of HTTP server requests.", "s", true, []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10})},
		{"httpconv/ServerResponseBodySize", constructorTest(httpconv.NewServerResponseBodySize, "http.server.response.body.size", "Size of HTTP server response bodies.", "By", true, nil)},
		{"hwconv/BatteryCharge", constructorTest(hwconv.NewBatteryCharge, "hw.battery.charge", "Remaining fraction of battery charge.", "1", false, nil)},
		{"hwconv/BatteryChargeObservable", constructorTest(hwconv.NewBatteryChargeObservable, "hw.battery.charge", "Remaining fraction of battery charge.", "1", false, nil)},
		{"hwconv/BatteryChargeLimit", constructorTest(hwconv.NewBatteryChargeLimit, "hw.battery.charge.limit", "Lower limit of battery charge fraction to ensure proper operation.", "1", false, nil)},
		{"hwconv/BatteryChargeLimitObservable", constructorTest(hwconv.NewBatteryChargeLimitObservable, "hw.battery.charge.limit", "Lower limit of battery charge fraction to ensure proper operation.", "1", false, nil)},
		{"hwconv/BatteryTimeLeft", constructorTest(hwconv.NewBatteryTimeLeft, "hw.battery.time_left", "Time left before battery is completely charged or discharged.", "s", false, nil)},
		{"hwconv/BatteryTimeLeftObservable", constructorTest(hwconv.NewBatteryTimeLeftObservable, "hw.battery.time_left", "Time left before battery is completely charged or discharged.", "s", false, nil)},
		{"hwconv/CPUSpeed", constructorTest(hwconv.NewCPUSpeed, "hw.cpu.speed", "CPU current frequency.", "Hz", false, nil)},
		{"hwconv/CPUSpeedObservable", constructorTest(hwconv.NewCPUSpeedObservable, "hw.cpu.speed", "CPU current frequency.", "Hz", false, nil)},
		{"hwconv/CPUSpeedLimit", constructorTest(hwconv.NewCPUSpeedLimit, "hw.cpu.speed.limit", "CPU maximum frequency.", "Hz", false, nil)},
		{"hwconv/CPUSpeedLimitObservable", constructorTest(hwconv.NewCPUSpeedLimitObservable, "hw.cpu.speed.limit", "CPU maximum frequency.", "Hz", false, nil)},
		{"hwconv/Energy", constructorTest(hwconv.NewEnergy, "hw.energy", "Energy consumed by the component.", "J", false, nil)},
		{"hwconv/EnergyObservable", constructorTest(hwconv.NewEnergyObservable, "hw.energy", "Energy consumed by the component.", "J", false, nil)},
		{"hwconv/Errors", constructorTest(hwconv.NewErrors, "hw.errors", "Number of errors encountered by the component.", "{error}", false, nil)},
		{"hwconv/ErrorsObservable", constructorTest(hwconv.NewErrorsObservable, "hw.errors", "Number of errors encountered by the component.", "{error}", false, nil)},
		{"hwconv/FanSpeed", constructorTest(hwconv.NewFanSpeed, "hw.fan.speed", "Fan speed in revolutions per minute.", "rpm", false, nil)},
		{"hwconv/FanSpeedObservable", constructorTest(hwconv.NewFanSpeedObservable, "hw.fan.speed", "Fan speed in revolutions per minute.", "rpm", false, nil)},
		{"hwconv/FanSpeedLimit", constructorTest(hwconv.NewFanSpeedLimit, "hw.fan.speed.limit", "Speed limit in rpm.", "rpm", false, nil)},
		{"hwconv/FanSpeedLimitObservable", constructorTest(hwconv.NewFanSpeedLimitObservable, "hw.fan.speed.limit", "Speed limit in rpm.", "rpm", false, nil)},
		{"hwconv/FanSpeedRatio", constructorTest(hwconv.NewFanSpeedRatio, "hw.fan.speed_ratio", "Fan speed expressed as a fraction of its maximum speed.", "1", false, nil)},
		{"hwconv/FanSpeedRatioObservable", constructorTest(hwconv.NewFanSpeedRatioObservable, "hw.fan.speed_ratio", "Fan speed expressed as a fraction of its maximum speed.", "1", false, nil)},
		{"hwconv/GpuIO", constructorTest(hwconv.NewGpuIO, "hw.gpu.io", "Received and transmitted bytes by the GPU.", "By", false, nil)},
		{"hwconv/GpuIOObservable", constructorTest(hwconv.NewGpuIOObservable, "hw.gpu.io", "Received and transmitted bytes by the GPU.", "By", false, nil)},
		{"hwconv/GpuMemoryLimit", constructorTest(hwconv.NewGpuMemoryLimit, "hw.gpu.memory.limit", "Size of the GPU memory.", "By", false, nil)},
		{"hwconv/GpuMemoryLimitObservable", constructorTest(hwconv.NewGpuMemoryLimitObservable, "hw.gpu.memory.limit", "Size of the GPU memory.", "By", false, nil)},
		{"hwconv/GpuMemoryUsage", constructorTest(hwconv.NewGpuMemoryUsage, "hw.gpu.memory.usage", "GPU memory used.", "By", false, nil)},
		{"hwconv/GpuMemoryUsageObservable", constructorTest(hwconv.NewGpuMemoryUsageObservable, "hw.gpu.memory.usage", "GPU memory used.", "By", false, nil)},
		{"hwconv/GpuMemoryUtilization", constructorTest(hwconv.NewGpuMemoryUtilization, "hw.gpu.memory.utilization", "Fraction of GPU memory used.", "1", false, nil)},
		{"hwconv/GpuMemoryUtilizationObservable", constructorTest(hwconv.NewGpuMemoryUtilizationObservable, "hw.gpu.memory.utilization", "Fraction of GPU memory used.", "1", false, nil)},
		{"hwconv/GpuUtilization", constructorTest(hwconv.NewGpuUtilization, "hw.gpu.utilization", "Fraction of time spent in a specific task.", "1", false, nil)},
		{"hwconv/GpuUtilizationObservable", constructorTest(hwconv.NewGpuUtilizationObservable, "hw.gpu.utilization", "Fraction of time spent in a specific task.", "1", false, nil)},
		{"hwconv/HostAmbientTemperature", constructorTest(hwconv.NewHostAmbientTemperature, "hw.host.ambient_temperature", "Ambient (external) temperature of the physical host.", "Cel", false, nil)},
		{"hwconv/HostAmbientTemperatureObservable", constructorTest(hwconv.NewHostAmbientTemperatureObservable, "hw.host.ambient_temperature", "Ambient (external) temperature of the physical host.", "Cel", false, nil)},
		{"hwconv/HostEnergy", constructorTest(hwconv.NewHostEnergy, "hw.host.energy", "Total energy consumed by the entire physical host, in joules.", "J", false, nil)},
		{"hwconv/HostEnergyObservable", constructorTest(hwconv.NewHostEnergyObservable, "hw.host.energy", "Total energy consumed by the entire physical host, in joules.", "J", false, nil)},
		{"hwconv/HostHeatingMargin", constructorTest(hwconv.NewHostHeatingMargin, "hw.host.heating_margin", "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors.", "Cel", false, nil)},
		{"hwconv/HostHeatingMarginObservable", constructorTest(hwconv.NewHostHeatingMarginObservable, "hw.host.heating_margin", "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors.", "Cel", false, nil)},
		{"hwconv/HostPower", constructorTest(hwconv.NewHostPower, "hw.host.power", "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred).", "W", false, nil)},
		{"hwconv/HostPowerObservable", constructorTest(hwconv.NewHostPowerObservable, "hw.host.power", "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred).", "W", false, nil)},
		{"hwconv/LogicalDiskLimit", constructorTest(hwconv.NewLogicalDiskLimit, "hw.logical_disk.limit", "Size of the logical disk.", "By", false, nil)},
		{"hwconv/LogicalDiskLimitObservable", constructorTest(hwconv.NewLogicalDiskLimitObservable, "hw.logical_disk.limit", "Size of the logical disk.", "By", false, nil)},
		{"hwconv/LogicalDiskUsage", constructorTest(hwconv.NewLogicalDiskUsage, "hw.logical_disk.usage", "Logical disk space usage.", "By", false, nil)},
		{"hwconv/LogicalDiskUsageObservable", constructorTest(hwconv.NewLogicalDiskUsageObservable, "hw.logical_disk.usage", "Logical disk space usage.", "By", false, nil)},
		{"hwconv/LogicalDiskUtilization", constructorTest(hwconv.NewLogicalDiskUtilization, "hw.logical_disk.utilization", "Logical disk space utilization as a fraction.", "1", false, nil)},
		{"hwconv/LogicalDiskUtilizationObservable", constructorTest(hwconv.NewLogicalDiskUtilizationObservable, "hw.logical_disk.utilization", "Logical disk space utilization as a fraction.", "1", false, nil)},
		{"hwconv/MemorySize", constructorTest(hwconv.NewMemorySize, "hw.memory.size", "Size of the memory module.", "By", false, nil)},
		{"hwconv/MemorySizeObservable", constructorTest(hwconv.NewMemorySizeObservable, "hw.memory.size", "Size of the memory module.", "By", false, nil)},
		{"hwconv/NetworkBandwidthLimit", constructorTest(hwconv.NewNetworkBandwidthLimit, "hw.network.bandwidth.limit", "Link speed.", "By/s", false, nil)},
		{"hwconv/NetworkBandwidthLimitObservable", constructorTest(hwconv.NewNetworkBandwidthLimitObservable, "hw.network.bandwidth.limit", "Link speed.", "By/s", false, nil)},
		{"hwconv/NetworkBandwidthUtilization", constructorTest(hwconv.NewNetworkBandwidthUtilization, "hw.network.bandwidth.utilization", "Utilization of the network bandwidth as a fraction.", "1", false, nil)},
		{"hwconv/NetworkBandwidthUtilizationObservable", constructorTest(hwconv.NewNetworkBandwidthUtilizationObservable, "hw.network.bandwidth.utilization", "Utilization of the network bandwidth as a fraction.", "1", false, nil)},
		{"hwconv/NetworkIO", constructorTest(hwconv.NewNetworkIO, "hw.network.io", "Received and transmitted network traffic in bytes.", "By", false, nil)},
		{"hwconv/NetworkIOObservable", constructorTest(hwconv.NewNetworkIOObservable, "hw.network.io", "Received and transmitted network traffic in bytes.", "By", false, nil)},
		{"hwconv/NetworkPackets", constructorTest(hwconv.NewNetworkPackets, "hw.network.packets", "Received and transmitted network traffic in packets (or frames).", "{packet}", false, nil)},
		{"hwconv/NetworkPacketsObservable", constructorTest(hwconv.NewNetworkPacketsObservable, "hw.network.packets", "Received and transmitted network traffic in packets (or frames).", "{packet}", false, nil)},
		{"hwconv/NetworkUp", constructorTest(hwconv.NewNetworkUp, "hw.network.up", "Link status: `1` (up) or `0` (down).", "1", false, nil)},
		{"hwconv/NetworkUpObservable", constructorTest(hwconv.NewNetworkUpObservable, "hw.network.up", "Link status: `1` (up) or `0` (down).", "1", false, nil)},
		{"hwconv/PhysicalDiskEnduranceUtilization", constructorTest(hwconv.NewPhysicalDiskEnduranceUtilization, "hw.physical_disk.endurance_utilization", "Endurance remaining for this SSD disk.", "1", false, nil)},
		{"hwconv/PhysicalDiskEnduranceUtilizationObservable", constructorTest(hwconv.NewPhysicalDiskEnduranceUtilizationObservable, "hw.physical_disk.endurance_utilization", "Endurance remaining for this SSD disk.", "1", false, nil)},
		{"hwconv/PhysicalDiskSize", constructorTest(hwconv.NewPhysicalDiskSize, "hw.physical_disk.size", "Size of the disk.", "By", false, nil)},
		{"hwconv/PhysicalDiskSizeObservable", constructorTest(hwconv.NewPhysicalDiskSizeObservable, "hw.physical_disk.size", "Size of the disk.", "By", false, nil)},
		{"hwconv/PhysicalDiskSmart", constructorTest(hwconv.NewPhysicalDiskSmart, "hw.physical_disk.smart", "Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute.", "1", false, nil)},
		{"hwconv/PhysicalDiskSmartObservable", constructorTest(hwconv.NewPhysicalDiskSmartObservable, "hw.physical_disk.smart", "Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute.", "1", false, nil)},
		{"hwconv/Power", constructorTest(hwconv.NewPower, "hw.power", "Instantaneous power consumed by the component.", "W", false, nil)},
		{"hwconv/PowerObservable", constructorTest(hwconv.NewPowerObservable, "hw.power", "Instantaneous power consumed by the component.", "W", false, nil)},
		{"hwconv/PowerSupplyLimit", constructorTest(hwconv.NewPowerSupplyLimit, "hw.power_supply.limit", "Maximum power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyLimitObservable", constructorTest(hwconv.NewPowerSupplyLimitObservable, "hw.power_supply.limit", "Maximum power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyUsage", constructorTest(hwconv.NewPowerSupplyUsage, "hw.power_supply.usage", "Current power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyUsageObservable", constructorTest(hwconv.NewPowerSupplyUsageObservable, "hw.power_supply.usage", "Current power output of the power supply.", "W", false, nil)},
		{"hwconv/PowerSupplyUtilization", constructorTest(hwconv.NewPowerSupplyUtilization, "hw.power_supply.utilization", "Utilization of the power supply as a fraction of its maximum output.", "1", false, nil)},
		{"hwconv/PowerSupplyUtilizationObservable", constructorTest(hwconv.NewPowerSupplyUtilizationObservable, "hw.power_supply.utilization", "Utilization of the power supply as a fraction of its maximum output.", "1", false, nil)},
		{"hwconv/Status", constructorTest(hwconv.NewStatus, "hw.status", "Operational status: `1` (true) or `0` (false) for each of the possible states.", "1", false, nil)},
		{"hwconv/StatusObservable", constructorTest(hwconv.NewStatusObservable, "hw.status", "Operational status: `1` (true) or `0` (false) for each of the possible states.", "1", false, nil)},
		{"hwconv/TapeDriveOperations", constructorTest(hwconv.NewTapeDriveOperations, "hw.tape_drive.operations", "Operations performed by the tape drive.", "{operation}", false, nil)},
		{"hwconv/TapeDriveOperationsObservable", constructorTest(hwconv.NewTapeDriveOperationsObservable, "hw.tape_drive.operations", "Operations performed by the tape drive.", "{operation}", false, nil)},
		{"hwconv/Temperature", constructorTest(hwconv.NewTemperature, "hw.temperature", "Temperature in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/TemperatureObservable", constructorTest(hwconv.NewTemperatureObservable, "hw.temperature", "Temperature in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/TemperatureLimit", constructorTest(hwconv.NewTemperatureLimit, "hw.temperature.limit", "Temperature limit in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/TemperatureLimitObservable", constructorTest(hwconv.NewTemperatureLimitObservable, "hw.temperature.limit", "Temperature limit in degrees Celsius.", "Cel", false, nil)},
		{"hwconv/Voltage", constructorTest(hwconv.NewVoltage, "hw.voltage", "Voltage measured by the sensor.", "V", false, nil)},
		{"hwconv/VoltageObservable", constructorTest(hwconv.NewVoltageObservable, "hw.voltage", "Voltage measured by the sensor.", "V", false, nil)},
		{"hwconv/VoltageLimit", constructorTest(hwconv.NewVoltageLimit, "hw.voltage.limit", "Voltage limit in Volts.", "V", false, nil)},
		{"hwconv/VoltageLimitObservable", constructorTest(hwconv.NewVoltageLimitObservable, "hw.voltage.limit", "Voltage limit in Volts.", "V", false, nil)},
		{"hwconv/VoltageNominal", constructorTest(hwconv.NewVoltageNominal, "hw.voltage.nominal", "Nominal (expected) voltage.", "V", false, nil)},
		{"hwconv/VoltageNominalObservable", constructorTest(hwconv.NewVoltageNominalObservable, "hw.voltage.nominal", "Nominal (expected) voltage.", "V", false, nil)},
		{"k8sconv/ContainerCPULimitCurrent", constructorTest(k8sconv.NewContainerCPULimitCurrent, "k8s.container.cpu.limit.current", "Maximum CPU resource limit currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitCurrentObservable", constructorTest(k8sconv.NewContainerCPULimitCurrentObservable, "k8s.container.cpu.limit.current", "Maximum CPU resource limit currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitDesired", constructorTest(k8sconv.NewContainerCPULimitDesired, "k8s.container.cpu.limit.desired", "Maximum CPU resource limit as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitDesiredObservable", constructorTest(k8sconv.NewContainerCPULimitDesiredObservable, "k8s.container.cpu.limit.desired", "Maximum CPU resource limit as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPULimitUtilization", constructorTest(k8sconv.NewContainerCPULimitUtilization, "k8s.container.cpu.limit.utilization", "The ratio of container CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/ContainerCPULimitUtilizationObservable", constructorTest(k8sconv.NewContainerCPULimitUtilizationObservable, "k8s.container.cpu.limit.utilization", "The ratio of container CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/ContainerCPURequestCurrent", constructorTest(k8sconv.NewContainerCPURequestCurrent, "k8s.container.cpu.request.current", "CPU resource requested currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestCurrentObservable", constructorTest(k8sconv.NewContainerCPURequestCurrentObservable, "k8s.container.cpu.request.current", "CPU resource requested currently configured for a running container.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestDesired", constructorTest(k8sconv.NewContainerCPURequestDesired, "k8s.container.cpu.request.desired", "CPU resource requested as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestDesiredObservable", constructorTest(k8sconv.NewContainerCPURequestDesiredObservable, "k8s.container.cpu.request.desired", "CPU resource requested as defined by the container spec.", "{cpu}", false, nil)},
		{"k8sconv/ContainerCPURequestUtilization", constructorTest(k8sconv.NewContainerCPURequestUtilization, "k8s.container.cpu.request.utilization", "The ratio of container CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/ContainerCPURequestUtilizationObservable", constructorTest(k8sconv.NewContainerCPURequestUtilizationObservable, "k8s.container.cpu.request.utilization", "The ratio of container CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/ContainerEphemeralStorageLimit", constructorTest(k8sconv.NewContainerEphemeralStorageLimit, "k8s.container.ephemeral_storage.limit", "Maximum ephemeral storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageLimitObservable", constructorTest(k8sconv.NewContainerEphemeralStorageLimitObservable, "k8s.container.ephemeral_storage.limit", "Maximum ephemeral storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageRequest", constructorTest(k8sconv.NewContainerEphemeralStorageRequest, "k8s.container.ephemeral_storage.request", "Ephemeral storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageRequestObservable", constructorTest(k8sconv.NewContainerEphemeralStorageRequestObservable, "k8s.container.ephemeral_storage.request", "Ephemeral storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageUsage", constructorTest(k8sconv.NewContainerEphemeralStorageUsage, "k8s.container.ephemeral_storage.usage", "The ephemeral storage used by a container.", "By", false, nil)},
		{"k8sconv/ContainerEphemeralStorageUsageObservable", constructorTest(k8sconv.NewContainerEphemeralStorageUsageObservable, "k8s.container.ephemeral_storage.usage", "The ephemeral storage used by a container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitCurrent", constructorTest(k8sconv.NewContainerMemoryLimitCurrent, "k8s.container.memory.limit.current", "Maximum memory resource limit currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitCurrentObservable", constructorTest(k8sconv.NewContainerMemoryLimitCurrentObservable, "k8s.container.memory.limit.current", "Maximum memory resource limit currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitDesired", constructorTest(k8sconv.NewContainerMemoryLimitDesired, "k8s.container.memory.limit.desired", "Maximum memory resource limit as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitDesiredObservable", constructorTest(k8sconv.NewContainerMemoryLimitDesiredObservable, "k8s.container.memory.limit.desired", "Maximum memory resource limit as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryLimitUtilization", constructorTest(k8sconv.NewContainerMemoryLimitUtilization, "k8s.container.memory.limit.utilization", "The ratio of container memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/ContainerMemoryLimitUtilizationObservable", constructorTest(k8sconv.NewContainerMemoryLimitUtilizationObservable, "k8s.container.memory.limit.utilization", "The ratio of container memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/ContainerMemoryRequestCurrent", constructorTest(k8sconv.NewContainerMemoryRequestCurrent, "k8s.container.memory.request.current", "Memory resource request currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestCurrentObservable", constructorTest(k8sconv.NewContainerMemoryRequestCurrentObservable, "k8s.container.memory.request.current", "Memory resource request currently configured for a running container.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestDesired", constructorTest(k8sconv.NewContainerMemoryRequestDesired, "k8s.container.memory.request.desired", "Memory resource requested as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestDesiredObservable", constructorTest(k8sconv.NewContainerMemoryRequestDesiredObservable, "k8s.container.memory.request.desired", "Memory resource requested as defined by the container spec.", "By", false, nil)},
		{"k8sconv/ContainerMemoryRequestUtilization", constructorTest(k8sconv.NewContainerMemoryRequestUtilization, "k8s.container.memory.request.utilization", "The ratio of container memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/ContainerMemoryRequestUtilizationObservable", constructorTest(k8sconv.NewContainerMemoryRequestUtilizationObservable, "k8s.container.memory.request.utilization", "The ratio of container memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/ContainerReady", constructorTest(k8sconv.NewContainerReady, "k8s.container.ready", "Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready).", "{container}", false, nil)},
		{"k8sconv/ContainerReadyObservable", constructorTest(k8sconv.NewContainerReadyObservable, "k8s.container.ready", "Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready).", "{container}", false, nil)},
		{"k8sconv/ContainerRestartCount", constructorTest(k8sconv.NewContainerRestartCount, "k8s.container.restart.count", "Describes how many times the container has restarted (since the last counter reset).", "{restart}", false, nil)},
		{"k8sconv/ContainerRestartCountObservable", constructorTest(k8sconv.NewContainerRestartCountObservable, "k8s.container.restart.count", "Describes how many times the container has restarted (since the last counter reset).", "{restart}", false, nil)},
		{"k8sconv/ContainerStatusReason", constructorTest(k8sconv.NewContainerStatusReason, "k8s.container.status.reason", "Describes the number of K8s containers that are currently in a state for a given reason.", "{container}", false, nil)},
		{"k8sconv/ContainerStatusReasonObservable", constructorTest(k8sconv.NewContainerStatusReasonObservable, "k8s.container.status.reason", "Describes the number of K8s containers that are currently in a state for a given reason.", "{container}", false, nil)},
		{"k8sconv/ContainerStatusState", constructorTest(k8sconv.NewContainerStatusState, "k8s.container.status.state", "Describes the number of K8s containers that are currently in a given state.", "{container}", false, nil)},
		{"k8sconv/ContainerStatusStateObservable", constructorTest(k8sconv.NewContainerStatusStateObservable, "k8s.container.status.state", "Describes the number of K8s containers that are currently in a given state.", "{container}", false, nil)},
		{"k8sconv/ContainerStorageLimit", constructorTest(k8sconv.NewContainerStorageLimit, "k8s.container.storage.limit", "Maximum storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerStorageLimitObservable", constructorTest(k8sconv.NewContainerStorageLimitObservable, "k8s.container.storage.limit", "Maximum storage resource limit set for the container.", "By", false, nil)},
		{"k8sconv/ContainerStorageRequest", constructorTest(k8sconv.NewContainerStorageRequest, "k8s.container.storage.request", "Storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/ContainerStorageRequestObservable", constructorTest(k8sconv.NewContainerStorageRequestObservable, "k8s.container.storage.request", "Storage resource requested for the container.", "By", false, nil)},
		{"k8sconv/CronJobJobActive", constructorTest(k8sconv.NewCronJobJobActive, "k8s.cronjob.job.active", "The number of actively running jobs for a cronjob.", "{job}", false, nil)},
		{"k8sconv/CronJobJobActiveObservable", constructorTest(k8sconv.NewCronJobJobActiveObservable, "k8s.cronjob.job.active", "The number of actively running jobs for a cronjob.", "{job}", false, nil)},
		{"k8sconv/DaemonSetNodeCurrentScheduled", constructorTest(k8sconv.NewDaemonSetNodeCurrentScheduled, "k8s.daemonset.node.current_scheduled", "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeCurrentScheduledObservable", constructorTest(k8sconv.NewDaemonSetNodeCurrentScheduledObservable, "k8s.daemonset.node.current_scheduled", "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeDesiredScheduled", constructorTest(k8sconv.NewDaemonSetNodeDesiredScheduled, "k8s.daemonset.node.desired_scheduled", "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod).", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeDesiredScheduledObservable", constructorTest(k8sconv.NewDaemonSetNodeDesiredScheduledObservable, "k8s.daemonset.node.desired_scheduled", "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod).", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeMisscheduled", constructorTest(k8sconv.NewDaemonSetNodeMisscheduled, "k8s.daemonset.node.misscheduled", "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeMisscheduledObservable", constructorTest(k8sconv.NewDaemonSetNodeMisscheduledObservable, "k8s.daemonset.node.misscheduled", "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeReady", constructorTest(k8sconv.NewDaemonSetNodeReady, "k8s.daemonset.node.ready", "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.", "{node}", false, nil)},
		{"k8sconv/DaemonSetNodeReadyObservable", constructorTest(k8sconv.NewDaemonSetNodeReadyObservable, "k8s.daemonset.node.ready", "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.", "{node}", false, nil)},
		{"k8sconv/DeploymentPodAvailable", constructorTest(k8sconv.NewDeploymentPodAvailable, "k8s.deployment.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment.", "{pod}", false, nil)},
		{"k8sconv/DeploymentPodAvailableObservable", constructorTest(k8sconv.NewDeploymentPodAvailableObservable, "k8s.deployment.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment.", "{pod}", false, nil)},
		{"k8sconv/DeploymentPodDesired", constructorTest(k8sconv.NewDeploymentPodDesired, "k8s.deployment.pod.desired", "Number of desired replica pods in this deployment.", "{pod}", false, nil)},
		{"k8sconv/DeploymentPodDesiredObservable", constructorTest(k8sconv.NewDeploymentPodDesiredObservable, "k8s.deployment.pod.desired", "Number of desired replica pods in this deployment.", "{pod}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageUtilization", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageUtilization, "k8s.hpa.metric.target.cpu.average_utilization", "Target average utilization, in percentage, for CPU resource in HPA config.", "1", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageUtilizationObservable", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageUtilizationObservable, "k8s.hpa.metric.target.cpu.average_utilization", "Target average utilization, in percentage, for CPU resource in HPA config.", "1", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageValue", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageValue, "k8s.hpa.metric.target.cpu.average_value", "Target average value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUAverageValueObservable", constructorTest(k8sconv.NewHPAMetricTargetCPUAverageValueObservable, "k8s.hpa.metric.target.cpu.average_value", "Target average value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUValue", constructorTest(k8sconv.NewHPAMetricTargetCPUValue, "k8s.hpa.metric.target.cpu.value", "Target value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAMetricTargetCPUValueObservable", constructorTest(k8sconv.NewHPAMetricTargetCPUValueObservable, "k8s.hpa.metric.target.cpu.value", "Target value for CPU resource in HPA config.", "{cpu}", false, nil)},
		{"k8sconv/HPAPodCurrent", constructorTest(k8sconv.NewHPAPodCurrent, "k8s.hpa.pod.current", "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodCurrentObservable", constructorTest(k8sconv.NewHPAPodCurrentObservable, "k8s.hpa.pod.current", "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodDesired", constructorTest(k8sconv.NewHPAPodDesired, "k8s.hpa.pod.desired", "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodDesiredObservable", constructorTest(k8sconv.NewHPAPodDesiredObservable, "k8s.hpa.pod.desired", "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMax", constructorTest(k8sconv.NewHPAPodMax, "k8s.hpa.pod.max", "The upper limit for the number of replica pods to which the autoscaler can scale up.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMaxObservable", constructorTest(k8sconv.NewHPAPodMaxObservable, "k8s.hpa.pod.max", "The upper limit for the number of replica pods to which the autoscaler can scale up.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMin", constructorTest(k8sconv.NewHPAPodMin, "k8s.hpa.pod.min", "The lower limit for the number of replica pods to which the autoscaler can scale down.", "{pod}", false, nil)},
		{"k8sconv/HPAPodMinObservable", constructorTest(k8sconv.NewHPAPodMinObservable, "k8s.hpa.pod.min", "The lower limit for the number of replica pods to which the autoscaler can scale down.", "{pod}", false, nil)},
		{"k8sconv/JobPodActive", constructorTest(k8sconv.NewJobPodActive, "k8s.job.pod.active", "The number of pending and actively running pods for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodActiveObservable", constructorTest(k8sconv.NewJobPodActiveObservable, "k8s.job.pod.active", "The number of pending and actively running pods for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodDesiredSuccessful", constructorTest(k8sconv.NewJobPodDesiredSuccessful, "k8s.job.pod.desired_successful", "The desired number of successfully finished pods the job should be run with.", "{pod}", false, nil)},
		{"k8sconv/JobPodDesiredSuccessfulObservable", constructorTest(k8sconv.NewJobPodDesiredSuccessfulObservable, "k8s.job.pod.desired_successful", "The desired number of successfully finished pods the job should be run with.", "{pod}", false, nil)},
		{"k8sconv/JobPodFailed", constructorTest(k8sconv.NewJobPodFailed, "k8s.job.pod.failed", "The number of pods which reached phase Failed for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodFailedObservable", constructorTest(k8sconv.NewJobPodFailedObservable, "k8s.job.pod.failed", "The number of pods which reached phase Failed for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodMaxParallel", constructorTest(k8sconv.NewJobPodMaxParallel, "k8s.job.pod.max_parallel", "The max desired number of pods the job should run at any given time.", "{pod}", false, nil)},
		{"k8sconv/JobPodMaxParallelObservable", constructorTest(k8sconv.NewJobPodMaxParallelObservable, "k8s.job.pod.max_parallel", "The max desired number of pods the job should run at any given time.", "{pod}", false, nil)},
		{"k8sconv/JobPodSuccessful", constructorTest(k8sconv.NewJobPodSuccessful, "k8s.job.pod.successful", "The number of pods which reached phase Succeeded for a job.", "{pod}", false, nil)},
		{"k8sconv/JobPodSuccessfulObservable", constructorTest(k8sconv.NewJobPodSuccessfulObservable, "k8s.job.pod.successful", "The number of pods which reached phase Succeeded for a job.", "{pod}", false, nil)},
		{"k8sconv/NamespacePhase", constructorTest(k8sconv.NewNamespacePhase, "k8s.namespace.phase", "Describes number of K8s namespaces that are currently in a given phase.", "{namespace}", false, nil)},
		{"k8sconv/NamespacePhaseObservable", constructorTest(k8sconv.NewNamespacePhaseObservable, "k8s.namespace.phase", "Describes number of K8s namespaces that are currently in a given phase.", "{namespace}", false, nil)},
		{"k8sconv/NodeConditionStatus", constructorTest(k8sconv.NewNodeConditionStatus, "k8s.node.condition.status", "Describes the condition of a particular Node.", "{node}", false, nil)},
		{"k8sconv/NodeConditionStatusObservable", constructorTest(k8sconv.NewNodeConditionStatusObservable, "k8s.node.condition.status", "Describes the condition of a particular Node.", "{node}", false, nil)},
		{"k8sconv/NodeCPUAllocatable", constructorTest(k8sconv.NewNodeCPUAllocatable, "k8s.node.cpu.allocatable", "Amount of CPU allocatable on the node.", "{cpu}", false, nil)},
		{"k8sconv/NodeCPUAllocatableObservable", constructorTest(k8sconv.NewNodeCPUAllocatableObservable, "k8s.node.cpu.allocatable", "Amount of CPU allocatable on the node.", "{cpu}", false, nil)},
		{"k8sconv/NodeCPUTime", constructorTest(k8sconv.NewNodeCPUTime, "k8s.node.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/NodeCPUTimeObservable", constructorTest(k8sconv.NewNodeCPUTimeObservable, "k8s.node.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/NodeCPUUsage", constructorTest(k8sconv.NewNodeCPUUsage, "k8s.node.cpu.usage", "Node\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeCPUUsageObservable", constructorTest(k8sconv.NewNodeCPUUsageObservable, "k8s.node.cpu.usage", "Node\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeEphemeralStorageAllocatable", constructorTest(k8sconv.NewNodeEphemeralStorageAllocatable, "k8s.node.ephemeral_storage.allocatable", "Amount of ephemeral-storage allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeEphemeralStorageAllocatableObservable", constructorTest(k8sconv.NewNodeEphemeralStorageAllocatableObservable, "k8s.node.ephemeral_storage.allocatable", "Amount of ephemeral-storage allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeFilesystemAvailable", constructorTest(k8sconv.NewNodeFilesystemAvailable, "k8s.node.filesystem.available", "Node filesystem available bytes.", "By", false, nil)},
		{"k8sconv/NodeFilesystemAvailableObservable", constructorTest(k8sconv.NewNodeFilesystemAvailableObservable, "k8s.node.filesystem.available", "Node filesystem available bytes.", "By", false, nil)},
		{"k8sconv/NodeFilesystemCapacity", constructorTest(k8sconv.NewNodeFilesystemCapacity, "k8s.node.filesystem.capacity", "Node filesystem capacity.", "By", false, nil)},
		{"k8sconv/NodeFilesystemCapacityObservable", constructorTest(k8sconv.NewNodeFilesystemCapacityObservable, "k8s.node.filesystem.capacity", "Node filesystem capacity.", "By", false, nil)},
		{"k8sconv/NodeFilesystemInodeCount", constructorTest(k8sconv.NewNodeFilesystemInodeCount, "k8s.node.filesystem.inode.count", "The total inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemInodeCountObservable", constructorTest(k8sconv.NewNodeFilesystemInodeCountObservable, "k8s.node.filesystem.inode.count", "The total inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemInodeFree", constructorTest(k8sconv.NewNodeFilesystemInodeFree, "k8s.node.filesystem.inode.free", "The free inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemInodeFreeObservable", constructorTest(k8sconv.NewNodeFilesystemInodeFreeObservable, "k8s.node.filesystem.inode.free", "The free inodes in the node\u0027s root filesystem.", "{inode}", false, nil)},
		{"k8sconv/NodeFilesystemUsage", constructorTest(k8sconv.NewNodeFilesystemUsage, "k8s.node.filesystem.usage", "Node filesystem usage.", "By", false, nil)},
		{"k8sconv/NodeFilesystemUsageObservable", constructorTest(k8sconv.NewNodeFilesystemUsageObservable, "k8s.node.filesystem.usage", "Node filesystem usage.", "By", false, nil)},
		{"k8sconv/NodeMemoryAllocatable", constructorTest(k8sconv.NewNodeMemoryAllocatable, "k8s.node.memory.allocatable", "Amount of memory allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeMemoryAllocatableObservable", constructorTest(k8sconv.NewNodeMemoryAllocatableObservable, "k8s.node.memory.allocatable", "Amount of memory allocatable on the node.", "By", false, nil)},
		{"k8sconv/NodeMemoryAvailable", constructorTest(k8sconv.NewNodeMemoryAvailable, "k8s.node.memory.available", "Node memory available.", "By", false, nil)},
		{"k8sconv/NodeMemoryAvailableObservable", constructorTest(k8sconv.NewNodeMemoryAvailableObservable, "k8s.node.memory.available", "Node memory available.", "By", false, nil)},
		{"k8sconv/NodeMemoryRss", constructorTest(k8sconv.NewNodeMemoryRss, "k8s.node.memory.rss", "Node memory RSS.", "By", false, nil)},
		{"k8sconv/NodeMemoryRssObservable", constructorTest(k8sconv.NewNodeMemoryRssObservable, "k8s.node.memory.rss", "Node memory RSS.", "By", false, nil)},
		{"k8sconv/NodeMemoryUsage", constructorTest(k8sconv.NewNodeMemoryUsage, "k8s.node.memory.usage", "Memory usage of the Node.", "By", false, nil)},
		{"k8sconv/NodeMemoryUsageObservable", constructorTest(k8sconv.NewNodeMemoryUsageObservable, "k8s.node.memory.usage", "Memory usage of the Node.", "By", false, nil)},
		{"k8sconv/NodeMemoryWorkingSet", constructorTest(k8sconv.NewNodeMemoryWorkingSet, "k8s.node.memory.working_set", "Node memory working set.", "By", false, nil)},
		{"k8sconv/NodeMemoryWorkingSetObservable", constructorTest(k8sconv.NewNodeMemoryWorkingSetObservable, "k8s.node.memory.working_set", "Node memory working set.", "By", false, nil)},
		{"k8sconv/NodeNetworkErrors", constructorTest(k8sconv.NewNodeNetworkErrors, "k8s.node.network.errors", "Node network errors.", "{error}", false, nil)},
		{"k8sconv/NodeNetworkErrorsObservable", constructorTest(k8sconv.NewNodeNetworkErrorsObservable, "k8s.node.network.errors", "Node network errors.", "{error}", false, nil)},
		{"k8sconv/NodeNetworkIO", constructorTest(k8sconv.NewNodeNetworkIO, "k8s.node.network.io", "Network bytes for the Node.", "By", false, nil)},
		{"k8sconv/NodeNetworkIOObservable", constructorTest(k8sconv.NewNodeNetworkIOObservable, "k8s.node.network.io", "Network bytes for the Node.", "By", false, nil)},
		{"k8sconv/NodePagingFaults", constructorTest(k8sconv.NewNodePagingFaults, "k8s.node.paging.faults", "Node memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/NodePagingFaultsObservable", constructorTest(k8sconv.NewNodePagingFaultsObservable, "k8s.node.paging.faults", "Node memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/NodePIDLimit", constructorTest(k8sconv.NewNodePIDLimit, "k8s.node.pid.limit", "Total number of processes/threads allowed by the operating system on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePIDLimitObservable", constructorTest(k8sconv.NewNodePIDLimitObservable, "k8s.node.pid.limit", "Total number of processes/threads allowed by the operating system on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePIDUsage", constructorTest(k8sconv.NewNodePIDUsage, "k8s.node.pid.usage", "Total number of existing processes and threads on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePIDUsageObservable", constructorTest(k8sconv.NewNodePIDUsageObservable, "k8s.node.pid.usage", "Total number of existing processes and threads on the node.", "{thread}", false, nil)},
		{"k8sconv/NodePodAllocatable", constructorTest(k8sconv.NewNodePodAllocatable, "k8s.node.pod.allocatable", "Amount of pods allocatable on the node.", "{pod}", false, nil)},
		{"k8sconv/NodePodAllocatableObservable", constructorTest(k8sconv.NewNodePodAllocatableObservable, "k8s.node.pod.allocatable", "Amount of pods allocatable on the node.", "{pod}", false, nil)},
		{"k8sconv/NodeSystemContainerCPUTime", constructorTest(k8sconv.NewNodeSystemContainerCPUTime, "k8s.node.system_container.cpu.time", "Node\u0027s system container CPU time.", "s", false, nil)},
		{"k8sconv/NodeSystemContainerCPUTimeObservable", constructorTest(k8sconv.NewNodeSystemContainerCPUTimeObservable, "k8s.node.system_container.cpu.time", "Node\u0027s system container CPU time.", "s", false, nil)},
		{"k8sconv/NodeSystemContainerCPUUsage", constructorTest(k8sconv.NewNodeSystemContainerCPUUsage, "k8s.node.system_container.cpu.usage", "Node\u0027s system container CPU usage, measured in CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeSystemContainerCPUUsageObservable", constructorTest(k8sconv.NewNodeSystemContainerCPUUsageObservable, "k8s.node.system_container.cpu.usage", "Node\u0027s system container CPU usage, measured in CPUs.", "{cpu}", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryUsage", constructorTest(k8sconv.NewNodeSystemContainerMemoryUsage, "k8s.node.system_container.memory.usage", "Node\u0027s system container memory usage.", "By", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryUsageObservable", constructorTest(k8sconv.NewNodeSystemContainerMemoryUsageObservable, "k8s.node.system_container.memory.usage", "Node\u0027s system container memory usage.", "By", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryWorkingSet", constructorTest(k8sconv.NewNodeSystemContainerMemoryWorkingSet, "k8s.node.system_container.memory.working_set", "The amount of working set memory.", "By", false, nil)},
		{"k8sconv/NodeSystemContainerMemoryWorkingSetObservable", constructorTest(k8sconv.NewNodeSystemContainerMemoryWorkingSetObservable, "k8s.node.system_container.memory.working_set", "The amount of working set memory.", "By", false, nil)},
		{"k8sconv/NodeUptime", constructorTest(k8sconv.NewNodeUptime, "k8s.node.uptime", "The time the Node has been running.", "s", false, nil)},
		{"k8sconv/NodeUptimeObservable", constructorTest(k8sconv.NewNodeUptimeObservable, "k8s.node.uptime", "The time the Node has been running.", "s", false, nil)},
		{"k8sconv/PersistentvolumeStatusPhase", constructorTest(k8sconv.NewPersistentvolumeStatusPhase, "k8s.persistentvolume.status.phase", "Number of PersistentVolumes in a given phase.", "{persistentvolume}", false, nil)},
		{"k8sconv/PersistentvolumeStatusPhaseObservable", constructorTest(k8sconv.NewPersistentvolumeStatusPhaseObservable, "k8s.persistentvolume.status.phase", "Number of PersistentVolumes in a given phase.", "{persistentvolume}", false, nil)},
		{"k8sconv/PersistentvolumeStorageCapacity", constructorTest(k8sconv.NewPersistentvolumeStorageCapacity, "k8s.persistentvolume.storage.capacity", "The storage capacity of the PersistentVolume.", "By", false, nil)},
		{"k8sconv/PersistentvolumeStorageCapacityObservable", constructorTest(k8sconv.NewPersistentvolumeStorageCapacityObservable, "k8s.persistentvolume.storage.capacity", "The storage capacity of the PersistentVolume.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStatusPhase", constructorTest(k8sconv.NewPersistentvolumeclaimStatusPhase, "k8s.persistentvolumeclaim.status.phase", "Number of PersistentVolumeClaims in a given phase.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/PersistentvolumeclaimStatusPhaseObservable", constructorTest(k8sconv.NewPersistentvolumeclaimStatusPhaseObservable, "k8s.persistentvolumeclaim.status.phase", "Number of PersistentVolumeClaims in a given phase.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageCapacity", constructorTest(k8sconv.NewPersistentvolumeclaimStorageCapacity, "k8s.persistentvolumeclaim.storage.capacity", "The actual storage capacity provisioned for the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageCapacityObservable", constructorTest(k8sconv.NewPersistentvolumeclaimStorageCapacityObservable, "k8s.persistentvolumeclaim.storage.capacity", "The actual storage capacity provisioned for the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageRequest", constructorTest(k8sconv.NewPersistentvolumeclaimStorageRequest, "k8s.persistentvolumeclaim.storage.request", "The storage requested by the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PersistentvolumeclaimStorageRequestObservable", constructorTest(k8sconv.NewPersistentvolumeclaimStorageRequestObservable, "k8s.persistentvolumeclaim.storage.request", "The storage requested by the PersistentVolumeClaim.", "By", false, nil)},
		{"k8sconv/PodCPULimitCurrent", constructorTest(k8sconv.NewPodCPULimitCurrent, "k8s.pod.cpu.limit.current", "Maximum CPU resource limit currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitCurrentObservable", constructorTest(k8sconv.NewPodCPULimitCurrentObservable, "k8s.pod.cpu.limit.current", "Maximum CPU resource limit currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitDesired", constructorTest(k8sconv.NewPodCPULimitDesired, "k8s.pod.cpu.limit.desired", "Maximum CPU resource limit as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitDesiredObservable", constructorTest(k8sconv.NewPodCPULimitDesiredObservable, "k8s.pod.cpu.limit.desired", "Maximum CPU resource limit as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPULimitUtilization", constructorTest(k8sconv.NewPodCPULimitUtilization, "k8s.pod.cpu.limit.utilization", "The ratio of Pod CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/PodCPULimitUtilizationObservable", constructorTest(k8sconv.NewPodCPULimitUtilizationObservable, "k8s.pod.cpu.limit.utilization", "The ratio of Pod CPU usage to its current CPU limit.", "1", false, nil)},
		{"k8sconv/PodCPURequestCurrent", constructorTest(k8sconv.NewPodCPURequestCurrent, "k8s.pod.cpu.request.current", "CPU resource requested currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestCurrentObservable", constructorTest(k8sconv.NewPodCPURequestCurrentObservable, "k8s.pod.cpu.request.current", "CPU resource requested currently configured for a running Pod.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestDesired", constructorTest(k8sconv.NewPodCPURequestDesired, "k8s.pod.cpu.request.desired", "CPU resource requested as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestDesiredObservable", constructorTest(k8sconv.NewPodCPURequestDesiredObservable, "k8s.pod.cpu.request.desired", "CPU resource requested as defined by the Pod spec.", "{cpu}", false, nil)},
		{"k8sconv/PodCPURequestUtilization", constructorTest(k8sconv.NewPodCPURequestUtilization, "k8s.pod.cpu.request.utilization", "The ratio of Pod CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/PodCPURequestUtilizationObservable", constructorTest(k8sconv.NewPodCPURequestUtilizationObservable, "k8s.pod.cpu.request.utilization", "The ratio of Pod CPU usage to its current CPU request.", "1", false, nil)},
		{"k8sconv/PodCPUTime", constructorTest(k8sconv.NewPodCPUTime, "k8s.pod.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/PodCPUTimeObservable", constructorTest(k8sconv.NewPodCPUTimeObservable, "k8s.pod.cpu.time", "Total CPU time consumed.", "s", false, nil)},
		{"k8sconv/PodCPUUsage", constructorTest(k8sconv.NewPodCPUUsage, "k8s.pod.cpu.usage", "Pod\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/PodCPUUsageObservable", constructorTest(k8sconv.NewPodCPUUsageObservable, "k8s.pod.cpu.usage", "Pod\u0027s CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.", "{cpu}", false, nil)},
		{"k8sconv/PodFilesystemAvailable", constructorTest(k8sconv.NewPodFilesystemAvailable, "k8s.pod.filesystem.available", "Pod filesystem available bytes.", "By", false, nil)},
		{"k8sconv/PodFilesystemAvailableObservable", constructorTest(k8sconv.NewPodFilesystemAvailableObservable, "k8s.pod.filesystem.available", "Pod filesystem available bytes.", "By", false, nil)},
		{"k8sconv/PodFilesystemCapacity", constructorTest(k8sconv.NewPodFilesystemCapacity, "k8s.pod.filesystem.capacity", "Pod filesystem capacity.", "By", false, nil)},
		{"k8sconv/PodFilesystemCapacityObservable", constructorTest(k8sconv.NewPodFilesystemCapacityObservable, "k8s.pod.filesystem.capacity", "Pod filesystem capacity.", "By", false, nil)},
		{"k8sconv/PodFilesystemUsage", constructorTest(k8sconv.NewPodFilesystemUsage, "k8s.pod.filesystem.usage", "Pod filesystem usage.", "By", false, nil)},
		{"k8sconv/PodFilesystemUsageObservable", constructorTest(k8sconv.NewPodFilesystemUsageObservable, "k8s.pod.filesystem.usage", "Pod filesystem usage.", "By", false, nil)},
		{"k8sconv/PodMemoryAvailable", constructorTest(k8sconv.NewPodMemoryAvailable, "k8s.pod.memory.available", "Pod memory available.", "By", false, nil)},
		{"k8sconv/PodMemoryAvailableObservable", constructorTest(k8sconv.NewPodMemoryAvailableObservable, "k8s.pod.memory.available", "Pod memory available.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitCurrent", constructorTest(k8sconv.NewPodMemoryLimitCurrent, "k8s.pod.memory.limit.current", "Maximum memory resource limit currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitCurrentObservable", constructorTest(k8sconv.NewPodMemoryLimitCurrentObservable, "k8s.pod.memory.limit.current", "Maximum memory resource limit currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitDesired", constructorTest(k8sconv.NewPodMemoryLimitDesired, "k8s.pod.memory.limit.desired", "Maximum memory resource limit as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitDesiredObservable", constructorTest(k8sconv.NewPodMemoryLimitDesiredObservable, "k8s.pod.memory.limit.desired", "Maximum memory resource limit as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryLimitUtilization", constructorTest(k8sconv.NewPodMemoryLimitUtilization, "k8s.pod.memory.limit.utilization", "The ratio of Pod memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/PodMemoryLimitUtilizationObservable", constructorTest(k8sconv.NewPodMemoryLimitUtilizationObservable, "k8s.pod.memory.limit.utilization", "The ratio of Pod memory working set to its current memory limit.", "1", false, nil)},
		{"k8sconv/PodMemoryRequestCurrent", constructorTest(k8sconv.NewPodMemoryRequestCurrent, "k8s.pod.memory.request.current", "Memory resource request currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestCurrentObservable", constructorTest(k8sconv.NewPodMemoryRequestCurrentObservable, "k8s.pod.memory.request.current", "Memory resource request currently configured for a running Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestDesired", constructorTest(k8sconv.NewPodMemoryRequestDesired, "k8s.pod.memory.request.desired", "Memory resource requested as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestDesiredObservable", constructorTest(k8sconv.NewPodMemoryRequestDesiredObservable, "k8s.pod.memory.request.desired", "Memory resource requested as defined by the Pod spec.", "By", false, nil)},
		{"k8sconv/PodMemoryRequestUtilization", constructorTest(k8sconv.NewPodMemoryRequestUtilization, "k8s.pod.memory.request.utilization", "The ratio of Pod memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/PodMemoryRequestUtilizationObservable", constructorTest(k8sconv.NewPodMemoryRequestUtilizationObservable, "k8s.pod.memory.request.utilization", "The ratio of Pod memory working set to its current memory request.", "1", false, nil)},
		{"k8sconv/PodMemoryRss", constructorTest(k8sconv.NewPodMemoryRss, "k8s.pod.memory.rss", "Pod memory RSS.", "By", false, nil)},
		{"k8sconv/PodMemoryRssObservable", constructorTest(k8sconv.NewPodMemoryRssObservable, "k8s.pod.memory.rss", "Pod memory RSS.", "By", false, nil)},
		{"k8sconv/PodMemoryUsage", constructorTest(k8sconv.NewPodMemoryUsage, "k8s.pod.memory.usage", "Memory usage of the Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryUsageObservable", constructorTest(k8sconv.NewPodMemoryUsageObservable, "k8s.pod.memory.usage", "Memory usage of the Pod.", "By", false, nil)},
		{"k8sconv/PodMemoryWorkingSet", constructorTest(k8sconv.NewPodMemoryWorkingSet, "k8s.pod.memory.working_set", "Pod memory working set.", "By", false, nil)},
		{"k8sconv/PodMemoryWorkingSetObservable", constructorTest(k8sconv.NewPodMemoryWorkingSetObservable, "k8s.pod.memory.working_set", "Pod memory working set.", "By", false, nil)},
		{"k8sconv/PodNetworkErrors", constructorTest(k8sconv.NewPodNetworkErrors, "k8s.pod.network.errors", "Pod network errors.", "{error}", false, nil)},
		{"k8sconv/PodNetworkErrorsObservable", constructorTest(k8sconv.NewPodNetworkErrorsObservable, "k8s.pod.network.errors", "Pod network errors.", "{error}", false, nil)},
		{"k8sconv/PodNetworkIO", constructorTest(k8sconv.NewPodNetworkIO, "k8s.pod.network.io", "Network bytes for the Pod.", "By", false, nil)},
		{"k8sconv/PodNetworkIOObservable", constructorTest(k8sconv.NewPodNetworkIOObservable, "k8s.pod.network.io", "Network bytes for the Pod.", "By", false, nil)},
		{"k8sconv/PodPagingFaults", constructorTest(k8sconv.NewPodPagingFaults, "k8s.pod.paging.faults", "Pod memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/PodPagingFaultsObservable", constructorTest(k8sconv.NewPodPagingFaultsObservable, "k8s.pod.paging.faults", "Pod memory paging faults.", "{fault}", false, nil)},
		{"k8sconv/PodStatusPhase", constructorTest(k8sconv.NewPodStatusPhase, "k8s.pod.status.phase", "Describes number of K8s Pods that are currently in a given phase.", "{pod}", false, nil)},
		{"k8sconv/PodStatusPhaseObservable", constructorTest(k8sconv.NewPodStatusPhaseObservable, "k8s.pod.status.phase", "Describes number of K8s Pods that are currently in a given phase.", "{pod}", false, nil)},
		{"k8sconv/PodStatusReason", constructorTest(k8sconv.NewPodStatusReason, "k8s.pod.status.reason", "Describes the number of K8s Pods that are currently in a state for a given reason.", "{pod}", false, nil)},
		{"k8sconv/PodStatusReasonObservable", constructorTest(k8sconv.NewPodStatusReasonObservable, "k8s.pod.status.reason", "Describes the number of K8s Pods that are currently in a state for a given reason.", "{pod}", false, nil)},
		{"k8sconv/PodUptime", constructorTest(k8sconv.NewPodUptime, "k8s.pod.uptime", "The time the Pod has been running.", "s", false, nil)},
		{"k8sconv/PodUptimeObservable", constructorTest(k8sconv.NewPodUptimeObservable, "k8s.pod.uptime", "The time the Pod has been running.", "s", false, nil)},
		{"k8sconv/PodVolumeAvailable", constructorTest(k8sconv.NewPodVolumeAvailable, "k8s.pod.volume.available", "Pod volume storage space available.", "By", false, nil)},
		{"k8sconv/PodVolumeAvailableObservable", constructorTest(k8sconv.NewPodVolumeAvailableObservable, "k8s.pod.volume.available", "Pod volume storage space available.", "By", false, nil)},
		{"k8sconv/PodVolumeCapacity", constructorTest(k8sconv.NewPodVolumeCapacity, "k8s.pod.volume.capacity", "Pod volume total capacity.", "By", false, nil)},
		{"k8sconv/PodVolumeCapacityObservable", constructorTest(k8sconv.NewPodVolumeCapacityObservable, "k8s.pod.volume.capacity", "Pod volume total capacity.", "By", false, nil)},
		{"k8sconv/PodVolumeInodeCount", constructorTest(k8sconv.NewPodVolumeInodeCount, "k8s.pod.volume.inode.count", "The total inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeCountObservable", constructorTest(k8sconv.NewPodVolumeInodeCountObservable, "k8s.pod.volume.inode.count", "The total inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeFree", constructorTest(k8sconv.NewPodVolumeInodeFree, "k8s.pod.volume.inode.free", "The free inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeFreeObservable", constructorTest(k8sconv.NewPodVolumeInodeFreeObservable, "k8s.pod.volume.inode.free", "The free inodes in the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeUsed", constructorTest(k8sconv.NewPodVolumeInodeUsed, "k8s.pod.volume.inode.used", "The inodes used by the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeInodeUsedObservable", constructorTest(k8sconv.NewPodVolumeInodeUsedObservable, "k8s.pod.volume.inode.used", "The inodes used by the filesystem of the Pod\u0027s volume.", "{inode}", false, nil)},
		{"k8sconv/PodVolumeUsage", constructorTest(k8sconv.NewPodVolumeUsage, "k8s.pod.volume.usage", "Pod volume usage.", "By", false, nil)},
		{"k8sconv/PodVolumeUsageObservable", constructorTest(k8sconv.NewPodVolumeUsageObservable, "k8s.pod.volume.usage", "Pod volume usage.", "By", false, nil)},
		{"k8sconv/ReplicaSetPodAvailable", constructorTest(k8sconv.NewReplicaSetPodAvailable, "k8s.replicaset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicaSetPodAvailableObservable", constructorTest(k8sconv.NewReplicaSetPodAvailableObservable, "k8s.replicaset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicaSetPodDesired", constructorTest(k8sconv.NewReplicaSetPodDesired, "k8s.replicaset.pod.desired", "Number of desired replica pods in this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicaSetPodDesiredObservable", constructorTest(k8sconv.NewReplicaSetPodDesiredObservable, "k8s.replicaset.pod.desired", "Number of desired replica pods in this replicaset.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodAvailable", constructorTest(k8sconv.NewReplicationControllerPodAvailable, "k8s.replicationcontroller.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodAvailableObservable", constructorTest(k8sconv.NewReplicationControllerPodAvailableObservable, "k8s.replicationcontroller.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodDesired", constructorTest(k8sconv.NewReplicationControllerPodDesired, "k8s.replicationcontroller.pod.desired", "Number of desired replica pods in this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ReplicationControllerPodDesiredObservable", constructorTest(k8sconv.NewReplicationControllerPodDesiredObservable, "k8s.replicationcontroller.pod.desired", "Number of desired replica pods in this replication controller.", "{pod}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitHard", constructorTest(k8sconv.NewResourceQuotaCPULimitHard, "k8s.resourcequota.cpu.limit.hard", "The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitHardObservable", constructorTest(k8sconv.NewResourceQuotaCPULimitHardObservable, "k8s.resourcequota.cpu.limit.hard", "The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitUsed", constructorTest(k8sconv.NewResourceQuotaCPULimitUsed, "k8s.resourcequota.cpu.limit.used", "The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPULimitUsedObservable", constructorTest(k8sconv.NewResourceQuotaCPULimitUsedObservable, "k8s.resourcequota.cpu.limit.used", "The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestHard", constructorTest(k8sconv.NewResourceQuotaCPURequestHard, "k8s.resourcequota.cpu.request.hard", "The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestHardObservable", constructorTest(k8sconv.NewResourceQuotaCPURequestHardObservable, "k8s.resourcequota.cpu.request.hard", "The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestUsed", constructorTest(k8sconv.NewResourceQuotaCPURequestUsed, "k8s.resourcequota.cpu.request.used", "The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaCPURequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaCPURequestUsedObservable, "k8s.resourcequota.cpu.request.used", "The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{cpu}", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitHard", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitHard, "k8s.resourcequota.ephemeral_storage.limit.hard", "The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitHardObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitHardObservable, "k8s.resourcequota.ephemeral_storage.limit.hard", "The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitUsed", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitUsed, "k8s.resourcequota.ephemeral_storage.limit.used", "The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageLimitUsedObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageLimitUsedObservable, "k8s.resourcequota.ephemeral_storage.limit.used", "The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestHard", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestHard, "k8s.resourcequota.ephemeral_storage.request.hard", "The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestHardObservable, "k8s.resourcequota.ephemeral_storage.request.hard", "The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestUsed", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestUsed, "k8s.resourcequota.ephemeral_storage.request.used", "The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaEphemeralStorageRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaEphemeralStorageRequestUsedObservable, "k8s.resourcequota.ephemeral_storage.request.used", "The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestHard", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestHard, "k8s.resourcequota.hugepage_count.request.hard", "The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestHardObservable, "k8s.resourcequota.hugepage_count.request.hard", "The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestUsed", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestUsed, "k8s.resourcequota.hugepage_count.request.used", "The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaHugepageCountRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaHugepageCountRequestUsedObservable, "k8s.resourcequota.hugepage_count.request.used", "The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{hugepage}", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitHard", constructorTest(k8sconv.NewResourceQuotaMemoryLimitHard, "k8s.resourcequota.memory.limit.hard", "The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitHardObservable", constructorTest(k8sconv.NewResourceQuotaMemoryLimitHardObservable, "k8s.resourcequota.memory.limit.hard", "The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitUsed", constructorTest(k8sconv.NewResourceQuotaMemoryLimitUsed, "k8s.resourcequota.memory.limit.used", "The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryLimitUsedObservable", constructorTest(k8sconv.NewResourceQuotaMemoryLimitUsedObservable, "k8s.resourcequota.memory.limit.used", "The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestHard", constructorTest(k8sconv.NewResourceQuotaMemoryRequestHard, "k8s.resourcequota.memory.request.hard", "The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaMemoryRequestHardObservable, "k8s.resourcequota.memory.request.hard", "The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestUsed", constructorTest(k8sconv.NewResourceQuotaMemoryRequestUsed, "k8s.resourcequota.memory.request.used", "The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaMemoryRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaMemoryRequestUsedObservable, "k8s.resourcequota.memory.request.used", "The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountHard", constructorTest(k8sconv.NewResourceQuotaObjectCountHard, "k8s.resourcequota.object_count.hard", "The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountHardObservable", constructorTest(k8sconv.NewResourceQuotaObjectCountHardObservable, "k8s.resourcequota.object_count.hard", "The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountUsed", constructorTest(k8sconv.NewResourceQuotaObjectCountUsed, "k8s.resourcequota.object_count.used", "The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaObjectCountUsedObservable", constructorTest(k8sconv.NewResourceQuotaObjectCountUsedObservable, "k8s.resourcequota.object_count.used", "The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "{object}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountHard", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountHard, "k8s.resourcequota.persistentvolumeclaim_count.hard", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountHardObservable", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountHardObservable, "k8s.resourcequota.persistentvolumeclaim_count.hard", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountUsed", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountUsed, "k8s.resourcequota.persistentvolumeclaim_count.used", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaPersistentvolumeclaimCountUsedObservable", constructorTest(k8sconv.NewResourceQuotaPersistentvolumeclaimCountUsedObservable, "k8s.resourcequota.persistentvolumeclaim_count.used", "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace.", "{persistentvolumeclaim}", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestHard", constructorTest(k8sconv.NewResourceQuotaStorageRequestHard, "k8s.resourcequota.storage.request.hard", "The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestHardObservable", constructorTest(k8sconv.NewResourceQuotaStorageRequestHardObservable, "k8s.resourcequota.storage.request.hard", "The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestUsed", constructorTest(k8sconv.NewResourceQuotaStorageRequestUsed, "k8s.resourcequota.storage.request.used", "The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ResourceQuotaStorageRequestUsedObservable", constructorTest(k8sconv.NewResourceQuotaStorageRequestUsedObservable, "k8s.resourcequota.storage.request.used", "The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace.", "By", false, nil)},
		{"k8sconv/ServiceEndpointCount", constructorTest(k8sconv.NewServiceEndpointCount, "k8s.service.endpoint.count", "Number of endpoints for a service by condition and address type.", "{endpoint}", false, nil)},
		{"k8sconv/ServiceEndpointCountObservable", constructorTest(k8sconv.NewServiceEndpointCountObservable, "k8s.service.endpoint.count", "Number of endpoints for a service by condition and address type.", "{endpoint}", false, nil)},
		{"k8sconv/ServiceLoadBalancerIngressCount", constructorTest(k8sconv.NewServiceLoadBalancerIngressCount, "k8s.service.load_balancer.ingress.count", "Number of load balancer ingress points (external IPs/hostnames) assigned to the service.", "{ingress}", false, nil)},
		{"k8sconv/ServiceLoadBalancerIngressCountObservable", constructorTest(k8sconv.NewServiceLoadBalancerIngressCountObservable, "k8s.service.load_balancer.ingress.count", "Number of load balancer ingress points (external IPs/hostnames) assigned to the service.", "{ingress}", false, nil)},
		{"k8sconv/StatefulSetPodAvailable", constructorTest(k8sconv.NewStatefulSetPodAvailable, "k8s.statefulset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodAvailableObservable", constructorTest(k8sconv.NewStatefulSetPodAvailableObservable, "k8s.statefulset.pod.available", "Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodCurrent", constructorTest(k8sconv.NewStatefulSetPodCurrent, "k8s.statefulset.pod.current", "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodCurrentObservable", constructorTest(k8sconv.NewStatefulSetPodCurrentObservable, "k8s.statefulset.pod.current", "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodDesired", constructorTest(k8sconv.NewStatefulSetPodDesired, "k8s.statefulset.pod.desired", "Number of desired replica pods in this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodDesiredObservable", constructorTest(k8sconv.NewStatefulSetPodDesiredObservable, "k8s.statefulset.pod.desired", "Number of desired replica pods in this statefulset.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodReady", constructorTest(k8sconv.NewStatefulSetPodReady, "k8s.statefulset.pod.ready", "The number of replica pods created for this statefulset with a Ready Condition.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodReadyObservable", constructorTest(k8sconv.NewStatefulSetPodReadyObservable, "k8s.statefulset.pod.ready", "The number of replica pods created for this statefulset with a Ready Condition.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodUpdated", constructorTest(k8sconv.NewStatefulSetPodUpdated, "k8s.statefulset.pod.updated", "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision.", "{pod}", false, nil)},
		{"k8sconv/StatefulSetPodUpdatedObservable", constructorTest(k8sconv.NewStatefulSetPodUpdatedObservable, "k8s.statefulset.pod.updated", "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision.", "{pod}", false, nil)},
		{"messagingconv/ClientConsumedMessages", constructorTest(messagingconv.NewClientConsumedMessages, "messaging.client.consumed.messages", "Number of messages that were delivered to the application.", "{message}", false, nil)},
		{"messagingconv/ClientConsumedMessagesObservable", constructorTest(messagingconv.NewClientConsumedMessagesObservable, "messaging.client.consumed.messages", "Number of messages that were delivered to the application.", "{message}", false, nil)},
		{"messagingconv/ClientOperationDuration", constructorTest(messagingconv.NewClientOperationDuration, "messaging.client.operation.duration", "Duration of messaging operation initiated by a producer or consumer client.", "s", true, nil)},
		{"messagingconv/ClientSentMessages", constructorTest(messagingconv.NewClientSentMessages, "messaging.client.sent.messages", "Number of messages producer attempted to send to the broker.", "{message}", false, nil)},
		{"messagingconv/ClientSentMessagesObservable", constructorTest(messagingconv.NewClientSentMessagesObservable, "messaging.client.sent.messages", "Number of messages producer attempted to send to the broker.", "{message}", false, nil)},
		{"messagingconv/ProcessDuration", constructorTest(messagingconv.NewProcessDuration, "messaging.process.duration", "Duration of processing operation.", "s", true, nil)},
		{"nfsconv/ClientNetCount", constructorTest(nfsconv.NewClientNetCount, "nfs.client.net.count", "Reports the count of kernel NFS client TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ClientNetCountObservable", constructorTest(nfsconv.NewClientNetCountObservable, "nfs.client.net.count", "Reports the count of kernel NFS client TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ClientNetTCPConnectionAccepted", constructorTest(nfsconv.NewClientNetTCPConnectionAccepted, "nfs.client.net.tcp.connection.accepted", "Reports the count of kernel NFS client TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ClientNetTCPConnectionAcceptedObservable", constructorTest(nfsconv.NewClientNetTCPConnectionAcceptedObservable, "nfs.client.net.tcp.connection.accepted", "Reports the count of kernel NFS client TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ClientOperationCount", constructorTest(nfsconv.NewClientOperationCount, "nfs.client.operation.count", "Reports the count of kernel NFSv4+ client operations.", "{operation}", false, nil)},
		{"nfsconv/ClientOperationCountObservable", constructorTest(nfsconv.NewClientOperationCountObservable, "nfs.client.operation.count", "Reports the count of kernel NFSv4+ client operations.", "{operation}", false, nil)},
		{"nfsconv/ClientProcedureCount", constructorTest(nfsconv.NewClientProcedureCount, "nfs.client.procedure.count", "Reports the count of kernel NFS client procedures.", "{procedure}", false, nil)},
		{"nfsconv/ClientProcedureCountObservable", constructorTest(nfsconv.NewClientProcedureCountObservable, "nfs.client.procedure.count", "Reports the count of kernel NFS client procedures.", "{procedure}", false, nil)},
		{"nfsconv/ClientRPCAuthrefreshCount", constructorTest(nfsconv.NewClientRPCAuthrefreshCount, "nfs.client.rpc.authrefresh.count", "Reports the count of kernel NFS client RPC authentication refreshes.", "{authrefresh}", false, nil)},
		{"nfsconv/ClientRPCAuthrefreshCountObservable", constructorTest(nfsconv.NewClientRPCAuthrefreshCountObservable, "nfs.client.rpc.authrefresh.count", "Reports the count of kernel NFS client RPC authentication refreshes.", "{authrefresh}", false, nil)},
		{"nfsconv/ClientRPCCount", constructorTest(nfsconv.NewClientRPCCount, "nfs.client.rpc.count", "Reports the count of kernel NFS client RPCs sent, regardless of whether they\u0027re accepted/rejected by the server.", "{request}", false, nil)},
		{"nfsconv/ClientRPCCountObservable", constructorTest(nfsconv.NewClientRPCCountObservable, "nfs.client.rpc.count", "Reports the count of kernel NFS client RPCs sent, regardless of whether they\u0027re accepted/rejected by the server.", "{request}", false, nil)},
		{"nfsconv/ClientRPCRetransmitCount", constructorTest(nfsconv.NewClientRPCRetransmitCount, "nfs.client.rpc.retransmit.count", "Reports the count of kernel NFS client RPC retransmits.", "{retransmit}", false, nil)},
		{"nfsconv/ClientRPCRetransmitCountObservable", constructorTest(nfsconv.NewClientRPCRetransmitCountObservable, "nfs.client.rpc.retransmit.count", "Reports the count of kernel NFS client RPC retransmits.", "{retransmit}", false, nil)},
		{"nfsconv/ServerFhStaleCount", constructorTest(nfsconv.NewServerFhStaleCount, "nfs.server.fh.stale.count", "Reports the count of kernel NFS server stale file handles.", "{fh}", false, nil)},
		{"nfsconv/ServerFhStaleCountObservable", constructorTest(nfsconv.NewServerFhStaleCountObservable, "nfs.server.fh.stale.count", "Reports the count of kernel NFS server stale file handles.", "{fh}", false, nil)},
		{"nfsconv/ServerIO", constructorTest(nfsconv.NewServerIO, "nfs.server.io", "Reports the count of kernel NFS server bytes returned to receive and transmit (read and write) requests.", "By", false, nil)},
		{"nfsconv/ServerIOObservable", constructorTest(nfsconv.NewServerIOObservable, "nfs.server.io", "Reports the count of kernel NFS server bytes returned to receive and transmit (read and write) requests.", "By", false, nil)},
		{"nfsconv/ServerNetCount", constructorTest(nfsconv.NewServerNetCount, "nfs.server.net.count", "Reports the count of kernel NFS server TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ServerNetCountObservable", constructorTest(nfsconv.NewServerNetCountObservable, "nfs.server.net.count", "Reports the count of kernel NFS server TCP segments and UDP datagrams handled.", "{record}", false, nil)},
		{"nfsconv/ServerNetTCPConnectionAccepted", constructorTest(nfsconv.NewServerNetTCPConnectionAccepted, "nfs.server.net.tcp.connection.accepted", "Reports the count of kernel NFS server TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ServerNetTCPConnectionAcceptedObservable", constructorTest(nfsconv.NewServerNetTCPConnectionAcceptedObservable, "nfs.server.net.tcp.connection.accepted", "Reports the count of kernel NFS server TCP connections accepted.", "{connection}", false, nil)},
		{"nfsconv/ServerOperationCount", constructorTest(nfsconv.NewServerOperationCount, "nfs.server.operation.count", "Reports the count of kernel NFSv4+ server operations.", "{operation}", false, nil)},
		{"nfsconv/ServerOperationCountObservable", constructorTest(nfsconv.NewServerOperationCountObservable, "nfs.server.operation.count", "Reports the count of kernel NFSv4+ server operations.", "{operation}", false, nil)},
		{"nfsconv/ServerProcedureCount", constructorTest(nfsconv.NewServerProcedureCount, "nfs.server.procedure.count", "Reports the count of kernel NFS server procedures.", "{procedure}", false, nil)},
		{"nfsconv/ServerProcedureCountObservable", constructorTest(nfsconv.NewServerProcedureCountObservable, "nfs.server.procedure.count", "Reports the count of kernel NFS server procedures.", "{procedure}", false, nil)},
		{"nfsconv/ServerRepcacheRequests", constructorTest(nfsconv.NewServerRepcacheRequests, "nfs.server.repcache.requests", "Reports the kernel NFS server reply cache request count by cache hit status.", "{request}", false, nil)},
		{"nfsconv/ServerRepcacheRequestsObservable", constructorTest(nfsconv.NewServerRepcacheRequestsObservable, "nfs.server.repcache.requests", "Reports the kernel NFS server reply cache request count by cache hit status.", "{request}", false, nil)},
		{"nfsconv/ServerRPCCount", constructorTest(nfsconv.NewServerRPCCount, "nfs.server.rpc.count", "Reports the count of kernel NFS server RPCs handled.", "{request}", false, nil)},
		{"nfsconv/ServerRPCCountObservable", constructorTest(nfsconv.NewServerRPCCountObservable, "nfs.server.rpc.count", "Reports the count of kernel NFS server RPCs handled.", "{request}", false, nil)},
		{"nfsconv/ServerThreadCount", constructorTest(nfsconv.NewServerThreadCount, "nfs.server.thread.count", "Reports the count of kernel NFS server available threads.", "{thread}", false, nil)},
		{"nfsconv/ServerThreadCountObservable", constructorTest(nfsconv.NewServerThreadCountObservable, "nfs.server.thread.count", "Reports the count of kernel NFS server available threads.", "{thread}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitHard", constructorTest(openshiftconv.NewClusterquotaCPULimitHard, "openshift.clusterquota.cpu.limit.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitHardObservable", constructorTest(openshiftconv.NewClusterquotaCPULimitHardObservable, "openshift.clusterquota.cpu.limit.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitUsed", constructorTest(openshiftconv.NewClusterquotaCPULimitUsed, "openshift.clusterquota.cpu.limit.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPULimitUsedObservable", constructorTest(openshiftconv.NewClusterquotaCPULimitUsedObservable, "openshift.clusterquota.cpu.limit.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestHard", constructorTest(openshiftconv.NewClusterquotaCPURequestHard, "openshift.clusterquota.cpu.request.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestHardObservable", constructorTest(openshiftconv.NewClusterquotaCPURequestHardObservable, "openshift.clusterquota.cpu.request.hard", "The enforced hard limit of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestUsed", constructorTest(openshiftconv.NewClusterquotaCPURequestUsed, "openshift.clusterquota.cpu.request.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaCPURequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaCPURequestUsedObservable, "openshift.clusterquota.cpu.request.used", "The current observed total usage of the resource across all projects.", "{cpu}", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitHard", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitHard, "openshift.clusterquota.ephemeral_storage.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitHardObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitHardObservable, "openshift.clusterquota.ephemeral_storage.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitUsed", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitUsed, "openshift.clusterquota.ephemeral_storage.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageLimitUsedObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageLimitUsedObservable, "openshift.clusterquota.ephemeral_storage.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestHard", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestHard, "openshift.clusterquota.ephemeral_storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestHardObservable, "openshift.clusterquota.ephemeral_storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestUsed", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestUsed, "openshift.clusterquota.ephemeral_storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaEphemeralStorageRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaEphemeralStorageRequestUsedObservable, "openshift.clusterquota.ephemeral_storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestHard", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestHard, "openshift.clusterquota.hugepage_count.request.hard", "The enforced hard limit of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestHardObservable, "openshift.clusterquota.hugepage_count.request.hard", "The enforced hard limit of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestUsed", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestUsed, "openshift.clusterquota.hugepage_count.request.used", "The current observed total usage of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaHugepageCountRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaHugepageCountRequestUsedObservable, "openshift.clusterquota.hugepage_count.request.used", "The current observed total usage of the resource across all projects.", "{hugepage}", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitHard", constructorTest(openshiftconv.NewClusterquotaMemoryLimitHard, "openshift.clusterquota.memory.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitHardObservable", constructorTest(openshiftconv.NewClusterquotaMemoryLimitHardObservable, "openshift.clusterquota.memory.limit.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitUsed", constructorTest(openshiftconv.NewClusterquotaMemoryLimitUsed, "openshift.clusterquota.memory.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryLimitUsedObservable", constructorTest(openshiftconv.NewClusterquotaMemoryLimitUsedObservable, "openshift.clusterquota.memory.limit.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestHard", constructorTest(openshiftconv.NewClusterquotaMemoryRequestHard, "openshift.clusterquota.memory.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaMemoryRequestHardObservable, "openshift.clusterquota.memory.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestUsed", constructorTest(openshiftconv.NewClusterquotaMemoryRequestUsed, "openshift.clusterquota.memory.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaMemoryRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaMemoryRequestUsedObservable, "openshift.clusterquota.memory.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountHard", constructorTest(openshiftconv.NewClusterquotaObjectCountHard, "openshift.clusterquota.object_count.hard", "The enforced hard limit of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountHardObservable", constructorTest(openshiftconv.NewClusterquotaObjectCountHardObservable, "openshift.clusterquota.object_count.hard", "The enforced hard limit of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountUsed", constructorTest(openshiftconv.NewClusterquotaObjectCountUsed, "openshift.clusterquota.object_count.used", "The current observed total usage of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaObjectCountUsedObservable", constructorTest(openshiftconv.NewClusterquotaObjectCountUsedObservable, "openshift.clusterquota.object_count.used", "The current observed total usage of the resource across all projects.", "{object}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountHard", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountHard, "openshift.clusterquota.persistentvolumeclaim_count.hard", "The enforced hard limit of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountHardObservable", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountHardObservable, "openshift.clusterquota.persistentvolumeclaim_count.hard", "The enforced hard limit of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountUsed", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountUsed, "openshift.clusterquota.persistentvolumeclaim_count.used", "The current observed total usage of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaPersistentvolumeclaimCountUsedObservable", constructorTest(openshiftconv.NewClusterquotaPersistentvolumeclaimCountUsedObservable, "openshift.clusterquota.persistentvolumeclaim_count.used", "The current observed total usage of the resource across all projects.", "{persistentvolumeclaim}", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestHard", constructorTest(openshiftconv.NewClusterquotaStorageRequestHard, "openshift.clusterquota.storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestHardObservable", constructorTest(openshiftconv.NewClusterquotaStorageRequestHardObservable, "openshift.clusterquota.storage.request.hard", "The enforced hard limit of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestUsed", constructorTest(openshiftconv.NewClusterquotaStorageRequestUsed, "openshift.clusterquota.storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"openshiftconv/ClusterquotaStorageRequestUsedObservable", constructorTest(openshiftconv.NewClusterquotaStorageRequestUsedObservable, "openshift.clusterquota.storage.request.used", "The current observed total usage of the resource across all projects.", "By", false, nil)},
		{"otelconv/SDKExporterLogExported", constructorTest(otelconv.NewSDKExporterLogExported, "otel.sdk.exporter.log.exported", "The number of log records for which the export has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKExporterLogExportedObservable", constructorTest(otelconv.NewSDKExporterLogExportedObservable, "otel.sdk.exporter.log.exported", "The number of log records for which the export has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKExporterLogInflight", constructorTest(otelconv.NewSDKExporterLogInflight, "otel.sdk.exporter.log.inflight", "The number of log records which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{log_record}", false, nil)},
		{"otelconv/SDKExporterLogInflightObservable", constructorTest(otelconv.NewSDKExporterLogInflightObservable, "otel.sdk.exporter.log.inflight", "The number of log records which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{log_record}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointExported", constructorTest(otelconv.NewSDKExporterMetricDataPointExported, "otel.sdk.exporter.metric_data_point.exported", "The number of metric data points for which the export has finished, either successful or failed.", "{data_point}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointExportedObservable", constructorTest(otelconv.NewSDKExporterMetricDataPointExportedObservable, "otel.sdk.exporter.metric_data_point.exported", "The number of metric data points for which the export has finished, either successful or failed.", "{data_point}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointInflight", constructorTest(otelconv.NewSDKExporterMetricDataPointInflight, "otel.sdk.exporter.metric_data_point.inflight", "The number of metric data points which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{data_point}", false, nil)},
		{"otelconv/SDKExporterMetricDataPointInflightObservable", constructorTest(otelconv.NewSDKExporterMetricDataPointInflightObservable, "otel.sdk.exporter.metric_data_point.inflight", "The number of metric data points which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{data_point}", false, nil)},
		{"otelconv/SDKExporterOperationDuration", constructorTest(otelconv.NewSDKExporterOperationDuration, "otel.sdk.exporter.operation.duration", "The duration of exporting a batch of telemetry records.", "s", true, nil)},
		{"otelconv/SDKExporterSpanExported", constructorTest(otelconv.NewSDKExporterSpanExported, "otel.sdk.exporter.span.exported", "The number of spans for which the export has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKExporterSpanExportedObservable", constructorTest(otelconv.NewSDKExporterSpanExportedObservable, "otel.sdk.exporter.span.exported", "The number of spans for which the export has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKExporterSpanInflight", constructorTest(otelconv.NewSDKExporterSpanInflight, "otel.sdk.exporter.span.inflight", "The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{span}", false, nil)},
		{"otelconv/SDKExporterSpanInflightObservable", constructorTest(otelconv.NewSDKExporterSpanInflightObservable, "otel.sdk.exporter.span.inflight", "The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed).", "{span}", false, nil)},
		{"otelconv/SDKLogCreated", constructorTest(otelconv.NewSDKLogCreated, "otel.sdk.log.created", "The number of log records submitted to an enabled `Logger`.", "{log_record}", false, nil)},
		{"otelconv/SDKLogCreatedObservable", constructorTest(otelconv.NewSDKLogCreatedObservable, "otel.sdk.log.created", "The number of log records submitted to an enabled `Logger`.", "{log_record}", false, nil)},
		{"otelconv/SDKMetricReaderCollectionDuration", constructorTest(otelconv.NewSDKMetricReaderCollectionDuration, "otel.sdk.metric_reader.collection.duration", "The duration of the collect operation of the metric reader.", "s", true, nil)},
		{"otelconv/SDKProcessorLogProcessed", constructorTest(otelconv.NewSDKProcessorLogProcessed, "otel.sdk.processor.log.processed", "The number of log records for which the processing has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorLogProcessedObservable", constructorTest(otelconv.NewSDKProcessorLogProcessedObservable, "otel.sdk.processor.log.processed", "The number of log records for which the processing has finished, either successful or failed.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorLogQueueCapacity", constructorTest(otelconv.NewSDKProcessorLogQueueCapacity, "otel.sdk.processor.log.queue.capacity", "The maximum number of log records the queue of a given instance of an SDK Log Record processor can hold.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorLogQueueSize", constructorTest(otelconv.NewSDKProcessorLogQueueSize, "otel.sdk.processor.log.queue.size", "The number of log records in the queue of a given instance of an SDK log processor.", "{log_record}", false, nil)},
		{"otelconv/SDKProcessorSpanProcessed", constructorTest(otelconv.NewSDKProcessorSpanProcessed, "otel.sdk.processor.span.processed", "The number of spans for which the processing has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKProcessorSpanProcessedObservable", constructorTest(otelconv.NewSDKProcessorSpanProcessedObservable, "otel.sdk.processor.span.processed", "The number of spans for which the processing has finished, either successful or failed.", "{span}", false, nil)},
		{"otelconv/SDKProcessorSpanQueueCapacity", constructorTest(otelconv.NewSDKProcessorSpanQueueCapacity, "otel.sdk.processor.span.queue.capacity", "The maximum number of spans the queue of a given instance of an SDK span processor can hold.", "{span}", false, nil)},
		{"otelconv/SDKProcessorSpanQueueSize", constructorTest(otelconv.NewSDKProcessorSpanQueueSize, "otel.sdk.processor.span.queue.size", "The number of spans in the queue of a given instance of an SDK span processor.", "{span}", false, nil)},
		{"otelconv/SDKSpanLive", constructorTest(otelconv.NewSDKSpanLive, "otel.sdk.span.live", "The number of created spans with `recording=true` for which the end operation has not been called yet.", "{span}", false, nil)},
		{"otelconv/SDKSpanLiveObservable", constructorTest(otelconv.NewSDKSpanLiveObservable, "otel.sdk.span.live", "The number of created spans with `recording=true` for which the end operation has not been called yet.", "{span}", false, nil)},
		{"otelconv/SDKSpanStarted", constructorTest(otelconv.NewSDKSpanStarted, "otel.sdk.span.started", "The number of created spans.", "{span}", false, nil)},
		{"otelconv/SDKSpanStartedObservable", constructorTest(otelconv.NewSDKSpanStartedObservable, "otel.sdk.span.started", "The number of created spans.", "{span}", false, nil)},
		{"processconv/ContextSwitches", constructorTest(processconv.NewContextSwitches, "process.context_switches", "Number of times the process has been context switched.", "{context_switch}", false, nil)},
		{"processconv/ContextSwitchesObservable", constructorTest(processconv.NewContextSwitchesObservable, "process.context_switches", "Number of times the process has been context switched.", "{context_switch}", false, nil)},
		{"processconv/CPUTime", constructorTest(processconv.NewCPUTime, "process.cpu.time", "Total CPU seconds broken down by different CPU modes.", "s", false, nil)},
		{"processconv/CPUUtilization", constructorTest(processconv.NewCPUUtilization, "process.cpu.utilization", "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process.", "1", false, nil)},
		{"processconv/CPUUtilizationObservable", constructorTest(processconv.NewCPUUtilizationObservable, "process.cpu.utilization", "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process.", "1", false, nil)},
		{"processconv/DiskIO", constructorTest(processconv.NewDiskIO, "process.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"processconv/DiskIOObservable", constructorTest(processconv.NewDiskIOObservable, "process.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"processconv/DiskOperations", constructorTest(processconv.NewDiskOperations, "process.disk.operations", "Number of disk operations performed by the process.", "{operation}", false, nil)},
		{"processconv/DiskOperationsObservable", constructorTest(processconv.NewDiskOperationsObservable, "process.disk.operations", "Number of disk operations performed by the process.", "{operation}", false, nil)},
		{"processconv/MemoryUsage", constructorTest(processconv.NewMemoryUsage, "process.memory.usage", "The amount of physical memory in use.", "By", false, nil)},
		{"processconv/MemoryUsageObservable", constructorTest(processconv.NewMemoryUsageObservable, "process.memory.usage", "The amount of physical memory in use.", "By", false, nil)},
		{"processconv/MemoryUtilization", constructorTest(processconv.NewMemoryUtilization, "process.memory.utilization", "Percentage of total physical memory that is used by the process.", "1", false, nil)},
		{"processconv/MemoryUtilizationObservable", constructorTest(processconv.NewMemoryUtilizationObservable, "process.memory.utilization", "Percentage of total physical memory that is used by the process.", "1", false, nil)},
		{"processconv/MemoryVirtual", constructorTest(processconv.NewMemoryVirtual, "process.memory.virtual", "The amount of committed virtual memory.", "By", false, nil)},
		{"processconv/MemoryVirtualObservable", constructorTest(processconv.NewMemoryVirtualObservable, "process.memory.virtual", "The amount of committed virtual memory.", "By", false, nil)},
		{"processconv/NetworkIO", constructorTest(processconv.NewNetworkIO, "process.network.io", "Network bytes transferred.", "By", false, nil)},
		{"processconv/NetworkIOObservable", constructorTest(processconv.NewNetworkIOObservable, "process.network.io", "Network bytes transferred.", "By", false, nil)},
		{"processconv/PagingFaults", constructorTest(processconv.NewPagingFaults, "process.paging.faults", "Number of page faults the process has made.", "{fault}", false, nil)},
		{"processconv/PagingFaultsObservable", constructorTest(processconv.NewPagingFaultsObservable, "process.paging.faults", "Number of page faults the process has made.", "{fault}", false, nil)},
		{"processconv/SignalsPending", constructorTest(processconv.NewSignalsPending, "process.signals_pending", "Number of pending signals for the process.", "{signal}", false, nil)},
		{"processconv/SignalsPendingObservable", constructorTest(processconv.NewSignalsPendingObservable, "process.signals_pending", "Number of pending signals for the process.", "{signal}", false, nil)},
		{"processconv/ThreadCount", constructorTest(processconv.NewThreadCount, "process.thread.count", "Process threads count.", "{thread}", false, nil)},
		{"processconv/ThreadCountObservable", constructorTest(processconv.NewThreadCountObservable, "process.thread.count", "Process threads count.", "{thread}", false, nil)},
		{"processconv/UnixFileDescriptorCount", constructorTest(processconv.NewUnixFileDescriptorCount, "process.unix.file_descriptor.count", "Number of UNIX file descriptors in use by the process.", "{file_descriptor}", false, nil)},
		{"processconv/UnixFileDescriptorCountObservable", constructorTest(processconv.NewUnixFileDescriptorCountObservable, "process.unix.file_descriptor.count", "Number of UNIX file descriptors in use by the process.", "{file_descriptor}", false, nil)},
		{"processconv/Uptime", constructorTest(processconv.NewUptime, "process.uptime", "The time the process has been running.", "s", false, nil)},
		{"processconv/UptimeObservable", constructorTest(processconv.NewUptimeObservable, "process.uptime", "The time the process has been running.", "s", false, nil)},
		{"processconv/WindowsHandleCount", constructorTest(processconv.NewWindowsHandleCount, "process.windows.handle.count", "Number of handles held by the process.", "{handle}", false, nil)},
		{"processconv/WindowsHandleCountObservable", constructorTest(processconv.NewWindowsHandleCountObservable, "process.windows.handle.count", "Number of handles held by the process.", "{handle}", false, nil)},
		{"rpcconv/ClientCallDuration", constructorTest(rpcconv.NewClientCallDuration, "rpc.client.call.duration", "Measures the duration of an outgoing Remote Procedure Call (RPC).", "s", true, []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 60})},
		{"rpcconv/ServerCallDuration", constructorTest(rpcconv.NewServerCallDuration, "rpc.server.call.duration", "Measures the duration of an incoming Remote Procedure Call (RPC).", "s", true, []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30, 60})},
		{"signalrconv/ServerActiveConnections", constructorTest(signalrconv.NewServerActiveConnections, "signalr.server.active_connections", "Number of connections that are currently active on the server.", "{connection}", false, nil)},
		{"signalrconv/ServerActiveConnectionsObservable", constructorTest(signalrconv.NewServerActiveConnectionsObservable, "signalr.server.active_connections", "Number of connections that are currently active on the server.", "{connection}", false, nil)},
		{"signalrconv/ServerConnectionDuration", constructorTest(signalrconv.NewServerConnectionDuration, "signalr.server.connection.duration", "The duration of connections on the server.", "s", true, nil)},
		{"systemconv/CPUFrequency", constructorTest(systemconv.NewCPUFrequency, "system.cpu.frequency", "Operating frequency of the logical CPU in Hertz.", "Hz", false, nil)},
		{"systemconv/CPUFrequencyObservable", constructorTest(systemconv.NewCPUFrequencyObservable, "system.cpu.frequency", "Operating frequency of the logical CPU in Hertz.", "Hz", false, nil)},
		{"systemconv/CPULogicalCount", constructorTest(systemconv.NewCPULogicalCount, "system.cpu.logical.count", "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking.", "{cpu}", false, nil)},
		{"systemconv/CPULogicalCountObservable", constructorTest(systemconv.NewCPULogicalCountObservable, "system.cpu.logical.count", "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking.", "{cpu}", false, nil)},
		{"systemconv/CPUPhysicalCount", constructorTest(systemconv.NewCPUPhysicalCount, "system.cpu.physical.count", "Reports the number of actual physical processor cores on the hardware.", "{cpu}", false, nil)},
		{"systemconv/CPUPhysicalCountObservable", constructorTest(systemconv.NewCPUPhysicalCountObservable, "system.cpu.physical.count", "Reports the number of actual physical processor cores on the hardware.", "{cpu}", false, nil)},
		{"systemconv/CPUTime", constructorTest(systemconv.NewCPUTime, "system.cpu.time", "Seconds each logical CPU spent on each mode.", "s", false, nil)},
		{"systemconv/CPUUtilization", constructorTest(systemconv.NewCPUUtilization, "system.cpu.utilization", "For each logical CPU, the utilization is calculated as the change in cumulative CPU time (cpu.time) over a measurement interval, divided by the elapsed time.", "1", false, nil)},
		{"systemconv/CPUUtilizationObservable", constructorTest(systemconv.NewCPUUtilizationObservable, "system.cpu.utilization", "For each logical CPU, the utilization is calculated as the change in cumulative CPU time (cpu.time) over a measurement interval, divided by the elapsed time.", "1", false, nil)},
		{"systemconv/DiskIO", constructorTest(systemconv.NewDiskIO, "system.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"systemconv/DiskIOObservable", constructorTest(systemconv.NewDiskIOObservable, "system.disk.io", "Disk bytes transferred.", "By", false, nil)},
		{"systemconv/DiskIOTime", constructorTest(systemconv.NewDiskIOTime, "system.disk.io_time", "Time disk spent activated.", "s", false, nil)},
		{"systemconv/DiskIOTimeObservable", constructorTest(systemconv.NewDiskIOTimeObservable, "system.disk.io_time", "Time disk spent activated.", "s", false, nil)},
		{"systemconv/DiskLimit", constructorTest(systemconv.NewDiskLimit, "system.disk.limit", "The total storage capacity of the disk.", "By", false, nil)},
		{"systemconv/DiskLimitObservable", constructorTest(systemconv.NewDiskLimitObservable, "system.disk.limit", "The total storage capacity of the disk.", "By", false, nil)},
		{"systemconv/DiskMerged", constructorTest(systemconv.NewDiskMerged, "system.disk.merged", "The number of disk reads/writes merged into single physical disk access operations.", "{operation}", false, nil)},
		{"systemconv/DiskMergedObservable", constructorTest(systemconv.NewDiskMergedObservable, "system.disk.merged", "The number of disk reads/writes merged into single physical disk access operations.", "{operation}", false, nil)},
		{"systemconv/DiskOperationTime", constructorTest(systemconv.NewDiskOperationTime, "system.disk.operation_time", "Sum of the time each operation took to complete.", "s", false, nil)},
		{"systemconv/DiskOperationTimeObservable", constructorTest(systemconv.NewDiskOperationTimeObservable, "system.disk.operation_time", "Sum of the time each operation took to complete.", "s", false, nil)},
		{"systemconv/DiskOperations", constructorTest(systemconv.NewDiskOperations, "system.disk.operations", "Disk operations count.", "{operation}", false, nil)},
		{"systemconv/DiskOperationsObservable", constructorTest(systemconv.NewDiskOperationsObservable, "system.disk.operations", "Disk operations count.", "{operation}", false, nil)},
		{"systemconv/FilesystemLimit", constructorTest(systemconv.NewFilesystemLimit, "system.filesystem.limit", "The total storage capacity of the filesystem.", "By", false, nil)},
		{"systemconv/FilesystemLimitObservable", constructorTest(systemconv.NewFilesystemLimitObservable, "system.filesystem.limit", "The total storage capacity of the filesystem.", "By", false, nil)},
		{"systemconv/FilesystemLockCount", constructorTest(systemconv.NewFilesystemLockCount, "system.filesystem.lock.count", "Filesystem lock counts.", "{lock}", false, nil)},
		{"systemconv/FilesystemLockCountObservable", constructorTest(systemconv.NewFilesystemLockCountObservable, "system.filesystem.lock.count", "Filesystem lock counts.", "{lock}", false, nil)},
		{"systemconv/FilesystemUsage", constructorTest(systemconv.NewFilesystemUsage, "system.filesystem.usage", "Reports a filesystem\u0027s space usage across different states.", "By", false, nil)},
		{"systemconv/FilesystemUsageObservable", constructorTest(systemconv.NewFilesystemUsageObservable, "system.filesystem.usage", "Reports a filesystem\u0027s space usage across different states.", "By", false, nil)},
		{"systemconv/FilesystemUtilization", constructorTest(systemconv.NewFilesystemUtilization, "system.filesystem.utilization", "Fraction of filesystem bytes used.", "1", false, nil)},
		{"systemconv/FilesystemUtilizationObservable", constructorTest(systemconv.NewFilesystemUtilizationObservable, "system.filesystem.utilization", "Fraction of filesystem bytes used.", "1", false, nil)},
		{"systemconv/MemoryLimit", constructorTest(systemconv.NewMemoryLimit, "system.memory.limit", "Total virtual memory available in the system.", "By", false, nil)},
		{"systemconv/MemoryLimitObservable", constructorTest(systemconv.NewMemoryLimitObservable, "system.memory.limit", "Total virtual memory available in the system.", "By", false, nil)},
		{"systemconv/MemoryLinuxAvailable", constructorTest(systemconv.NewMemoryLinuxAvailable, "system.memory.linux.available", "An estimate of how much memory is available for starting new applications, without causing swapping.", "By", false, nil)},
		{"systemconv/MemoryLinuxAvailableObservable", constructorTest(systemconv.NewMemoryLinuxAvailableObservable, "system.memory.linux.available", "An estimate of how much memory is available for starting new applications, without causing swapping.", "By", false, nil)},
		{"systemconv/MemoryLinuxHugepagesLimit", constructorTest(systemconv.NewMemoryLinuxHugepagesLimit, "system.memory.linux.hugepages.limit", "Total number of hugepages available.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesLimitObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesLimitObservable, "system.memory.linux.hugepages.limit", "Total number of hugepages available.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesPageSize", constructorTest(systemconv.NewMemoryLinuxHugepagesPageSize, "system.memory.linux.hugepages.page_size", "System hugepage size in bytes.", "By", false, nil)},
		{"systemconv/MemoryLinuxHugepagesPageSizeObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesPageSizeObservable, "system.memory.linux.hugepages.page_size", "System hugepage size in bytes.", "By", false, nil)},
		{"systemconv/MemoryLinuxHugepagesReserved", constructorTest(systemconv.NewMemoryLinuxHugepagesReserved, "system.memory.linux.hugepages.reserved", "Number of reserved hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesReservedObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesReservedObservable, "system.memory.linux.hugepages.reserved", "Number of reserved hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesSurplus", constructorTest(systemconv.NewMemoryLinuxHugepagesSurplus, "system.memory.linux.hugepages.surplus", "Number of surplus hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesSurplusObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesSurplusObservable, "system.memory.linux.hugepages.surplus", "Number of surplus hugepages.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUsage", constructorTest(systemconv.NewMemoryLinuxHugepagesUsage, "system.memory.linux.hugepages.usage", "Number of hugepages in use by state.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUsageObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesUsageObservable, "system.memory.linux.hugepages.usage", "Number of hugepages in use by state.", "{page}", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUtilization", constructorTest(systemconv.NewMemoryLinuxHugepagesUtilization, "system.memory.linux.hugepages.utilization", "Percentage of hugepages in use by state.", "1", false, nil)},
		{"systemconv/MemoryLinuxHugepagesUtilizationObservable", constructorTest(systemconv.NewMemoryLinuxHugepagesUtilizationObservable, "system.memory.linux.hugepages.utilization", "Percentage of hugepages in use by state.", "1", false, nil)},
		{"systemconv/MemoryLinuxShared", constructorTest(systemconv.NewMemoryLinuxShared, "system.memory.linux.shared", "Shared memory used (mostly by tmpfs).", "By", false, nil)},
		{"systemconv/MemoryLinuxSharedObservable", constructorTest(systemconv.NewMemoryLinuxSharedObservable, "system.memory.linux.shared", "Shared memory used (mostly by tmpfs).", "By", false, nil)},
		{"systemconv/MemoryLinuxSlabUsage", constructorTest(systemconv.NewMemoryLinuxSlabUsage, "system.memory.linux.slab.usage", "Reports the memory used by the Linux kernel for managing caches of frequently used objects.", "By", false, nil)},
		{"systemconv/MemoryLinuxSlabUsageObservable", constructorTest(systemconv.NewMemoryLinuxSlabUsageObservable, "system.memory.linux.slab.usage", "Reports the memory used by the Linux kernel for managing caches of frequently used objects.", "By", false, nil)},
		{"systemconv/MemoryUsage", constructorTest(systemconv.NewMemoryUsage, "system.memory.usage", "Reports memory by state.", "By", false, nil)},
		{"systemconv/MemoryUtilization", constructorTest(systemconv.NewMemoryUtilization, "system.memory.utilization", "Fraction of memory bytes by state.", "1", false, nil)},
		{"systemconv/NetworkConnectionCount", constructorTest(systemconv.NewNetworkConnectionCount, "system.network.connection.count", "The number of connections.", "{connection}", false, nil)},
		{"systemconv/NetworkConnectionCountObservable", constructorTest(systemconv.NewNetworkConnectionCountObservable, "system.network.connection.count", "The number of connections.", "{connection}", false, nil)},
		{"systemconv/NetworkErrors", constructorTest(systemconv.NewNetworkErrors, "system.network.errors", "Count of network errors detected.", "{error}", false, nil)},
		{"systemconv/NetworkErrorsObservable", constructorTest(systemconv.NewNetworkErrorsObservable, "system.network.errors", "Count of network errors detected.", "{error}", false, nil)},
		{"systemconv/NetworkIO", constructorTest(systemconv.NewNetworkIO, "system.network.io", "The number of bytes transmitted and received.", "By", false, nil)},
		{"systemconv/NetworkPacketCount", constructorTest(systemconv.NewNetworkPacketCount, "system.network.packet.count", "The number of packets transferred.", "{packet}", false, nil)},
		{"systemconv/NetworkPacketCountObservable", constructorTest(systemconv.NewNetworkPacketCountObservable, "system.network.packet.count", "The number of packets transferred.", "{packet}", false, nil)},
		{"systemconv/NetworkPacketDropped", constructorTest(systemconv.NewNetworkPacketDropped, "system.network.packet.dropped", "Count of packets that are dropped or discarded even though there was no error.", "{packet}", false, nil)},
		{"systemconv/NetworkPacketDroppedObservable", constructorTest(systemconv.NewNetworkPacketDroppedObservable, "system.network.packet.dropped", "Count of packets that are dropped or discarded even though there was no error.", "{packet}", false, nil)},
		{"systemconv/PagingFaults", constructorTest(systemconv.NewPagingFaults, "system.paging.faults", "The number of page faults.", "{fault}", false, nil)},
		{"systemconv/PagingFaultsObservable", constructorTest(systemconv.NewPagingFaultsObservable, "system.paging.faults", "The number of page faults.", "{fault}", false, nil)},
		{"systemconv/PagingOperations", constructorTest(systemconv.NewPagingOperations, "system.paging.operations", "The number of paging operations.", "{operation}", false, nil)},
		{"systemconv/PagingOperationsObservable", constructorTest(systemconv.NewPagingOperationsObservable, "system.paging.operations", "The number of paging operations.", "{operation}", false, nil)},
		{"systemconv/PagingUsage", constructorTest(systemconv.NewPagingUsage, "system.paging.usage", "UNIX swap or windows pagefile usage.", "By", false, nil)},
		{"systemconv/PagingUsageObservable", constructorTest(systemconv.NewPagingUsageObservable, "system.paging.usage", "UNIX swap or windows pagefile usage.", "By", false, nil)},
		{"systemconv/PagingUtilization", constructorTest(systemconv.NewPagingUtilization, "system.paging.utilization", "Swap (UNIX) or pagefile (windows) utilization.", "1", false, nil)},
		{"systemconv/PagingUtilizationObservable", constructorTest(systemconv.NewPagingUtilizationObservable, "system.paging.utilization", "Swap (UNIX) or pagefile (windows) utilization.", "1", false, nil)},
		{"systemconv/ProcessCount", constructorTest(systemconv.NewProcessCount, "system.process.count", "Total number of processes in each state.", "{process}", false, nil)},
		{"systemconv/ProcessCountObservable", constructorTest(systemconv.NewProcessCountObservable, "system.process.count", "Total number of processes in each state.", "{process}", false, nil)},
		{"systemconv/ProcessCreated", constructorTest(systemconv.NewProcessCreated, "system.process.created", "Total number of processes created over uptime of the host.", "{process}", false, nil)},
		{"systemconv/ProcessCreatedObservable", constructorTest(systemconv.NewProcessCreatedObservable, "system.process.created", "Total number of processes created over uptime of the host.", "{process}", false, nil)},
		{"systemconv/ProcessLimit", constructorTest(systemconv.NewProcessLimit, "system.process.limit", "The maximum number of concurrent processes/tasks allowed by the operating system.", "{thread}", false, nil)},
		{"systemconv/ProcessLimitObservable", constructorTest(systemconv.NewProcessLimitObservable, "system.process.limit", "The maximum number of concurrent processes/tasks allowed by the operating system.", "{thread}", false, nil)},
		{"systemconv/Uptime", constructorTest(systemconv.NewUptime, "system.uptime", "The time the system has been running.", "s", false, nil)},
		{"systemconv/UptimeObservable", constructorTest(systemconv.NewUptimeObservable, "system.uptime", "The time the system has been running.", "s", false, nil)},
		{"vcsconv/ChangeCount", constructorTest(vcsconv.NewChangeCount, "vcs.change.count", "The number of changes (pull requests/merge requests/changelists) in a repository, categorized by their state (e.g. open or merged).", "{change}", false, nil)},
		{"vcsconv/ChangeCountObservable", constructorTest(vcsconv.NewChangeCountObservable, "vcs.change.count", "The number of changes (pull requests/merge requests/changelists) in a repository, categorized by their state (e.g. open or merged).", "{change}", false, nil)},
		{"vcsconv/ChangeDuration", constructorTest(vcsconv.NewChangeDuration, "vcs.change.duration", "The time duration a change (pull request/merge request/changelist) has been in a given state.", "s", false, nil)},
		{"vcsconv/ChangeDurationObservable", constructorTest(vcsconv.NewChangeDurationObservable, "vcs.change.duration", "The time duration a change (pull request/merge request/changelist) has been in a given state.", "s", false, nil)},
		{"vcsconv/ChangeTimeToApproval", constructorTest(vcsconv.NewChangeTimeToApproval, "vcs.change.time_to_approval", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get the first approval.", "s", false, nil)},
		{"vcsconv/ChangeTimeToApprovalObservable", constructorTest(vcsconv.NewChangeTimeToApprovalObservable, "vcs.change.time_to_approval", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get the first approval.", "s", false, nil)},
		{"vcsconv/ChangeTimeToMerge", constructorTest(vcsconv.NewChangeTimeToMerge, "vcs.change.time_to_merge", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get merged into the target(base) ref.", "s", false, nil)},
		{"vcsconv/ChangeTimeToMergeObservable", constructorTest(vcsconv.NewChangeTimeToMergeObservable, "vcs.change.time_to_merge", "The amount of time since its creation it took a change (pull request/merge request/changelist) to get merged into the target(base) ref.", "s", false, nil)},
		{"vcsconv/ContributorCount", constructorTest(vcsconv.NewContributorCount, "vcs.contributor.count", "The number of unique contributors to a repository.", "{contributor}", false, nil)},
		{"vcsconv/ContributorCountObservable", constructorTest(vcsconv.NewContributorCountObservable, "vcs.contributor.count", "The number of unique contributors to a repository.", "{contributor}", false, nil)},
		{"vcsconv/RefCount", constructorTest(vcsconv.NewRefCount, "vcs.ref.count", "The number of refs of type branch or tag in a repository.", "{ref}", false, nil)},
		{"vcsconv/RefCountObservable", constructorTest(vcsconv.NewRefCountObservable, "vcs.ref.count", "The number of refs of type branch or tag in a repository.", "{ref}", false, nil)},
		{"vcsconv/RefLinesDelta", constructorTest(vcsconv.NewRefLinesDelta, "vcs.ref.lines_delta", "The number of lines added/removed in a ref (branch) relative to the ref from the `vcs.ref.base.name` attribute.", "{line}", false, nil)},
		{"vcsconv/RefLinesDeltaObservable", constructorTest(vcsconv.NewRefLinesDeltaObservable, "vcs.ref.lines_delta", "The number of lines added/removed in a ref (branch) relative to the ref from the `vcs.ref.base.name` attribute.", "{line}", false, nil)},
		{"vcsconv/RefRevisionsDelta", constructorTest(vcsconv.NewRefRevisionsDelta, "vcs.ref.revisions_delta", "The number of revisions (commits) a ref (branch) is ahead/behind the branch from the `vcs.ref.base.name` attribute.", "{revision}", false, nil)},
		{"vcsconv/RefRevisionsDeltaObservable", constructorTest(vcsconv.NewRefRevisionsDeltaObservable, "vcs.ref.revisions_delta", "The number of revisions (commits) a ref (branch) is ahead/behind the branch from the `vcs.ref.base.name` attribute.", "{revision}", false, nil)},
		{"vcsconv/RefTime", constructorTest(vcsconv.NewRefTime, "vcs.ref.time", "Time a ref (branch) created from the default branch (trunk) has existed. The `ref.type` attribute will always be `branch`.", "s", false, nil)},
		{"vcsconv/RefTimeObservable", constructorTest(vcsconv.NewRefTimeObservable, "vcs.ref.time", "Time a ref (branch) created from the default branch (trunk) has existed. The `ref.type` attribute will always be `branch`.", "s", false, nil)},
		{"vcsconv/RepositoryCount", constructorTest(vcsconv.NewRepositoryCount, "vcs.repository.count", "The number of repositories in an organization.", "{repository}", false, nil)},
		{"vcsconv/RepositoryCountObservable", constructorTest(vcsconv.NewRepositoryCountObservable, "vcs.repository.count", "The number of repositories in an organization.", "{repository}", false, nil)},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func constructorTest[O, I any](
	newInstrument func(metric.Meter, ...O) (I, error),
	name, description, unit string,
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
				if _, err := newInstrument(m, options...); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(storage, before) {
					t.Error("constructor changed caller option storage")
				}
				if m.config == nil {
					t.Fatal("constructor did not create an instrument")
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
