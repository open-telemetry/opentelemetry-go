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

package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// Handle CTRL+C gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setup OpenTelemetry.
	serviceName := "dice"
	serviceVersion := "0.1.0"
	otelShutdown, err := setupOTelSDK(ctx, serviceName, serviceVersion)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Create HTTP server with requests' base context canceled on server shutdown.
	rCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return rCtx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		// Add HTTP instrumentation for the whole server.
		Handler: otelhttp.NewHandler(http.DefaultServeMux, "/"),
	}
	srv.RegisterOnShutdown(cancel)

	// Start the server.
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}

// handleFunc is a replacement for [net/http.HandleFunc]
// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
func handleFunc(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	// Configure the "http.route" for the HTTP instrumentation.
	handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
	http.Handle(pattern, handler)
}
