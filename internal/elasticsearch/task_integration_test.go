//go:build integration


package elasticsearch_test

import (
	"context"
	"testing"
	"time"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/elasticsearch"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	esTask "github.com/MarioCarrion/todo-api-microservice-example/internal/elasticsearch"
)

func TestTask_Index_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// Start Elasticsearch container
	esContainer, err := elasticsearch.Run(ctx, "docker.elastic.co/elasticsearch/elasticsearch:7.17.10")
	if err != nil {
		t.Fatalf("failed to start elasticsearch container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(esContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := esContainer.Address(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Create Elasticsearch client
	cfg := esv7.Config{
		Addresses: []string{connStr},
	}
	client, err := esv7.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create elasticsearch client: %v", err)
	}

	// Create task repository
	taskRepo := esTask.NewTask(client)

	// Test Index method
	now := time.Now()
	task := internal.Task{
		ID:          "test-123",
		Description: "Test task for elasticsearch",
		Priority:    internal.PriorityHigh.Pointer(),
		IsDone:      false,
		Dates: &internal.Dates{
			Start: &now,
			Due:   &now,
		},
	}

	err = taskRepo.Index(ctx, task)
	if err != nil {
		t.Fatalf("Failed to index task: %v", err)
	}

	// Give ES a moment to index
	time.Sleep(500 * time.Millisecond)

	// Test Search to verify it was indexed
	results, err := taskRepo.Search(ctx, internal.SearchParams{
		Description: internal.ValueToPointer("Test"),
		From:        0,
		Size:        10,
	})
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results.Tasks) == 0 {
		t.Error("Expected to find indexed task")
	}
}

func TestTask_Delete_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	esContainer, err := elasticsearch.Run(ctx, "docker.elastic.co/elasticsearch/elasticsearch:7.17.10")
	if err != nil {
		t.Fatalf("failed to start elasticsearch container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(esContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := esContainer.Address(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	cfg := esv7.Config{
		Addresses: []string{connStr},
	}
	client, err := esv7.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create elasticsearch client: %v", err)
	}

	taskRepo := esTask.NewTask(client)

	// First index a task
	task := internal.Task{
		ID:          "test-delete",
		Description: "Task to delete",
		IsDone:      false,
	}

	err = taskRepo.Index(ctx, task)
	if err != nil {
		t.Fatalf("Failed to index task: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Now delete it
	err = taskRepo.Delete(ctx, task.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
}

func TestTask_Search_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	esContainer, err := elasticsearch.Run(ctx, "docker.elastic.co/elasticsearch/elasticsearch:7.17.10")
	if err != nil {
		t.Fatalf("failed to start elasticsearch container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(esContainer); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := esContainer.Address(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	cfg := esv7.Config{
		Addresses: []string{connStr},
	}
	client, err := esv7.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create elasticsearch client: %v", err)
	}

	taskRepo := esTask.NewTask(client)

	// Index multiple tasks
	tasks := []internal.Task{
		{ID: "1", Description: "High priority task", Priority: internal.PriorityHigh.Pointer(), IsDone: false},
		{ID: "2", Description: "Low priority task", Priority: internal.PriorityLow.Pointer(), IsDone: false},
		{ID: "3", Description: "Completed task", Priority: internal.PriorityMedium.Pointer(), IsDone: true},
	}

	for _, task := range tasks {
		if err := taskRepo.Index(ctx, task); err != nil {
			t.Fatalf("Failed to index task: %v", err)
		}
	}

	time.Sleep(1 * time.Second)

	// Search for high priority tasks
	results, err := taskRepo.Search(ctx, internal.SearchParams{
		Priority: internal.PriorityHigh.Pointer(),
		From:     0,
		Size:     10,
	})
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results.Tasks) == 0 {
		t.Error("Expected to find high priority tasks")
	}

	// Search for completed tasks
	results, err = taskRepo.Search(ctx, internal.SearchParams{
		IsDone: internal.ValueToPointer(true),
		From:   0,
		Size:   10,
	})
	if err != nil {
		t.Fatalf("Failed to search: %v", err)
	}

	if len(results.Tasks) == 0 {
		t.Error("Expected to find completed tasks")
	}
}
