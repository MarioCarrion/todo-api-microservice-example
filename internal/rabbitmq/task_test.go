package rabbitmq_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/streadway/amqp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	rabbitmqtask "github.com/MarioCarrion/todo-api-microservice-example/internal/rabbitmq"
)

// dockerImage must match the docker image listed in `compose.rabbitmq.yml`.
const dockerImage = "rabbitmq:3.11.10-management-alpine"

func TestTask_Created_Integration(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	rmqContainer, err := rabbitmq.Run(ctx, dockerImage)
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %v", err)
	}

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(rmqContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

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

	t.Cleanup(func() { conn.Close() })

	channel, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}

	t.Cleanup(func() { channel.Close() })

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
	taskPub := rabbitmqtask.NewTask(channel)

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
	t.Parallel()

	ctx := t.Context()

	rmqContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine")
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %v", err)
	}

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(rmqContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	connStr, err := rmqContainer.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	conn, err := amqp.Dial(connStr)
	if err != nil {
		t.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	channel, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}

	t.Cleanup(func() { channel.Close() })

	err = channel.ExchangeDeclare("tasks", "topic", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("failed to declare exchange: %v", err)
	}

	taskPub := rabbitmqtask.NewTask(channel)

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
	t.Parallel()

	ctx := t.Context()

	rmqContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12-management-alpine")
	if err != nil {
		t.Fatalf("failed to start rabbitmq container: %v", err)
	}

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(rmqContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	connStr, err := rmqContainer.AmqpURL(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	conn, err := amqp.Dial(connStr)
	if err != nil {
		t.Fatalf("failed to connect to rabbitmq: %v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	channel, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}

	t.Cleanup(func() { channel.Close() })

	err = channel.ExchangeDeclare("tasks", "topic", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("failed to declare exchange: %v", err)
	}

	taskPub := rabbitmqtask.NewTask(channel)

	err = taskPub.Deleted(ctx, "test-789")
	if err != nil {
		t.Fatalf("Failed to publish deleted event: %v", err)
	}
}

//-

var setupClient = sync.OnceValue(func() RabbitMQClient { //nolint: gochecknoglobals
	var res RabbitMQClient

	ctx := context.Background()

	rmqContainer, err := rabbitmq.Run(ctx, dockerImage)
	if err != nil {
		res.err = fmt.Errorf("failed to start rabbitmq container: %w", err)

		return res
	}

	res.container = rmqContainer

	connStr, err := rmqContainer.AmqpURL(ctx)
	if err != nil {
		res.err = fmt.Errorf("failed to get connection string: %w", err)

		return res
	}

	//- RabbitMQ connection
	conn, err := amqp.Dial(connStr)
	if err != nil {
		res.err = fmt.Errorf("failed to connect to rabbitmq: %w", err)

		return res
	}

	res.connection = conn

	//- RabbitMQ Channel

	channel, err := conn.Channel()
	if err != nil {
		res.err = fmt.Errorf("failed to open channel: %w", err)

		return res
	}

	res.channel = channel

	return res

})

type RabbitMQClient struct {
	container  *rabbitmq.RabbitMQContainer
	channel    *amqp.Channel
	connection *amqp.Connection
	err        error
}

func (r *RabbitMQClient) Teardown() error {
	var err error

	if r.channel != nil {
		err = r.channel.Close()
	}

	if r.connection != nil {
		err = errors.Join(err, fmt.Errorf("failed to close connection: %w", r.connection.Close()))
	}

	if r.container != nil {
		if err1 := testcontainers.TerminateContainer(r.container); err1 != nil {
			err = errors.Join(err, fmt.Errorf("failed to terminate container: %w", err1))
		}
	}

	return err
}
