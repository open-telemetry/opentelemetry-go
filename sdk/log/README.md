# Log SDK

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/log)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/log)

`BatchProcessor` batches by item count by default. Use
`WithExportBatchSizeUnit(BatchExportSizerTypeBytes)` to interpret
`WithExportMaxBatchSize` in serialized bytes when the exporter supports it.
