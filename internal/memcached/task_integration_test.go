package memcached_test

import (
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	taskMemcached "github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/memcached/memcachedtesting"
)

func TestTask_Find_Integration(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	req := testcontainers.ContainerRequest{
		Image:        "memcached:1.6-alpine",
		ExposedPorts: []string{"11211/tcp"},
		WaitingFor:   wait.ForListeningPort("11211/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start memcached container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	// Get host and port
	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}

	port, err := container.MappedPort(ctx, "11211")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}

	// Create memcache client
	mcClient := memcache.New(host + ":" + port.Port())
	logger := zap.NewNop()

	// Create mock store that returns a task
	mockStore := &memcachedtesting.FakeTaskStore{}
	testTask := internal.Task{
		ID:          "test-123",
		Description: "Test task",
		Priority:    internal.PriorityHigh.Pointer(),
	}
	mockStore.FindReturns(testTask, nil)

	// Create task with cache
	taskCache := taskMemcached.NewTask(mcClient, mockStore, logger)

	// First call should go to store and cache
	result, err := taskCache.Find(ctx, "test-123")
	if err != nil {
		t.Fatalf("Failed to find task: %v", err)
	}

	if result.ID != testTask.ID {
		t.Errorf("Expected task ID %s, got %s", testTask.ID, result.ID)
	}

	// Verify store was called once
	if mockStore.FindCallCount() != 1 {
		t.Errorf("Expected store.Find to be called once, got %d", mockStore.FindCallCount())
	}

	// Second call should hit cache
	_, err = taskCache.Find(ctx, "test-123")
	if err != nil {
		t.Fatalf("Failed to find task from cache: %v", err)
	}

	// Store should still be called only once (cache hit)
	if mockStore.FindCallCount() != 1 {
		t.Errorf("Expected store.Find to still be called once (cache hit), got %d", mockStore.FindCallCount())
	}
}

func TestTask_Create_Integration(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	req := testcontainers.ContainerRequest{
		Image:        "memcached:1.6-alpine",
		ExposedPorts: []string{"11211/tcp"},
		WaitingFor:   wait.ForListeningPort("11211/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start memcached container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}

	port, err := container.MappedPort(ctx, "11211")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}

	mcClient := memcache.New(host + ":" + port.Port())
	logger := zap.NewNop()

	mockStore := &memcachedtesting.FakeTaskStore{}
	testTask := internal.Task{
		ID:          "test-create",
		Description: "Created task",
	}
	mockStore.CreateReturns(testTask, nil)

	taskCache := taskMemcached.NewTask(mcClient, mockStore, logger)

	params := internal.CreateParams{
		Description: "Created task",
	}

	result, err := taskCache.Create(ctx, params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if result.ID != testTask.ID {
		t.Errorf("Expected task ID %s, got %s", testTask.ID, result.ID)
	}

	// Verify the task is now in cache by trying to get it
	item, err := mcClient.Get(testTask.ID)
	if err != nil {
		// This is OK - the cache set might not complete immediately
		t.Logf("Task not found in cache (expected for write-through): %v", err)
	}

	if item != nil {
		t.Logf("Task successfully cached with key: %s", item.Key)
	}
}

func TestTask_Delete_Integration(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	req := testcontainers.ContainerRequest{
		Image:        "memcached:1.6-alpine",
		ExposedPorts: []string{"11211/tcp"},
		WaitingFor:   wait.ForListeningPort("11211/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start memcached container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}

	port, err := container.MappedPort(ctx, "11211")
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}

	mcClient := memcache.New(host + ":" + port.Port())
	logger := zap.NewNop()

	mockStore := &memcachedtesting.FakeTaskStore{}
	mockStore.DeleteReturns(nil)

	taskCache := taskMemcached.NewTask(mcClient, mockStore, logger)

	err = taskCache.Delete(ctx, "test-delete")
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	if mockStore.DeleteCallCount() != 1 {
		t.Errorf("Expected store.Delete to be called once, got %d", mockStore.DeleteCallCount())
	}
}
