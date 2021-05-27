package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"

	"github.com/MarioCarrion/todo-api/pkg/openapi3"
)

func main() {
	initTracer()

	//-

	clientOA3 := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	client, err := openapi3.NewClientWithResponses("http://0.0.0.0:9234", openapi3.WithHTTPClient(&clientOA3))
	if err != nil {
		log.Fatalf("Couldn't instantiate client: %s", err)
	}

	newPtrStr := func(s string) *string {
		return &s
	}

	newPtrTime := func(t time.Time) *time.Time {
		return &t
	}

	//- Create

	priority := openapi3.Priority_low

	respC, err := client.CreateTaskWithResponse(context.Background(),
		openapi3.CreateTaskJSONRequestBody{
			Dates: &openapi3.Dates{
				Start: newPtrTime(time.Now()),
				Due:   newPtrTime(time.Now().Add(time.Hour * 24)),
			},
			Description: newPtrStr("Sleep early"),
			Priority:    &priority,
		})
	if err != nil {
		log.Fatalf("Couldn't create task %s", err)
	}

	fmt.Printf("New Task\n\tID: %s\n", *respC.JSON201.Task.Id)
	fmt.Printf("\tPriority: %s\n", *respC.JSON201.Task.Priority)
	fmt.Printf("\tDescription: %s\n", *respC.JSON201.Task.Description)
	fmt.Printf("\tStart: %s\n", *respC.JSON201.Task.Dates.Start)
	fmt.Printf("\tDue: %s\n", *respC.JSON201.Task.Dates.Due)

	//- Update

	priority = openapi3.Priority_high
	done := true

	_, err = client.UpdateTaskWithResponse(context.Background(),
		*respC.JSON201.Task.Id,
		openapi3.UpdateTaskJSONRequestBody{
			Dates: &openapi3.Dates{
				Start: respC.JSON201.Task.Dates.Start,
				Due:   respC.JSON201.Task.Dates.Due,
			},
			Description: newPtrStr("Sleep early..."),
			Priority:    &priority,
			IsDone:      &done,
		})
	if err != nil {
		log.Fatalf("Couldn't create task %s", err)
	}

	//- Read

	respR, err := client.ReadTaskWithResponse(context.Background(), *respC.JSON201.Task.Id)
	if err != nil {
		log.Fatalf("Couldn't create task %s", err)
	}

	fmt.Printf("Updated Task\n\tID: %s\n", *respR.JSON200.Task.Id)
	fmt.Printf("\tPriority: %s\n", *respR.JSON200.Task.Priority)
	fmt.Printf("\tDescription: %s\n", *respR.JSON200.Task.Description)
	fmt.Printf("\tStart: %s\n", *respR.JSON200.Task.Dates.Start)
	fmt.Printf("\tDue: %s\n", *respR.JSON200.Task.Dates.Due)
	fmt.Printf("\tDone: %t\n", *respR.JSON200.Task.IsDone)

	time.Sleep(10 * time.Second)
}

func initTracer() {
	jaegerEndpoint := "http://localhost:14268/api/traces"

	jaegerExporter, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
	)
	if err != nil {
		log.Fatalln("Couldn't initialize exporter", err)
	}

	// Create stdout exporter to be able to retrieve the collected spans.
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatalln("Couldn't initialize exporter", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithSyncer(jaegerExporter),
		sdktrace.WithResource(resource.NewWithAttributes(attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue("cli"),
		})),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
