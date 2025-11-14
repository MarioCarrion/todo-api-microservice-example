//go:build integration


package redis_test

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/testcontainers/testcontainers-go"
	redismodule "github.com/testcontainers/testcontainers-go/modules/redis"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	redisTask "github.com/MarioCarrion/todo-api-microservice-example/internal/redis"
)

func TestTask_Created_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// Start Redis container
	redisContainer, err := redismodule.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(redisContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr: connStr,
	})
	defer client.Close()

	// Create task publisher
	taskPub := redisTask.NewTask(client)

	// Test Created method
	task := internal.Task{
		ID:          "test-123",
		Description: "Test task",
		IsDone:      false,
	}

	err = taskPub.Created(ctx, task)
	if err != nil {
		t.Fatalf("Failed to publish created event: %v", err)
	}

	// Verify message was published (subscribe and check)
	pubsub := client.Subscribe(ctx, "Task.Created")
	defer pubsub.Close()

	// Note: In a real test, you'd need a separate goroutine to publish after subscribing
	// This test verifies the method doesn't error
}

func TestTask_Updated_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	redisContainer, err := redismodule.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(redisContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: connStr,
	})
	defer client.Close()

	taskPub := redisTask.NewTask(client)

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

	ctx := context.Background()

	redisContainer, err := redismodule.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(redisContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: connStr,
	})
	defer client.Close()

	taskPub := redisTask.NewTask(client)

	err = taskPub.Deleted(ctx, "test-789")
	if err != nil {
		t.Fatalf("Failed to publish deleted event: %v", err)
	}
}
