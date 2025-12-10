package kafka_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	kafkatest "github.com/testcontainers/testcontainers-go/modules/kafka"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	kafkatask "github.com/MarioCarrion/todo-api-microservice-example/internal/kafka"
)

const (
	// dockerImage version must match the docker image version listed in `compose.kafka.yml`.
	dockerImage = "confluentinc/confluent-local:7.6.2"

	// topicName is the Kafka topic used in tests.
	topicName = "test-tasks-topic"

	// consumerReadTimeout is the timeout to read messages from Kafka consumer.
	consumerReadTimeout = 5 * time.Second
)

func TestMain(m *testing.M) {
	client := setupClient()
	if client.err != nil {
		panic(fmt.Sprintf("Failed to set up Kafka client: %v", client.err))
	}

	code := m.Run()

	if err := client.Teardown(); err != nil {
		panic(fmt.Sprintf("Failed to close Kafka client: %v", err))
	}

	os.Exit(code)
}

func TestTask_All(t *testing.T) {
	t.Parallel()

	client := setupClient()
	if client.err != nil {
		t.Fatalf("Failed to setupClient: %v", client.err)
	}

	type Event struct {
		Type  string
		Value internal.Task
	}

	taskPub := kafkatask.NewTask(client.producer, topicName)

	tests := []struct {
		name   string
		call   func(t *testing.T, client *kafkatask.Task)
		verify func(t *testing.T, evnt Event)
	}{
		{
			name: "Created",
			call: func(t *testing.T, client *kafkatask.Task) {
				task := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.PriorityHigh.Pointer(),
					IsDone:      true,
				}

				if err := taskPub.Created(t.Context(), task); err != nil {
					t.Fatalf("Failed to publish created event: %v", err)
				}
			},
			verify: func(t *testing.T, evnt Event) {
				if evnt.Type != kafkatask.TaskCreatedMessageType {
					t.Fatalf("Expected created event type, got %s", evnt.Type)
				}

				expected := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.PriorityHigh.Pointer(),
					IsDone:      true,
				}

				if diff := cmp.Diff(evnt.Value, expected); diff != "" {
					t.Fatalf("Received created task event is not the same as the published one: %s", diff)
				}
			},
		},
		{
			name: "Deleted",
			call: func(t *testing.T, client *kafkatask.Task) {
				if err := taskPub.Deleted(t.Context(), "task-id-123-456"); err != nil {
					t.Fatalf("Failed to publish deleted event: %v", err)
				}
			},
			verify: func(t *testing.T, evnt Event) {
				if evnt.Type != kafkatask.TaskDeletedMessageType {
					t.Fatalf("Expected deleted event type, got %s", evnt.Type)
				}

				expected := internal.Task{
					ID: "task-id-123-456",
				}

				if diff := cmp.Diff(evnt.Value, expected); diff != "" {
					t.Fatalf("Received deleted task event is not the same as the published one: %s", diff)
				}
			},
		},
		{
			name: "Updated",
			call: func(t *testing.T, client *kafkatask.Task) {
				task := internal.Task{
					ID:          "test-123-updated",
					Description: "Test task updated",
					Priority:    internal.PriorityHigh.Pointer(),
					IsDone:      true,
				}

				if err := taskPub.Created(t.Context(), task); err != nil {
					t.Fatalf("Failed to publish updated event: %v", err)
				}
			},
			verify: func(t *testing.T, evnt Event) {
				if evnt.Type != kafkatask.TaskCreatedMessageType {
					t.Fatalf("Expected updated event type, got %s", evnt.Type)
				}

				expected := internal.Task{
					ID:          "test-123-updated",
					Description: "Test task updated",
					Priority:    internal.PriorityHigh.Pointer(),
					IsDone:      true,
				}

				if diff := cmp.Diff(evnt.Value, expected); diff != "" {
					t.Fatalf("Received updated task event is not the same as the published one: %s", diff)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.call(t, taskPub)

			producedEvent := <-client.producer.Events()
			msg, ok := producedEvent.(*kafka.Message)
			if !ok {
				t.Fatalf("Failed to receive produced event: %T - %v", msg, ok)
			}

			msg, err := client.consumer.ReadMessage(consumerReadTimeout)
			if err != nil {
				t.Fatalf("Failed to read message: %v", err)
			}

			var evt Event

			if err := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&evt); err != nil {
				t.Errorf("Failed to decode message: %v", err)
				if _, err := client.consumer.CommitMessage(msg); err != nil {
					t.Fatalf("Failed to commit message: %v", err)
				}
			}

			tt.verify(t, evt)
		})
	}
}

var setupClient = sync.OnceValue(func() KafkaClient {
	var res KafkaClient

	ctx := context.Background()

	container, err := kafkatest.Run(ctx, dockerImage)
	if err != nil {
		res.err = fmt.Errorf("failed to start kafka container: %w", err)

		return res
	}

	res.container = container

	brokers, err := container.Brokers(ctx)
	if err != nil {
		res.err = fmt.Errorf("failed to get kafka brokers: %w", err)

		return res
	}

	//- Producer

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
	})
	if err != nil {
		res.err = fmt.Errorf("failed to instantiate kafka producer: %w", err)

		return res
	}

	res.producer = producer

	//- Consumer

	config := kafka.ConfigMap{
		"bootstrap.servers":  brokers[0],
		"group.id":           "tests-group-id",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		res.err = fmt.Errorf("failed to instantiate kafka consumer: %w", err)

		return res
	}

	res.consumer = consumer

	if err := consumer.Subscribe(topicName, nil); err != nil {
		res.err = fmt.Errorf("failed to subscribe kafka consumer: %w", err)

		return res
	}

	return res
})

type KafkaClient struct {
	producer  *kafka.Producer
	consumer  *kafka.Consumer
	container *kafkatest.KafkaContainer
	err       error
}

func (r *KafkaClient) Teardown() error {
	var err error

	if r.consumer != nil {
		if err1 := r.consumer.Close(); err1 != nil {
			err = err1
		}
	}

	if r.producer != nil {
		r.producer.Close()
	}

	if r.container != nil {
		if err1 := testcontainers.TerminateContainer(r.container); err1 != nil {
			err = errors.Join(err, fmt.Errorf("failed to terminate container: %w", err1))
		}
	}

	if r.err != nil {
		err = errors.Join(err, r.err)
	}

	return err
}
