package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/service"
)

// mockTaskRepository is a mock implementation of TaskRepository for testing
type mockTaskRepository struct {
	createFn func(ctx context.Context, params internal.CreateParams) (internal.Task, error)
	deleteFn func(ctx context.Context, id string) error
	findFn   func(ctx context.Context, id string) (internal.Task, error)
	updateFn func(ctx context.Context, id string, params internal.UpdateParams) error
}

func (m *mockTaskRepository) Create(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
	if m.createFn != nil {
		return m.createFn(ctx, params)
	}
	return internal.Task{}, nil
}

func (m *mockTaskRepository) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockTaskRepository) Find(ctx context.Context, id string) (internal.Task, error) {
	if m.findFn != nil {
		return m.findFn(ctx, id)
	}
	return internal.Task{}, nil
}

func (m *mockTaskRepository) Update(ctx context.Context, id string, params internal.UpdateParams) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, params)
	}
	return nil
}

// mockTaskSearchRepository is a mock implementation of TaskSearchRepository
type mockTaskSearchRepository struct {
	searchFn func(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error)
}

func (m *mockTaskSearchRepository) Search(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
	if m.searchFn != nil {
		return m.searchFn(ctx, args)
	}
	return internal.SearchResults{}, nil
}

// mockTaskMessageBrokerPublisher is a mock implementation of TaskMessageBrokerPublisher
type mockTaskMessageBrokerPublisher struct {
	createdFn func(ctx context.Context, task internal.Task) error
	deletedFn func(ctx context.Context, id string) error
	updatedFn func(ctx context.Context, task internal.Task) error
}

func (m *mockTaskMessageBrokerPublisher) Created(ctx context.Context, task internal.Task) error {
	if m.createdFn != nil {
		return m.createdFn(ctx, task)
	}
	return nil
}

func (m *mockTaskMessageBrokerPublisher) Deleted(ctx context.Context, id string) error {
	if m.deletedFn != nil {
		return m.deletedFn(ctx, id)
	}
	return nil
}

func (m *mockTaskMessageBrokerPublisher) Updated(ctx context.Context, task internal.Task) error {
	if m.updatedFn != nil {
		return m.updatedFn(ctx, task)
	}
	return nil
}

