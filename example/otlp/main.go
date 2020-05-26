// dummy application for testing opentelemetry Go agent + collector

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// "github.com/open-telemetry/opentelemetry-collector/translator/conventions"

	"go.opentelemetry.io/otel/api/global"
	// "go.opentelemetry.io/otel/api/kv"
	// "go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// tracestdout "go.opentelemetry.io/otel/exporters/trace/stdout"
)

func initExporter() {
	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:55680"),
	)
	// exp, err := tracestdout.NewExporter(tracestdout.Options{PrettyPrint: true})
	handleErr(err, "Failed to create exporter: $v")

	// defer handleErr(exp.Stop(), "Failed to stop exporter: %v")

	provider, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
        sdktrace.WithSyncer(exp),
	)
	handleErr(err, "Failed to create trace provider: %v")

	global.SetTraceProvider(provider)
}

func main() {
	initExporter()
	tracer := global.Tracer("mage-sense")

	ctx, span := tracer.Start(context.Background(), "Archmage-Overlord")
	for i := 0; i < 10; i++ {
		_, innerSpan := tracer.Start(ctx, fmt.Sprintf("Minion-%d", i))
		log.Println("Minions hard at work, scribing...")
		<-time.After(time.Second)
		innerSpan.End()
	}

	span.End()
	<-time.After(time.Second)

    log.Println("Spell-scroll scribed!")
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
