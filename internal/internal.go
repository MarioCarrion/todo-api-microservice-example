package internal

// PointerToValue converts a pointer of type *T to the value of type T it points to.
// It returns the zero value of T if the pointer is nil.
func PointerToValue[T any](ptr *T) T { //nolint: ireturn
	if ptr == nil {
		var zero T

		return zero
	}

	return *ptr
}
