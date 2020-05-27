# Example
Here is a collection of OpenTelemtry-Go examples showcasing basic functionality.

## OTLP
This example demonstrates how to export trace and metric data from an
application using OpenTelemetry's own wire protocal OTLP. We will also walk
you through configuring a collector to accept OTLP exports.

### How to run?

#### Prequisites
- go >=1.14 installed
- `GOPATH` is configured
- (optional) OpenTelemetry collector is available

#### Start the Application
An example application is included in `example/otlp`. It simulates the process
of scribing a spell scroll (e.g. in [D&D](https://roll20.net/compendium/dnd5e/Spell%20Scroll#content)).
The application has been instrumented and exports both trace and metric data
via OTLP to any listening receiver. To run it:

```sh
go get -d go.opentelemetry.io/otel
cd $GOPATH/go.opentelemetry.io/otel/example/otlp
go run main.go
```

The application is currently configured to transmit exported data to
`localhost:55680`. See [main.go](example/otlp/main.go) for full details.

Note, if you don't have a receiver configured to take in metric data, the
application will complain about being unable to connect.

#### (optional) Configure the Collector
Follow the instructions [on the
website](https://opentelemetry.io/docs/collector/about/) to install a working
instance of the collector. This example assumes you have the collector installed
locally.

To configure the collector to accept OTLP traffic from our application,
ensure that it has the following configs:

```yaml
receivers:
    otlp:
        endpoint: 0.0.0.0:55680   # listens to localhost:55680

    # potentially other receivers

service:
    pipelines:

        traces:
            receivers:
                - otlp
                # potentially other receivers
            processors: # whatever processors you need
            exporters: # wherever you want your data to go

        metrics:
            receivers:
                -otlp
                # potentially other receivers
            processors: etc
            exporters: etc

    # other services
```

An example config has been provided at
[example-otlp-config.yaml](example/otlp/example-otlp-config.yaml).

Then to run:
```sh
./[YOUR_COLLECTOR_BINARY]  --config [PATH_TO_CONFIG]
```

If you use the example config, it's set to export to `stdout`. If you run
the collector on the same machine as the example application, you should
see trace and metric outputs from the collector.



## HTTP
This is a simple example that demonstrates tracing http request from client to server. The example
shows key aspects of tracing such as:

- Root Span (on Client)
- Child Span (on Client)
- Child Span from a Remote Parent (on Server)
- SpanContext Propagation (from Client to Server)
- Span Events
- Span Attributes

Example uses
- open-telemetry SDK as trace instrumentation provider,
- httptrace plugin to facilitate tracing http request on client and server
- http trace_context propagation to propagate SpanContext on the wire.
- stdout exporter to print information about spans in the terminal

### How to run?

#### Prequisites

- go 1.13 installed
- GOPATH is configured.

#### 1 Download git repo
```
GO111MODULE="" go get -d go.opentelemetry.io/otel
```

#### 2 Start Server
```
cd $GOPATH/src/go.opentelemetry.io/otel/example/http/
go run ./server/server.go
```

#### 3 Start Client
```
cd $GOPATH/src/go.opentelemetry.io/otel/example/http/
go run ./client/client.go
```

#### 4 Check traces in stdout

The spans should be visible in stdout in the order that they were exported.
