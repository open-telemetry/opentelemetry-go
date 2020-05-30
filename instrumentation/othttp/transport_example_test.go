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

package othttp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ExampleNewTransport() {
	// Start with a working trace provider
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)

	// Create an http.Client that uses the othttp.Transport
	// wrapped around the http.DefaultTransport
	client := http.Client{
		Transport: NewTransport(http.DefaultTransport),
	}

	// Make a request with our tracing client
	response, err := client.Get("https://postman-echo.com/get")
	if err != nil {
		log.Fatal(err)
	}

	// Read the whole body and close it. The span created by the
	// othttp.Transport does not end until a read from the response
	// body returns io.EOF or the response body is closed.
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", body)
	// body should look like this, with a different "traceparent" value:
	// {"args":{},"headers":{"x-forwarded-proto":"https","host":"postman-echo.com","accept-encoding":"gzip","traceparent":"00-fb1d6775b94db561d9b51adbb3640de5-919c41073ec08f50-01","user-agent":"Go-http-client/1.1","x-forwarded-port":"443"},"url":"https://postman-echo.com/get"}
}
