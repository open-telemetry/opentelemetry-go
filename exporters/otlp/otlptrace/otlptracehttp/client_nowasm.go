// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !wasm

package otlptracehttp // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

import (
	"net"
	"net/http"
	"time"
)

// Keep it in sync with golang's DefaultTransport from net/http! We
// have our own copy to avoid handling a situation where the
// DefaultTransport is overwritten with some different implementation
// of http.RoundTripper or it's modified by other package.
var ourTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
