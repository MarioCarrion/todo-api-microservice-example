package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MarioCarrion/todo-api/internal"
)

//go:generate counterfeiter -generate

//counterfeiter:generate -o resttesting/task_service.gen.go . TaskService

// TaskService ...
type TaskService interface {
	By(ctx context.Context, args internal.SearchParams) (internal.SearchResults, error)
	Create(ctx context.Context, params internal.CreateParams) (internal.Task, error)
	Delete(ctx context.Context, id string) error
	Task(ctx context.Context, id string) (internal.Task, error)
	Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error
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

// Register connects the handlers to the router.
func (t *TaskHandler) Register(r *gin.Engine) {
	r.POST("/tasks", t.create)
	r.GET("/tasks/:id", t.task)
	r.PUT("/tasks/:id", t.update)
	r.DELETE("/tasks/:id", t.delete)
	r.POST("/search/tasks", t.search)
}

// Task is an activity that needs to be completed within a period of time.
//
//nolint:tagliatelle
type Task struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	Dates       Dates    `json:"dates"`
	IsDone      bool     `json:"is_done"`
}

// CreateTasksRequest defines the request used for creating tasks.
type CreateTasksRequest struct {
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	Dates       Dates    `json:"dates"`
}

// CreateTasksResponse defines the response returned back after creating tasks.
type CreateTasksResponse struct {
	Task Task `json:"task"`
}

func (t *TaskHandler) create(router *gin.Context) {
	var req CreateTasksRequest
	if err := router.ShouldBindJSON(&req); err != nil {
		renderErrorResponse(router, "invalid request",
			internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))

		return
	}

	task, err := t.svc.Create(router, internal.CreateParams{
		Description: req.Description,
		Priority:    req.Priority.Convert(),
		Dates:       req.Dates.Convert(),
	})
	if err != nil {
		renderErrorResponse(router, "create failed", err)

		return
	}

	router.JSON(http.StatusCreated,
		&CreateTasksResponse{
			Task: Task{
				ID:          task.ID,
				Description: task.Description,
				Priority:    NewPriority(task.Priority),
				Dates:       NewDates(task.Dates),
			},
		})
}

func (t *TaskHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := t.svc.Delete(c, id); err != nil {
		renderErrorResponse(c, "delete failed", err)

		return
	}

	c.JSON(http.StatusOK, &struct{}{})
}

// ReadTasksResponse defines the response returned back after searching one task.
type ReadTasksResponse struct {
	Task Task `json:"task"`
}

func (t *TaskHandler) task(c *gin.Context) {
	id := c.Param("id")

	task, err := t.svc.Task(c, id)
	if err != nil {
		renderErrorResponse(c, "find failed", err)

		return
	}

	c.JSON(http.StatusOK,
		&ReadTasksResponse{
			Task: Task{
				ID:          task.ID,
				Description: task.Description,
				Priority:    NewPriority(task.Priority),
				Dates:       NewDates(task.Dates),
				IsDone:      task.IsDone,
			},
		})
}

// UpdateTasksRequest defines the request used for updating a task.
//
//nolint:tagliatelle
type UpdateTasksRequest struct {
	Description string   `json:"description"`
	IsDone      bool     `json:"is_done"`
	Priority    Priority `json:"priority"`
	Dates       Dates    `json:"dates"`
}

func (t *TaskHandler) update(router *gin.Context) {
	var req UpdateTasksRequest
	if err := router.ShouldBindJSON(&req); err != nil {
		renderErrorResponse(router, "invalid request",
			internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))

		return
	}

	id := router.Param("id")

	err := t.svc.Update(router, id, req.Description, req.Priority.Convert(), req.Dates.Convert(), req.IsDone)
	if err != nil {
		renderErrorResponse(router, "update failed", err)

		return
	}

	router.JSON(http.StatusOK, &struct{}{})
}

// SearchTasksRequest defines the request used for searching tasks.
//
//nolint:tagliatelle
type SearchTasksRequest struct {
	Description *string   `json:"description"`
	Priority    *Priority `json:"priority"`
	IsDone      *bool     `json:"is_done"`
	From        int64     `json:"from"`
	Size        int64     `json:"size"`
}

// SearchTasksResponse defines the response returned back after searching for any task.
type SearchTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Total int64  `json:"total"`
}

func (t *TaskHandler) search(router *gin.Context) {
	var req SearchTasksRequest
	if err := router.ShouldBindJSON(&req); err != nil {
		renderErrorResponse(router, "invalid request",
			internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))

		return
	}

	var priority *internal.Priority

	if req.Priority != nil {
		res := req.Priority.Convert()
		priority = &res
	}

	res, err := t.svc.By(router, internal.SearchParams{
		Description: req.Description,
		Priority:    priority,
		IsDone:      req.IsDone,
		From:        req.From,
		Size:        req.Size,
	})
	if err != nil {
		renderErrorResponse(router, "search failed", err)

		return
	}

	tasks := make([]Task, len(res.Tasks))

	for i, task := range res.Tasks {
		tasks[i].ID = task.ID
		tasks[i].Description = task.Description
		tasks[i].Priority = NewPriority(task.Priority)
		tasks[i].Dates = NewDates(task.Dates)
	}

	router.JSON(http.StatusOK,
		&SearchTasksResponse{
			Tasks: tasks,
			Total: res.Total,
		})
}
