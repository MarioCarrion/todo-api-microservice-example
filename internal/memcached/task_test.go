package memcached_test

import (
	"errors"
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
		name           string
		params         internal.CreateParams
		setupMock      func(*memcachedtesting.FakeTaskStore)
		expectedTask   internal.Task
		expectedErrMsg string
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
			expectedTask: internal.Task{
				ID:          "123",
				Description: "test task",
				Priority:    internal.PriorityHigh.Pointer(),
			},
			expectedErrMsg: "",
		},
		{
			name: "store error",
			params: internal.CreateParams{
				Description: "test task",
			},
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.CreateReturns(internal.Task{}, errors.New("create error"))
			},
			expectedTask:   internal.Task{},
			expectedErrMsg: "orig.Create",
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
				if diff := cmp.Diff(tt.expectedTask, result); diff != "" {
					t.Errorf("task mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestTask_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*memcachedtesting.FakeTaskStore)
		expectedErrMsg string
	}{
		{
			name: "successful delete",
			id:   "123",
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.DeleteReturns(nil)
			},
			expectedErrMsg: "",
		},
		{
			name: "store error",
			id:   "123",
			setupMock: func(m *memcachedtesting.FakeTaskStore) {
				m.DeleteReturns(errors.New("delete error"))
			},
			expectedErrMsg: "orig.Delete",
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

func TestTask_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		id             string
		params         internal.UpdateParams
		setupMock      func(*memcachedtesting.FakeTaskStore)
		expectedErrMsg string
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
			expectedErrMsg: "",
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
			expectedErrMsg: "orig.Update",
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

// TestTask_Find is skipped because it requires a running memcached instance.
// The Find method attempts to get from cache first, and without a running
// memcached server, it returns a connection error, making it difficult to
// unit test without infrastructure. The method implements Cache-Aside pattern:
// 1. Try to get from cache
// 2. On cache miss, fetch from origin store
// 3. Cache the result
//
// Integration tests with testcontainers would be more appropriate.
