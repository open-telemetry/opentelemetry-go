# Semantic Convention Changes

The `go.opentelemetry.io/otel/semconv/v1.28.0` should be a drop-in replacement for `go.opentelemetry.io/otel/semconv/v1.27.0` with the following exceptions.

## Dropped deprecations

The following declarations have been deprecated in the [OpenTelemetry Semantic Conventions].
Refer to the respective documentation in that repository for deprecation instructions for each type.

- `AndroidStateBackground`
- `AndroidStateCreated`
- `AndroidStateForeground`
- `AndroidStateKey`
- `DBCosmosDBStatusCode`
- `DBCosmosDBStatusCodeKey`
- `GenAICompletion`
- `GenAICompletionKey`
- `GenAIPrompt`
- `GenAIPromptKey`

### Dropping the `aspnetcore` namespace

The following declarations have been dropped in `go.opentelemetry.io/otel/semconv/v1.28.0`.
These semantic conventions were determined to have not have use in Go applications.

- `ASPNETCoreDiagnosticsExceptionResultAborted`
- `ASPNETCoreDiagnosticsExceptionResultHandled`
- `ASPNETCoreDiagnosticsExceptionResultKey`
- `ASPNETCoreDiagnosticsExceptionResultSkipped`
- `ASPNETCoreDiagnosticsExceptionResultUnhandled`
- `ASPNETCoreDiagnosticsHandlerType`
- `ASPNETCoreDiagnosticsHandlerTypeKey`
- `ASPNETCoreRateLimitingPolicy`
- `ASPNETCoreRateLimitingPolicyKey`
- `ASPNETCoreRateLimitingResultAcquired`
- `ASPNETCoreRateLimitingResultEndpointLimiter`
- `ASPNETCoreRateLimitingResultGlobalLimiter`
- `ASPNETCoreRateLimitingResultKey`
- `ASPNETCoreRateLimitingResultRequestCanceled`
- `ASPNETCoreRequestIsUnhandled`
- `ASPNETCoreRequestIsUnhandledKey`
- `ASPNETCoreRoutingIsFallback`
- `ASPNETCoreRoutingIsFallbackKey`
- `ASPNETCoreRoutingMatchStatusFailure`
- `ASPNETCoreRoutingMatchStatusKey`
- `ASPNETCoreRoutingMatchStatusSuccess`
- `AspnetcoreDiagnosticsExceptionsDescription`
- `AspnetcoreDiagnosticsExceptionsName`
- `AspnetcoreDiagnosticsExceptionsUnit`
- `AspnetcoreRateLimitingActiveRequestLeasesDescription`
- `AspnetcoreRateLimitingActiveRequestLeasesName`
- `AspnetcoreRateLimitingActiveRequestLeasesUnit`
- `AspnetcoreRateLimitingQueuedRequestsDescription`
- `AspnetcoreRateLimitingQueuedRequestsName`
- `AspnetcoreRateLimitingQueuedRequestsUnit`
- `AspnetcoreRateLimitingRequestLeaseDurationDescription`
- `AspnetcoreRateLimitingRequestLeaseDurationName`
- `AspnetcoreRateLimitingRequestLeaseDurationUnit`
- `AspnetcoreRateLimitingRequestTimeInQueueDescription`
- `AspnetcoreRateLimitingRequestTimeInQueueName`
- `AspnetcoreRateLimitingRequestTimeInQueueUnit`
- `AspnetcoreRateLimitingRequestsDescription`
- `AspnetcoreRateLimitingRequestsName`
- `AspnetcoreRateLimitingRequestsUnit`
- `AspnetcoreRoutingMatchAttemptsDescription`
- `AspnetcoreRoutingMatchAttemptsName`
- `AspnetcoreRoutingMatchAttemptsUnit`

If you use any of these declarations in your Go application, please [open an issue] describing your use-case.

### Dropping the `jvm` namespace

The following declarations have been dropped in `go.opentelemetry.io/otel/semconv/v1.28.0`.
These semantic conventions were determined to have not have use in Go applications.

