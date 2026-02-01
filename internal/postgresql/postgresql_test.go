package postgresql

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/postgresql/db"
)

func Test_convertPriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  db.Priority
		verify func(t *testing.T, priority internal.Priority, err error)
	}{
		{
			name:  "PriorityNone",
			input: db.PriorityNone,
			verify: func(t *testing.T, priority internal.Priority, err error) {
				t.Helper()

				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if priority != internal.PriorityNone {
					t.Errorf("expected %v, got %v", internal.PriorityNone, priority)
				}
			},
		},
		{
			name:  "PriorityLow",
			input: db.PriorityLow,
			verify: func(t *testing.T, priority internal.Priority, err error) {
				t.Helper()

				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if priority != internal.PriorityLow {
					t.Errorf("expected %v, got %v", internal.PriorityLow, priority)
				}
			},
		},
		{
			name:  "PriorityMedium",
			input: db.PriorityMedium,
			verify: func(t *testing.T, priority internal.Priority, err error) {
				t.Helper()

				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if priority != internal.PriorityMedium {
					t.Errorf("expected %v, got %v", internal.PriorityMedium, priority)
				}
			},
		},
		{
			name:  "PriorityHigh",
			input: db.PriorityHigh,
			verify: func(t *testing.T, priority internal.Priority, err error) {
				t.Helper()

				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if priority != internal.PriorityHigh {
					t.Errorf("expected %v, got %v", internal.PriorityHigh, priority)
				}
			},
		},
		{
			name:  "error",
			input: db.Priority("invalid"),
			verify: func(t *testing.T, priority internal.Priority, err error) {
				t.Helper()

				if err == nil {
					t.Errorf("expected error, got nothing")
				}

				if priority != internal.Priority(-1) {
					t.Errorf("expected %v, got %v", internal.Priority(-1), priority)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			priority, err := convertPriority(tt.input)
			tt.verify(t, priority, err)
		})
	}
}

func Test_newTimestamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *time.Time
		expected pgtype.Timestamp
	}{
		{
			name:     "nil",
			input:    nil,
			expected: pgtype.Timestamp{Valid: false},
		},
		{
			name:     "zero",
			input:    &time.Time{},
			expected: pgtype.Timestamp{Valid: false},
		},
		{
			name:  "truncated by minute",
			input: internal.ValueToPointer(time.Date(2026, 1, 1, 3, 2, 0, 0, time.UTC)),
			expected: pgtype.Timestamp{
				Time:  time.Date(2026, 1, 1, 3, 2, 0, 0, time.UTC),
				Valid: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if result := newTimestamp(tt.input); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func Test_newPriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *internal.Priority
		expected db.Priority
	}{
		{
			name:     "nil",
			input:    nil,
			expected: db.PriorityNone,
		},
		{
			name:     "PriorityNone",
			input:    internal.ValueToPointer(internal.PriorityNone),
			expected: db.PriorityNone,
		},
		{
			name:     "PriorityLow",
			input:    internal.ValueToPointer(internal.PriorityLow),
			expected: db.PriorityLow,
		},
		{
			name:     "PriorityMedium",
			input:    internal.ValueToPointer(internal.PriorityMedium),
			expected: db.PriorityMedium,
		},
		{
			name:     "PriorityHigh",
			input:    internal.ValueToPointer(internal.PriorityHigh),
			expected: db.PriorityHigh,
		},
		{
			name:     "invalid",
			input:    internal.ValueToPointer(internal.Priority(-1)),
			expected: db.Priority("invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if result := newPriority(tt.input); result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
