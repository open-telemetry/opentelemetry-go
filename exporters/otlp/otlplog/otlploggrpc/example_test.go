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

// ExampleWithTLSCredentials demonstrates how to configure the exporter with certificates, including self-signed certificates.
func ExampleWithTLSCredentials() {
	ctx := context.Background()
	var grpcExpOpt []otlploggrpc.Option
	// the filepath to the server's CA certificate
	caFile := os.Getenv("CUSTOM_SERVER_CA_CERTIFICATE")
	// the filepath to the client's certificate
	clientCert := os.Getenv("CUSTOM_CLIENT_CERTIFICATE")
	// the filepath to the client's private key
	clientKey := os.Getenv("CUSTOM_CLIENT_KEY")
	if caFile != "" && clientCert != "" && clientKey != "" {
		// mTLS connection
		tlsCfg := tls.Config{
			InsecureSkipVerify: false,
		}
		// loads CA certificate
		pool, _ := x509.SystemCertPool()
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
	} else if caFile != "" {
		// TLS connection
		creds, err := credentials.NewClientTLSFromFile(caFile, "")
		if err != nil {
			panic(err)
		}
		option := otlploggrpc.WithTLSCredentials(creds)
		grpcExpOpt = append(grpcExpOpt, option)
	}
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
