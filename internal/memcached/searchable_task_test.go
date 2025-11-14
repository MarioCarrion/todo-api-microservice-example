package memcached_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	taskMemcached "github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
)

// mockSearchableTaskStore is a mock implementation of SearchableTaskStore for testing
type mockSearchableTaskStore struct {
	deleteFn func(ctx context.Context, id string) error
	indexFn  func(ctx context.Context, task internal.Task) error
	searchFn func(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error)
}

func (m *mockSearchableTaskStore) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockSearchableTaskStore) Index(ctx context.Context, task internal.Task) error {
	if m.indexFn != nil {
		return m.indexFn(ctx, task)
	}
	return nil
}

func (m *mockSearchableTaskStore) Search(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, args)
	}
	return internal.SearchResults{}, nil
}

func TestNewSearchableTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "creates new searchable task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Use a dummy memcache client address - we won't actually connect
			client := memcache.New("localhost:11211")
			store := &mockSearchableTaskStore{}

			result := taskMemcached.NewSearchableTask(client, store)

			if result == nil {
				t.Fatal("expected non-nil SearchableTask")
			}
		})
	}
}

func TestSearchableTask_Index(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		task           internal.Task
		mockStore      *mockSearchableTaskStore
		expectedErrMsg string
	}{
		{
			name: "successful index",
			task: internal.Task{
				ID:          "123",
				Description: "test task",
			},
			mockStore: &mockSearchableTaskStore{
				indexFn: func(ctx context.Context, task internal.Task) error {
					return nil
				},
			},
			expectedErrMsg: "",
		},
		{
			name: "store error",
			task: internal.Task{
				ID:          "123",
				Description: "test task",
			},
			mockStore: &mockSearchableTaskStore{
				indexFn: func(ctx context.Context, task internal.Task) error {
					return errors.New("index error")
				},
			},
			expectedErrMsg: "orig.Index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			searchable := taskMemcached.NewSearchableTask(client, tt.mockStore)

			err := searchable.Index(context.Background(), tt.task)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedErrMsg)
				}
				if !containsString(err.Error(), tt.expectedErrMsg) {
					t.Errorf("expected error containing %q, got %q", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSearchableTask_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		id             string
		mockStore      *mockSearchableTaskStore
		expectedErrMsg string
	}{
		{
			name: "successful delete",
			id:   "123",
			mockStore: &mockSearchableTaskStore{
				deleteFn: func(ctx context.Context, id string) error {
					return nil
				},
			},
			expectedErrMsg: "",
		},
		{
			name: "store error",
			id:   "123",
			mockStore: &mockSearchableTaskStore{
				deleteFn: func(ctx context.Context, id string) error {
					return errors.New("delete error")
				},
			},
			expectedErrMsg: "orig.Delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			searchable := taskMemcached.NewSearchableTask(client, tt.mockStore)

			err := searchable.Delete(context.Background(), tt.id)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedErrMsg)
				}
				if !containsString(err.Error(), tt.expectedErrMsg) {
					t.Errorf("expected error containing %q, got %q", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestSearchableTask_Search is skipped because it requires a running memcached instance.
// The Search method attempts to connect to memcache first (getTask), and without a running
// memcached server, it returns a connection error instead of ErrCacheMiss, making it
// difficult to unit test without infrastructure. Integration tests with testcontainers
// would be more appropriate for testing the caching behavior.

// containsString checks if a string contains a substring (simple helper for test)
func containsString(s, substr string) bool {
	// Simple contains check
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