- `JVMBufferPoolName`
- `JVMBufferPoolNameKey`
- `JVMGCAction`
- `JVMGCActionKey`
- `JVMGCName`
- `JVMGCNameKey`
- `JVMMemoryPoolName`
- `JVMMemoryPoolNameKey`
- `JVMMemoryTypeHeap`
- `JVMMemoryTypeKey`
- `JVMMemoryTypeNonHeap`
- `JVMThreadDaemon`
- `JVMThreadDaemonKey`
- `JVMThreadStateBlocked`
- `JVMThreadStateKey`
- `JVMThreadStateNew`
- `JVMThreadStateRunnable`
- `JVMThreadStateTerminated`
- `JVMThreadStateTimedWaiting`
- `JVMThreadStateWaiting`
- `JvmBufferCountDescription`
- `JvmBufferCountName`
- `JvmBufferCountUnit`
- `JvmBufferMemoryLimitDescription`
- `JvmBufferMemoryLimitName`
- `JvmBufferMemoryLimitUnit`
- `JvmBufferMemoryUsageDescription`
- `JvmBufferMemoryUsageName`
- `JvmBufferMemoryUsageUnit`
- `JvmBufferMemoryUsedDescription`
- `JvmBufferMemoryUsedName`
- `JvmBufferMemoryUsedUnit`
- `JvmCPUCountDescription`
- `JvmCPUCountName`
- `JvmCPUCountUnit`
- `JvmCPURecentUtilizationDescription`
- `JvmCPURecentUtilizationName`
- `JvmCPURecentUtilizationUnit`
- `JvmCPUTimeDescription`
- `JvmCPUTimeName`
- `JvmCPUTimeUnit`
- `JvmClassCountDescription`
- `JvmClassCountName`
- `JvmClassCountUnit`
- `JvmClassLoadedDescription`
- `JvmClassLoadedName`
- `JvmClassLoadedUnit`
- `JvmClassUnloadedDescription`
- `JvmClassUnloadedName`
- `JvmClassUnloadedUnit`
- `JvmGcDurationDescription`
- `JvmGcDurationName`
- `JvmGcDurationUnit`
- `JvmMemoryCommittedDescription`
- `JvmMemoryCommittedName`
- `JvmMemoryCommittedUnit`
- `JvmMemoryInitDescription`
- `JvmMemoryInitName`
- `JvmMemoryInitUnit`
- `JvmMemoryLimitDescription`
- `JvmMemoryLimitName`
- `JvmMemoryLimitUnit`
- `JvmMemoryUsedAfterLastGcDescription`
- `JvmMemoryUsedAfterLastGcName`
- `JvmMemoryUsedAfterLastGcUnit`
- `JvmMemoryUsedDescription`
- `JvmMemoryUsedName`
- `JvmMemoryUsedUnit`
- `JvmSystemCPULoad1mDescription`
- `JvmSystemCPULoad1mName`
- `JvmSystemCPULoad1mUnit`
- `JvmSystemCPUUtilizationDescription`
- `JvmSystemCPUUtilizationName`
- `JvmSystemCPUUtilizationUnit`
- `JvmThreadCountDescription`
- `JvmThreadCountName`
- `JvmThreadCountUnit`

If you use any of these declarations in your Go application, please [open an issue] describing your use-case.

### Dropping the `nodejs` namespace

The following declarations have been dropped in `go.opentelemetry.io/otel/semconv/v1.28.0`.
These semantic conventions were determined to have not have use in Go applications.

- `NodejsEventloopDelayMaxDescription`
- `NodejsEventloopDelayMaxName`
- `NodejsEventloopDelayMaxUnit`
- `NodejsEventloopDelayMeanDescription`
- `NodejsEventloopDelayMeanName`
- `NodejsEventloopDelayMeanUnit`
- `NodejsEventloopDelayMinDescription`
- `NodejsEventloopDelayMinName`
- `NodejsEventloopDelayMinUnit`
- `NodejsEventloopDelayP50Description`
- `NodejsEventloopDelayP50Name`
- `NodejsEventloopDelayP50Unit`
- `NodejsEventloopDelayP90Description`
- `NodejsEventloopDelayP90Name`
- `NodejsEventloopDelayP90Unit`
- `NodejsEventloopDelayP99Description`
- `NodejsEventloopDelayP99Name`
- `NodejsEventloopDelayP99Unit`
- `NodejsEventloopDelayStddevDescription`
- `NodejsEventloopDelayStddevName`
- `NodejsEventloopDelayStddevUnit`
- `NodejsEventloopUtilizationDescription`
- `NodejsEventloopUtilizationName`
- `NodejsEventloopUtilizationUnit`

