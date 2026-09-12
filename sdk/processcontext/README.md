# processcontext

[![Go Reference](https://pkg.go.dev/badge/go.opentelemetry.io/otel/sdk/processcontext.svg)](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/processcontext)

Package `processcontext` implements the SDK publishing side of the OTEL_CTX
process context mechanism defined in
[OTEP 4719](https://github.com/open-telemetry/opentelemetry-specification/blob/main/oteps/profiles/4719-process-ctx.md).

It makes the process's [resource attributes](https://opentelemetry.io/docs/specs/otel/resource/sdk/)
discoverable by external readers (e.g. OBI or eBPF profiler) via a named
memory-mapped region visible in `/proc/<pid>/maps`.

## Supported platforms

**Linux only.** `NewPublisher` returns an error on other platforms.

## How it works

When `NewPublisher` is called it:

1. Creates a memory-mapped region (using `memfd_create` when available,
   otherwise an anonymous private mapping).
2. Serializes the resource attributes as a `ProcessContext` protobuf payload
   (OTLP wire format).
3. Writes a 32-byte header containing the signature `OTEL_CTX`, format
   version, payload size, a monotonic timestamp, and a pointer to the payload.
4. Names the mapping `OTEL_CTX` via `prctl(PR_SET_VMA_ANON_NAME)` so it
   appears in `/proc/<pid>/maps` for external readers to discover.

The timestamp field uses a seqlock-compatible protocol: it is set to zero
while the payload is being written and restored to a strictly increasing
non-zero value once the write is complete.

## Usage

```go
import (
    "context"

    "go.opentelemetry.io/otel/sdk/processcontext"
    "go.opentelemetry.io/otel/sdk/resource"
)

func main() {
    res, _ := resource.New(context.Background(),
        resource.WithTelemetrySDK(),
        resource.WithProcess(),
    )

    pub, err := processcontext.NewPublisher(res)
    if err != nil {
        // Not fatal; external readers simply won't see the attributes.
        _ = err
    } else {
        defer pub.Shutdown(context.Background())
    }

    // ...run the application...

    // When resource attributes change (e.g., late K8s label injection):
    _ = pub.Update(updatedRes)
}
```
