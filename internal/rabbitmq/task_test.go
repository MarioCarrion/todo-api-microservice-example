package rabbitmq_test

import (
	"testing"

	"github.com/streadway/amqp"

	rabbitmqTask "github.com/MarioCarrion/todo-api-microservice-example/internal/rabbitmq"
)

// TestNewTask tests the constructor.
func TestNewTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "creates new task publisher",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock channel (note: this will be nil but tests constructor logic)
			// In a real integration test, we'd connect to RabbitMQ
			var channel *amqp.Channel

			taskPub := rabbitmqTask.NewTask(channel)

			if taskPub == nil {
				t.Fatal("expected non-nil Task publisher")
			}
		})
	}
}

// Note: Testing the Created, Updated, and Deleted methods requires either:
// 1. A running RabbitMQ broker for integration testing
// 2. Mocking the amqp.Channel interface (complex due to internal implementation)
// 3. Using testcontainers for RabbitMQ
//
// These are integration tests rather than unit tests. For true unit testing, the Task
// struct would need to accept an interface for the channel rather than *amqp.Channel.
//
// The methods themselves:
// - Created: publishes "Task.Created" event with task data encoded using gob
// - Updated: publishes "Task.Updated" event with task data encoded using gob
// - Deleted: publishes "Task.Deleted" event with task ID encoded using gob
//
// Integration tests would be more appropriate for validating the actual message publishing behavior.
