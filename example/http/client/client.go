// Copyright 2019, OpenTelemetry Authors
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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/api/trace"
	"github.com/open-telemetry/opentelemetry-go/plugin/httptrace"

	"google.golang.org/grpc/codes"

	_ "github.com/open-telemetry/opentelemetry-go/exporter/loader"
)

var (
	tracer = trace.GlobalTracer().
		WithService("client").
		WithComponent("main").
		WithResources(
			tag.New("whatevs").String("yesss"),
		)
)

func main() {
	client := http.DefaultClient

	ctx := tag.NewContext(context.Background(),
		tag.Insert(tag.New("username").String("donuts")),
	)

	var body []byte

	err := tracer.WithSpan(ctx, "say hello",
		func(ctx context.Context) error {
			req, _ := http.NewRequest("GET", "http://localhost:7777/hello", nil)

			ctx, req, inj := httptrace.W3C(ctx, req)

			trace.Inject(ctx, inj)

			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			body, err = ioutil.ReadAll(res.Body)
			res.Body.Close()
			trace.Active(ctx).SetStatus(codes.OK)

			return err
		})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", body)
}
