package kafka_test

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	kafkaTask "github.com/MarioCarrion/todo-api-microservice-example/internal/kafka"
)

// TestNewTask tests the constructor
func TestNewTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		topicName string
	}{
		{
			name:      "creates new task publisher",
			topicName: "test-topic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock producer configuration (won't actually connect)
			producer, err := kafka.NewProducer(&kafka.ConfigMap{
				"bootstrap.servers": "localhost:9092",
				"client.id":         "test-client",
			})
			if err != nil {
				t.Skip("Skipping test: requires Kafka producer setup")
			}
			defer producer.Close()

			taskPub := kafkaTask.NewTask(producer, tt.topicName)

			if taskPub == nil {
				t.Fatal("expected non-nil Task publisher")
			}
		})
	}
}

// Note: Testing the Created, Updated, and Deleted methods requires either:
// 1. A running Kafka broker for integration testing
// 2. Mocking the kafka.Producer interface (which is complex due to its internal implementation)
// 3. Using testcontainers for Kafka
//
// These are integration tests rather than unit tests. For true unit testing, the Task
// struct would need to accept an interface for the producer rather than *kafka.Producer.
//
// The methods themselves are simple wrappers around the publish method that:
// - Encode the task data as JSON
// - Call producer.Produce()
// - Return any errors
//
// Integration tests would be more appropriate for validating the actual message publishing behavior.
