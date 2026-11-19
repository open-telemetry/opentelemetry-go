// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package otlpprofilegrpc provides an experimental OTLP profile exporter that
transports OTLP Profiles v1development protobuf payloads over gRPC.

All Exporter values must be created with [New].

The environment variables described below can be used for configuration.

OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT (default: "https://localhost:4317") -
target to which the exporter sends telemetry. OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT
takes precedence over OTEL_EXPORTER_OTLP_ENDPOINT. Overridden by [WithEndpoint],
[WithEndpointURL], [WithInsecure], and [WithGRPCConn].

OTEL_EXPORTER_OTLP_INSECURE, OTEL_EXPORTER_OTLP_PROFILES_INSECURE (default: "false") -
setting "true" disables client transport security for the gRPC connection.

OTEL_EXPORTER_OTLP_HEADERS, OTEL_EXPORTER_OTLP_PROFILES_HEADERS (default: none) -
key-value pairs sent as gRPC metadata (e.g. "k1=v1,k2=v2").

OTEL_EXPORTER_OTLP_TIMEOUT, OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT (default: "10000") -
maximum time in milliseconds the exporter waits for each export.

OTEL_EXPORTER_OTLP_COMPRESSION, OTEL_EXPORTER_OTLP_PROFILES_COMPRESSION (default: none) -
gRPC compressor ("none" or "gzip").

OTEL_EXPORTER_OTLP_CERTIFICATE, OTEL_EXPORTER_OTLP_PROFILES_CERTIFICATE (default: none) -
filepath to the trusted CA certificate PEM used to verify the server.

OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE, OTEL_EXPORTER_OTLP_PROFILES_CLIENT_CERTIFICATE (default: none) -
filepath to the client certificate PEM for mTLS.

OTEL_EXPORTER_OTLP_CLIENT_KEY, OTEL_EXPORTER_OTLP_PROFILES_CLIENT_KEY (default: none) -
filepath to the client private key PEM for mTLS.
*/
package otlpprofilegrpc
