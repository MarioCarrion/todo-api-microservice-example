package vault_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/testcontainers/testcontainers-go"
	tvault "github.com/testcontainers/testcontainers-go/modules/vault"

	"github.com/MarioCarrion/todo-api-microservice-example/internal/envvar/vault"
)

type vaultClient struct {
	Token   string
	Address string
	Client  *api.Client
}

//nolint:paralleltest,tparallel
func TestProvider_Get(t *testing.T) {
	t.Parallel()

	type output struct {
		res     string
		withErr bool
	}

	tests := []struct {
		name   string
		setup  func(*vaultClient) error
		input  string
		output output
	}{
		{
			"OK",
			func(v *vaultClient) error {
				if _, err := v.Client.Logical().Write("/secret/data/ok",
					map[string]any{
						"data": map[string]any{
							"one": "1",
							"two": "2",
						},
					}); err != nil {
					return fmt.Errorf("couldn't write: %w", err)
				}

				return nil
			},
			"/ok:one",
			output{
				res: "1",
			},
		},
		{
			"OK: cached",
			func(_ *vaultClient) error { return nil },
			"/ok:two",
			output{
				res: "2",
			},
		},
		{
			"ERR: missing key value",
			func(_ *vaultClient) error { return nil },
			"/ok",
			output{
				withErr: true,
			},
		},
		{
			"ERR: key not found in cached data",
			func(_ *vaultClient) error { return nil },
			"/ok:three",
			output{
				withErr: true,
			},
		},
		{
			"ERR: secret not found",
			func(_ *vaultClient) error { return nil },
			"/not:found",
			output{
				withErr: true,
			},
		},
		{
			"ERR: key not found in retrieved data",
			func(v *vaultClient) error {
				if _, err := v.Client.Logical().Write("/secret/data/err",
					map[string]any{
						"data": map[string]any{
							"hello": "world",
						},
					}); err != nil {
					return fmt.Errorf("couldn't write: %w", err)
				}

				return nil
			},
			"/err:something",
			output{
				withErr: true,
			},
		},
	}

	// Provider is not local to the subtest because we want to test the local caching logic

	client := newVault(t)
	provider, err := vault.New(client.Token, client.Address, "/secret")

	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { //nolint: wsl
			// Not calling t.Parallel() because vault.Provider is not goroutine safe.

			if err := tt.setup(client); err != nil {
				t.Fatalf("could not set up %s", err)
			}

			actualRes, actualErr := provider.Get(tt.input)

			if actualRes != tt.output.res {
				t.Fatalf("expected %s, got %s", tt.output.res, actualRes)
			}

			if (actualErr != nil) != tt.output.withErr {
				t.Fatalf("expected error %t, got %s", tt.output.withErr, actualErr)
			}
		})
	}
}

func newVault(tb testing.TB) *vaultClient {
	tb.Helper()

	token := "myroot"

	vaultContainer, err := tvault.Run(
		tb.Context(),
		"hashicorp/vault:1.21.0",
		tvault.WithToken(token),
	)

	tb.Cleanup(func() {
		if err := testcontainers.TerminateContainer(vaultContainer); err != nil {
			tb.Logf("Failed to terminate container: %s", err)
		}
	})

	if err != nil {
		tb.Fatalf("Failed run container: %s", err)
	}

	host, err := vaultContainer.HttpHostAddress(tb.Context())
	if err != nil {
		tb.Fatalf("Failed to get host address: %s", err)
	}

	config := &api.Config{
		Address: host,
	}

	client, err := api.NewClient(config)
	if err != nil {
		tb.Fatalf("Failed to instantiate client: %s", err)
	}

	client.SetToken(token)

	_, err = client.Logical().Write("/secret/data/example",
		map[string]any{
			"data": map[string]any{
				"one": "two",
			},
		})
	if err != nil {
		tb.Fatalf("Failed to write secret: %s", err)
	}

	secret, err := client.Logical().Read("/secret/data/example")
	if err != nil {
		tb.Fatalf("Failed to read write secret: %s", err)
	}

	if secret == nil {
		tb.Fatalf("Secret was nil: %s", err)
	}

	return &vaultClient{
		Client:  client,
		Token:   token,
		Address: host,
	}
}
