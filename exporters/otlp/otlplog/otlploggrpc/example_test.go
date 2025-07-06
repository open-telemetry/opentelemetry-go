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
	ctx := context.Background()
	var grpcExpOpt []otlploggrpc.Option
	// the filepath to the server's CA certificate
	var caFile string
	// TLS connection
	creds, err := credentials.NewClientTLSFromFile(caFile, "")
	if err != nil {
		panic(err)
	}
	option := otlploggrpc.WithTLSCredentials(creds)
	grpcExpOpt = append(grpcExpOpt, option)
	exp, err := otlploggrpc.New(ctx, grpcExpOpt...)
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
	ctx := context.Background()
	var grpcExpOpt []otlploggrpc.Option
	// the filepath to the server's CA certificate
	var caFile string
	// the filepath to the client's certificate
	var clientCert string
	// the filepath to the client's private key
	var clientKey string
	// mTLS connection
	tlsCfg := tls.Config{}
	// loads CA certificate
	pool := x509.NewCertPool()
	data, err := os.ReadFile(caFile)
	if err != nil {
		panic(err)
	}
	if !pool.AppendCertsFromPEM(data) {
		panic(errors.New("failed to add CA certificate to root CA pool"))
	}
	tlsCfg.RootCAs = pool
	// load client cert and key
	keypair, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		panic(err)
	}
	tlsCfg.Certificates = []tls.Certificate{keypair}
	creds := credentials.NewTLS(&tlsCfg)
	option := otlploggrpc.WithTLSCredentials(creds)
	grpcExpOpt = append(grpcExpOpt, option)
	exp, err := otlploggrpc.New(ctx, grpcExpOpt...)
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
