# "Passthrough" setup for OpenTelemetry

Some Go programs may wish to propagate context without recording spans. To do this in OpenTelemetry, simply install `TextMapPropagators`, but do not install a TracerProvider using the SDK. This works because the default TracerProvider implementation returns a "Non-Recording" span that keeps the context of the caller but does not record spans.

For example, when you initialize your global settings, the following will propagate context without recording spans:

```golang
// Setup Propagators only
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
```

But the following will propagate context _and_ create new, potentially recorded spans:

```golang
// Setup SDK
exp, _ := stdout.New(stdout.WithPrettyPrint())
tp = sdktrace.NewTracerProvider(
    sdktrace.WithBatcher(exp),
)
otel.SetTracerProvider(tp)
// Setup Propagators
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
```

## The Demo

The demo has the following call structure:

`Outer -> Passthrough -> Inner`

If all components had both an SDK and propagators registered, we would expect the trace to look like:

```
|-------outer---------|
 |-Passthrough recv-|
  |Passthrough send|
    |---inner---|
```

However, in this demo, only the outer and inner have TracerProvider backed by the SDK. All components have Propagators set. In this case, we expect to see:

```
|-------outer---------|
    |---inner---|
```

Run the demo as follows:

```console
$ go run .
2022/07/14 16:40:03 Register a global TextMapPropagator, but do not register a global TracerProvider to be in "passthrough" mode.
2022/07/14 16:40:03 The "passthrough" mode propagates the TraceContext and Baggage, but does not record spans.
2022/07/14 16:40:03 The "make outer request" span should be recorded, because it is recorded with a Tracer from the SDK TracerProvider
2022/07/14 16:40:03 The "handle passthrough request" span should NOT be recorded, because it is recorded by a TracerProvider not backed by the SDK.
2022/07/14 16:40:04 The "make outgoing request from passthrough" span should NOT be recorded, because it is recorded by a TracerProvider not backed by the SDK.
2022/07/14 16:40:04 The "handle inner request" span should be recorded, because it is recorded with a Tracer from the SDK TracerProvider
{
	"Name": "handle inner request",
	"SpanContext": {
		"TraceID": "3296fd65d8383262134e53eac727c3e5",
		"SpanID": "253a29450291d34b",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "3296fd65d8383262134e53eac727c3e5",
		"SpanID": "5becf3c6908649e2",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": true
	},
	"SpanKind": 1,
	"StartTime": "2022-07-14T16:40:04.217495-04:00",
	"EndTime": "2022-07-14T16:40:05.217610898-04:00",
	"Attributes": null,
	"Events": null,
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:passthrough"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.8.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "example/passthrough/inner",
		"Version": "",
		"SchemaURL": ""
	}
}
{
	"Name": "make outer request",
	"SpanContext": {
		"TraceID": "3296fd65d8383262134e53eac727c3e5",
		"SpanID": "5becf3c6908649e2",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "00000000000000000000000000000000",
		"SpanID": "0000000000000000",
		"TraceFlags": "00",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 1,
	"StartTime": "2022-07-14T16:40:03.216172-04:00",
	"EndTime": "2022-07-14T16:40:05.217618421-04:00",
	"Attributes": null,
	"Events": null,
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:passthrough"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.8.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "example/passthrough/outer",
		"Version": "",
		"SchemaURL": ""
	}
}
```
