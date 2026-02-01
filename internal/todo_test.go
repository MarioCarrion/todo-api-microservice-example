package internal_test

import (
	"errors"
	"testing"
	"time"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

func TestPriority_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   internal.Priority
		withErr bool
	}{
		{
			"OK: PriorityNone",
			internal.PriorityNone,
			false,
		},
		{
			"OK: PriorityLow",
			internal.PriorityLow,
			false,
		},
		{
			"OK: PriorityMedium",
			internal.PriorityMedium,
			false,
		},
		{
			"OK: PriorityHigh",
			internal.PriorityHigh,
			false,
		},
		{
			"ERR: unknown value",
			internal.Priority(-1),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualErr := tt.input.Validate()
			if (actualErr != nil) != tt.withErr {
				t.Fatalf("expected error %t, got %s", tt.withErr, actualErr)
			}

			var ierr *internal.Error
			if tt.withErr && !errors.As(actualErr, &ierr) {
				t.Fatalf("expected %T error, got %T", ierr, actualErr)
			}
		})
	}
}

func TestPriority_Pointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		priority internal.Priority
	}{
		{
			name:     "PriorityNone",
			priority: internal.PriorityNone,
		},
		{
			name:     "PriorityLow",
			priority: internal.PriorityLow,
		},
		{
			name:     "PriorityMedium",
			priority: internal.PriorityMedium,
		},
		{
			name:     "PriorityHigh",
			priority: internal.PriorityHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ptr := internal.ValueToPointer(tt.priority)
			if ptr == nil {
				t.Fatal("expected non-nil pointer")
			} else if *ptr != tt.priority {
				t.Errorf("expected *%v, got *%v", tt.priority, *ptr)
			}
		})
	}
}

func TestDates_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   internal.Dates
		withErr bool
	}{
		{
			"OK: Start.IsZero",
			internal.Dates{
				Due: internal.ValueToPointer(time.Now()),
			},
			false,
		},
		{
			"OK: Due.IsZero",
			internal.Dates{
				Start: internal.ValueToPointer(time.Now()),
			},
			false,
		},
		{
			"OK: Start < Due",
			internal.Dates{
				Start: internal.ValueToPointer(time.Now()),
				Due:   internal.ValueToPointer(time.Now().Add(2 * time.Hour)),
			},
			false,
		},
		{
			"ERR: Start > Due",
			internal.Dates{
				Start: internal.ValueToPointer(time.Now().Add(2 * time.Hour)),
				Due:   internal.ValueToPointer(time.Now()),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualErr := tt.input.Validate()
			if (actualErr != nil) != tt.withErr {
				t.Fatalf("expected error %t, got %s", tt.withErr, actualErr)
			}

			var ierr *internal.Error
			if tt.withErr && !errors.As(actualErr, &ierr) {
				t.Fatalf("expected %T error, got %T", ierr, actualErr)
			}
		})
	}
}

func TestTask_Validate(t *testing.T) {
	t.Parallel()

	newDate := func(start time.Time, due time.Time) *internal.Dates {
		res := internal.Dates{
			Start: &start,
			Due:   &due,
		}

		return &res
	}

	tests := []struct {
		name    string
		input   internal.Task
		withErr bool
	}{
		{
			"OK",
			internal.Task{
				Description: "complete this microservice",
				Priority:    internal.ValueToPointer(internal.PriorityHigh),
				Dates:       newDate(time.Now(), time.Now().Add(time.Hour)),
			},
			false,
		},
		{
			"ERR: Description",
			internal.Task{
				Priority: internal.ValueToPointer(internal.PriorityHigh),
				Dates:    newDate(time.Now(), time.Now().Add(time.Hour)),
			},
			true,
		},
		{
			"ERR: Priority",
			internal.Task{
				Description: "complete this microservice",
				Priority:    internal.ValueToPointer(internal.Priority(-1)),
				Dates:       newDate(time.Now(), time.Now().Add(time.Hour)),
			},
			true,
		},
		{
			"ERR: Dates",
			internal.Task{
				Description: "complete this microservice",
				Priority:    internal.ValueToPointer(internal.PriorityHigh),
				Dates:       newDate(time.Now().Add(time.Hour), time.Now()),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualErr := tt.input.Validate()
			if (actualErr != nil) != tt.withErr {
				t.Fatalf("expected error %t, got %s", tt.withErr, actualErr)
			}

			var ierr *internal.Error
			if tt.withErr && !errors.As(actualErr, &ierr) {
				t.Fatalf("expected %T error, got %T", ierr, actualErr)
			}
		})
	}
}
