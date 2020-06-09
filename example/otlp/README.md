# OTLP Example
This example demonstrates how to export trace and metric data from an
application using OpenTelemetry's own wire protocol
[OTLP](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/protocol/README.md).
We will also walk you through configuring a collector to accept OTLP exports.

### How to run?

#### Prequisites
- go >=1.13 installed
- OpenTelemetry collector is available

#### Configure the Collector
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
[example-otlp-config.yaml](otlp/example-otlp-config.yaml).

Then to run:
```sh
./[YOUR_COLLECTOR_BINARY]  --config [PATH_TO_CONFIG]
```

If you use the example config, it's set to export to `stdout`. If you run
the collector on the same machine as the example application, you should
see trace and metric outputs from the collector.

#### Start the Application
An example application is included in this directory. It simulates the process
of scribing a spell scroll (e.g. in [D&D](https://roll20.net/compendium/dnd5e/Spell%20Scroll#content)).
The application has been instrumented and exports both trace and metric data
via OTLP to any listening receiver. To run it:

```sh
go get -d go.opentelemetry.io/otel
cd $GOPATH/go.opentelemetry.io/otel/example/otlp
go run main.go
```

The application is currently configured to transmit exported data to
`localhost:55680`. See [main.go](otlp/main.go) for full details.

After starting the application, you should see trace and metric log output
on the collector.

Note, if the receiver is incorrectly configured to take in metric data, the
application may complain about being unable to connect.
