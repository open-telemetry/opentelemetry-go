# Experimental Exemplar Reservoirs

[![pkg.go.dev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/metric/exemplar/x.svg)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric/exemplar/x)

This package contains experimental exemplar reservoirs for the OpenTelemetry Go SDK.

## FixedSizeRoundRobinReservoir

`FixedSizeRoundRobinReservoir` is an experimental reservoir that samples at most a fixed number of exemplars using a round-robin strategy to distribute measurements across independent buckets, each using Algorithm L for sampling. This can be used as a higher-performance drop-in replacement for `FixedSizeReservoir` when some bias is acceptable.
