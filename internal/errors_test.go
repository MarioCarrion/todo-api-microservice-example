package internal_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

func TestWrapErrorf(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("original error")

	tests := []struct {
		name         string
		orig         error
		code         internal.ErrorCode
		format       string
		args         []any
		expectedMsg  string
		expectedCode internal.ErrorCode
	}{
		{
			name:         "wrap error with message",
			orig:         originalErr,
			code:         internal.ErrorCodeNotFound,
			format:       "failed to find: %s",
			args:         []any{"resource"},
			expectedMsg:  "failed to find: resource: original error",
			expectedCode: internal.ErrorCodeNotFound,
		},
		{
			name:         "wrap error with multiple args",
			orig:         originalErr,
			code:         internal.ErrorCodeInvalidArgument,
			format:       "validation failed for %s on field %s",
			args:         []any{"Task", "description"},
			expectedMsg:  "validation failed for Task on field description: original error",
			expectedCode: internal.ErrorCodeInvalidArgument,
		},
		{
			name:         "wrap nil error",
			orig:         nil,
			code:         internal.ErrorCodeUnknown,
			format:       "some operation failed",
			args:         nil,
			expectedMsg:  "some operation failed",
			expectedCode: internal.ErrorCodeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := internal.WrapErrorf(tt.orig, tt.code, tt.format, tt.args...)

			if err.Error() != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, err.Error())
			}

			var ierr *internal.Error
			if !errors.As(err, &ierr) {
				t.Fatalf("expected error to be *internal.Error, got %T", err)
			}

			if ierr.Code() != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, ierr.Code())
			}

			if tt.orig != nil {
				if !errors.Is(err, tt.orig) {
					t.Errorf("expected error to wrap original error")
				}

				if !errors.Is(ierr, tt.orig) {
					t.Errorf("expected Unwrap() to return original error")
				}
			}
		})
	}
}

func TestNewErrorf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		code         internal.ErrorCode
		format       string
		args         []any
		expectedMsg  string
		expectedCode internal.ErrorCode
	}{
		{
			name:         "create error with not found code",
			code:         internal.ErrorCodeNotFound,
			format:       "task %s not found",
			args:         []any{"123"},
			expectedMsg:  "task 123 not found",
			expectedCode: internal.ErrorCodeNotFound,
		},
		{
			name:         "create error with invalid argument code",
			code:         internal.ErrorCodeInvalidArgument,
			format:       "invalid value: %d",
			args:         []any{42},
			expectedMsg:  "invalid value: 42",
			expectedCode: internal.ErrorCodeInvalidArgument,
		},
		{
			name:         "create error without args",
			code:         internal.ErrorCodeUnknown,
			format:       "unknown error occurred",
			args:         nil,
			expectedMsg:  "unknown error occurred",
			expectedCode: internal.ErrorCodeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := internal.NewErrorf(tt.code, tt.format, tt.args...)

			if err.Error() != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, err.Error())
			}

			var ierr *internal.Error
			if !errors.As(err, &ierr) {
				t.Fatalf("expected error to be *internal.Error, got %T", err)
			}

			if ierr.Code() != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, ierr.Code())
			}

			if ierr.Unwrap() != nil {
				t.Errorf("expected Unwrap() to return nil, got %v", ierr.Unwrap())
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		err         error
		expectedMsg string
	}{
		{
			name:        "error without wrapped error",
			err:         internal.NewErrorf(internal.ErrorCodeNotFound, "not found"),
			expectedMsg: "not found",
		},
		{
			name:        "error with wrapped error",
			err:         internal.WrapErrorf(fmt.Errorf("database error"), internal.ErrorCodeUnknown, "operation failed"),
			expectedMsg: "operation failed: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err.Error() != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, tt.err.Error())
			}
		})
	}
}

func TestErrorCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code internal.ErrorCode
	}{
		{
			name: "ErrorCodeUnknown",
			code: internal.ErrorCodeUnknown,
		},
		{
			name: "ErrorCodeNotFound",
			code: internal.ErrorCodeNotFound,
		},
		{
			name: "ErrorCodeInvalidArgument",
			code: internal.ErrorCodeInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := internal.NewErrorf(tt.code, "test error")
			var ierr *internal.Error
			if !errors.As(err, &ierr) {
				t.Fatalf("expected error to be *internal.Error")
			}

			if ierr.Code() != tt.code {
				t.Errorf("expected code %d, got %d", tt.code, ierr.Code())
			}
		})
	}
}
