package memcached_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	memcachedtask "github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/memcached/memcachedtesting"
)

func TestSearchableTask_All(t *testing.T) {
	t.Parallel()

	client := setupClient()

	tests := []struct {
		name           string
		setUpMockStore func() *memcachedtesting.FakeSearchableTaskStore
		callAndVerify  func(t *testing.T, task *memcachedtask.SearchableTask, store *memcachedtesting.FakeSearchableTaskStore)
	}{
		{
			name: "Index",
			setUpMockStore: func() *memcachedtesting.FakeSearchableTaskStore {
				return &memcachedtesting.FakeSearchableTaskStore{}
			},
			callAndVerify: func(t *testing.T, task *memcachedtask.SearchableTask, store *memcachedtesting.FakeSearchableTaskStore) {
				t.Helper()

				expected := internal.Task{
					ID:          "test-123",
					Description: "Test task",
					Priority:    internal.ValueToPointer(internal.PriorityHigh),
				}

				if err := task.Index(t.Context(), expected); err != nil {
					t.Fatalf("Failed to index task: %v", err)
				}

				_, got := store.IndexArgsForCall(0)

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Received indexed task is not the same as the mocked one: %s", diff)
				}

				// Verify store was called once
				if count := store.IndexCallCount(); count != 1 {
					t.Errorf("Expected store.Index to be called once, got %d", count)
				}
			},
		},
		{
			name: "Delete",
			setUpMockStore: func() *memcachedtesting.FakeSearchableTaskStore {
				return &memcachedtesting.FakeSearchableTaskStore{}
			},
			callAndVerify: func(t *testing.T, task *memcachedtask.SearchableTask, store *memcachedtesting.FakeSearchableTaskStore) {
				t.Helper()

				expected := "test-123"

				if err := task.Delete(t.Context(), "test-123"); err != nil {
					t.Fatalf("Failed to index task: %v", err)
				}

				_, got := store.DeleteArgsForCall(0)

				if got != expected {
					t.Fatalf("Received %s task id not the same as expected one: %s", got, expected)
				}

				// Verify store was called once
				if count := store.DeleteCallCount(); count != 1 {
					t.Errorf("Expected store.Delete to be called once, got %d", count)
				}
			},
		},
		{
			name: "Search",
			setUpMockStore: func() *memcachedtesting.FakeSearchableTaskStore {
				mock := &memcachedtesting.FakeSearchableTaskStore{}

				mock.SearchReturns(internal.SearchResults{
					Tasks: []internal.Task{
						{
							ID:          "test-123",
							Description: "Test task",
							Priority:    internal.ValueToPointer(internal.PriorityHigh),
						},
					},
					Total: 1,
				}, nil)

				return mock
			},
			callAndVerify: func(t *testing.T, task *memcachedtask.SearchableTask, store *memcachedtesting.FakeSearchableTaskStore) {
				t.Helper()

				got, err := task.Search(t.Context(), internal.SearchParams{
					Description: internal.ValueToPointer("description"),
					Priority:    internal.ValueToPointer(internal.PriorityHigh),
					IsDone:      internal.ValueToPointer(true),
					From:        0,
					Size:        2,
				})
				if err != nil {
					t.Fatalf("Failed to search tasks: %v", err)
				}
				// return fmt.Sprintf("%s_%d_%t_%d_%d", description, priority, isDone, args.From, args.Size)

				if count := store.SearchCallCount(); count != 1 {
					t.Errorf("Expected store.Search to be called once, got %d", count)
				}

				expected := internal.SearchResults{
					Tasks: []internal.Task{
						{
							ID:          "test-123",
							Description: "Test task",
							Priority:    internal.ValueToPointer(internal.PriorityHigh),
						},
					},
					Total: 1,
				}

				if diff := cmp.Diff(got, expected); diff != "" {
					t.Fatalf("Found tasks are not the same as the mocked one: %s", diff)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStore := tt.setUpMockStore()
			task := memcachedtask.NewSearchableTask(client.client, mockStore)

			tt.callAndVerify(t, task, mockStore)
		})
	}
}
