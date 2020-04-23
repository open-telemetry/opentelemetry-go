// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlp

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	DefaultCollectorPort uint16 = 55680
	DefaultCollectorHost string = "localhost"
	DefaultNumWorkers    uint   = 1
)

type ExporterOption func(*Config)

type Config struct {
	canDialInsecure    bool
	collectorAddr      string
	compressor         string
	reconnectionPeriod time.Duration
	grpcDialOptions    []grpc.DialOption
	headers            map[string]string
	clientCredentials  credentials.TransportCredentials
	numWorkers         uint
}

// WorkerCount sets the number of Goroutines to use when processing telemetry.
func WorkerCount(n uint) ExporterOption {
	if n == 0 {
		n = DefaultNumWorkers
	}
	return func(cfg *Config) {
		cfg.numWorkers = n
	}
}

// WithInsecure disables client transport security for the exporter's gRPC connection
// just like grpc.WithInsecure() https://pkg.go.dev/google.golang.org/grpc#WithInsecure
// does. Note, by default, client security is required unless WithInsecure is used.
func WithInsecure() ExporterOption {
	return func(cfg *Config) {
		cfg.canDialInsecure = true
	}
}

// WithAddress allows one to set the address that the exporter will
// connect to the collector on. If unset, it will instead try to use
// connect to DefaultCollectorHost:DefaultCollectorPort.
func WithAddress(addr string) ExporterOption {
	return func(cfg *Config) {
		cfg.collectorAddr = addr
	}
}

// WithReconnectionPeriod allows one to set the delay between next connection attempt
// after failing to connect with the collector.
func WithReconnectionPeriod(rp time.Duration) ExporterOption {
	return func(cfg *Config) {
		cfg.reconnectionPeriod = rp
	}
}

// WithCompressor will set the compressor for the gRPC client to use when sending requests.
// It is the responsibility of the caller to ensure that the compressor set has been registered
// with google.golang.org/grpc/encoding. This can be done by encoding.RegisterCompressor. Some
// compressors auto-register on import, such as gzip, which can be registered by calling
// `import _ "google.golang.org/grpc/encoding/gzip"`
func WithCompressor(compressor string) ExporterOption {
	return func(cfg *Config) {
		cfg.compressor = compressor
	}
}

// WithHeaders will send the provided headers when the gRPC stream connection
// is instantiated.
func WithHeaders(headers map[string]string) ExporterOption {
	return func(cfg *Config) {
		cfg.headers = headers
	}
}

// WithTLSCredentials allows the connection to use TLS credentials
// when talking to the server. It takes in grpc.TransportCredentials instead
// of say a Certificate file or a tls.Certificate, because the retrieving
// these credentials can be done in many ways e.g. plain file, in code tls.Config
// or by certificate rotation, so it is up to the caller to decide what to use.
func WithTLSCredentials(creds credentials.TransportCredentials) ExporterOption {
	return func(cfg *Config) {
		cfg.clientCredentials = creds
	}
}

// WithGRPCDialOption opens support to any grpc.DialOption to be used. If it conflicts
// with some other configuration the GRPC specified via the collector the ones here will
// take preference since they are set last.
func WithGRPCDialOption(opts ...grpc.DialOption) ExporterOption {
	return func(cfg *Config) {
		cfg.grpcDialOptions = opts
	}
}
