// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package otlpprofilehttp provides an experimental OTLP profile exporter that
transports OTLP Profiles v1development protobuf payloads over HTTP.

Exporter values should be created using [New].

The environment variables described below can be used for configuration.

OTEL_EXPORTER_OTLP_ENDPOINT (default: "https://localhost:4318") -
target base URL ("/v1development/profiles" is appended) to which the exporter
sends telemetry. The configuration can be overridden by
OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT and by [WithEndpoint], [WithEndpointURL],
and [WithInsecure].

OTEL_EXPORTER_OTLP_PROFILES_ENDPOINT (default: "https://localhost:4318/v1development/profiles") -
target URL to which the exporter sends telemetry. The configuration can be
overridden by [WithEndpoint], [WithEndpointURL], [WithInsecure], and [WithURLPath].

OTEL_EXPORTER_OTLP_INSECURE, OTEL_EXPORTER_OTLP_PROFILES_INSECURE (default: "false") -
setting "true" disables client transport security for the exporter's HTTP connection.

OTEL_EXPORTER_OTLP_HEADERS, OTEL_EXPORTER_OTLP_PROFILES_HEADERS (default: none) -
key-value pairs used as headers associated with HTTP requests (e.g. "k1=v1,k2=v2").

OTEL_EXPORTER_OTLP_TIMEOUT, OTEL_EXPORTER_OTLP_PROFILES_TIMEOUT (default: "10000") -
maximum time in milliseconds the OTLP exporter waits for each export.

OTEL_EXPORTER_OTLP_COMPRESSION, OTEL_EXPORTER_OTLP_PROFILES_COMPRESSION (default: none) -
compression strategy used for the HTTP body ("none" or "gzip").

OTEL_EXPORTER_OTLP_CERTIFICATE, OTEL_EXPORTER_OTLP_PROFILES_CERTIFICATE (default: none) -
filepath to the trusted CA certificate PEM used to verify the server.

OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE, OTEL_EXPORTER_OTLP_PROFILES_CLIENT_CERTIFICATE (default: none) -
filepath to the client certificate PEM for mTLS.

OTEL_EXPORTER_OTLP_CLIENT_KEY, OTEL_EXPORTER_OTLP_PROFILES_CLIENT_KEY (default: none) -
filepath to the client private key PEM for mTLS.
*/
package otlpprofilehttp
