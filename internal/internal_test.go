package internal_test

import (
	"testing"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

func Test_PointerToValue(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		if val := internal.PointerToValue[int](nil); val != 0 {
			t.Errorf("expected zero value, got %v", val)
		}
	})

	t.Run("non nil", func(t *testing.T) {
		t.Parallel()

		value := 42

		if val := internal.PointerToValue(&value); val != 42 {
			t.Errorf("expected zero value, got %v", val)
		}
	})
}
