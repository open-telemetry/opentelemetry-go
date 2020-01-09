package stdout_test

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/stdout"
)

func ExampleNewExportPipeline() {
	// Create a meter
	pusher, err := stdout.NewExportPipeline(stdout.Config{
		PrettyPrint:    true,
		DoNotPrintTime: true,
	})
	if err != nil {
		log.Fatal("Could not initialize stdout exporter:", err)
	}
	defer pusher.Stop()

	ctx := context.Background()

	key := key.New("key")
	meter := scope.NamedMeter(pusher.Meter(), "test")

	// Create and update a single counter:
	counter := meter.NewInt64Counter("a.counter", metric.WithKeys(key))

	counter.Add(ctx, 100, key.String("value"))

	// Output:
	// {
	// 	"updates": [
	// 		{
	// 			"name": "test.a.counter{key=value}",
	// 			"sum": 100
	// 		}
	// 	]
	// }
}
