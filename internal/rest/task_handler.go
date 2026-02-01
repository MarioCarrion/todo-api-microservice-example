package rest

import (
	"context"

	"github.com/google/uuid"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

//go:generate counterfeiter -generate

//counterfeiter:generate -o resttesting/task_service.gen.go . TaskService

// TaskService ...
type TaskService interface {
	By(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error)
	Create(ctx context.Context, params internal.CreateParams) (internal.Task, error)
	Delete(ctx context.Context, id string) error
	ByID(ctx context.Context, id string) (internal.Task, error)
	Update(ctx context.Context, id string, args internal.UpdateParams) error
}

// TaskHandler ...
type TaskHandler struct {
	svc TaskService
}

// NewTaskHandler ...
func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{
		svc: svc,
	}
}

func (t *TaskHandler) CreateTask(ctx context.Context, req CreateTaskRequestObject) (CreateTaskResponseObject, error) {
	var priority *internal.Priority
	if req.Body.Priority != nil {
		priority = req.Body.Priority.ToDomain()
	}

	var dates *internal.Dates
	if req.Body.Dates != nil {
		dates = internal.ValueToPointer(req.Body.Dates.ToDomain())
	}

	task, err := t.svc.Create(ctx, internal.CreateParams{
		Description: req.Body.Description,
		Priority:    priority,
		Dates:       dates,
	})

	// TODO: determine if this a validation error or a different kind of error, and use "CreateTask400JSONResponse"
	if err != nil {
		return CreateTask500JSONResponse{ //nolint: nilerr
			Error: err.Error(),
		}, nil
	}

	id, err := uuid.Parse(task.ID)
	if err != nil {
		return CreateTask500JSONResponse{ //nolint: nilerr
			Error: err.Error(),
		}, nil
	}

	resp := CreateTask201JSONResponse{}
	resp.Task = Task{
		ID:          id,
		Description: task.Description,
	}

	if task.Dates != nil {
		resp.Task.Dates = internal.ValueToPointer(NewDates(*task.Dates))
	}

	if task.Priority != nil {
		resp.Task.Priority = internal.ValueToPointer(NewPriority(*task.Priority))
	}

	return resp, nil
}

func (t *TaskHandler) DeleteTask(ctx context.Context, request DeleteTaskRequestObject) (DeleteTaskResponseObject, error) {
	if err := t.svc.Delete(ctx, request.Id.String()); err != nil {
		resp := DeleteTask500JSONResponse{}
		resp.Error = err.Error()

		return resp, nil //nolint: nilerr
	}

	// TODO: Consider "DeleteTask404Response"

	return DeleteTask200Response{}, nil
}

func (t *TaskHandler) ReadTask(ctx context.Context, request ReadTaskRequestObject) (ReadTaskResponseObject, error) {
	task, err := t.svc.ByID(ctx, request.Id.String())
	// TODO: determine if this is a validation error or a different kind of error, and use "CreateTask400JSONResponse"
	if err != nil {
		resp := ReadTask500JSONResponse{}
		resp.Error = err.Error()

		return resp, nil //nolint: nilerr
	}

	id, err := uuid.Parse(task.ID)
	if err != nil {
		resp := ReadTask500JSONResponse{}
		resp.Error = err.Error()

		return resp, nil //nolint: nilerr
	}

	resp := ReadTask200JSONResponse{}
	resp.Task = &Task{
		ID:          id,
		Description: task.Description,
	}

	if task.Dates != nil {
		resp.Task.Dates = internal.ValueToPointer(NewDates(*task.Dates))
	}

	if task.Priority != nil {
		resp.Task.Priority = internal.ValueToPointer(NewPriority(*task.Priority))
	}

	return resp, nil
}

func (t *TaskHandler) UpdateTask(ctx context.Context, req UpdateTaskRequestObject) (UpdateTaskResponseObject, error) {
	var priority *internal.Priority
	if req.Body.Priority != nil {
		priority = req.Body.Priority.ToDomain()
	}

	var dates *internal.Dates
	if req.Body.Dates != nil {
		dates = internal.ValueToPointer(req.Body.Dates.ToDomain())
	}

	if err := t.svc.Update(ctx, req.Id.String(), internal.UpdateParams{
		Description: req.Body.Description,
		Priority:    priority,
		Dates:       dates,
		IsDone:      req.Body.IsDone,
	}); err != nil {
		// TODO: determine if this is a validation error or a different kind of error, and use "CreateTask400JSONResponse"
		return UpdateTask500JSONResponse{Error: err.Error()}, nil //nolint: nilerr
	}

	return UpdateTask200Response{}, nil
}

func (t *TaskHandler) SearchTask(ctx context.Context, req SearchTaskRequestObject) (SearchTaskResponseObject, error) {
	var priority *internal.Priority

	if req.Body.Priority != nil {
		priority = req.Body.Priority.ToDomain()
	}

	res, err := t.svc.By(ctx, internal.SearchParams{
		Description: req.Body.Description,
		Priority:    priority,
		IsDone:      req.Body.IsDone,
		From:        req.Body.From,
		Size:        req.Body.Size,
	})
	// TODO: determine if this is a validation error or a different kind of error, and use "CreateTask400JSONResponse"
	if err != nil {
		return SearchTask500JSONResponse{ //nolint: nilerr
			Error: err.Error(),
		}, nil
	}

	tasks := make([]Task, len(res.Tasks))

	for i, task := range res.Tasks { //nolint: varnamelen
		id, err := uuid.Parse(task.ID)
		if err != nil {
			return SearchTask500JSONResponse{ //nolint: nilerr
				Error: err.Error(),
			}, nil
		}

		tasks[i].ID = id
		tasks[i].Description = task.Description

		if task.Priority != nil {
			tasks[i].Priority = internal.ValueToPointer(NewPriority(*task.Priority))
		}

		if task.Dates != nil {
			tasks[i].Dates = internal.ValueToPointer(NewDates(*task.Dates))
		}

		tasks[i].IsDone = &task.IsDone
	}

	resp := SearchTask200JSONResponse{}
	resp.Tasks = &tasks

	return resp, nil
}
