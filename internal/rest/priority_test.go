package rest_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
	"github.com/MarioCarrion/todo-api-microservice-example/internal/rest"
)

func TestNewPriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  internal.Priority
		output rest.Priority
	}{
		{
			"OK: none",
			internal.PriorityNone,
			rest.Priority("none"),
		},
		{
			"OK: low",
			internal.PriorityLow,
			rest.Priority("low"),
		},
		{
			"OK: medium",
			internal.PriorityMedium,
			rest.Priority("medium"),
		},
		{
			"OK: high",
			internal.PriorityHigh,
			rest.Priority("high"),
		},
		{
			"OK: unknonwn",
			internal.Priority(-1),
			rest.Priority("none"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes := rest.NewPriority(tt.input)

			if !cmp.Equal(tt.output, actualRes) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output, actualRes))
			}
		})
	}
}

func TestPriority_ToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  rest.Priority
		output *internal.Priority
	}{
		{
			name:  "OK: none",
			input: rest.Priority("none"),
			output: func() *internal.Priority {
				res := internal.PriorityNone

				return &res
			}(),
		},
		{
			name:  "OK: none",
			input: rest.Priority("none"),
			output: func() *internal.Priority {
				res := internal.PriorityNone

				return &res
			}(),
		},
		{
			name:  "OK: low",
			input: rest.Priority("low"),
			output: func() *internal.Priority {
				res := internal.PriorityLow

				return &res
			}(),
		},
		{
			name:  "OK: medium",
			input: rest.Priority("medium"),
			output: func() *internal.Priority {
				res := internal.PriorityMedium

				return &res
			}(),
		},
		{
			name:  "OK: high",
			input: rest.Priority("high"),
			output: func() *internal.Priority {
				res := internal.PriorityHigh

				return &res
			}(),
		},
		{
			name:  "ERR",
			input: rest.Priority("err"),
			output: func() *internal.Priority {
				res := internal.PriorityNone

				return &res
			}(),
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes := tt.input.ToDomain()

			if diff := cmp.Diff(tt.output, actualRes); diff != "" {
				t.Fatalf("expected output do not match\n%s", diff)
			}
		})
	}
}

func TestPriority_MarshalJSON(t *testing.T) {
	t.Parallel()

	type output struct {
		res     []byte
		withErr bool
	}

	tests := []struct {
		name   string
		input  rest.Priority
		output output
	}{
		{
			"OK",
			rest.Priority("none"),
			output{
				res: []byte(`"none"`),
			},
		},
		{
			"ERR",
			rest.Priority("unknown"),
			output{
				withErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes, actualErr := json.Marshal(&tt.input)

			if (actualErr != nil) != tt.output.withErr {
				t.Fatalf("expected error %t, actual %s", tt.output.withErr, actualErr)
			}

			if !cmp.Equal(tt.output.res, actualRes) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output.res, actualRes))
			}
		})
	}
}

func TestPriority_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type output struct {
		res     rest.Priority
		withErr bool
	}

	tests := []struct {
		name   string
		input  []byte
		output output
	}{
		{
			"OK",
			[]byte(`"none"`),
			output{
				res: rest.Priority("none"),
			},
		},
		{
			"ERR: conver",
			[]byte(`"unknown"`),
			output{
				withErr: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var actualRes rest.Priority

			actualErr := json.Unmarshal(tt.input, &actualRes)

			if (actualErr != nil) != tt.output.withErr {
				t.Fatalf("expected error %t, actual %s", tt.output.withErr, actualErr)
			}

			if !cmp.Equal(tt.output.res, actualRes) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output.res, actualRes))
			}
		})
	}
}
