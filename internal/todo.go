// Package internal defines the types used to create Tasks and their corresponding attributes.
package internal

import (
	"time"

	"github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// PriorityNone indicates the task needs to be prioritized.
	PriorityNone Priority = iota

	// PriorityLow indicates a non urgent task.
	PriorityLow

	// PriorityMedium indicates a task that should be completed soon.
	PriorityMedium

	// PriorityHigh indicates an urgent task that must be completed as soon as possible.
	PriorityHigh
)

// Priority indicates how important a Task is.
type Priority int8

// Validate ...
func (p Priority) Validate() error {
	switch p {
	case PriorityNone, PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	}

	return NewErrorf(ErrorCodeInvalidArgument, "unknown value")
}

// Pointer returns the pointer of p
func (p Priority) Pointer() *Priority {
	return &p
}

// Category is human readable value meant to be used to organize your tasks. Category values are unique.
type Category string

// Dates indicates a point in time where a task starts or completes, dates are not enforced on Tasks.
type Dates struct {
	Start *time.Time
	Due   *time.Time
}

// Validate ...
func (d Dates) Validate() error {
	if d.Start == nil || d.Due == nil {
		return nil
	}

	if !d.Start.IsZero() && !d.Due.IsZero() && d.Start.After(*d.Due) {
		return NewErrorf(ErrorCodeInvalidArgument, "start dates should be before end date")
	}

	return nil
}

// Pointer returns the pointer of d
func (d Dates) Pointer() *Dates {
	return &d
}

// Task is an activity that needs to be completed within a period of time.
type Task struct {
	ID          uuid.UUID
	IsDone      bool
	Priority    *Priority
	Description string
	Dates       *Dates
	SubTasks    []Task
	Categories  []Category
}

// Validate ...
func (t Task) Validate() error {
	if err := validation.ValidateStruct(&t,
		validation.Field(&t.Description, validation.Required),
		validation.Field(&t.Priority),
		validation.Field(&t.Dates),
	); err != nil {
		return WrapErrorf(err, ErrorCodeInvalidArgument, "invalid values")
	}

	return nil
}
