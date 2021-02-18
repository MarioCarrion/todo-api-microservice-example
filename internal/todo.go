// Package internal defines the types used to create Tasks and their corresponding attributes.
package internal

import (
	"errors"
	"fmt"
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

	return errors.New("unknown value")
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
		return errors.New("start dates should be before end date")
	}

	return nil
}

// Task is an activity that needs to be completed within a period of time.
type Task struct {
	ID          string
	Description string
	Priority    Priority
	Dates       Dates
	SubTasks    []Task
	Categories  []Category
	IsDone      bool
}

// Validate ...
func (t Task) Validate() error {
	if t.Description == "" {
		return errors.New("description is required")
	}

	if err := t.Priority.Validate(); err != nil {
		return fmt.Errorf("priority is invalid: %w", err)
	}

	if err := t.Dates.Validate(); err != nil {
		return fmt.Errorf("dates are invalid: %w", err)
	}

	return nil
}
