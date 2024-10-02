Instructions on how to run instrumented and uninstrumented examples.

## Prerequisites

- [Go](https://golang.org/dl/) installed on your system.
- Necessary permissions to execute shell scripts.

## Usage

The `run.sh` script accepts one argument to determine which example to run:

- `instrumented`
- `uninstrumented`

### Running the Instrumented Example

The instrumented example includes OpenTelemetry instrumentation for collecting telemetry data like traces and metrics. 

To run the instrumented example, execute:

```bash
./run.sh instrumented
```

### Running the Uninstrumented Example

The uninstrumented example is the exact same application, without OTEL instrumentation.

To run the instrumented example, execute:

```bash
./run.sh uninstrumented
```





