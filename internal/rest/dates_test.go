package rest_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/rest"
)

func TestNewDates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  internal.Dates
		output rest.Dates
	}{
		{
			"OK",
			internal.Dates{
				Start: internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
				Due:   internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour)),
			},
			rest.Dates{
				Start: internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
				Due:   internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes := rest.NewDates(tt.input)
			if !cmp.Equal(tt.output, actualRes, cmpopts.IgnoreUnexported(time.Time{})) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output, actualRes, cmpopts.IgnoreUnexported(time.Time{})))
			}
		})
	}
}

func TestDates_ToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  rest.Dates
		output internal.Dates
	}{
		{
			"OK",
			rest.Dates{
				Start: internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
				Due:   internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour)),
			},
			internal.Dates{
				Start: internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
				Due:   internal.ValueToPointer(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes := tt.input.ToDomain()
			if !cmp.Equal(tt.output, actualRes, cmpopts.IgnoreUnexported(time.Time{})) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output, actualRes, cmpopts.IgnoreUnexported(time.Time{})))
			}
		})
	}
}