If you use any of these declarations in your Go application, please [open an issue] describing your use-case.

### Dropping the `v8js` namespace

The following declarations have been dropped in `go.opentelemetry.io/otel/semconv/v1.28.0`.
These semantic conventions were determined to have not have use in Go applications.

- `V8JSGCTypeIncremental`
- `V8JSGCTypeKey`
- `V8JSGCTypeMajor`
- `V8JSGCTypeMinor`
- `V8JSGCTypeWeakcb`
- `V8JSHeapSpaceNameCodeSpace`
- `V8JSHeapSpaceNameKey`
- `V8JSHeapSpaceNameLargeObjectSpace`
- `V8JSHeapSpaceNameMapSpace`
- `V8JSHeapSpaceNameNewSpace`
- `V8JSHeapSpaceNameOldSpace`
- `V8jsGcDurationDescription`
- `V8jsGcDurationName`
- `V8jsGcDurationUnit`
- `V8jsHeapSpaceAvailableSizeDescription`
- `V8jsHeapSpaceAvailableSizeName`
- `V8jsHeapSpaceAvailableSizeUnit`
- `V8jsHeapSpacePhysicalSizeDescription`
- `V8jsHeapSpacePhysicalSizeName`
- `V8jsHeapSpacePhysicalSizeUnit`
- `V8jsMemoryHeapLimitDescription`
- `V8jsMemoryHeapLimitName`
- `V8jsMemoryHeapLimitUnit`
- `V8jsMemoryHeapUsedDescription`
- `V8jsMemoryHeapUsedName`
- `V8jsMemoryHeapUsedUnit`

If you use any of these declarations in your Go application, please [open an issue] describing your use-case.

### Renames

The following renames have been introduced to better match Go and industry naming standards.
Be sure to update any use from `go.opentelemetry.io/otel/semconv/v1.27.0` with the equivalent in `go.opentelemetry.io/otel/semconv/v1.28.0`.

| `v1.28.0` | `v1.27.0` |
| --- | --- |
| DBSystemDb2 | DBSystemDB2|
| GenAIClientOperationDurationDescription | GenAiClientOperationDurationDescription|
| GenAIClientOperationDurationName | GenAiClientOperationDurationName|
| GenAIClientOperationDurationUnit | GenAiClientOperationDurationUnit|
| GenAIClientTokenUsageDescription | GenAiClientTokenUsageDescription|
| GenAIClientTokenUsageName | GenAiClientTokenUsageName|
| GenAIClientTokenUsageUnit | GenAiClientTokenUsageUnit|
| GenAIServerRequestDurationDescription | GenAiServerRequestDurationDescription|
| GenAIServerRequestDurationName | GenAiServerRequestDurationName|
| GenAIServerRequestDurationUnit | GenAiServerRequestDurationUnit|
| GenAIServerTimePerOutputTokenDescription | GenAiServerTimePerOutputTokenDescription|
| GenAIServerTimePerOutputTokenName | GenAiServerTimePerOutputTokenName|
| GenAIServerTimePerOutputTokenUnit | GenAiServerTimePerOutputTokenUnit|
| GenAIServerTimeToFirstTokenDescription | GenAiServerTimeToFirstTokenDescription|
| GenAIServerTimeToFirstTokenName | GenAiServerTimeToFirstTokenName|
| GenAIServerTimeToFirstTokenUnit | GenAiServerTimeToFirstTokenUnit|
| GoMemoryGCGoalDescription | GoMemoryGcGoalDescription|
| GoMemoryGCGoalName | GoMemoryGcGoalName|
| GoMemoryGCGoalUnit | GoMemoryGcGoalUnit|

[OpenTelemetry Semantic Conventions]: https://github.com/open-telemetry/semantic-conventions
[open an issue]: https://github.com/open-telemetry/opentelemetry-go/issues/new?template=Blank+issue
