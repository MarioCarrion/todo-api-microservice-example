package memcached_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	taskMemcached "github.com/MarioCarrion/todo-api-microservice-example/internal/memcached"
)

// mockTaskStore is a mock implementation of TaskStore for testing
type mockTaskStore struct {
	createFn func(ctx context.Context, params internal.CreateParams) (internal.Task, error)
	deleteFn func(ctx context.Context, id string) error
	findFn   func(ctx context.Context, id string) (internal.Task, error)
	updateFn func(ctx context.Context, id string, params internal.UpdateParams) error
}

func (m *mockTaskStore) Create(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
	if m.createFn != nil {
		return m.createFn(ctx, params)
	}
	return internal.Task{}, nil
}

func (m *mockTaskStore) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockTaskStore) Find(ctx context.Context, id string) (internal.Task, error) {
	if m.findFn != nil {
		return m.findFn(ctx, id)
	}
	return internal.Task{}, nil
}

func (m *mockTaskStore) Update(ctx context.Context, id string, params internal.UpdateParams) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, params)
	}
	return nil
}

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
			store := &mockTaskStore{}
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
		mockStore      *mockTaskStore
		expectedTask   internal.Task
		expectedErrMsg string
	}{
		{
			name: "successful create",
			params: internal.CreateParams{
				Description: "test task",
				Priority:    internal.PriorityHigh.Pointer(),
			},
			mockStore: &mockTaskStore{
				createFn: func(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
					return internal.Task{
						ID:          "123",
						Description: params.Description,
						Priority:    params.Priority,
					}, nil
				},
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
			mockStore: &mockTaskStore{
				createFn: func(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
					return internal.Task{}, errors.New("create error")
				},
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
			task := taskMemcached.NewTask(client, tt.mockStore, logger)

			result, err := task.Create(context.Background(), tt.params)

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
		mockStore      *mockTaskStore
		expectedErrMsg string
	}{
		{
			name: "successful delete",
			id:   "123",
			mockStore: &mockTaskStore{
				deleteFn: func(ctx context.Context, id string) error {
					return nil
				},
			},
			expectedErrMsg: "",
		},
		{
			name: "store error",
			id:   "123",
			mockStore: &mockTaskStore{
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
			logger := zap.NewNop()
			task := taskMemcached.NewTask(client, tt.mockStore, logger)

			err := task.Delete(context.Background(), tt.id)

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
		mockStore      *mockTaskStore
		expectedErrMsg string
	}{
		{
			name: "successful update",
			id:   "123",
			params: internal.UpdateParams{
				Description: internal.ValueToPointer("updated"),
			},
			mockStore: &mockTaskStore{
				updateFn: func(ctx context.Context, id string, params internal.UpdateParams) error {
					return nil
				},
				findFn: func(ctx context.Context, id string) (internal.Task, error) {
					return internal.Task{ID: id, Description: "updated"}, nil
				},
			},
			expectedErrMsg: "",
		},
		{
			name: "update error",
			id:   "123",
			params: internal.UpdateParams{
				Description: internal.ValueToPointer("updated"),
			},
			mockStore: &mockTaskStore{
				updateFn: func(ctx context.Context, id string, params internal.UpdateParams) error {
					return errors.New("update error")
				},
			},
			expectedErrMsg: "orig.Update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := memcache.New("localhost:11211")
			logger := zap.NewNop()
			task := taskMemcached.NewTask(client, tt.mockStore, logger)

			err := task.Update(context.Background(), tt.id, tt.params)

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


