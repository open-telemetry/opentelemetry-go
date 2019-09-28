# Example

## HTTP
This is a simple example that demonstrates tracing http request from client to server. The example
shows key aspects of tracing such as 
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
- jaeger exporter to export spans to visualize and store them.

### How to run?

#### Prequisites

- go 1.12 installed 
- GOPATH is configured.

#### 1 Download git repo
```
GO111MODULE="" go get -d go.opentelemetry.io
```

#### 2 Start Server
```
cd $GOPATH/src/go.opentelemetry.io/example/http/
go run ./server/server.go
``` 

#### 3 Start Client
```
cd $GOPATH/src/go.opentelemetry.io/example/http/
go run ./client/client.go
``` 

#### 4 Check traces in stdout

The spans should be visible in stdout in the order that they were exported.


