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
		name      string
		task      internal.Task
		setupMock func(*memcachedtesting.FakeSearchableTaskStore)
		verify    func(*testing.T, error)
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
			verify: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
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
			verify: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", "orig.Index")
				}
				if !strings.Contains(err.Error(), "orig.Index") {
					t.Errorf("expected error containing %q, got %q", "orig.Index", err.Error())
				}
			},
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
			tt.verify(t, err)
		})
	}
}

func TestSearchableTask_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		setupMock func(*memcachedtesting.FakeSearchableTaskStore)
		verify    func(*testing.T, error)
	}{
		{
			name: "successful delete",
			id:   "123",
			setupMock: func(m *memcachedtesting.FakeSearchableTaskStore) {
				m.DeleteReturns(nil)
			},
			verify: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "store error",
			id:   "123",
			setupMock: func(m *memcachedtesting.FakeSearchableTaskStore) {
				m.DeleteReturns(errors.New("delete error"))
			},
			verify: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", "orig.Delete")
				}
				if !strings.Contains(err.Error(), "orig.Delete") {
					t.Errorf("expected error containing %q, got %q", "orig.Delete", err.Error())
				}
			},
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
			tt.verify(t, err)
		})
	}
}

// TestSearchableTask_Search is skipped because it requires a running memcached instance.
// The Search method attempts to connect to memcache first (getTask), and without a running
// memcached server, it returns a connection error instead of ErrCacheMiss, making it
// difficult to unit test without infrastructure. Integration tests with testcontainers
// would be more appropriate for testing the caching behavior.
