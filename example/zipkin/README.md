# Zipkin Exporter Example

Sends spans to zipkin collector.

### Run collector

```sh
docker run -d -p 9411:9411 openzipkin/zipkin
```

### Run client

```sh
go build .
./zipkin
```
