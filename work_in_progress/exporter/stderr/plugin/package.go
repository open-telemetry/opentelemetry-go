package main

import "github.com/lightstep/opentelemetry-golang-prototype/exporter/stderr"

var (
	Observer = stderr.New()
)

func main() {
}
