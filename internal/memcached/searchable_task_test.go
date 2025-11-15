package memcached_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	taskMemcached "github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/memcached/memcachedtesting"
)

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
			store := &memcachedtesting.FakeSearchableTaskStore{}

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
		setupMock      func(*memcachedtesting.FakeSearchableTaskStore)
		expectedErrMsg string
	}{
		{
			name: "successful index",
			task: internal.Task{
				ID:          "123",
				Description: "test task",
			},
			setupMock: func(m *memcachedtesting.FakeSearchableTaskStore) {
				m.IndexReturns(nil)
			},
			expectedErrMsg: "",
		},
		{
			name: "store error",
			task: internal.Task{
				ID:          "123",
				Description: "test task",
			},
			setupMock: func(m *memcachedtesting.FakeSearchableTaskStore) {
				m.IndexReturns(errors.New("index error"))
			},
			expectedErrMsg: "orig.Index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			mockStore := &memcachedtesting.FakeSearchableTaskStore{}
			tt.setupMock(mockStore)
			searchable := taskMemcached.NewSearchableTask(client, mockStore)

			err := searchable.Index(t.Context(), tt.task)

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
		setupMock      func(*memcachedtesting.FakeSearchableTaskStore)
		expectedErrMsg string
	}{
		{
			name: "successful delete",
			id:   "123",
			setupMock: func(m *memcachedtesting.FakeSearchableTaskStore) {
				m.DeleteReturns(nil)
			},
			expectedErrMsg: "",
		},
		{
			name: "store error",
			id:   "123",
			setupMock: func(m *memcachedtesting.FakeSearchableTaskStore) {
				m.DeleteReturns(errors.New("delete error"))
			},
			expectedErrMsg: "orig.Delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			mockStore := &memcachedtesting.FakeSearchableTaskStore{}
			tt.setupMock(mockStore)
			searchable := taskMemcached.NewSearchableTask(client, mockStore)

			err := searchable.Delete(t.Context(), tt.id)

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
	return strings.Contains(s, substr)
}
