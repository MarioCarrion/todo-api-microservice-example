package internal

// ValueToPointer converts a value of type T to a pointer of type *T.
func ValueToPointer[T any](value T) *T {
	return &value
}

// PointerToValue converts a pointer of type *T to the value of type T it points to.
// It returns the zero value of T if the pointer is nil.
func PointerToValue[T any](ptr *T) T { //nolint: ireturn
	if ptr == nil {
		var zero T

		return zero
	}

	return *ptr
}
