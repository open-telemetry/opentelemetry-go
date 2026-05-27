# SDK Trace

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/trace)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace)

`BatchSpanProcessor` batches by item count by default. Use
`WithExportBatchSizeUnit(BatchSpanProcessorSizerTypeBytes)` to interpret
`WithMaxExportBatchSize` in serialized bytes when the exporter supports it.
