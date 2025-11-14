package rest_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/rest"
)

// mockTaskService is a mock implementation of TaskService for testing
type mockTaskService struct {
	byFn     func(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error)
	createFn func(ctx context.Context, params internal.CreateParams) (internal.Task, error)
	deleteFn func(ctx context.Context, id string) error
	byIDFn   func(ctx context.Context, id string) (internal.Task, error)
	updateFn func(ctx context.Context, id string, args internal.UpdateParams) error
}

func (m *mockTaskService) By(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
	if m.byFn != nil {
		return m.byFn(ctx, args)
	}
	return internal.SearchResults{}, nil
}

func (m *mockTaskService) Create(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
	if m.createFn != nil {
		return m.createFn(ctx, params)
	}
	return internal.Task{}, nil
}

func (m *mockTaskService) Delete(ctx context.Context, id string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockTaskService) ByID(ctx context.Context, id string) (internal.Task, error) {
	if m.byIDFn != nil {
		return m.byIDFn(ctx, id)
	}
	return internal.Task{}, nil
}

func (m *mockTaskService) Update(ctx context.Context, id string, args internal.UpdateParams) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, args)
	}
	return nil
}

func TestNewTaskHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "creates new task handler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := &mockTaskService{}
			handler := rest.NewTaskHandler(svc)

			if handler == nil {
				t.Fatal("expected non-nil handler")
			}
		})
	}
}

func TestTaskHandler_CreateTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New()

	tests := []struct {
		name         string
		request      rest.CreateTaskRequestObject
		mockService  *mockTaskService
		expectError  bool
		validateResp func(t *testing.T, resp rest.CreateTaskResponseObject)
	}{
		{
			name: "successful creation",
			request: rest.CreateTaskRequestObject{
				Body: &rest.CreateTaskJSONRequestBody{
					Description: "test task",
					Priority:    (*rest.Priority)(internal.ValueToPointer("high")),
				},
			},
			mockService: &mockTaskService{
				createFn: func(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
					return internal.Task{
						ID:          taskID.String(),
						Description: params.Description,
						Priority:    params.Priority,
					}, nil
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.CreateTaskResponseObject) {
				t.Helper()
				r, ok := resp.(rest.CreateTask201JSONResponse)
				if !ok {
					t.Fatalf("expected CreateTask201JSONResponse, got %T", resp)
				}
				if r.Task.ID != taskID {
					t.Errorf("expected task ID %v, got %v", taskID, r.Task.ID)
				}
			},
		},
		{
			name: "service error",
			request: rest.CreateTaskRequestObject{
				Body: &rest.CreateTaskJSONRequestBody{
					Description: "test task",
				},
			},
			mockService: &mockTaskService{
				createFn: func(ctx context.Context, params internal.CreateParams) (internal.Task, error) {
					return internal.Task{}, errors.New("service error")
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.CreateTaskResponseObject) {
				t.Helper()
				_, ok := resp.(rest.CreateTask500JSONResponse)
				if !ok {
					t.Fatalf("expected CreateTask500JSONResponse, got %T", resp)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := rest.NewTaskHandler(tt.mockService)
			resp, err := handler.CreateTask(t.Context(), tt.request)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, resp)
			}
		})
	}
}

