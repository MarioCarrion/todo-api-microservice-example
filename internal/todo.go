// Package internal defines the types used to create Tasks and their corresponding attributes.
package internal

import (
	"time"
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

// Category is human readable value meant to be used to organize your tasks. Category values are unique.
type Category string

// Dates indicates a point in time where a task starts or completes, dates are not enforced on Tasks.
type Dates struct {
	Start time.Time
	Due   time.Time
}

// Validate ...
func (d Dates) Validate() error {
	if !d.Start.IsZero() && !d.Due.IsZero() && d.Start.After(d.Due) {
		return NewErrorf(ErrorCodeInvalidArgument, "start dates should be before end date")
	}

	return nil
}

// Task is an activity that needs to be completed within a period of time.
type Task struct {
	IsDone      bool
	Priority    Priority
	ID          string
	Description string
	Dates       Dates
	SubTasks    []Task
	Categories  []Category
}

// Validate ...
func (t Task) Validate() error {
	if t.Description == "" {
		return NewErrorf(ErrorCodeInvalidArgument, "description is required")
	}

	if err := t.Priority.Validate(); err != nil {
		return WrapErrorf(err, ErrorCodeInvalidArgument, "priority is invalid")
	}

	if err := t.Dates.Validate(); err != nil {
		return WrapErrorf(err, ErrorCodeInvalidArgument, "dates are invalid")
	}

	return nil
}
