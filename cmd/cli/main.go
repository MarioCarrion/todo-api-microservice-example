package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MarioCarrion/todo-api-microservice-example/pkg/openapi3"
)

func main() {
	clientOA3 := http.Client{}

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

	count := 1

	for count < 101 {
		priority := openapi3.Low

		_, err := client.CreateTaskWithResponse(context.Background(),
			openapi3.CreateTaskJSONRequestBody{
				Dates: &openapi3.Dates{
					Start: newPtrTime(time.Now()),
					Due:   newPtrTime(time.Now().Add(time.Hour * 24)),
				},
				Description: newPtrStr(fmt.Sprintf("Searchable Task %d", count)),
				Priority:    &priority,
			})
		if err != nil {
			log.Fatalf("Couldn't create task %s", err)
		}

		count++
	}
}
