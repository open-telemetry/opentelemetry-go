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

// dice is the "Roll the dice" getting started example application.
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Create HTTP request multiplexer.
	mux := http.NewServeMux()
	registerHandleFunc(mux, "/rolldice", handle)

	// Create HTTP server with requests' base context canceled on server shutdown.
	rCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return rCtx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		// Add HTTP instrumentation for the whole server.
		Handler: otelhttp.NewHandler(mux, "/"),
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
	return srv.Shutdown(context.Background())
}

func registerHandleFunc(mux *http.ServeMux, pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	// Configure the "http.route" for the HTTP instrumentation.
	handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
	mux.Handle(pattern, handler)
}

func handle(w http.ResponseWriter, r *http.Request) {
	roll := 1 + rand.Intn(6)
	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}
