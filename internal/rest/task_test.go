package rest_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/rest"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/rest/resttesting"
)

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

			svc := &resttesting.FakeTaskService{}
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
		setupMock    func(*resttesting.FakeTaskService)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.CreateReturns(internal.Task{
					ID:          taskID.String(),
					Description: "test task",
					Priority:    internal.ValueToPointer(internal.PriorityHigh),
				}, nil)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.CreateReturns(internal.Task{}, errors.New("service error"))
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

			mockService := &resttesting.FakeTaskService{}
			tt.setupMock(mockService)
			handler := rest.NewTaskHandler(mockService)
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
		setupMock    func(*resttesting.FakeTaskService)
		expectError  bool
		validateResp func(t *testing.T, resp rest.ReadTaskResponseObject)
	}{
		{
			name: "successful read",
			request: rest.ReadTaskRequestObject{
				Id: taskID,
			},
			setupMock: func(m *resttesting.FakeTaskService) {
				m.ByIDReturns(internal.Task{
					ID:          taskID.String(),
					Description: "test task",
				}, nil)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.ByIDReturns(internal.Task{}, errors.New("not found"))
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

			mockService := &resttesting.FakeTaskService{}
			tt.setupMock(mockService)
			handler := rest.NewTaskHandler(mockService)
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
		setupMock    func(*resttesting.FakeTaskService)
		expectError  bool
		validateResp func(t *testing.T, resp rest.DeleteTaskResponseObject)
	}{
		{
			name: "successful delete",
			request: rest.DeleteTaskRequestObject{
				Id: taskID,
			},
			setupMock: func(m *resttesting.FakeTaskService) {
				m.DeleteReturns(nil)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.DeleteReturns(errors.New("delete error"))
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

			mockService := &resttesting.FakeTaskService{}
			tt.setupMock(mockService)
			handler := rest.NewTaskHandler(mockService)
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
		setupMock    func(*resttesting.FakeTaskService)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.UpdateReturns(nil)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.UpdateReturns(errors.New("update error"))
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

			mockService := &resttesting.FakeTaskService{}
			tt.setupMock(mockService)
			handler := rest.NewTaskHandler(mockService)
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
		setupMock    func(*resttesting.FakeTaskService)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.ByReturns(internal.SearchResults{
					Tasks: []internal.Task{
						{ID: taskID1.String(), Description: "test task 1"},
						{ID: taskID2.String(), Description: "test task 2"},
					},
					Total: 2,
				}, nil)
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
			setupMock: func(m *resttesting.FakeTaskService) {
				m.ByReturns(internal.SearchResults{}, errors.New("search error"))
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

			mockService := &resttesting.FakeTaskService{}
			tt.setupMock(mockService)
			handler := rest.NewTaskHandler(mockService)
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
