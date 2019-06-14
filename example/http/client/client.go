package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/api/trace"
	"github.com/open-telemetry/opentelemetry-go/plugin/httptrace"

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

			return err
		})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", body)
}
