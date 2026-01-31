package rabbitmq_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	rabbitmqtask "github.com/MarioCarrion/todo-api-microservice-example/internal/rabbitmq"
)

const (
	// dockerImage must match the docker image listed in `compose.rabbitmq.yml`.
	dockerImage = "rabbitmq:3.11.10-management-alpine"

	rabbitMQConsumerName = "test-rabbitmq-consumer"
)

func TestMain(m *testing.M) {
	client := setupClient()
	if client.err != nil {
		panic(fmt.Sprintf("Failed to set up RabbitMQ client: %v", client.err))
	}

	code := m.Run()

	if err := client.Teardown(); err != nil {
		panic(fmt.Sprintf("Failed to close RabbitMQ client: %v", err))
	}

	os.Exit(code)
}

func TestTask_All(t *testing.T) { //nolint: tparallel
	t.Parallel()

	ctx := t.Context()

	client := setupClient()
	if client.err != nil {
		t.Fatalf("Failed to setupClient: %v", client.err)
	}

	tests := []struct {
		name   string
		call   func(t *testing.T, channel *amqp.Channel)
		verify func(t *testing.T, routingKey string, body []byte)
	}{
		{
			name: "Created",
			call: func(t *testing.T, channel *amqp.Channel) {
				t.Helper()

				taskPub := rabbitmqtask.NewTask(channel)

				task := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.ValueToPointer(internal.PriorityHigh),
					IsDone:      true,
				}

				if err := taskPub.Created(ctx, task); err != nil {
					t.Fatalf("Failed to publish created event: %v", err)
				}
			},
			verify: func(t *testing.T, routingKey string, body []byte) {
				t.Helper()

				if routingKey != rabbitmqtask.TaskCreatedMessageType {
					t.Fatalf("Expected routing key %s, got %s", rabbitmqtask.TaskCreatedMessageType, routingKey)
				}

				var got internal.Task

				if err := gob.NewDecoder(bytes.NewReader(body)).Decode(&got); err != nil {
					t.Fatalf("Failed to decode body: %v", err)
				}

				expected := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.ValueToPointer(internal.PriorityHigh),
					IsDone:      true,
				}

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Received task is not the same as the created one: %s", diff)
				}
			},
		},
		{
			name: "Updated",
			call: func(t *testing.T, channel *amqp.Channel) {
				t.Helper()

				taskPub := rabbitmqtask.NewTask(channel)

				task := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.ValueToPointer(internal.PriorityHigh),
					IsDone:      true,
				}

				if err := taskPub.Updated(ctx, task); err != nil {
					t.Fatalf("Failed to publish created event: %v", err)
				}
			},
			verify: func(t *testing.T, routingKey string, body []byte) {
				t.Helper()

				if routingKey != rabbitmqtask.TaskUpdatedMessageType {
					t.Fatalf("Expected routing key %s, got %s", rabbitmqtask.TaskUpdatedMessageType, routingKey)
				}

				var got internal.Task

				if err := gob.NewDecoder(bytes.NewReader(body)).Decode(&got); err != nil {
					t.Fatalf("Failed to decode body: %v", err)
				}

				expected := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.ValueToPointer(internal.PriorityHigh),
					IsDone:      true,
				}

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Received task is not the same as the created one: %s", diff)
				}
			},
		},
		{
			name: "Deleted",
			call: func(t *testing.T, channel *amqp.Channel) {
				t.Helper()

				taskPub := rabbitmqtask.NewTask(channel)

				if err := taskPub.Deleted(ctx, "test-123"); err != nil {
					t.Fatalf("Failed to publish created event: %v", err)
				}
			},
			verify: func(t *testing.T, routingKey string, body []byte) {
				t.Helper()

				if routingKey != rabbitmqtask.TaskDeletedMessageType {
					t.Fatalf("Expected routing key %s, got %s", rabbitmqtask.TaskDeletedMessageType, routingKey)
				}

				var got string

				if err := gob.NewDecoder(bytes.NewReader(body)).Decode(&got); err != nil {
					t.Fatalf("Failed to decode body: %v", err)
				}

				expected := "test-123"

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Received task is not the same as the created one: %s", diff)
				}
			},
		},
	}

	for _, tt := range tests { //nolint: paralleltest
		t.Run(tt.name, func(t *testing.T) {
			tt.call(t, client.producerChannel)

			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			t.Cleanup(cancel)

			select {
			case msg := <-client.msgs:
				tt.verify(t, msg.RoutingKey, msg.Body)
			case <-ctx.Done():
				t.Fatal("Did not receive message in time")
			}
		})
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

	//- Publisher

	pubChannel, err := conn.Channel()
	if err != nil {
		res.err = fmt.Errorf("failed to open publisher channel: %w", err)

		return res
	}

	res.producerChannel = pubChannel

	err = res.producerChannel.ExchangeDeclare(
		rabbitmqtask.ExchangeName, // name
		"topic",                   // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		res.err = fmt.Errorf("failed to exchange declare (producer): %w", err)

		return res
	}

	queue, err := res.producerChannel.QueueDeclare(
		rabbitmqtask.ExchangeName, // name
		true,                      // durable
		false,                     // delete when unused
		true,                      // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		res.err = fmt.Errorf("failed to queue declare (producer): %w", err)

		return res
	}

	if err := pubChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		res.err = fmt.Errorf("failed to channel Qos: %w", err)

		return res
	}

	if err := res.producerChannel.QueueBind(
		queue.Name,                // queue name
		"Tasks.*",                 // routing key
		rabbitmqtask.ExchangeName, // exchange
		false,
		nil,
	); err != nil {
		res.err = fmt.Errorf("failed to queue bind (producer): %w", err)

		return res
	}

	//- Consumer

	conChannel, err := conn.Channel()
	if err != nil {
		res.err = fmt.Errorf("failed to open consumer channel: %w", err)

		return res
	}

	res.consumerChannel = conChannel

	err = res.consumerChannel.ExchangeDeclare(
		rabbitmqtask.ExchangeName, // name
		"topic",                   // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		res.err = fmt.Errorf("failed to exchange declare (consumer): %w", err)

		return res
	}

	queue, err = res.consumerChannel.QueueDeclare(
		rabbitmqtask.ExchangeName, // name
		true,                      // durable
		false,                     // delete when unused
		true,                      // exclusive
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		res.err = fmt.Errorf("failed to queue declare (consumer): %w", err)

		return res
	}

	if err := res.consumerChannel.QueueBind(
		queue.Name,                // queue name
		"Task.*",                  // routing key
		rabbitmqtask.ExchangeName, // exchange
		false,
		nil,
	); err != nil {
		res.err = fmt.Errorf("failed to queue bind (consumer): %w", err)

		return res
	}

	msgs, err := res.consumerChannel.Consume(
		queue.Name,           // queue
		rabbitMQConsumerName, // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		res.err = fmt.Errorf("failed to consume msgs: %w", err)

		return res
	}

	res.msgs = msgs

	return res

})

type RabbitMQClient struct {
	container       *rabbitmq.RabbitMQContainer
	producerChannel *amqp.Channel
	consumerChannel *amqp.Channel
	connection      *amqp.Connection
	msgs            <-chan amqp.Delivery
	err             error
}

func (r *RabbitMQClient) Teardown() error {
	var err error

	if r.consumerChannel != nil {
		if err1 := r.consumerChannel.Close(); err1 != nil {
			err = fmt.Errorf("failed to close consumer channel: %w", err1)
		}
	}

	if r.producerChannel != nil {
		if err1 := r.producerChannel.Close(); err1 != nil {
			err = errors.Join(err, fmt.Errorf("failed to close producer channel: %w", err1))
		}
	}

	if r.connection != nil {
		if err1 := r.connection.Close(); err1 != nil {
			err = errors.Join(err, fmt.Errorf("failed to close connection: %w", err1))
		}
	}

	if r.container != nil {
		if err1 := testcontainers.TerminateContainer(r.container); err1 != nil {
			err = errors.Join(err, fmt.Errorf("failed to terminate container: %w", err1))
		}
	}

	return err
}
