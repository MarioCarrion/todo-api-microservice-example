package internal_test

import (
	"testing"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

func TestValueToPointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		check func(t *testing.T, result any)
	}{
		{
			name:  "int value",
			value: 42,
			check: func(t *testing.T, result any) {
				t.Helper()
				ptr := result.(*int)
				if *ptr != 42 {
					t.Errorf("expected *42, got *%d", *ptr)
				}
			},
		},
		{
			name:  "string value",
			value: "test",
			check: func(t *testing.T, result any) {
				t.Helper()
				ptr := result.(*string)
				if *ptr != "test" {
					t.Errorf("expected *test, got *%s", *ptr)
				}
			},
		},
		{
			name:  "bool value",
			value: true,
			check: func(t *testing.T, result any) {
				t.Helper()
				ptr := result.(*bool)
				if *ptr != true {
					t.Errorf("expected *true, got *%t", *ptr)
				}
			},
		},
		{
			name:  "priority value",
			value: internal.PriorityHigh,
			check: func(t *testing.T, result any) {
				t.Helper()
				ptr := result.(*internal.Priority)
				if *ptr != internal.PriorityHigh {
					t.Errorf("expected *PriorityHigh, got *%v", *ptr)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch v := tt.value.(type) {
			case int:
				result := internal.ValueToPointer(v)
				tt.check(t, result)
			case string:
				result := internal.ValueToPointer(v)
				tt.check(t, result)
			case bool:
				result := internal.ValueToPointer(v)
				tt.check(t, result)
			case internal.Priority:
				result := internal.ValueToPointer(v)
				tt.check(t, result)
			}
		})
	}
}

func TestPointerToValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "non-nil int pointer",
			input:    internal.ValueToPointer(42),
			expected: 42,
		},
		{
			name:     "nil int pointer",
			input:    (*int)(nil),
			expected: 0,
		},
		{
			name:     "non-nil string pointer",
			input:    internal.ValueToPointer("test"),
			expected: "test",
		},
		{
			name:     "nil string pointer",
			input:    (*string)(nil),
			expected: "",
		},
		{
			name:     "non-nil bool pointer",
			input:    internal.ValueToPointer(true),
			expected: true,
		},
		{
			name:     "nil bool pointer",
			input:    (*bool)(nil),
			expected: false,
		},
		{
			name:     "non-nil priority pointer",
			input:    internal.ValueToPointer(internal.PriorityHigh),
			expected: internal.PriorityHigh,
		},
		{
			name:     "nil priority pointer",
			input:    (*internal.Priority)(nil),
			expected: internal.PriorityNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch ptr := tt.input.(type) {
			case *int:
				result := internal.PointerToValue(ptr)
				if result != tt.expected.(int) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			case *string:
				result := internal.PointerToValue(ptr)
				if result != tt.expected.(string) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			case *bool:
				result := internal.PointerToValue(ptr)
				if result != tt.expected.(bool) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			case *internal.Priority:
				result := internal.PointerToValue(ptr)
				if result != tt.expected.(internal.Priority) {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestPointerToValue_RoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{
			name:  "int round trip",
			value: 42,
		},
		{
			name:  "string round trip",
			value: "hello",
		},
		{
			name:  "bool round trip",
			value: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch v := tt.value.(type) {
			case int:
				ptr := internal.ValueToPointer(v)
				result := internal.PointerToValue(ptr)
				if result != v {
					t.Errorf("round trip failed: expected %v, got %v", v, result)
				}
			case string:
				ptr := internal.ValueToPointer(v)
				result := internal.PointerToValue(ptr)
				if result != v {
					t.Errorf("round trip failed: expected %v, got %v", v, result)
				}
			case bool:
				ptr := internal.ValueToPointer(v)
				result := internal.PointerToValue(ptr)
				if result != v {
					t.Errorf("round trip failed: expected %v, got %v", v, result)
				}
			}
		})
	}
}
