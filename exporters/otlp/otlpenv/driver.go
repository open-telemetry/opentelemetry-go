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

package otlpenv

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	grpccreds "google.golang.org/grpc/credentials"
	grpcgzip "google.golang.org/grpc/encoding/gzip"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlphttp"
)

type protocol int

const (
	protocolGRPC protocol = iota
	protocolHTTP
)

type compression int

const (
	compressionNone compression = iota
	compressionGzip
)

const (
	defaultEndpoint = "localhost:4317"
	envVarPrefix    = "OTEL_EXPORTER_OTLP_"
)

type endpointRawSetup struct {
	protocol    string
	compression string
	endpoint    string
	insecure    string
	certFile    string
	headers     string
	timeout     string
}

type endpointSetup struct {
	protocol    protocol
	compression compression
	endpoint    string
	insecure    bool
	certFile    string
	headers     map[string]string
	timeout     time.Duration
}

func (rs endpointRawSetup) process() endpointSetup {
	if rs.protocol == "" {
		rs.protocol = "gzip"
	}
	if rs.endpoint == "" {
		rs.endpoint = defaultEndpoint
	}
	if rs.insecure == "" {
		rs.insecure = "false"
	}
	return endpointSetup{
		protocol:    stringToProtocol(rs.protocol),
		compression: stringToCompression(rs.compression),
		endpoint:    rs.endpoint,
		insecure:    stringToBoolean(rs.insecure),
		certFile:    rs.certFile,
		headers:     stringToHeaders(rs.headers),
		timeout:     stringToDuration(rs.timeout),
	}
}

type driverRawSetup struct {
	span   endpointRawSetup
	metric endpointRawSetup
}

func NewDriver(opts ...Option) otlp.ProtocolDriver {
	cfg := newConfig(opts...)
	envMap := envToMap(cfg.env)
	setup := driverRawSetup{}

	type entry struct {
		name      string
		forSpan   *string
		forMetric *string
	}
	for _, e := range []entry{
		{
			name:      "ENDPOINT",
			forSpan:   &setup.span.endpoint,
			forMetric: &setup.metric.endpoint,
		},
		{
			name:      "PROTOCOL",
			forSpan:   &setup.span.protocol,
			forMetric: &setup.metric.protocol,
		},
		{
			name:      "COMPRESSION",
			forSpan:   &setup.span.compression,
			forMetric: &setup.metric.compression,
		},
		{
			name:      "INSECURE",
			forSpan:   &setup.span.insecure,
			forMetric: &setup.metric.insecure,
		},
		{
			name:      "CERTIFICATE",
			forSpan:   &setup.span.certFile,
			forMetric: &setup.metric.certFile,
		},
		{
			name:      "HEADERS",
			forSpan:   &setup.span.headers,
			forMetric: &setup.metric.headers,
		},
		{
			name:      "TIMEOUT",
			forSpan:   &setup.span.timeout,
			forMetric: &setup.metric.timeout,
		},
	} {
		generalName := fmt.Sprintf("%s%s", envVarPrefix, e.name)
		if value, ok := envMap[generalName]; ok {
			*e.forSpan = value
			*e.forMetric = value
			continue
		}
		spanName := fmt.Sprintf("%sSPAN_%s", envVarPrefix, e.name)
		metricName := fmt.Sprintf("%sMETRIC_%s", envVarPrefix, e.name)
		if value, ok := envMap[spanName]; ok {
			*e.forSpan = value
		}
		if value, ok := envMap[metricName]; ok {
			*e.forMetric = value
		}
	}

	if setup.span == setup.metric {
		return endpointSetupToDriver(setup.span.process())
	}

	spanDriver := endpointSetupToDriver(setup.span.process())
	metricDriver := endpointSetupToDriver(setup.metric.process())
	splitCfg := otlp.SplitConfig{
		ForTraces:  spanDriver,
		ForMetrics: metricDriver,
	}
	return otlp.NewSplitDriver(splitCfg)
}

func envToMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			otel.Handle(fmt.Errorf(`otlpenv: invalid environment variable entry "%s" (should be in form of "key=value"), ignoring it`, e))
			continue
		}
		if strings.HasPrefix(parts[0], envVarPrefix) {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func stringToProtocol(s string) protocol {
	switch s {
	case "grpc":
		return protocolGRPC
	case "http":
		return protocolHTTP
	default:
		otel.Handle(fmt.Errorf(`otlpenv: invalid protocol "%s" (should be either "grpc" or "http"), falling back to "grpc"`, s))
		return protocolGRPC
	}
}

func stringToCompression(s string) compression {
	switch s {
	case "gzip":
		return compressionGzip
	case "":
		return compressionNone
	default:
		otel.Handle(fmt.Errorf(`otlpenv: invalid compression "%s" (should be either "gzip" or unset/empty), falling back to "gzip"`, s))
		return compressionNone
	}
}

func stringToBoolean(s string) bool {
	switch s {
	case "true":
		return true
	case "false":
		return false
	default:
		otel.Handle(fmt.Errorf(`otlpenv: invalid boolean value "%s" (should be either "true" or "false"), falling back to "false"`, s))
		return false
	}
}

func stringToHeaders(s string) map[string]string {
	kvs := strings.Split(s, ",")
	m := make(map[string]string)
	for _, kv := range kvs {
		kvPair := strings.Split(kv, "=")
		if len(kvPair) != 2 {
			otel.Handle(fmt.Errorf(`otlpenv: invalid header entry "%s" (should be in form of "key=value"), ignoring it`, s))
			continue
		}
		name := strings.TrimSpace(kvPair[0])
		value := strings.TrimSpace(kvPair[1])
		m[name] = value
	}
	if len(m) == 0 {
		return nil
	}
	return m
}

func endpointSetupToDriver(eps endpointSetup) otlp.ProtocolDriver {
	switch eps.protocol {
	case protocolGRPC:
		return endpointSetupToGRPCDriver(eps)
	case protocolHTTP:
		return endpointSetupToHTTPDriver(eps)
	default:
		panic(fmt.Sprintf("otlpenv: bug, invalid protocol after processing it (%d)", eps.protocol))
	}
}

func endpointSetupToGRPCDriver(eps endpointSetup) otlp.ProtocolDriver {
	var opts []otlpgrpc.Option
	switch eps.compression {
	case compressionNone:
		// nothing to do
	case compressionGzip:
		opts = append(opts, otlpgrpc.WithCompressor(grpcgzip.Name))
	}
	if eps.insecure {
		opts = append(opts, otlpgrpc.WithInsecure())
	} else if eps.certFile != "" {
		tlsCfg, err := tlsConfigFromCertFile(eps.certFile)
		if err != nil {
			otel.Handle(err)
		} else {
			opts = append(opts, otlpgrpc.WithTLSCredentials(grpccreds.NewTLS(tlsCfg)))
		}
	}
	if eps.headers != nil {
		opts = append(opts, otlpgrpc.WithHeaders(eps.headers))
	}
	if eps.timeout > 0 {
		opts = append(opts, otlpgrpc.WithTimeout(eps.timeout))
	}
	opts = append(opts, otlpgrpc.WithEndpoint(eps.endpoint))
	return otlpgrpc.NewDriver(opts...)
}

func endpointSetupToHTTPDriver(eps endpointSetup) otlp.ProtocolDriver {
	var (
		opts            []otlphttp.Option
		httpCompression otlphttp.Compression
	)
	switch eps.compression {
	case compressionNone:
		httpCompression = otlphttp.NoCompression
	case compressionGzip:
		httpCompression = otlphttp.GzipCompression
	}
	opts = append(opts, otlphttp.WithCompression(httpCompression))
	if eps.insecure {
		opts = append(opts, otlphttp.WithInsecure())
	} else if eps.certFile != "" {
		tlsCfg, err := tlsConfigFromCertFile(eps.certFile)
		if err != nil {
			otel.Handle(err)
		} else {
			opts = append(opts, otlphttp.WithTLSClientConfig(tlsCfg))
		}
	}
	if eps.headers != nil {
		opts = append(opts, otlphttp.WithHeaders(eps.headers))
	}
	if eps.timeout > 0 {
		opts = append(opts, otlphttp.WithTimeout(eps.timeout))
	}
	opts = append(opts, otlphttp.WithEndpoint(eps.endpoint))
	return otlphttp.NewDriver(opts...)
}

func tlsConfigFromCertFile(certFile string) (*tls.Config, error) {
	// stolen from grpc
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("otlpenv: failed to read a certificate file %s: %w", certFile, err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("otlpenv: failed to append certificates from file %s", certFile)
	}
	return &tls.Config{RootCAs: cp}, nil
}

func stringToDuration(s string) time.Duration {
	if s == "" {
		return 0
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		otel.Handle(fmt.Errorf(`otlpenv: failed to parse duration string "%s": %w`, s, err))
		return 0
	}
	if d < 0 {
		otel.Handle(fmt.Errorf(`otlpenv: duration "%s" is negative, ignoring`, s))
		return 0
	}
	return d
}
