package memcached_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	taskMemcached "github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/memcached/memcachedtesting"
)

func TestNewTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "creates new task with cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			store := &memcachedtesting.FakeTaskStore{}
			logger := zap.NewNop()

			task := taskMemcached.NewTask(client, store, logger)

			if task == nil {
				t.Fatal("expected non-nil Task")
			}
		})
	}
}

func TestTask_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    internal.CreateParams
		setupMock func(*memcachedtesting.FakeTaskStore)
		verify    func(*testing.T, internal.Task, error)
	}{
		{
			name: "successful create",
			params: internal.CreateParams{
				Description: "test task",
				Priority:    internal.PriorityHigh.Pointer(),
			},
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.CreateReturns(internal.Task{
					ID:          "123",
					Description: "test task",
					Priority:    internal.PriorityHigh.Pointer(),
				}, nil)
			},
			verify: func(t *testing.T, result internal.Task, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				expectedTask := internal.Task{
					ID:          "123",
					Description: "test task",
					Priority:    internal.PriorityHigh.Pointer(),
				}
				if diff := cmp.Diff(expectedTask, result); diff != "" {
					t.Errorf("task mismatch (-want +got):\n%s", diff)
				}
			},
		},
		{
			name: "store error",
			params: internal.CreateParams{
				Description: "test task",
			},
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.CreateReturns(internal.Task{}, errors.New("create error"))
			},
			verify: func(t *testing.T, _ internal.Task, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), "orig.Create") {
					t.Errorf("expected error to contain %q, got %q", "orig.Create", err.Error())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			logger := zap.NewNop()
			mockStore := &memcachedtesting.FakeTaskStore{}
			tt.setupMock(mockStore)
			task := taskMemcached.NewTask(client, mockStore, logger)

			result, err := task.Create(t.Context(), tt.params)
			tt.verify(t, result, err)
		})
	}
}

func TestTask_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		setupMock func(*memcachedtesting.FakeTaskStore)
		verify    func(*testing.T, error)
	}{
		{
			name: "successful delete",
			id:   "123",
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
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
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.DeleteReturns(errors.New("delete error"))
			},
			verify: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), "orig.Delete") {
					t.Errorf("expected error to contain %q, got %q", "orig.Delete", err.Error())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			logger := zap.NewNop()
			mockStore := &memcachedtesting.FakeTaskStore{}
			tt.setupMock(mockStore)
			task := taskMemcached.NewTask(client, mockStore, logger)

			err := task.Delete(t.Context(), tt.id)
			tt.verify(t, err)
		})
	}
}

func TestTask_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		params    internal.UpdateParams
		setupMock func(*memcachedtesting.FakeTaskStore)
		verify    func(*testing.T, error)
	}{
		{
			name: "successful update",
			id:   "123",
			params: internal.UpdateParams{
				Description: internal.ValueToPointer("updated"),
			},
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.UpdateReturns(nil)
				m.FindReturns(internal.Task{ID: "123", Description: "updated"}, nil)
			},
			verify: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "update error",
			id:   "123",
			params: internal.UpdateParams{
				Description: internal.ValueToPointer("updated"),
			},
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.UpdateReturns(errors.New("update error"))
			},
			verify: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), "orig.Update") {
					t.Errorf("expected error to contain %q, got %q", "orig.Update", err.Error())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			logger := zap.NewNop()
			mockStore := &memcachedtesting.FakeTaskStore{}
			tt.setupMock(mockStore)
			task := taskMemcached.NewTask(client, mockStore, logger)

			err := task.Update(t.Context(), tt.id, tt.params)
			tt.verify(t, err)
		})
	}
}

// TestTask_Find is skipped because it requires a running memcached instance.
// The Find method attempts to get from cache first, and without a running
// memcached server, it returns a connection error, making it difficult to
// unit test without infrastructure. The method implements Cache-Aside pattern:
// 1. Try to get from cache
// 2. On cache miss, fetch from origin store
// 3. Cache the result
//
// Integration tests with testcontainers would be more appropriate.
