package internal_test

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

func TestCreateParams_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   internal.CreateParams
		withErr bool
	}{
		{
			"OK",
			internal.CreateParams{
				Description: "Description",
				Priority:    internal.ValueToPointer(internal.PriorityLow),
			},
			false,
		},
		{
			"ERR",
			internal.CreateParams{},
			true,
		},
		{
			"ERR",
			internal.CreateParams{},
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

			var ierr validation.Errors
			if tt.withErr && !errors.As(actualErr, &ierr) {
				t.Fatalf("expected %T error, got %T", ierr, actualErr)
			}
		})
	}
}

func TestSearchParams_IsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  internal.SearchParams
		output bool
	}{
		{
			"OK",
			internal.SearchParams{
				Description: internal.ValueToPointer("description"),
				Priority:    internal.ValueToPointer(internal.PriorityHigh),
				IsDone:      internal.ValueToPointer(false),
			},
			false,
		},
		{
			"OK: Description",
			internal.SearchParams{
				Priority: internal.ValueToPointer(internal.PriorityHigh),
				IsDone:   internal.ValueToPointer(false),
			},
			false,
		},
		{
			"OK: Priority",
			internal.SearchParams{
				Description: internal.ValueToPointer("description"),
				IsDone:      internal.ValueToPointer(false),
			},
			false,
		},
		{
			"OK: IsDone",
			internal.SearchParams{
				Description: internal.ValueToPointer("description"),
				Priority:    internal.ValueToPointer(internal.PriorityHigh),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if actual := tt.input.IsZero(); actual != tt.output {
				t.Fatalf("expected %t, got %t", tt.output, actual)
			}
		})
	}
}
