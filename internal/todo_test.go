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

			ptr := tt.priority.Pointer()
			if ptr == nil {
				t.Fatal("expected non-nil pointer")
			} else if *ptr != tt.priority {
				t.Errorf("expected *%v, got *%v", tt.priority, *ptr)
			}
		})
	}
}

func TestPriority_ValueOrDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *internal.Priority
		expected internal.Priority
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: internal.PriorityNone,
		},
		{
			name:     "PriorityNone pointer",
			input:    internal.PriorityNone.Pointer(),
			expected: internal.PriorityNone,
		},
		{
			name:     "PriorityLow pointer",
			input:    internal.PriorityLow.Pointer(),
			expected: internal.PriorityLow,
		},
		{
			name:     "PriorityMedium pointer",
			input:    internal.PriorityMedium.Pointer(),
			expected: internal.PriorityMedium,
		},
		{
			name:     "PriorityHigh pointer",
			input:    internal.PriorityHigh.Pointer(),
			expected: internal.PriorityHigh,
		},
		{
			name:     "invalid priority pointer returns PriorityNone",
			input:    internal.Priority(-1).Pointer(),
			expected: internal.PriorityNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.input.ValueOrDefault()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
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

	tests := []struct {
		name    string
		input   internal.Task
		withErr bool
	}{
		{
			"OK",
			internal.Task{
				Description: "complete this microservice",
				Priority:    internal.PriorityHigh.Pointer(),
				Dates: internal.Dates{
					Start: internal.ValueToPointer(time.Now()),
					Due:   internal.ValueToPointer(time.Now().Add(time.Hour)),
				}.Pointer(),
			},
			false,
		},
		{
			"ERR: Description",
			internal.Task{
				Priority: internal.PriorityHigh.Pointer(),
				Dates: internal.Dates{
					Start: internal.ValueToPointer(time.Now()),
					Due:   internal.ValueToPointer(time.Now().Add(time.Hour)),
				}.Pointer(),
			},
			true,
		},
		{
			"ERR: Priority",
			internal.Task{
				Description: "complete this microservice",
				Priority:    internal.Priority(-1).Pointer(),
				Dates: internal.Dates{
					Start: internal.ValueToPointer(time.Now()),
					Due:   internal.ValueToPointer(time.Now().Add(time.Hour)),
				}.Pointer(),
			},
			true,
		},
		{
			"ERR: Dates",
			internal.Task{
				Description: "complete this microservice",
				Priority:    internal.PriorityHigh.Pointer(),
				Dates: internal.Dates{
					Start: internal.ValueToPointer(time.Now().Add(time.Hour)),
					Due:   internal.ValueToPointer(time.Now()),
				}.Pointer(),
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
