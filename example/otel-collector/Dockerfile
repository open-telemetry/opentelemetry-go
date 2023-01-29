FROM golang:1.18-alpine AS base

COPY . /src

WORKDIR /src/example/otel-collector

RUN go build -o main .

FROM alpine:latest

COPY --from=base /src/example/otel-collector/main /app/main

CMD ["/app/main"]