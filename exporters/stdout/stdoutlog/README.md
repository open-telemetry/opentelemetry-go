# STDOUT Log Exporter

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/otel/exporters/stdout/stdoutlog)](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/stdout/stdoutlog)

Exports honor context cancellation and deadlines while waiting for the configured
writer. Go's `io.Writer` interface cannot interrupt a write that has already
started, so that write may finish after the export returns. The exporter retains
only that record's encoded JSON and allows at most one write in progress. Later
exports wait for it using their own contexts.

`ForceFlush` and the first `Shutdown` call wait for an in-progress write with
their supplied contexts. After shutdown starts, further `Shutdown` and
`ForceFlush` calls do nothing, even if the write is still running. The application
owns the writer; the exporter does not close it, call its `Flush` method, or change
its deadlines. Use a writer that returns from `Write` to ensure all background
work can finish.
