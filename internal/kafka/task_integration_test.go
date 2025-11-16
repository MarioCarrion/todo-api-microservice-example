package kafka_test

import (
	"context"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/testcontainers/testcontainers-go"
	kafkamodule "github.com/testcontainers/testcontainers-go/modules/kafka"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	kafkaTask "github.com/MarioCarrion/todo-api-microservice-example/internal/kafka"
)

// setupKafkaProducer starts a Kafka container and returns a configured producer.
func setupKafkaProducer(ctx context.Context, t *testing.T) (*kafka.Producer, func()) {
	t.Helper()

	kafkaContainer, err := kafkamodule.Run(ctx, "confluentinc/confluent-local:7.5.0")
	if err != nil {
		t.Fatalf("failed to start kafka container: %v", err)
	}

	cleanup := func() {
		if err := testcontainers.TerminateContainer(kafkaContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		cleanup()
		t.Fatalf("failed to get brokers: %v", err)
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
	})
	if err != nil {
		cleanup()
		t.Fatalf("failed to create producer: %v", err)
	}

	cleanupAll := func() {
		producer.Close()
		cleanup()
	}

	return producer, cleanupAll
}

func TestTask_Created_Integration(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := t.Context()
	producer, cleanup := setupKafkaProducer(ctx, t)

	defer cleanup()

	// Create task publisher
	taskPub := kafkaTask.NewTask(producer, "test-tasks")

	// Test Created method
	task := internal.Task{
		ID:          "test-123",
		Description: "Test task",
		Priority:    internal.PriorityHigh.Pointer(),
		IsDone:      false,
	}

	err := taskPub.Created(ctx, task)
	if err != nil {
		t.Fatalf("Failed to publish created event: %v", err)
	}

	// Flush to ensure message is sent
	producer.Flush(5000)
}

func TestTask_Updated_Integration(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := t.Context()
	producer, cleanup := setupKafkaProducer(ctx, t)

	defer cleanup()

	taskPub := kafkaTask.NewTask(producer, "test-tasks")

	task := internal.Task{
		ID:          "test-456",
		Description: "Updated task",
		IsDone:      true,
	}

	err := taskPub.Updated(ctx, task)
	if err != nil {
		t.Fatalf("Failed to publish updated event: %v", err)
	}

	producer.Flush(5000)
}

func TestTask_Deleted_Integration(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := t.Context()
	producer, cleanup := setupKafkaProducer(ctx, t)

	defer cleanup()

	taskPub := kafkaTask.NewTask(producer, "test-tasks")

	err := taskPub.Deleted(ctx, "test-789")
	if err != nil {
		t.Fatalf("Failed to publish deleted event: %v", err)
	}

	producer.Flush(5000)
}
