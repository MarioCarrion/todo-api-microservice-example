package envvar_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/MarioCarrion/todo-api/internal/envvar"
	"github.com/MarioCarrion/todo-api/internal/envvar/envvartesting"
)

func TestConfiguration_Get(t *testing.T) {
	t.Parallel()

	type output struct {
		val     string
		withErr bool
	}

	tests := []struct {
		name   string
		setup  func(p *envvartesting.FakeProvider) (teardown func())
		input  string
		output output
		arg    string
	}{
		{
			"OK: no secret",
			func(_ *envvartesting.FakeProvider) func() {
				return func() {
					os.Setenv("ENVVAR_OK", "")
				}
			},
			"ENVVAR_OK",
			output{
				val: "value",
			},
			"",
		},
		{
			"OK: secure",
			func(p *envvartesting.FakeProvider) func() {
				os.Setenv("ENVVAR_OK1_SECURE", "/secret/value")

				p.GetReturns("provider value", nil)

				return func() {
					os.Setenv("ENVVAR_OK1", "")
					os.Setenv("ENVVAR_OK1_SECURE", "")
				}
			},
			"ENVVAR_OK1",
			output{
				val: "provider value",
			},
			"/secret/value",
		},
		{
			"ERR: provider failed",
			func(p *envvartesting.FakeProvider) func() {
				os.Setenv("ENVVAR_ERR_SECURE", "/failed")

				p.GetReturns("", errors.New("failed"))

				return func() {
					os.Setenv("ENVVAR_ERR_SECURE", "")
				}
			},
			"ENVVAR_ERR",
			output{
				withErr: true,
			},
			"/failed",
		},
	}

	_ = envvar.Load(path.Join("fixtures", "env"))

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			provider := envvartesting.FakeProvider{}

			teardown := tt.setup(&provider)
			t.Cleanup(teardown)

			actual, actualErr := envvar.New(&provider).Get(tt.input)

			if (actualErr != nil) != tt.output.withErr {
				t.Fatalf("expected error %t, got %s", tt.output.withErr, actualErr)
			}

			if !cmp.Equal(tt.output.val, actual) {
				t.Fatalf("expected result does not match: %s", cmp.Diff(tt.output.val, actual))
			}

			//- Provider Args

			if provider.GetCallCount() > 0 {
				if arg := provider.GetArgsForCall(0); arg != tt.arg {
					t.Fatalf("expected arg %s, got %s", arg, tt.arg)
				}
			}
		})
	}
}
