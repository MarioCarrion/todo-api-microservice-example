//go:build integration

package rabbitmq_test

import (
	"context"
	"testing"

	"github.com/streadway/amqp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	rabbitmqTask "github.com/MarioCarrion/todo-api-microservice-example/internal/rabbitmq"
)

func TestTask_Created_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := t.Context()

	// Start RabbitMQ container
	rmqContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine")
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(rmqContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := rmqContainer.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Create RabbitMQ connection
	conn, err := amqp.Dial(connStr)
	if err != nil {
		t.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}
	defer channel.Close()

	// Declare the exchange
	err = channel.ExchangeDeclare(
		"tasks", // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		t.Fatalf("failed to declare exchange: %v", err)
	}

	// Create task publisher
	taskPub := rabbitmqTask.NewTask(channel)

	// Test Created method
	task := internal.Task{
		ID:          "test-123",
		Description: "Test task",
		Priority:    internal.PriorityHigh.Pointer(),
		IsDone:      false,
	}

	err = taskPub.Created(ctx, task)
	if err != nil {
		t.Fatalf("Failed to publish created event: %v", err)
	}
}

func TestTask_Updated_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := t.Context()

	rmqContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine")
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(rmqContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := rmqContainer.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	conn, err := amqp.Dial(connStr)
	if err != nil {
		t.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}
	defer channel.Close()

	err = channel.ExchangeDeclare("tasks", "topic", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("failed to declare exchange: %v", err)
	}

	taskPub := rabbitmqTask.NewTask(channel)

	task := internal.Task{
		ID:          "test-456",
		Description: "Updated task",
		IsDone:      true,
	}

	err = taskPub.Updated(ctx, task)
	if err != nil {
		t.Fatalf("Failed to publish updated event: %v", err)
	}
}

func TestTask_Deleted_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := t.Context()

	rmqContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine")
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(rmqContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := rmqContainer.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	conn, err := amqp.Dial(connStr)
	if err != nil {
		t.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}
	defer channel.Close()

	err = channel.ExchangeDeclare("tasks", "topic", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("failed to declare exchange: %v", err)
	}

	taskPub := rabbitmqTask.NewTask(channel)

	err = taskPub.Deleted(ctx, "test-789")
	if err != nil {
		t.Fatalf("Failed to publish deleted event: %v", err)
	}
}
