// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploggrpc_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

func Example() {
	ctx := context.Background()
	exp, err := otlploggrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	global.SetLoggerProvider(provider)

	// From here, the provider can be used by instrumentation to collect
	// telemetry.
}

// Demonstrates how to configure the exporter using self-signed certificates for TLS connections.
func Example_selfSignedCertificates_TLS() {
	// Variables provided by the user.
	var (
		caFile string // The filepath to the server's CA certificate.
	)

	ctx := context.Background()

	// Configure TLS connection.
	creds, err := credentials.NewClientTLSFromFile(caFile, "")
	if err != nil {
		panic(err)
	}
	exp, err := otlploggrpc.New(ctx, otlploggrpc.WithTLSCredentials(creds))
	if err != nil {
		panic(err)
	}

	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	global.SetLoggerProvider(provider)

	// From here, the provider can be used by instrumentation to collect
	// telemetry.
}

// Demonstrates how to configure the exporter using self-signed certificates for mutual TLS (mTLS) connections.
func Example_selfSignedCertificates_mTLS() {
	// Variables provided by the user.
	var (
		caFile     string // The filepath to the server's CA certificate.
		clientCert string // The filepath to the client's certificate.
		clientKey  string // The filepath to the client's private key.
	)

	ctx := context.Background()

	// Configure mTLS connection.
	tlsCfg := tls.Config{}
	pool := x509.NewCertPool()
	data, err := os.ReadFile(caFile)
	if err != nil {
		panic(err)
	}
	if !pool.AppendCertsFromPEM(data) {
		panic(errors.New("failed to add CA certificate to root CA pool"))
	}
	tlsCfg.RootCAs = pool
	keypair, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		panic(err)
	}
	tlsCfg.Certificates = []tls.Certificate{keypair}
	creds := credentials.NewTLS(&tlsCfg)
	exp, err := otlploggrpc.New(ctx, otlploggrpc.WithTLSCredentials(creds))
	if err != nil {
		panic(err)
	}

	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	global.SetLoggerProvider(provider)

	// From here, the provider can be used by instrumentation to collect
	// telemetry.
}
