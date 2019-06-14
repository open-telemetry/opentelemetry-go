package main

import "github.com/open-telemetry/opentelemetry-go/exporter/stderr"

var (
	Observer = stderr.New()
)

func main() {
}
