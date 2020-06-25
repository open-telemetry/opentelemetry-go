# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- This CHANGELOG file to track all changes in the project going forward.

### Changed

- 

### Removed

- 

## [0.6.0] - 2020-05-21
## [0.5.0] - 2020-05-13
## [0.4.3] - 2020-04-24
## [0.4.2] - 2020-03-31
## [0.4.1] - 2020-03-31
## [0.4.0] - 2020-03-30
## [0.3.0] - 2020-03-21
## [0.2.3] - 2020-03-04

### Added

- `RecordError` method on `Span`s in the trace API to Simplify adding error events to spans. (#473)
- Configurable push frequency for exporters setup pipeline. (#504)

### Changed

- Rename the `exporter` directory to `exporters`.
   The `go.opentelemetry.io/otel/exporter/trace/jaeger` package was mistakenly released with a `v1.0.0` tag instead of `v0.1.0`.
   This resulted in all subsequent releases not becoming the default latest.
   A consequence of this was that all `go get`s pulled in the incompatible `v0.1.0` release of that package when pulling in more recent packages from other otel packages.
   Renaming the `exporter` directory to `exporters` fixes this issue by renaming the package and therefore clearing any existing dependency tags.
   Consequentially, this action also renames *all* exporter packages.

### Removed

- The `CorrelationContextHeader` constant in the `correlation` package is no longer exported. (#503)

## [0.2.2] - 2020-02-27

### Added

- `HTTPSupplier` interface in the propagation API to specify methods to retrieve and store a single value for a key to be associated with a carrier. (#467)
- `HTTPExtractor` interface in the propagation API to extract information from an `HTTPSupplier` into a context. (#467)
- `HTTPInjector` interface in the propagation API to inject information into an `HTTPSupplier.` (#467)
- `Config` and configuring `Option` to the propagator API. (#467)
- `Propagators` interface in the propagation API to contain the set of injectors and extractors for all supported carrier formats. (#467)
- `HTTPPropagator` interface in the propagation API to inject and extract from an `HTTPSupplier.` (#467)
- `WithInjectors` and `WithExtractors` functions to the propagator API to configure injectors and extractors to use. (#467)
- `ExtractHTTP` and `InjectHTTP` functions to apply configured HTTP extractors and injectors to a passed context. (#467)
- Histogram aggregator. (#433)
- `DefaultPropagator` function and have it return `trace.TraceContext` as the default context propagator. (#456)
- `AlwaysParentSample` sampler to the trace API. (#455)
- `WithNewRoot` option function to the trace API to specify the created span should be considered a root span. (#451)


### Changed

- Renamed `WithMap` to `ContextWithMap` in the correlation package. (#481)
- Renamed `FromContext` to `MapFromContext` in the correlation package. (#481)
- Move correlation context propagation to correlation package. (#479)
- Do not default to putting remote span context into links. (#480)
- Propagators extrac
- `Tracer.WithSpan` updated to accept `StartOptions`. (#472)
- Renamed `MetricKind` to `Kind` to not stutter in the type usage. (#432)
- Renamed the `export` package to `metric` to match directory structure. (#432)
- Rename the `api/distributedcontext` package to `api/correlation`. (#444)
- Rename the `api/propagators` package to `api/propagation`. (#444)
- Move the propagators from the `propagators` package into the `trace` API package. (#444)
- Update `Float64Gauge`, `Int64Gauge`, `Float64Counter`, `Int64Counter`, `Float64Measure`, and `Int64Measure` metric methods to use value receivers instead of pointers. (#462)
- Moved all dependencies of tools package to a tools directory. (#466)

### Removed

- Binary propagators. (#467)
- NOOP propagator. (#467)

### Fixed

- Upgraded `github.com/golangci/golangci-lint` from `v1.21.0` to `v1.23.6` in `tools/`. (#492)
- Fix a possible nil-dereference crash (#478)
- Correct comments for `InstallNewPipeline` in the stdout exporter. (#483)
- Correct comments for `InstallNewPipeline` in the dogstatsd exporter. (#484)
- Correct comments for `InstallNewPipeline` in the prometheus exporter. (#482)
- Initialize `onError` based on `Config` in prometheus exporter. (#486)
- Correct module name in prometheus exporter README. (#475)
- Remove tracer name prefix from span names. (#430)
- Fix `aggregator_test.go` import package comment. (#431)
- Improved detail in stdout exporter. (#436)
- Fix a dependency issue (generate target should depend on stringer, not lint target) in Makefile. (#442)
- Reorders the Makefile targets within `precommit` target so we generate files and build the code before doing linting, so we can get much nicer errors about syntax errors from the compiler. (#442)
- Reword function documentation in gRPC plugin. (#446)
- Send the `span.kind` tag to Jaeger from the jaeger exporter. (#441)
- Fix `metadataSupplier` in the jaeger exporter to overwrite the header if existing instead of appending to it. (#441)
- Upgraded to Go 1.13 in CI. (#465)
- Correct opentelemetry.io URL in trace SDK documentation. (#464)
- Refactored reference counting logic in SDK determination of stale records. (#468)
- Add call to `runtime.Gosched` in instrument `acquireHandle` logic to not block the collector. (#469)

## [0.2.1.1] - 2020-01-13

### Fixes

- Use stateful batcher on Prometheus exporter fixing regresion introduced in #395. (#428)

## [0.2.1] - 2020-01-08

### Added

- Global meter forwarding implementation.
   This enables deferred initialization for metric instruments registered before the first Meter SDK is installed. (#392)
- Global trace forwarding implementation.
   This enables deferred initialization for tracers registered before the first Trace SDK is installed. (#406)
- Standardize export pipeline creation in all exporters. (#395)
- A testing, organization, and comments for 64-bit field alignment. (#418)
- Script to tag all modules in the project. (#414)

### Changed

- Renamed `propagation` package to `propagators`. (#362)
- Renamed `B3Propagator` propagator to `B3`. (#362)
- Renamed `TextFormatPropagator` propagator to `TextFormat`. (#362)
- Renamed `BinaryPropagator` propagator to `Binary`. (#362)
- Renamed `BinaryFormatPropagator` propagator to `BinaryFormat`. (#362)
- Renamed `NoopTextFormatPropagator` propagator to `NoopTextFormat`. (#362)
- Renamed `TraceContextPropagator` propagator to `TraceContext`. (#362)
- Renamed `SpanOption` to `StartOption` in the trace API. (#369)
- Renamed `StartOptions` to `StartConfig` in the trace API. (#369)
- Renamed `EndOptions` to `EndConfig` in the trace API. (#369)
- `Number` now has a pointer receiver for its methods. (#375)
- Renamed `CurrentSpan` to `SpanFromContext` in the trace API. (#379)
- Renamed `SetCurrentSpan` to `ContextWithSpan` in the trace API. (#379)
- Renamed `Message` in Event to `Name` in the trace API. (#389)
- Prometheus exporter no longer aggregates metrics, instead it only exports them. (#385)
- Renamed `HandleImpl` to `BoundInstrumentImpl` in the metric API. (#400)
- Renamed `Float64CounterHandle` to `Float64CounterBoundInstrument` in the metric API. (#400)
- Renamed `Int64CounterHandle` to `Int64CounterBoundInstrument` in the metric API. (#400)
- Renamed `Float64GaugeHandle` to `Float64GaugeBoundInstrument` in the metric API. (#400)
- Renamed `Int64GaugeHandle` to `Int64GaugeBoundInstrument` in the metric API. (#400)
- Renamed `Float64MeasureHandle` to `Float64MeasureBoundInstrument` in the metric API. (#400)
- Renamed `Int64MeasureHandle` to `Int64MeasureBoundInstrument` in the metric API. (#400)
- Renamed `Release` method for bound instruments in the metric API to `Unbind`. (#400)
- Renamed `AcquireHandle` method for bound instruments in the metric API to `Bind`. (#400)
- Renamed the `File` option in the stdout exporter to `Writer`. (#404)
- Renamed all `Options` to `Config` for all metric exports where this wasn't already the case.

### Fixed

- Aggregator import path corrected. (#421)
- Correct links in README. (#368)
- The README was updated to match latest code changes in its examples. (#374)
- Don't capitalize error statements. (#375)
- Fix ignored errors. (#375)
- Fix ambiguous variable naming. (#375)
- Remove unnecessary type casting. (#375)
- Use named parameters. (#375)
- Updated release schedule. (#378)
- Correct http-stackdriver example module name. (#394)
- Removed the `http.request` span in `httptrace` package. (#397)
- Add comments in the metrics SDK (#399)
- Initialize checkpoint when creating ddsketch aggregator to prevent panic when merging into a empty one. (#402) (#403)
- Add documentation of compatible exporters in the README. (#405)
- Typo fix. (#408)
- Simplify span check logic in SDK tracer implementation. (#419)

## [0.2.0] - 2019-12-03

### Added

- Unary gRPC tracing example. (#351)
- Prometheus exporter. (#334)
- Dogstatsd metrics exporter. (#326)

### Changed

- Rename `MaxSumCount` aggregation to `MinMaxSumCount` and add the `Min` interface for this aggregation. (#352)
- Rename `GetMeter` to `Meter`. (#357)
- Rename `HTTPTraceContextPropagator` to `TraceContextPropagator`. (#355)
- Rename `HTTPB3Propagator` to `B3Propagator`. (#355)
- Rename `HTTPTraceContextPropagator` to `TraceContextPropagator`. (#355)
- Move `/global` package to `/api/global`. (#356)
- Rename `GetTracer` to `Tracer`. (#347)

### Removed

- `SetAttribute` from the `Span` interface in the trace API. (#361)
- `AddLink` from the `Span` interface in the trace API. (#349)
- `Link` from the `Span` interface in the trace API. (#349)

### Fixed

- Exclude example directories from coverage report. (#365)
- Lint make target now implements automatic fixes with `golangci-lint` before a second run to report the remaining issues. (#360)
- Drop `GO111MODULE` environment variable in Makefile as Go 1.13 is the project specified minimum version and this is environment variable is not needed for that version of Go. (#359)
- Run the race checker for all test. (#354)
- Redundant commands in the Makefile are removed. (#354)
- Split the `generate` and `lint` targets of the Makefile. (#354)
- Renames `circle-ci` target to more generic `ci` in Makefile. (#354)
- Add example Prometheus binary to gitignore. (#358)
- Support negative numbers with the `MaxSumCount`. (#335)
- Resolve race conditions in `push_test.go` identified in #339. (#340)
- Use `/usr/bin/env bash` as a shebang in scripts rather than `/bin/bash`. (#336)
- Trace benchmark now tests both `AlwaysSample` and `NeverSample`.
   Previously it was testing `AlwaysSample` twice. (#325)
- Trace benchmark now uses a `[]byte` for `TraceID` to fix failing test. (#325)
- Added a trace benchmark to test variadic functions in `setAttribute` vs `setAttributes` (#325)
- The `defaultkeys` batcher was only using the encoded label set as its map key while building a checkpoint.
   This allowed distinct label sets through, but any metrics sharing a label set could be overwritten or merged incorrectly.
   This was corrected. (#333)


## [0.1.2] - 2019-11-18

### Fixed

- Optimized the `simplelru` map for attributes to reduce the number of allocations. (#328)
- Removed unnecessary unslicing of parameters that are already a slice. (#324)

## [0.1.1] - 2019-11-18

This release contains a Metrics SDK with stdout exporter and supports basic aggregations such as counter, gauges, array, maxsumcount, and ddsketch.

### Added

- Metrics stdout export pipeline. (#265)
- Array aggregation for raw measure metrics. (#282)
- The core.Value now have a `MarshalJSON` method. (#281)

### Removed

- `WithService`, `WithResources`, and `WithComponent` methods of tracers. (#314)
- Prefix slash in `Tracer.Start()` for the Jaeger example. (#292)

### Changed

- Allocation in LabelSet construction to reduce GC overhead. (#318)
- `trace.WithAttributes` to append values instead of replacing (#315)
- Use a formula for tolerance in sampling tests. (#298)
- Move export types into trace and metric-specific sub-directories. (#289)
- `SpanKind` back to being based on an `int` type. (#288)

### Fixed

- URL to OpenTelemetry website in README. (#323)
- Name of othttp default tracer. (#321)
- `ExportSpans` for the stackdriver exporter now handles `nil` context. (#294)
- CI modules cache to correctly restore/save from/to the cache. (#316)
- Fix metric SDK race condition between `LoadOrStore` and the assignment `rec.recorder = i.meter.exporter.AggregatorFor(rec)`. (#293)
- README now reflects the new code structure introduced with these changes. (#291)
- Make the basic example work. (#279)

## [0.1.0] - 2019-11-04

This is the first release of open-telemetry go library.
It contains api and sdk for trace and meter.

### Added

- Initial OpenTelemetry trace and metric API prototypes.
- Initial OpenTelemetry trace, metric, and export SDK packages.
- A wireframe bridge to support compatibility with OpenTracing.
- Example code for a basic, http-stackdriver, http, jaeger, and named tracer setup.
- Exporters for Jaeger, Stackdriver, and stdout.
- Propagators for binary, B3, and trace-context protocols.
- Project information and guidelines in the form of a README and CONTRIBUTING.
- Tools to build the project and a Makefile to automate the process.
- Apache-2.0 license.
- CircleCI build CI manifest files.
- CODEOWNERS file to track owners of this project.


[Unreleased]: https://github.com/open-telemetry/opentelemetry-go/compare/v0.6.0...HEAD
[0.6.0]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.6.0
[0.5.0]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.5.0
[0.4.3]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.4.3
[0.4.2]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.4.2
[0.4.1]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.4.1
[0.4.0]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.4.0
[0.3.0]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.3.0
[0.2.3]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.2.3
[0.2.2]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.2.2
[0.2.1.1]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.2.1.1
[0.2.1]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.2.1
[0.2.0]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.2.0
[0.1.2]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.1.2
[0.1.1]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.1.1
[0.1.0]: https://github.com/open-telemetry/opentelemetry-go/releases/tag/v0.1.0
