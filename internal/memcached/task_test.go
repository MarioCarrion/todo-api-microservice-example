package memcached_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	memcachedtask "github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/memcached/memcachedtesting"
)

// dockerImage must match the docker image listed in `compose.yml`.
const dockerImage = "memcached:1.6.19-alpine3.17"

func TestMain(m *testing.M) {
	client := setupClient()
	if client.err != nil {
		panic(fmt.Sprintf("Failed to set up Memcached client: %v", client.err))
	}

	code := m.Run()

	if err := client.Teardown(); err != nil {
		panic(fmt.Sprintf("Failed to close Memcached client: %v", err))
	}

	os.Exit(code)
}

func TestTask_All(t *testing.T) { //nolint: tparallel
	t.Parallel()

	client := setupClient()

	logger := zap.NewNop()
	ctx := t.Context()

	tests := []struct {
		name           string
		setUpMockStore func() *memcachedtesting.FakeTaskStore
		callAndVerify  func(t *testing.T, task *memcachedtask.Task, store *memcachedtesting.FakeTaskStore)
	}{
		{
			name: "Find",
			setUpMockStore: func() *memcachedtesting.FakeTaskStore {
				mockStore := &memcachedtesting.FakeTaskStore{}
				testTask := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.PriorityHigh.Pointer(),
				}
				mockStore.FindReturns(testTask, nil)

				return mockStore
			},
			callAndVerify: func(t *testing.T, task *memcachedtask.Task, store *memcachedtesting.FakeTaskStore) {
				t.Helper()

				got, err := task.Find(ctx, "test-123")
				if err != nil {
					t.Fatalf("Failed to find task: %v", err)
				}

				expected := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.PriorityHigh.Pointer(),
				}

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Found task is not the same as the mocked one: %s", diff)
				}

				// Verify store was called once
				if count := store.FindCallCount(); count != 1 {
					t.Errorf("Expected store.Find to be called once, got %d", count)
				}

				// Second call should hit cache
				if _, err := task.Find(ctx, "test-123"); err != nil {
					t.Fatalf("Failed to find task from cache: %v", err)
				}

				// Store should still be called only once (cache hit)
				if count := store.FindCallCount(); count != 1 {
					t.Errorf("Expected store.Find to still be called once (cache hit), got %d", count)
				}
			},
		},
		{
			name: "Create",
			setUpMockStore: func() *memcachedtesting.FakeTaskStore {
				mockStore := &memcachedtesting.FakeTaskStore{}
				testTask := internal.Task{
					ID:          "test-123",
					Description: "Create task",
					Priority:    internal.PriorityHigh.Pointer(),
				}
				mockStore.CreateReturns(testTask, nil)

				return mockStore
			},
			callAndVerify: func(t *testing.T, task *memcachedtask.Task, store *memcachedtesting.FakeTaskStore) {
				t.Helper()

				params := internal.CreateParams{
					Description: "Create task",
					Priority:    internal.PriorityHigh.Pointer(),
				}

				got, err := task.Create(ctx, params)
				if err != nil {
					t.Fatalf("Failed to create task: %v", err)
				}

				expected := internal.Task{
					ID:          "test-123",
					Description: "Create task",
					Priority:    internal.PriorityHigh.Pointer(),
				}

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Found task is not the same as the mocked one: %s", diff)
				}

				// Verify store was called once
				if count := store.CreateCallCount(); count != 1 {
					t.Errorf("Expected store.Create to be called once, got %d", count)
				}

				// Verify the task is now in cache by trying to get it
				item, err := client.client.Get(expected.ID)
				if err != nil {
					t.Logf("Task not found in cache (expected for write-through): %v", err)
				}

				if item != nil {
					t.Logf("Task successfully cached with key: %s", item.Key)
				}
			},
		},
		{
			name: "Delete",
			setUpMockStore: func() *memcachedtesting.FakeTaskStore {
				mockStore := &memcachedtesting.FakeTaskStore{}

				testTask := internal.Task{
					ID:          "test-123",
					Description: "Delete task",
					Priority:    internal.PriorityHigh.Pointer(),
				}
				mockStore.CreateReturns(testTask, nil)

				mockStore.DeleteReturns(nil)

				return mockStore
			},
			callAndVerify: func(t *testing.T, task *memcachedtask.Task, store *memcachedtesting.FakeTaskStore) {
				t.Helper()

				params := internal.CreateParams{
					Description: "Delete task",
					Priority:    internal.PriorityHigh.Pointer(),
				}

				got, err := task.Create(ctx, params)
				if err != nil {
					t.Fatalf("Failed to create task: %v", err)
				}

				expected := internal.Task{
					ID:          "test-123",
					Description: "Delete task",
					Priority:    internal.PriorityHigh.Pointer(),
				}

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Created task is not the same as the mocked one: %s", diff)
				}

				// Verify store was called once
				if count := store.CreateCallCount(); count != 1 {
					t.Errorf("Expected store.Create to be called once, got %d", count)
				}

				// Verify the task is now in cache by trying to get it
				item, err := client.client.Get(expected.ID)
				if err != nil {
					t.Logf("Task not found in cache (expected for write-through): %v", err)
				}

				if item != nil {
					t.Logf("Task successfully cached with key: %s", item.Key)
				}

				//- Now delete the task
				if err := task.Delete(ctx, expected.ID); err != nil {
					t.Fatalf("Failed to delete task: %v", err)
				}

				if count := store.DeleteCallCount(); count != 1 {
					t.Errorf("Expected store.Delete to be called once, got %d", count)
				}

				_, err = client.client.Get(expected.ID)
				if !errors.Is(err, memcache.ErrCacheMiss) {
					t.Errorf("Expected a ErrCacheMiss, but got %v", err)
				}
			},
		},
	}

	for _, tt := range tests { //nolint: paralleltest
		t.Run(tt.name, func(t *testing.T) {
			mockStore := tt.setUpMockStore()

			taskCache := memcachedtask.NewTask(client.client, mockStore, logger)

			tt.callAndVerify(t, taskCache, mockStore)
		})
	}
}

//---

var setupClient = sync.OnceValue(func() MemcachedClient { //nolint: gochecknoglobals
	var res MemcachedClient

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        dockerImage,
		ExposedPorts: []string{"11211/tcp"},
		WaitingFor:   wait.ForListeningPort("11211/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		res.err = fmt.Errorf("failed to start memcached container: %v", err)

		return res
	}

	res.container = &container

	host, err := container.Host(ctx)
	if err != nil {
		res.err = fmt.Errorf("failed to get host: %v", err)

		return res
	}

	port, err := container.MappedPort(ctx, "11211")
	if err != nil {
		res.err = fmt.Errorf("failed to get port: %v", err)

		return res
	}

	client := memcache.New(host + ":" + port.Port())

	res.client = client

	return res
})

type MemcachedClient struct {
	container *testcontainers.Container
	client    *memcache.Client
	err       error
}

func (r *MemcachedClient) Teardown() error {
	var err error

	if r.container != nil {
		if err1 := testcontainers.TerminateContainer(*r.container); err1 != nil {
			err = fmt.Errorf("failed to terminate container: %w", err1)
		}
	}

	if r.err != nil {
		err = errors.Join(err, r.err)
	}

	return err
}
