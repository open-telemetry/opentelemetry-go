# gRPC Tracing Example

Traces client and server calls via interceptors.

### Compile .proto

Only required if the service definition (.proto) changes.

```sh
cd ./example/grpc

# protobuf v1.3.2
protoc -I api --go_out=plugins=grpc,paths=source_relative:./api api/hello-service.proto
```

### Run server

```sh
cd ./example/grpc

go run ./server
```

### Run client

```sh
go run ./client
```