func TestTask_Create(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name           string
		params         internal.CreateParams
		mockRepo       *mockTaskRepository
		mockMsgBroker  *mockTaskMessageBrokerPublisher
		expectedTask   internal.Task
		expectedErrMsg string
	}{
		{
			name: "successful create",
			params: internal.CreateParams{
				Description: "test task",
				Priority:    internal.PriorityHigh.Pointer(),
			},
			mockRepo: &mockTaskRepository{
				createFn: func(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
					return internal.Task{
						ID:          "123",
						Description: params.Description,
						Priority:    params.Priority,
					}, nil
				},
			},
			mockMsgBroker: &mockTaskMessageBrokerPublisher{},
			expectedTask: internal.Task{
				ID:          "123",
				Description: "test task",
				Priority:    internal.PriorityHigh.Pointer(),
			},
			expectedErrMsg: "",
		},
		{
			name: "validation error",
			params: internal.CreateParams{
				Description: "", // Invalid - empty description
			},
			mockRepo:       &mockTaskRepository{},
			mockMsgBroker:  &mockTaskMessageBrokerPublisher{},
			expectedTask:   internal.Task{},
			expectedErrMsg: "params.Validate",
		},
		{
			name: "repository error",
			params: internal.CreateParams{
				Description: "test task",
				Priority:    internal.PriorityHigh.Pointer(),
			},
			mockRepo: &mockTaskRepository{
				createFn: func(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
					return internal.Task{}, errors.New("database error")
				},
			},
			mockMsgBroker:  &mockTaskMessageBrokerPublisher{},
			expectedTask:   internal.Task{},
			expectedErrMsg: "repo.Create",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewTask(logger, tt.mockRepo, &mockTaskSearchRepository{}, tt.mockMsgBroker)

			task, err := svc.Create(t.Context(), tt.params)

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
				if diff := cmp.Diff(tt.expectedTask, task); diff != "" {
					t.Errorf("task mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestTask_Delete(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name           string
		id             string
		mockRepo       *mockTaskRepository
		mockMsgBroker  *mockTaskMessageBrokerPublisher
		expectedErrMsg string
	}{
		{
			name: "successful delete",
			id:   "123",
			mockRepo: &mockTaskRepository{
				deleteFn: func(ctx context.Context, id string) error {
					return nil
				},
			},
			mockMsgBroker:  &mockTaskMessageBrokerPublisher{},
			expectedErrMsg: "",
		},
		{
			name: "repository error",
			id:   "123",
			mockRepo: &mockTaskRepository{
				deleteFn: func(ctx context.Context, id string) error {
					return errors.New("database error")
				},
			},
			mockMsgBroker:  &mockTaskMessageBrokerPublisher{},
			expectedErrMsg: "Delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewTask(logger, tt.mockRepo, &mockTaskSearchRepository{}, tt.mockMsgBroker)

			err := svc.Delete(t.Context(), tt.id)

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

func TestTask_ByID(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name           string
		id             string
		mockRepo       *mockTaskRepository
		expectedTask   internal.Task
		expectedErrMsg string
	}{
		{
			name: "successful find",
			id:   "123",
			mockRepo: &mockTaskRepository{
				findFn: func(ctx context.Context, id string) (internal.Task, error) {
					return internal.Task{
						ID:          id,
						Description: "test task",
					}, nil
				},
			},
			expectedTask: internal.Task{
				ID:          "123",
				Description: "test task",
			},
			expectedErrMsg: "",
		},
		{
			name: "repository error",
			id:   "123",
			mockRepo: &mockTaskRepository{
				findFn: func(ctx context.Context, id string) (internal.Task, error) {
					return internal.Task{}, errors.New("not found")
				},
			},
			expectedTask:   internal.Task{},
			expectedErrMsg: "Find",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewTask(logger, tt.mockRepo, &mockTaskSearchRepository{}, &mockTaskMessageBrokerPublisher{})

			task, err := svc.ByID(t.Context(), tt.id)

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
				if diff := cmp.Diff(tt.expectedTask, task); diff != "" {
					t.Errorf("task mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestTask_Update(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name           string
		id             string
		params         internal.UpdateParams
		mockRepo       *mockTaskRepository
		mockMsgBroker  *mockTaskMessageBrokerPublisher
		expectedErrMsg string
	}{
		{
			name: "successful update",
			id:   "123",
			params: internal.UpdateParams{
				Description: internal.ValueToPointer("updated task"),
			},
			mockRepo: &mockTaskRepository{
				updateFn: func(ctx context.Context, id string, params internal.UpdateParams) error {
					return nil
				},
				findFn: func(ctx context.Context, id string) (internal.Task, error) {
					return internal.Task{ID: id, Description: "updated task"}, nil
				},
			},
			mockMsgBroker:  &mockTaskMessageBrokerPublisher{},
			expectedErrMsg: "",
		},
		{
			name: "repository update error",
			id:   "123",
			params: internal.UpdateParams{
				Description: internal.ValueToPointer("updated task"),
			},
			mockRepo: &mockTaskRepository{
				updateFn: func(ctx context.Context, id string, params internal.UpdateParams) error {
					return errors.New("database error")
				},
			},
			mockMsgBroker:  &mockTaskMessageBrokerPublisher{},
			expectedErrMsg: "repo.Update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewTask(logger, tt.mockRepo, &mockTaskSearchRepository{}, tt.mockMsgBroker)

			err := svc.Update(t.Context(), tt.id, tt.params)

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

func TestTask_By(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name           string
		params         internal.SearchParams
		mockSearch     *mockTaskSearchRepository
		expectedResult internal.SearchResults
		expectedErrMsg string
	}{
		{
			name: "successful search",
			params: internal.SearchParams{
				Description: internal.ValueToPointer("test"),
				From:        0,
				Size:        10,
			},
			mockSearch: &mockTaskSearchRepository{
				searchFn: func(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
					return internal.SearchResults{
						Tasks: []internal.Task{
							{ID: "1", Description: "test task 1"},
							{ID: "2", Description: "test task 2"},
						},
						Total: 2,
					}, nil
				},
			},
			expectedResult: internal.SearchResults{
				Tasks: []internal.Task{
					{ID: "1", Description: "test task 1"},
					{ID: "2", Description: "test task 2"},
				},
				Total: 2,
			},
			expectedErrMsg: "",
		},
		{
			name: "search error",
			params: internal.SearchParams{
				Description: internal.ValueToPointer("test"),
			},
			mockSearch: &mockTaskSearchRepository{
				searchFn: func(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
					return internal.SearchResults{}, errors.New("search error")
				},
			},
			expectedResult: internal.SearchResults{},
			expectedErrMsg: "search",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := service.NewTask(logger, &mockTaskRepository{}, tt.mockSearch, &mockTaskMessageBrokerPublisher{})

			result, err := svc.By(t.Context(), tt.params)

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
				if diff := cmp.Diff(tt.expectedResult, result); diff != "" {
					t.Errorf("result mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestNewTask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "creates new task service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop()
			repo := &mockTaskRepository{}
			search := &mockTaskSearchRepository{}
			msgBroker := &mockTaskMessageBrokerPublisher{}

			svc := service.NewTask(logger, repo, search, msgBroker)

			if svc == nil {
				t.Fatal("expected non-nil service")
			}
		})
	}
}
