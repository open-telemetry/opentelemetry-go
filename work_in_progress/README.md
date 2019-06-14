This is not a complete implementation of the OpenTracing and OpenCensus API surface areas. I'm posting this here, now, to have as a point of reference for several of the issues in the specification repo. Some of the high-lights of the approach taken here:

* Always associate current span context with stats and metrics events
* Introduce a low-level "observer" exporter
* Avoid excessive memory allocations
* Avoid buffering state with API objects
* Use `context.Context` to propagate context tags and active scope
* Introduce a "reader" implementation to interpret "observer"-exported events and build state
* Use a common `KeyValue` construct for span attributes, context tags, resource definitions, log fields, and metric fields
* Support logging API w/o a current span context
* Support for golang `plugin` package to load implementations
* Example use of golang's `net/http/httptrace` w/ @iredelmeier's [tracecontext.go](https://github.com/lightstep/tracecontext.go) package

The first bullet about associating current span context and stats/metrics events bridges the tracing data model with the metrics data model. The APIs here would make this association not an option, as the `stats.Record` API takes a context, which passes through to the observer, which could choose to use the span-context association. The prototype includes a stderr exporter that writes a debugging log of events to the console. One of the critical features here, enabled by the low-level observer exporter, is that Span start events can be logged in chronological order, not as the span finishes.

To run the examples, first build the stderr tracer plugin (requires Linux or OS X):

```
(cd ./exporter/stderr/plugin && make)
```

then set the `OPENTELEMETRY_LIB` environment variable to the .so file in that directory:

```
OPENTELEMETRY_LIB=./exporter/stderr/plugin/stderr.so go run ./example/client/client.go
```

The output of this program reads:

```
2019/06/04 16-51-32.075834 start say hello, a root span [ component=main whatevs=yesss service=client username=donuts span_id=4d6..d52 trace_id=786..21d ]
2019/06/04 16-51-32.076073 start http.request < parent_span_id=b80..c03 service=client component=main whatevs=yesss > [ username=donuts span_id=b80..c03 trace_id=786..21d ]
2019/06/04 16-51-32.076208 start http.getconn < parent_span_id=365..2d1 service=client component=main whatevs=yesss > [ http.host=localhost:7777 username=donuts span_id=365..2d1 trace_id=786..21d ]
2019/06/04 16-51-32.076399 start http.dns < parent_span_id=57e..8d8 service=client component=main whatevs=yesss > [ username=donuts span_id=57e..8d8 trace_id=786..21d ]
2019/06/04 16-51-32.172726 finish http.dns (96.32742ms) [ username=donuts span_id=57e..8d8 trace_id=786..21d ]
2019/06/04 16-51-32.172871 start http.connect < parent_span_id=886..01e service=client component=main whatevs=yesss > [ username=donuts span_id=886..01e trace_id=786..21d ]
2019/06/04 16-51-32.173570 finish http.connect (698.912µs) [ username=donuts span_id=886..01e trace_id=786..21d ]
2019/06/04 16-51-32.173654 modify attr [ http.host=localhost:7777 http.remote=127.0.0.1:7777 username=donuts span_id=365..2d1 trace_id=786..21d ]
2019/06/04 16-51-32.173757 modify attr [ http.host=localhost:7777 http.remote=127.0.0.1:7777 http.local=127.0.0.1:61376 username=donuts span_id=365..2d1 trace_id=786..21d ]
2019/06/04 16-51-32.173858 finish http.getconn (97.649927ms) [ http.host=localhost:7777 http.remote=127.0.0.1:7777 http.local=127.0.0.1:61376 username=donuts span_id=365..2d1 trace_id=786..21d ]
2019/06/04 16-51-32.174024 start http.headers < parent_span_id=940..294 service=client component=main whatevs=yesss > [ username=donuts span_id=940..294 trace_id=786..21d ]
2019/06/04 16-51-32.174050 modify attr [ http.host=localhost:7777 username=donuts span_id=b80..c03 trace_id=786..21d ]
2019/06/04 16-51-32.174066 modify attr [ http.user-agent=Go-http-client/1.1 http.host=localhost:7777 username=donuts span_id=b80..c03 trace_id=786..21d ]
2019/06/04 16-51-32.174117 modify attr [ http.host=localhost:7777 http.user-agent=Go-http-client/1.1 http.traceparent=00-78629a0f5f3f164fd5104dc76695721d-4d65822107fcfd52-01 username=donuts span_id=b80..c03 trace_id=786..21d ]
2019/06/04 16-51-32.174139 modify attr [ http.host=localhost:7777 http.user-agent=Go-http-client/1.1 http.traceparent=00-78629a0f5f3f164fd5104dc76695721d-4d65822107fcfd52-01 http.tracestate=ot@username=donuts username=donuts span_id=b80..c03 trace_id=786..21d ]
2019/06/04 16-51-32.174165 modify attr [ http.user-agent=Go-http-client/1.1 http.traceparent=00-78629a0f5f3f164fd5104dc76695721d-4d65822107fcfd52-01 http.tracestate=ot@username=donuts http.accept-encoding=gzip http.host=localhost:7777 username=donuts span_id=b80..c03 trace_id=786..21d ]
2019/06/04 16-51-32.174183 start http.send < parent_span_id=0c6..7a0 whatevs=yesss service=client component=main > [ username=donuts span_id=0c6..7a0 trace_id=786..21d ]
2019/06/04 16-51-32.174201 finish http.send (18.179µs) [ username=donuts span_id=0c6..7a0 trace_id=786..21d ]
2019/06/04 16-51-32.175226 start http.receive < parent_span_id=a68..b99 service=client component=main whatevs=yesss > [ username=donuts span_id=a68..b99 trace_id=786..21d ]
2019/06/04 16-51-32.175384 finish http.receive (157.839µs) [ username=donuts span_id=a68..b99 trace_id=786..21d ]
2019/06/04 16-51-32.175438 finish say hello (99.603385ms) [ service=client component=main whatevs=yesss username=donuts span_id=4d6..d52 trace_id=786..21d ]
```
