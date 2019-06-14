package main

import "github.com/open-telemetry/opentelemetry-go/exporter/stdout"

var (
	Observer = stdout.New()
)

func main() {
}
