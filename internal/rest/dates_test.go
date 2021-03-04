package rest_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MarioCarrion/todo-api/internal"
	"github.com/MarioCarrion/todo-api/internal/rest"
)

func TestDates_Marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  rest.Dates
		output []byte
	}{
		{
			"OK",
			rest.Dates{
				Start: rest.Time(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
			},
			[]byte(`{"start":"2009-11-10T23:00:00Z","due":"0001-01-01T00:00:00Z"}`),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes, actualErr := json.Marshal(&tt.input)
			if actualErr != nil {
				t.Fatalf("expected no error, got %s", actualErr)
			}

			if !cmp.Equal(tt.output, actualRes) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output, actualRes))
			}
		})
	}
}

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
				Start: time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC),
				Due:   time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour),
			},
			rest.Dates{
				Start: rest.Time(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
				Due:   rest.Time(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour)),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes := rest.NewDates(tt.input)
			if !cmp.Equal(tt.output, actualRes, cmpopts.IgnoreUnexported(rest.Time{})) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output, actualRes, cmpopts.IgnoreUnexported(rest.Time{})))
			}
		})
	}
}

func TestDates_Convert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  rest.Dates
		output internal.Dates
	}{
		{
			"OK",
			rest.Dates{
				Start: rest.Time(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
				Due:   rest.Time(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour)),
			},
			internal.Dates{
				Start: time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC),
				Due:   time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC).Add(time.Hour),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualRes := tt.input.Convert()
			if !cmp.Equal(tt.output, actualRes, cmpopts.IgnoreUnexported(time.Time{})) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output, actualRes, cmpopts.IgnoreUnexported(time.Time{})))
			}
		})
	}
}

func TestTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	type output struct {
		res     []byte
		withErr bool
	}

	tests := []struct {
		name   string
		input  rest.Time
		output output
	}{
		{
			"OK",
			rest.Time(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
			output{
				res: []byte(`"2009-11-10T23:00:00Z"`),
			},
		},
	}

	for _, tt := range tests {
		tt := tt

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

func TestTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type output struct {
		res     rest.Time
		withErr bool
	}

	tests := []struct {
		name   string
		input  []byte
		output output
	}{
		{
			"OK",
			[]byte(`"2009-11-10T23:00:00Z"`),
			output{
				res: rest.Time(time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC)),
			},
		},
		{
			"ERR",
			[]byte(`2009-`),
			output{
				withErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var actualRes rest.Time

			actualErr := json.Unmarshal(tt.input, &actualRes)

			if (actualErr != nil) != tt.output.withErr {
				t.Fatalf("expected error %t, actual %s", tt.output.withErr, actualErr)
			}

			if !cmp.Equal(tt.output.res, actualRes, cmpopts.IgnoreUnexported(rest.Time{})) {
				t.Fatalf("expected output do not match\n%s", cmp.Diff(tt.output.res, actualRes, cmpopts.IgnoreUnexported(rest.Time{})))
			}
		})
	}
}