func TestTaskHandler_ReadTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New()

	tests := []struct {
		name         string
		request      rest.ReadTaskRequestObject
		mockService  *mockTaskService
		expectError  bool
		validateResp func(t *testing.T, resp rest.ReadTaskResponseObject)
	}{
		{
			name: "successful read",
			request: rest.ReadTaskRequestObject{
				Id: taskID,
			},
			mockService: &mockTaskService{
				byIDFn: func(ctx context.Context, id string) (internal.Task, error) {
					return internal.Task{
						ID:          taskID.String(),
						Description: "test task",
					}, nil
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.ReadTaskResponseObject) {
				t.Helper()
				r, ok := resp.(rest.ReadTask200JSONResponse)
				if !ok {
					t.Fatalf("expected ReadTask200JSONResponse, got %T", resp)
				}
				if r.Task.ID != taskID {
					t.Errorf("expected task ID %v, got %v", taskID, r.Task.ID)
				}
			},
		},
		{
			name: "service error",
			request: rest.ReadTaskRequestObject{
				Id: taskID,
			},
			mockService: &mockTaskService{
				byIDFn: func(ctx context.Context, id string) (internal.Task, error) {
					return internal.Task{}, errors.New("not found")
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.ReadTaskResponseObject) {
				t.Helper()
				_, ok := resp.(rest.ReadTask500JSONResponse)
				if !ok {
					t.Fatalf("expected ReadTask500JSONResponse, got %T", resp)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := rest.NewTaskHandler(tt.mockService)
			resp, err := handler.ReadTask(t.Context(), tt.request)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, resp)
			}
		})
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New()

	tests := []struct {
		name         string
		request      rest.DeleteTaskRequestObject
		mockService  *mockTaskService
		expectError  bool
		validateResp func(t *testing.T, resp rest.DeleteTaskResponseObject)
	}{
		{
			name: "successful delete",
			request: rest.DeleteTaskRequestObject{
				Id: taskID,
			},
			mockService: &mockTaskService{
				deleteFn: func(ctx context.Context, id string) error {
					return nil
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.DeleteTaskResponseObject) {
				t.Helper()
				_, ok := resp.(rest.DeleteTask200Response)
				if !ok {
					t.Fatalf("expected DeleteTask200Response, got %T", resp)
				}
			},
		},
		{
			name: "service error",
			request: rest.DeleteTaskRequestObject{
				Id: taskID,
			},
			mockService: &mockTaskService{
				deleteFn: func(ctx context.Context, id string) error {
					return errors.New("delete error")
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.DeleteTaskResponseObject) {
				t.Helper()
				_, ok := resp.(rest.DeleteTask500JSONResponse)
				if !ok {
					t.Fatalf("expected DeleteTask500JSONResponse, got %T", resp)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := rest.NewTaskHandler(tt.mockService)
			resp, err := handler.DeleteTask(t.Context(), tt.request)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, resp)
			}
		})
	}
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	t.Parallel()

	taskID := uuid.New()

	tests := []struct {
		name         string
		request      rest.UpdateTaskRequestObject
		mockService  *mockTaskService
		expectError  bool
		validateResp func(t *testing.T, resp rest.UpdateTaskResponseObject)
	}{
		{
			name: "successful update",
			request: rest.UpdateTaskRequestObject{
				Id: taskID,
				Body: &rest.UpdateTaskJSONRequestBody{
					Description: internal.ValueToPointer("updated task"),
				},
			},
			mockService: &mockTaskService{
				updateFn: func(ctx context.Context, id string, args internal.UpdateParams) error {
					return nil
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.UpdateTaskResponseObject) {
				t.Helper()
				_, ok := resp.(rest.UpdateTask200Response)
				if !ok {
					t.Fatalf("expected UpdateTask200Response, got %T", resp)
				}
			},
		},
		{
			name: "service error",
			request: rest.UpdateTaskRequestObject{
				Id: taskID,
				Body: &rest.UpdateTaskJSONRequestBody{
					Description: internal.ValueToPointer("updated task"),
				},
			},
			mockService: &mockTaskService{
				updateFn: func(ctx context.Context, id string, args internal.UpdateParams) error {
					return errors.New("update error")
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.UpdateTaskResponseObject) {
				t.Helper()
				_, ok := resp.(rest.UpdateTask500JSONResponse)
				if !ok {
					t.Fatalf("expected UpdateTask500JSONResponse, got %T", resp)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := rest.NewTaskHandler(tt.mockService)
			resp, err := handler.UpdateTask(t.Context(), tt.request)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, resp)
			}
		})
	}
}

func TestTaskHandler_SearchTask(t *testing.T) {
	t.Parallel()

	taskID1 := uuid.New()
	taskID2 := uuid.New()

	tests := []struct {
		name         string
		request      rest.SearchTaskRequestObject
		mockService  *mockTaskService
		expectError  bool
		validateResp func(t *testing.T, resp rest.SearchTaskResponseObject)
	}{
		{
			name: "successful search",
			request: rest.SearchTaskRequestObject{
				Body: &rest.SearchTaskJSONRequestBody{
					Description: internal.ValueToPointer("test"),
					From:        0,
					Size:        10,
				},
			},
			mockService: &mockTaskService{
				byFn: func(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
					return internal.SearchResults{
						Tasks: []internal.Task{
							{ID: taskID1.String(), Description: "test task 1"},
							{ID: taskID2.String(), Description: "test task 2"},
						},
						Total: 2,
					}, nil
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.SearchTaskResponseObject) {
				t.Helper()
				r, ok := resp.(rest.SearchTask200JSONResponse)
				if !ok {
					t.Fatalf("expected SearchTask200JSONResponse, got %T", resp)
				}
				if r.Tasks == nil || len(*r.Tasks) != 2 {
					t.Errorf("expected 2 tasks, got %v", r.Tasks)
				}
			},
		},
		{
			name: "service error",
			request: rest.SearchTaskRequestObject{
				Body: &rest.SearchTaskJSONRequestBody{
					Description: internal.ValueToPointer("test"),
				},
			},
			mockService: &mockTaskService{
				byFn: func(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error) {
					return internal.SearchResults{}, errors.New("search error")
				},
			},
			expectError: false,
			validateResp: func(t *testing.T, resp rest.SearchTaskResponseObject) {
				t.Helper()
				_, ok := resp.(rest.SearchTask500JSONResponse)
				if !ok {
					t.Fatalf("expected SearchTask500JSONResponse, got %T", resp)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := rest.NewTaskHandler(tt.mockService)
			resp, err := handler.SearchTask(t.Context(), tt.request)

			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, resp)
			}
		})
	}
}
